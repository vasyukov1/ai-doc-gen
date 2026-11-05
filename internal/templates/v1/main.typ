#import "rich.typ": rich
#import "fkn.typ": render-fkn
#let d = json("document.json")
#let l = d.labels
#let m = d.metadata
#let plain(value) = {
  show regex("(?i)TODO[^\\n]*"): text.with(fill: rgb("#c62828"), weight: "bold")
  text(value)
}
#let field(key) = plain(m.at(key))
#set document(title: d.title, author: m.author, date: none)
#set text(font: "Liberation Serif", size: 12pt, lang: d.language)
#set page(paper: "a4", margin: (left: 20mm, right: 10mm, top: 25mm, bottom: 25mm))
#set par(justify: true, spacing: 0.7em, leading: 0.55em)
#let title = [#text(size: 14pt, weight: "bold", plain(d.title))]
#let body = [
#for s in d.sections {
 heading(level: s.level, numbering: if s.unnumbered {none} else {"1.1.1"}, plain(s.heading))
 if s.at("blocks", default: none) != none { rich(s.blocks, plain) } else if s.body != "" { for paragraph in s.body.split("\n\n") { plain(paragraph); parbreak() } }
 if s.at("rich_tables", default: none) != none { for t in s.rich_tables {
  if t.caption.len() > 0 { block(strong(rich(t.caption, plain))) }
  table(columns: t.columns.len(), inset: 5pt,
    table.header(..t.columns.map(x => strong(rich(x, plain)))),
    ..t.rows.map(row => row.map(x => rich(x, plain))).flatten())
  parbreak()
 } } else if s.at("tables", default: none) != none { for t in s.tables {
  if t.caption != "" { block(text(weight: "bold", plain(t.caption))) }
  table(columns: t.columns.len(), inset: 5pt,
    table.header(..t.columns.map(x => text(weight: "bold", plain(x)))),
    ..t.rows.flatten().map(x => plain(x)))
  parbreak()
 } }
 for warning in s.missing { block(inset: 5pt, stroke: 0.4pt + rgb("#c62828"))[#plain(warning)] }
}
#heading(level: 1, numbering: none, l.specification)
#for paragraph in d.specification.split("\n\n") { plain(paragraph); parbreak() }

]
#if d.academic {
 render-fkn(d, plain, body)
} else {

  page()[
    #align(center)[#text(size: 20pt, weight: "bold", field("project")) #parbreak() #title]
    #v(1fr)
    #l.contributors #field("author") #parbreak()
    #field("organization") #parbreak() #field("year")
  ]

 set page(header: align(center)[#context counter(page).display() #h(8mm) #plain(d.document_code)])
 set heading(numbering: "1.1.1")
 show heading.where(level: 1): it => { pagebreak(weak: true); block(above: 12pt, below: 12pt, text(weight: "bold", it)) }
 heading(level: 1, numbering: none, outlined: false, l.annotation)
 plain(d.annotation)
 parbreak()
 plain(l.review_warning)
 outline(title: l.contents, indent: 5mm, depth: 2)
 body
}
