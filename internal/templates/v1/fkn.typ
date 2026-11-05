// Adapted FKN layout. See FKN-LICENSE for attribution.
#let render-fkn(d, plain, body) = {
  let m = d.metadata
  let l = d.labels
  let faculty-name = plain(m.faculty)
  let program-name = plain(m.program)
  let current-year = plain(m.year)
  let un(n) = "_" * n
  let doc-code = plain(d.document_code)
  let doc-title = plain(d.title)
  let approval-page-code = plain(d.approval_code)

  let storage-table = {
    set text(size: 10pt)
    place(
      bottom + left,
      dx: 5mm,
      dy: -10mm,
      rotate(
        -90deg,
        reflow: true,
        table(
          columns: (25mm, 35mm, 25mm, 25mm, 35mm),
          rows: (5mm, 7mm),
          align: center,
          [ #plain(l.inventory.at(0)) ], [ #plain(l.inventory.at(1)) ], [ #plain(l.inventory.at(2)) ], [ #plain(l.inventory.at(3)) ], [ #plain(l.inventory.at(1)) ],
        ),
      ),
    )
  }

  let approval-page = {
    let top-banner = [
      #set par(spacing: 0.65em)

      #text(weight: "bold", plain(m.organization))

      #faculty-name

      #program-name
    ]

    let approve-table = grid(
      columns: (1fr, 1fr),
      inset: (
        x: 10mm,
        y: 3mm,
      ),
      align: center,

      [
        #plain(l.agreed)

        #plain(m.supervisor_position)
      ],
      [
        #plain(l.approved)

        #plain(m.approver_position)
      ],

      box[
        #plain(l.signature_todo) #plain(m.supervisor)

        #plain(l.date_todo)
      ],
      [
        #plain(l.signature_todo) #plain(m.approver)

        #plain(l.date_todo)
      ],
    )

    let center-banner = [
      #set text(size: 14pt, weight: "bold")
      #set par(spacing: 2em)

      #par(spacing: 0.65em, plain(upper(m.project)))

      #doc-title

      #plain(l.approval_sheet)

      #approval-page-code
    ]

    let student-info = align(right)[
      #set par(spacing: 1em)

      #plain(l.contributors)

      #plain(l.student_group) #plain(m.group)

      #plain(l.signature_todo) / #plain(m.author) /

      #plain(l.date_todo)
    ]

    let bottom-banner = [
      #set text(weight: "bold")

      #current-year
    ]

    page(
      header: none,
      footer: none,
      margin: (
        left: 20mm,
        right: 10mm,
        top: 25mm,
        bottom: 15mm,
      ),
      foreground: storage-table,
    )[
      #set align(center)

      #grid(
        columns: (1fr),
        row-gutter: 1fr,
        top-banner,
        approve-table,
        center-banner,
        student-info,
        bottom-banner,
      )
    ]

    counter(page).update(1)
  }

  let title-page = {
    let top-banner = [
      #set align(left)
      #set par(spacing: 2em)

      #plain(l.approved_document)

      #approval-page-code
    ]

    let center-banner = [
      #set text(size: 14pt, weight: "bold")
      #set par(spacing: 2em)

      #par(spacing: 0.65em, plain(upper(m.project)))

      #doc-title

      #doc-code

      #plain(l.sheets) #context { counter(page).final().at(0) }
    ]

    let bottom-banner = [
      #set text(weight: "bold")

      #current-year
    ]

    page(
      header: none,
      footer: none,
      margin: (
        left: 20mm,
        right: 10mm,
        top: 25mm,
        bottom: 15mm,
      ),
      foreground: storage-table,
    )[
      #set align(center)

      #grid(
        columns: (1fr),
        row-gutter: 1fr,
        top-banner,
        center-banner,
        bottom-banner,
      )
    ]
  }

  let list-registration-changes = {
    let list-name = [
      #set text(size: 14pt, weight: "bold")
      #set par(spacing: 2em)
      #plain(l.changes)
    ]

    let changes-table = figure(
      table(
        columns: (10mm,) + (15mm,) * 4 + (20mm,) * 2 + (auto,) + (auto,) * 2,
        rows: (auto,) * 3 + (9.5mm,) * 21,

        table.cell(colspan: 10, align: horizon, [#plain(l.changes)]),
        table.cell(colspan: 5, align: horizon, [#plain(l.page_numbers)]),
        table.cell(rowspan: 2, align: horizon, [#plain(l.total_pages)]),
        table.cell(rowspan: 2, align: horizon, [#plain(l.document_number)]),
        table.cell(rowspan: 2, align: horizon, [#plain(l.incoming_document)]),
        table.cell(rowspan: 2, align: horizon, [#plain(l.signature_short)]),
        table.cell(rowspan: 2, align: horizon, [#plain(l.date)]),
        table.cell(align: horizon, rotate(-90deg, reflow: true)[#plain(l.change)]),
        table.cell(align: horizon, rotate(-90deg, reflow: true)[#plain(l.changed)]),
        table.cell(align: horizon, rotate(-90deg, reflow: true)[#plain(l.replaced)]),
        table.cell(align: horizon, rotate(-90deg, reflow: true)[#plain(l.new_pages)]),
        table.cell(align: horizon, rotate(-90deg, reflow: true)[#plain(l.cancelled_pages)]),
        ..range(210).map(_ => []),
      ),
    )

    page(
      header: none,
      footer: none,
      margin: (
        left: 20mm,
        right: 10mm,
        top: 25mm,
        bottom: 15mm,
      ),
    )[
      #set align(center)

      #grid(
        columns: (1fr),
        row-gutter: 1fr,
        list-name,
        changes-table,
      )
    ]
  }

  let outline-block = {
    pagebreak(weak: true)

    {
      set align(center)
      set text(weight: "bold")

      [#plain(l.contents)]
    }

    outline(
      title: none,
      indent: 5mm,
    )
  }

  let normal-pages = {
    set page(
      margin: (
        top: 25mm,
        left: 20mm,
        right: 10mm,
        bottom: 47mm,
      ),
      header: [
        #set align(center)
        #set text(weight: "bold")

        #context counter(page).display()

        #doc-code
      ],
      footer: [
        #set text(size: 8pt)
        #set par(first-line-indent: 0pt)
        #table(
          columns: (2fr, 1fr, 1fr, 1fr, 1fr),
          align: center,
          rows: 7mm,

          [], [], [], [], [],
          [#plain(l.change)], [#plain(l.sheet)], [#plain(l.document_number)], [#plain(l.signature_short)], [#plain(l.date)],
          doc-code, [], [], [], [],
          [#plain(l.inventory_original)], [ #plain(l.inventory.at(1)) ], [#plain(l.inventory_replacement)], [#plain(l.inventory_duplicate)], [ #plain(l.inventory.at(1)) ],
        )
      ],
    )

    set par(
      first-line-indent: (
        amount: 2em,
        all: true,
      ),
      justify: true,
      spacing: 0.65em,
      leading: 0.65em,
    )

    set list(
      indent: 2em,
      spacing: 0.65em,
      marker: "-",
    )

    set enum(
      indent: 2em,
      spacing: 0.65em,
    )

    set heading(numbering: "1.")

    show heading.where(level: 1): h => {
      set align(center)
      set text(
        weight: "bold",
        size: 12pt,
      )

      pagebreak(weak: true)
      if h.numbering != none [
        #counter(heading).display(h.numbering) #h.body
      ] else [
        #h.body
      ]
    }

    show heading.where(level: 2): h => {
      set text(
        weight: "bold",
        size: 12pt,
      )

      block(inset: (left: 1em))[#if h.numbering != none {counter(heading).display()} #h.body]
    }

    show heading.where(level: 3): h => {
      set text(weight: "bold", size: 12pt)

      block(inset: (left: 3em))[#if h.numbering != none {counter(heading).display()} #h.body]
    }

    pagebreak(weak: true)
    align(center, text(weight: "bold", size: 12pt, [#plain(l.annotation)]))
    plain(d.annotation)
    parbreak()
    plain(l.review_warning)

    outline-block

    body
  }

  set text(
    lang: d.language,
    size: 12pt,
    font: "Liberation Serif",
  )

  approval-page
  title-page
  normal-pages
  list-registration-changes
}
