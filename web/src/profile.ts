const profiles: Record<string, [string, string, string?]> = {
  "client-handover": ["Client handover: scope, installation, operating workflows, acceptance checks and maintenance. Add recipient constraints to tailor the draft.", "Передача заказчику: состав поставки, установка, работа, приёмка и сопровождение. Добавьте ограничения получателя.", "Lliurament al client: abast, instal·lació, funcionament, acceptació i manteniment. Afegiu les necessitats del destinatari."],
  tz: ["FKN technical specification: what the software must do. Source requirements are mandatory.", "ТЗ ФКН: что должна делать программа. Исходные требования обязательны."],
  pz: ["FKN explanatory note: implementation choices and rationale. Supply the source specification.", "ПЗ ФКН: реализация и обоснование решений. Приложите исходное ТЗ."],
  ro: ["FKN operator manual: setup and everyday use. Supply the source specification.", "РО ФКН: установка и использование программы. Приложите исходное ТЗ."],
  pmi: ["FKN test plan: proposed checks and acceptance procedures, not proof of successful tests.", "ПМИ ФКН: план проверок и приёмки, а не подтверждение успешных испытаний."],
  "pmi-team": ["Team FKN test plan: shared acceptance procedures. Supply the source specification.", "Командная ПМИ ФКН: общие процедуры приёмки. Приложите исходное ТЗ."],
  srs: ["SRS: structured software requirements aligned with ISO/IEC/IEEE 29148; not a certification of compliance. Source requirements improve the draft.", "SRS: спецификация требований с ориентиром на ISO/IEC/IEEE 29148, без сертификации соответствия. Исходные требования улучшат результат.", "SRS: requisits de programari estructurats amb referència a ISO/IEC/IEEE 29148, sense certificació de conformitat. Els requisits d’origen milloren l’esborrany."],
  arc42: ["arc42: architecture, building blocks, runtime scenarios, decisions and risks. Add project context for a more useful draft.", "arc42: архитектура, компоненты, сценарии, решения и риски. Добавьте контекст проекта для содержательного документа.", "arc42: arquitectura, components, escenaris d’execució, decisions i riscos. Afegiu context del projecte per obtenir un esborrany més útil."],
  "technical-overview": ["Technical overview: a compact explanation of the implementation. A short project description helps identify the intended scope.", "Технический обзор: краткое описание реализации. Короткое описание проекта поможет определить предмет документа.", "Visió tècnica: explicació compacta de la implementació. Una descripció breu ajuda a delimitar l’abast del document."],
};

export function initProfileGuidance(): void {
  const profile = document.querySelector<HTMLSelectElement>('[name="profile_id"]');
  const language = document.querySelector<HTMLSelectElement>('[name="language"]');
  const guide = document.querySelector<HTMLElement>("[data-profile-guide]");
  if (!profile || !guide) return;

  const update = () => {
    const description = profiles[profile.value];
    const catalan = language?.value === "ca";
    const supported = Boolean(description?.[2]);
    language?.setCustomValidity(catalan && !supported ? "Catalan is available for neutral technical profiles only. Choose a supported profile or English/Russian." : "");
    guide.textContent = catalan
      ? description?.[2] ?? "Aquest perfil acadèmic només està disponible en anglès i rus. Trieu un perfil tècnic neutre per generar el document en català."
      : description?.[language?.value === "ru" ? 1 : 0] ?? "";
  };
  profile.addEventListener("change", update);
  language?.addEventListener("change", update);
  update();
}

