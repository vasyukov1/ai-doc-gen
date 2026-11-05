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
