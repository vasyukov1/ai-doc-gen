# Document profiles and standards

Europe does not have one universal project-documentation template. ISO/IEC/IEEE standards can guide specific document types, while universities, employers, and countries may impose their own formatting and content rules. GOST is a separate standards family; it should be an optional, explicitly versioned profile rather than described as a European standard.

| Candidate profile | Basis and intended content | MVP status |
|---|---|---|
| Software project overview | Product purpose, system context, modules, setup, evidence and limitations | Good first generic profile; no compliance claim |
| Requirements specification | [ISO/IEC/IEEE 29148:2018](https://www.iso.org/standard/72089.html) concerns requirements engineering and information items | Requires user input: code alone cannot prove original requirements |
| Architecture description | [ISO/IEC/IEEE 42010:2022](https://www.iso.org/standard/74393.html) concerns architecture descriptions and viewpoints | Candidate after stakeholder/viewpoint questions are designed |
| User guide | [ISO/IEC/IEEE 26514:2022](https://www.iso.org/standard/77451.html) covers information for users of software | Candidate after target-user/task data is available |
| Quality or test report | [ISO/IEC 25023:2016](https://www.iso.org/standard/35747.html) defines measures used with ISO/IEC 25010 | Needs actual test/measurement evidence; do not infer results |
| Academic project report | Institution and programme rubric, supplied by the user | Choose a specific institution/format before implementation |
| GOST-oriented set | Named standard numbers and editions supplied or licensed for review | Later, if validated demand warrants it |

The first selected families are an existing institution-specific academic document set, adapted for the website, and generic technical project documentation. The academic set is already available as Russian Typst templates; [template_adaptation.md](template_adaptation.md) records the five variants and required changes. Implement the technical project overview as a generic profile. Requirements, architecture, user-guide, and quality-report profiles are later extensions based on the standards above. A language choice must select a complete locale-specific template and generation policy, rather than translate only the generated paragraphs.

**Template registry proposal:** `profile_id`, `edition`, `locale`, `language`, `required_metadata`, `section_schema`, `evidence_policy`, `typst_template_version`, `validation_rules`, and `claim_level` (`generic`, `aligned`, or `verified-against-supplied-rules`). Display the exact basis and edition in the UI. Avoid copying protected standard text into templates without appropriate rights.

Typst is suitable for the editable source and PDF artifact: its [official documentation](https://typst.app/docs/reference/pdf/) describes PDF compilation. PDF standards such as PDF/A are separate output-format choices and do not establish that the document content complies with an ISO or university profile.
