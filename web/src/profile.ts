const profiles: Record<string, [string, string]> = {
  tz: ["FKN technical specification: what the software must do. Source requirements are mandatory.", "ТЗ ФКН: что должна делать программа. Исходные требования обязательны."],
  pz: ["FKN explanatory note: implementation choices and rationale. Supply the source specification.", "ПЗ ФКН: реализация и обоснование решений. Приложите исходное ТЗ."],
  ro: ["FKN operator manual: setup and everyday use. Supply the source specification.", "РО ФКН: установка и использование программы. Приложите исходное ТЗ."],
  pmi: ["FKN test plan: proposed checks and acceptance procedures, not proof of successful tests.", "ПМИ ФКН: план проверок и приёмки, а не подтверждение успешных испытаний."],
  "pmi-team": ["Team FKN test plan: shared acceptance procedures. Supply the source specification.", "Командная ПМИ ФКН: общие процедуры приёмки. Приложите исходное ТЗ."],
  srs: ["SRS: structured software requirements aligned with ISO/IEC/IEEE 29148; not a certification of compliance. Source requirements improve the draft.", "SRS: спецификация требований с ориентиром на ISO/IEC/IEEE 29148, без сертификации соответствия. Исходные требования улучшат результат."],
  arc42: ["arc42: architecture, building blocks, runtime scenarios, decisions and risks. Add project context for a more useful draft.", "arc42: архитектура, компоненты, сценарии, решения и риски. Добавьте контекст проекта для содержательного документа."],
  "technical-overview": ["Technical overview: a compact explanation of the implementation. A short project description helps identify the intended scope.", "Технический обзор: краткое описание реализации. Короткое описание проекта поможет определить предмет документа."],
};

export function initProfileGuidance(): void {
  const profile = document.querySelector<HTMLSelectElement>('[name="profile_id"]');
  const language = document.querySelector<HTMLSelectElement>('[name="language"]');
  const guide = document.querySelector<HTMLElement>("[data-profile-guide]");
  if (!profile || !guide) return;

  const update = () => {
    const description = profiles[profile.value];
    guide.textContent = description?.[language?.value === "ru" ? 1 : 0] ?? "";
  };
  profile.addEventListener("change", update);
  language?.addEventListener("change", update);
  update();
}
