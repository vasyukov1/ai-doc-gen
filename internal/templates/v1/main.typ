#let d = json("document.json")
#let l = d.labels
#let m = d.metadata
#let plain(value) = text(value)
#let field(key) = plain(m.at(key))
#set document(title: d.title, author: m.author, date: none)
#set text(font: "Liberation Serif", size: 12pt, lang: d.language)
#set page(paper: "a4", margin: (left: 20mm, right: 10mm, top: 25mm, bottom: 25mm))
#set par(justify: true, spacing: 0.7em, leading: 0.55em)
#let inventory = table(columns: (1fr, 1fr, 1fr, 1fr, 1fr), inset: 3pt,
  ..l.inventory.map(x => text(size: 8pt, x)), [], [], [], [], [])
#let title = [#text(size: 14pt, weight: "bold", plain(d.title))]
#if d.academic {
  page(footer: inventory)[
    #align(center)[#text(weight: "bold", field("organization")) #parbreak()
    #field("faculty") #parbreak() #field("program")]
    #v(12mm)
    #grid(columns: (1fr, 1fr), gutter: 8mm,
      [#l.agreed #parbreak() #field("supervisor_position") #parbreak() #field("supervisor") #parbreak() #l.signature],
      [#l.approved #parbreak() #field("approver_position") #parbreak() #field("approver") #parbreak() #l.signature])
    #v(18mm)
    #align(center)[#text(size: 14pt, weight: "bold", field("project")) #parbreak() #title #parbreak()
    #l.approval_sheet #parbreak() #plain(d.approval_code)]
    #v(16mm)
    #align(right)[#l.contributors #parbreak() #l.group #field("group") #parbreak() #field("author") #parbreak() #l.signature]
    #v(1fr)
    #align(center, field("year"))
  ]
  page(footer: inventory)[
    #l.approved_document #parbreak() #plain(d.approval_code)
    #v(1fr)
    #align(center)[#text(size: 14pt, weight: "bold", field("project")) #parbreak() #title #parbreak()
    #plain(d.document_code) #parbreak() #l.sheets #context counter(page).final().first()]
    #v(1fr)
    #align(center, field("year"))
  ]
} else {
  page()[
    #align(center)[#text(size: 20pt, weight: "bold", field("project")) #parbreak() #title]
    #v(1fr)
    #l.contributors #field("author") #parbreak()
    #field("organization") #parbreak() #field("year")
  ]
}
#set page(header: align(center)[#context counter(page).display() #h(8mm) #plain(d.document_code)], footer: if d.academic { inventory } else {none})
#set heading(numbering: "1.1.1")
#show heading.where(level: 1): it => { pagebreak(weak: true); block(above: 12pt, below: 12pt, text(weight: "bold", it)) }
#heading(level: 1, numbering: none, outlined: false, l.annotation)
#plain(d.annotation)
#parbreak()
#text(weight: "bold", l.review_warning)
#outline(title: l.contents, indent: 5mm, depth: 2)
#for s in d.sections {
 heading(level: s.level, numbering: if s.unnumbered {none} else {"1.1.1"}, plain(s.heading))
 if s.body != "" { for paragraph in s.body.split("\n\n") { plain(paragraph); parbreak() } }
 for warning in s.missing { block(inset: 5pt, stroke: 0.4pt + gray)[#text(weight: "bold", l.missing) #plain(warning)] }
}
#heading(level: 1, numbering: none, l.specification)
#for paragraph in d.specification.split("\n\n") { plain(paragraph); parbreak() }
#if d.academic {
 pagebreak()
 heading(level: 1, numbering: none, outlined: false, l.changes)
 table(columns: (0.7fr, 1.2fr, 1.2fr, 1.2fr, 1.2fr, 1fr, 1fr), inset: 4pt,
   ..l.change_columns.map(x => text(size: 9pt, x)),
   ..range(105).map(_ => box(height: 7mm)))
}
