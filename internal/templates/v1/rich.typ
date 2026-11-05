// Render only allowlisted structure; all leaf text stays literal.
#let rich(nodes, plain) = {
 for n in nodes {
  let content = if ("list", "ordered", "table").contains(n.kind) { [] } else { rich(n.children, plain) }
  if n.kind == "text" { plain(n.text) }
  else if n.kind == "strong" { strong(content) }
  else if n.kind == "emphasis" { emph(content) }
  else if n.kind == "strike" { strike(content) }
  else if n.kind == "code" { text(size: 9pt, content) }
  else if n.kind == "codeblock" { block(inset: 6pt, fill: luma(95%), text(size: 9pt, plain(n.text))) }
  else if n.kind == "heading" { block(above: 8pt, below: 5pt, strong(content)) }
  else if n.kind == "paragraph" { content; parbreak() }
  else if n.kind == "list" { list(..n.children.map(x => rich(x.children, plain))) }
  else if n.kind == "ordered" { enum(..n.children.map(x => rich(x.children, plain))) }
  else if n.kind == "quote" { block(inset: (left: 10pt), stroke: (left: 1pt + gray), content) }
  else if n.kind == "link" and n.url != "" { link(n.url, content) }
  else if n.kind == "table" {
   if n.children.len() > 0 {
    let cols = n.children.first().children.len()
    if cols > 0 {
     table(columns: cols, inset: 4pt,
      table.header(..n.children.first().children.map(x => strong(rich(x.children, plain)))),
      ..n.children.slice(1).map(row => row.children.map(x => rich(x.children, plain))).flatten())
     parbreak()
    }
   }
  }
  else { content }
 }
}
