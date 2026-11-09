// Explicit opt-in only: credentials, confirmations and uploaded files never enter storage.
const DRAFT_KEY = "doc-gen.creation-draft.v1";
const allowed = new Set([
  "generation_intent", "github_url", "description", "specification", "profile_id", "language",
  "metadata.project", "metadata.organization", "metadata.faculty", "metadata.program",
  "metadata.group", "metadata.author", "metadata.supervisor", "metadata.supervisor_position",
  "metadata.approver", "metadata.approver_position", "metadata.year", "metadata.code",
  "gemini_model", "groq_model", "openrouter_model", "local_model",
]);

export function initSetup(): void {
  const form = document.querySelector<HTMLFormElement>("[data-create-document]");
  if (!form) return;
  const remember = form.querySelector<HTMLInputElement>("[data-remember-draft]")!;
  const feedback = form.querySelector<HTMLElement>("[data-draft-feedback]")!;
  form.querySelector<HTMLElement>("[data-draft-settings]")!.hidden = false;
  form.querySelector<HTMLElement>("[data-purpose-shortcuts]")!.hidden = false;

  const field = (name: string) => form.elements.namedItem(name) as HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement | null;
  let restoring = false;
  const restore = () => {
    restoring = true;
    try {
      const raw = localStorage.getItem(DRAFT_KEY);
      if (!raw) {
        feedback.textContent = "No saved form on this device. / Сохранённых данных нет.";
        return;
      }
      const values: unknown = JSON.parse(raw);
      if (!values || typeof values !== "object" || Array.isArray(values)) throw Error("Invalid draft");
      for (const [name, value] of Object.entries(values)) {
        const element = allowed.has(name) ? field(name) : null;
        if (!element || typeof value !== "string" || value.length > 200000) continue;
        if (element instanceof HTMLSelectElement && !Array.from(element.options).some(option => option.value === value)) continue;
        element.value = value;
        element.dispatchEvent(new Event("change", { bubbles: true }));
      }
      form.dispatchEvent(new Event("change", { bubbles: true }));
      remember.checked = true;
      feedback.textContent = "Saved text restored; you can edit every field. Files and keys stay unchanged. / Текст восстановлен, его можно изменить. Файлы и ключи не восстанавливаются.";
    } catch {
      feedback.textContent = "Saved form could not be restored. / Не удалось восстановить форму. Можно заполнить её заново.";
    } finally {
      restoring = false;
    }
  };
  try { if (localStorage.getItem(DRAFT_KEY)) restore(); } catch { feedback.textContent = "Browser storage is unavailable. / Хранилище браузера недоступно."; }

  const save = () => {
    if (!remember.checked || restoring) return;
    const values: Record<string, string> = {};
    for (const name of allowed) {
      const element = field(name);
      if (element) values[name] = element.value;
    }
    try {
      localStorage.setItem(DRAFT_KEY, JSON.stringify(values));
      feedback.textContent = "Form saved on this device. Keys and files are not saved.";
    } catch {
      feedback.textContent = "Browser storage is unavailable or full. This form has not been saved.";
    }
  };
  const forget = () => {
    remember.checked = false;
    try {
      localStorage.removeItem(DRAFT_KEY);
      feedback.textContent = "Saved form text removed. Current inputs remain. Profile and language preferences are remembered separately.";
    } catch {
      feedback.textContent = "Could not remove saved form. Clear this site's browser storage to remove it.";
    }
  };
  remember.addEventListener("change", () => remember.checked ? save() : forget());
  form.querySelector("[data-forget-draft]")!.addEventListener("click", forget);
  form.querySelector("[data-restore-draft]")!.addEventListener("click", restore);
  form.querySelector("[data-reset-draft]")!.addEventListener("click", () => {
    remember.checked = false;
    form.reset();
    remember.checked = false;
    for (const name of allowed) {
      const element = field(name);
      if (element && !(element instanceof HTMLSelectElement) && !name.endsWith("_model")) element.value = "";
    }
    field("metadata.year")!.value = String(new Date().getFullYear());
    form.querySelectorAll<HTMLInputElement>("input[type=password], input[type=file]").forEach(input => { input.value = ""; });
    form.querySelectorAll<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>("input[name], select[name], textarea[name]").forEach(input => {
      input.dispatchEvent(new Event("change", { bubbles: true }));
    });
    form.querySelectorAll<HTMLElement>(".form-feedback").forEach(box => { box.hidden = true; });
    feedback.textContent = "Current inputs cleared. Saved text remains available through Restore. / Поля очищены. Сохранённый текст можно восстановить кнопкой выше.";
    field("profile_id")!.focus();
  });
  form.addEventListener("input", save);
  form.addEventListener("change", save);

  form.querySelectorAll<HTMLButtonElement>("[data-profile-choice]").forEach(button => {
    button.addEventListener("click", () => {
      const profile = field("profile_id")!;
      profile.value = button.dataset.profileChoice!;
      profile.dispatchEvent(new Event("change", { bubbles: true }));
    });
  });

  form.querySelectorAll<HTMLElement>("[data-model-picker]").forEach(container => {
    const input = container.querySelector<HTMLInputElement>("input")!;
    const preset = input.defaultValue;
    const label = container.querySelector("label")!;
    const picker = document.createElement("select");
    picker.id = `${input.id}_picker`;
    picker.add(new Option(input.dataset.modelLabel!, preset));
    for (const entry of (input.dataset.modelOptions ?? "").split(";")) {
      const [id, title] = entry.split("|");
      if (id && title) picker.add(new Option(title, id));
    }
    const isPreset = (value: string) => Array.from(picker.options).some(option => option.value === value && value !== "custom");
    picker.add(new Option("Other model — enter ID / Другая модель", "custom"));
    picker.value = isPreset(input.value) ? input.value : "custom";
    label.htmlFor = picker.id;
    label.textContent = "Model / Модель";
    container.insertBefore(picker, input);
    input.setAttribute("aria-label", "Custom model ID");
    const update = () => {
      const custom = picker.value === "custom";
      input.hidden = !custom;
      if (!custom) input.value = picker.value;
      if (custom) input.focus();
    };
    picker.addEventListener("change", () => { update(); save(); });
    input.addEventListener("change", () => { picker.value = isPreset(input.value) ? input.value : "custom"; input.hidden = picker.value !== "custom"; });
    input.hidden = picker.value !== "custom";
  });
}