const academicFields = new Set(["faculty", "program", "group", "supervisor", "supervisor_position", "approver", "approver_position", "code"]);
const metadataHints: Record<string, string> = {
  organization: "University or client/company name; for independent work, your team name. / Вуз, компания заказчика или название вашей команды.",
  faculty: "Your faculty, for example Faculty of Computer Science. / Факультет, например «Факультет компьютерных наук».",
  program: "Your degree programme, for example Software Engineering. / Образовательная программа, например «Программная инженерия».",
  group: "Your student group identifier, for example SE-231. / Номер учебной группы, например БПИ231.",
  author: "Names as they should appear on the cover; separate team members with commas. / ФИО автора или участников команды через запятую.",
  supervisor: "Your academic supervisor's full name. / ФИО научного руководителя.",
  supervisor_position: "The supervisor's official role or academic title. / Должность или учёное звание руководителя.",
  approver: "Person who approves the document, as specified by your course. / Кто утверждает документ — уточните в требованиях курса.",
  approver_position: "The approver's official role. / Должность утверждающего документ.",
  year: "Year shown on the cover. / Год на титульном листе.",
  code: "Document/project code from your course requirements; ask your supervisor if unknown. / Шифр из требований курса; если неизвестен, уточните у руководителя.",
};

export function initAdaptiveInputs(): void {
  const form = document.querySelector<HTMLFormElement>("[data-create-document]");
  const profile = form?.querySelector<HTMLSelectElement>('[name="profile_id"]');
  if (!form || !profile) return;
  for (const [name, text] of Object.entries(metadataHints)) {
    const input = form.elements.namedItem(`metadata.${name}`) as HTMLInputElement;
    const hint = document.createElement("p");
    hint.id = `${input.id}-help`; hint.className = "form-note"; hint.textContent = text;
    input.setAttribute("aria-describedby", hint.id);
    input.after(hint);
  }
  const update = () => {
    const academic = ["tz", "pz", "ro", "pmi", "pmi-team"].includes(profile.value);
    for (const name of academicFields) {
      const input = form.elements.namedItem(`metadata.${name}`) as HTMLInputElement;
      input.disabled = !academic;
      input.closest<HTMLElement>(".field")!.hidden = !academic;
    }
    form.querySelectorAll<HTMLElement>(".form-row").forEach(row => {
      row.hidden = Array.from(row.children).every(child => (child as HTMLElement).hidden);
    });
    const details = form.querySelector<HTMLDetailsElement>("[data-spec-group]")!;
    if (academic) details.open = true;
    form.querySelector<HTMLElement>("[data-spec-label]")!.textContent = academic
      ? "Source specification / Исходное ТЗ" : "Additional requirements / Дополнительные требования";
    form.querySelector<HTMLElement>("[data-spec-help]")!.textContent = academic
      ? "Paste the actual assignment: purpose, functions, constraints and acceptance criteria. Alternatively attach UTF-8 .txt/.md. Copy readable text from PDF or Typst. / Вставьте реальное ТЗ или задание: цель, функции, ограничения и приёмка. Вместо текста можно прикрепить .txt/.md. Из PDF/Typst скопируйте текст."
      : "Optional: requirements that cannot be found in the code. Leave empty if the repository and description are sufficient. Use text OR a UTF-8 .txt/.md file. / Необязательно: требования, которых нет в коде. Если исходников и описания достаточно, оставьте пустым. Используйте текст ИЛИ файл .txt/.md.";
    const examples: Record<string, string> = {
      "client-handover": "Who receives the project, where they will run it, and what they need to maintain. / Кто принимает проект, где будет запускать и что должен сопровождать.",
      "technical-overview": "What the project does and who will read the overview. / Что делает проект и для кого нужен обзор.",
      srs: "Users, intended functions and constraints not visible in the code. / Пользователи, ожидаемые функции и ограничения, которых нет в коде.",
      arc42: "Architecture questions, deployment environment and important trade-offs. / Вопросы по архитектуре, среда развёртывания и важные компромиссы.",
    };
    form.querySelector<HTMLElement>("[data-description-help]")!.textContent = examples[profile.value]
      ?? "Briefly explain the coursework topic and what this document should emphasize. / Кратко опишите тему работы и на чём сделать акцент в документе.";
  };
  profile.addEventListener("change", update);
  update();
}
