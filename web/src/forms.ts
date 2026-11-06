import { initProfileGuidance, initAdaptiveInputs } from "./profile";

function reveal(element: HTMLElement): void {
  let parent = element.parentElement;
  while (parent) {
    if (parent instanceof HTMLDetailsElement) parent.open = true;
    parent = parent.parentElement;
  }
}

const PROFILE_KEY = "docgen.profile";
const LANGUAGE_KEY = "docgen.language";

function stored(key: string): string {
  try {
    return localStorage.getItem(key) ?? "";
  } catch {
    return "";
  }
}

function remember(key: string, value: string): void {
  try {
    localStorage.setItem(key, value);
  } catch {
    // Private browsing may refuse storage; the choice still works for this page.
  }
}

function initFileInputs(): void {
  for (const input of document.querySelectorAll<HTMLInputElement>("[data-file-input]")) {
    const zone = input.closest<HTMLElement>(".dropzone");
    const title = zone?.querySelector<HTMLElement>("[data-file-name]");
    const empty = title?.textContent ?? "";
    const form = input.form;
    const repo = form?.elements.namedItem("github_url") as HTMLInputElement | null;
    const clear = document.createElement("button");
    clear.type = "button"; clear.className = "button button--quiet"; clear.textContent = "Remove file / Убрать файл"; clear.hidden = !input.files?.length;
    zone?.after(clear);
    const reset = (focus: boolean): void => {
      input.value = "";
      input.dispatchEvent(new Event("change", {bubbles: true}));
      if (focus) input.focus();
    };
    clear.addEventListener("click", () => reset(true));
    zone?.addEventListener("dragover", event => { event.preventDefault(); });
    zone?.addEventListener("drop", event => {
      event.preventDefault();
      if (event.dataTransfer?.files.length === 1) { input.files = event.dataTransfer.files; input.dispatchEvent(new Event("change", {bubbles:true})); }
    });
    input.addEventListener("change", () => {
      clear.hidden = !input.files?.length;
      zone?.classList.toggle("is-set", Boolean(input.files?.length));
      if (title) title.textContent = input.files?.[0]?.name ?? empty;
      // One source per submission: a chosen archive replaces a typed repository.
      if (input.name === "archive" && input.files?.length && repo?.value) {
        repo.value = "";
        repo.dispatchEvent(new Event("input", {bubbles: true}));
      }
    });
    if (input.name === "archive" && repo) {
      // Typing a repository replaces a chosen archive instead of blocking submit.
      repo.addEventListener("input", () => {
        if (repo.value.trim() && input.files?.length) reset(false);
      });
    }
  }
}

// Fewer decisions on first visit: reuse the previous profile and language,
// and derive the project title from the source.
function initCreateDefaults(): void {
  const form = document.querySelector<HTMLFormElement>("form[data-create-document]");
  if (!form) {
    return;
  }

  const profile = form.elements.namedItem("profile_id") as HTMLSelectElement | null;
  const language = form.elements.namedItem("language") as HTMLSelectElement | null;
  const repo = form.elements.namedItem("github_url") as HTMLInputElement | null;
  const archive = form.elements.namedItem("archive") as HTMLInputElement | null;
  const project = form.elements.namedItem("metadata.project") as HTMLInputElement | null;
  const year = form.elements.namedItem("metadata.year") as HTMLInputElement | null;

  if (profile) {
    const saved = stored(PROFILE_KEY);
    if (saved && Array.from(profile.options).some(option => option.value === saved)) {
      profile.value = saved;
      profile.dispatchEvent(new Event("change", {bubbles: true}));
    }
    profile.addEventListener("change", () => remember(PROFILE_KEY, profile.value));
  }

  if (language) {
    // The user chooses the language; only their own earlier choice is restored.
    const saved = stored(LANGUAGE_KEY);
    if (saved === "en" || saved === "ru" || saved === "ca") {
      language.value = saved;
      language.dispatchEvent(new Event("change", {bubbles: true}));
    }
    language.addEventListener("change", () => remember(LANGUAGE_KEY, language.value));
  }

  if (year && !year.value.trim()) {
    year.value = String(new Date().getFullYear());
  }

  let suggestedTitle = "";
  const suggest = (): void => {
    if (!project || (project.value.trim() && project.value !== suggestedTitle)) {
      return;
    }
    const name = repo?.value.trim()
      ? repo.value.trim().replace(/\/+$/, "").split("/").pop() ?? ""
      : (archive?.files?.[0]?.name ?? "").replace(/\.zip$/i, "");
    if (name) {
      suggestedTitle = name.replace(/[-_]+/g, " ").trim();
      project.value = suggestedTitle;
    }
  };
  repo?.addEventListener("input", suggest);
  archive?.addEventListener("change", suggest);
}

const dirtyForms = new Set<HTMLFormElement>();
let pageEditRevision = 0;

function fileProblem(input: HTMLInputElement | null, extensions: string[], maxBytes: number): string {
  const file = input?.files?.[0];
  if (!file) return "";
  if (!extensions.some(ext => file.name.toLowerCase().endsWith(ext))) return `Choose ${extensions.join(" or ")} / Выберите ${extensions.join(" или ")}.`;
  if (!file.size || file.size > maxBytes) return `File must contain 1–${maxBytes.toLocaleString()} bytes. / Проверьте размер файла.`;
  return "";
}

function initBusyForms(): void {
  window.addEventListener("beforeunload", event => {
    if (dirtyForms.size) { event.preventDefault(); event.returnValue = ""; }
  });
  for (const form of document.querySelectorAll<HTMLFormElement>("form[data-busy]")) {
    const create = form.hasAttribute("data-create-document");
    const editable = form.hasAttribute("data-edit-form");
    let editRevision = 0;
    const markDirty = () => {
      if (!editable) return;
      editRevision++;
      pageEditRevision++;
      dirtyForms.add(form);
      form.dataset.unsaved = "true";
      const status = form.querySelector<HTMLElement>("[data-save-status]");
      if (status) status.textContent = "Unsaved changes / Несохранённые изменения";
    };
    const field = (name: string) => form.elements.namedItem(name) as HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement | null;
    const box = document.createElement("div");
    box.className = "form-feedback";
    box.setAttribute("role", "alert");
    box.tabIndex = -1;
    box.hidden = true;
    form.prepend(box);
    let approved = false;
    let sending = false;
    const academic = () => ["tz", "pz", "ro", "pmi", "pmi-team"].includes(field("profile_id")?.value ?? form.dataset.profile ?? "");
    const update = () => {
      const spec = field("specification");
      const file = field("specification_file") as HTMLInputElement | null;
      const structuredValues = ["purpose", "functions", "constraints", "acceptance"].map(n => field("spec." + n)?.value.trim() ?? "");
      const structured = structuredValues.every(Boolean);
      const hasStructured = structuredValues.some(Boolean);
      const specSources = Number(Boolean(spec?.value.trim())) + Number(Boolean(file?.files?.length)) + Number(hasStructured);
      const specError = fileProblem(file, [".txt", ".md"], 64000)
        || (specSources > 1 ? "Use one specification source: text, file OR structured fields. / Оставьте один источник ТЗ." : "")
        || (hasStructured && !structured ? "Complete all four structured fields. / Заполните все четыре поля ТЗ." : "")
        || (spec && new TextEncoder().encode(spec.value).length > 64000 ? "Specification exceeds 64,000 UTF-8 bytes. / Сократите текст ТЗ." : "");
      if (spec) {
        spec.setCustomValidity(specError || (academic() && !spec.value.trim() && !file?.files?.length && !structured ? "Добавьте техническое задание: текст или файл .txt/.md. / Provide specification text or a file." : ""));
        spec.setAttribute("aria-required", String(academic()));
      }
      const badge = form.querySelector<HTMLElement>("[data-spec-required]");
      if (badge) { badge.textContent = academic() ? "Required · Обязательно" : "Optional · Необязательно"; badge.className = academic() ? "required" : "optional"; }
      if (create) {
        const repo = field("github_url");
        const archive = field("archive") as HTMLInputElement;
        const hasRepo = Boolean(repo?.value.trim());
        const hasArchive = Boolean(archive.files?.length);
        archive.setCustomValidity(fileProblem(archive, [".zip"], 100 * 1024 * 1024));
        repo?.setCustomValidity(hasRepo === hasArchive ? "Выберите один источник: GitHub или ZIP. / Choose either GitHub or ZIP." : "");
      }
    };
    form.addEventListener("input", () => { approved = false; markDirty(); update(); });
    form.addEventListener("change", () => { approved = false; markDirty(); update(); });
    form.addEventListener("invalid", event => {
      const target = event.target as HTMLElement;
      reveal(target);
      box.classList.remove("form-feedback--success"); box.setAttribute("role", "alert");
      box.textContent = (target as HTMLInputElement).validationMessage;
      box.hidden = false;
    }, true);
    update();
    form.addEventListener("submit", async event => {
      event.preventDefault();
      update();
      if (sending || !form.reportValidity()) return;
      if (create && !approved) {
        const names = ["author", "year", ...(academic() ? ["group", "supervisor", "supervisor_position", "approver", "approver_position", "code"] : ["organization"])];
        const missing = names.filter(n => { const input = field("metadata." + n); return input && !input.disabled && !input.value.trim(); });
        if (missing.length) {
          box.replaceChildren();
          const note = document.createElement("p");
          note.textContent = "Не заполнены / Missing: " + missing.map(n => { const input = field("metadata." + n); return input?.labels?.[0]?.textContent?.trim() || n; }).join(", ") + ". В PDF будут красные TODO. / The PDF will contain red TODOs.";
          const proceed = document.createElement("button");
          proceed.type = "button"; proceed.className = "button"; proceed.textContent = "Продолжить с TODO / Continue with TODOs";
          proceed.onclick = () => { approved = true; form.requestSubmit(); };
          const fill = document.createElement("button");
          fill.type = "button"; fill.className = "button button--quiet"; fill.textContent = "Fill missing fields / Заполнить поля";
          fill.onclick = () => { const target = field("metadata." + missing[0]); if (target) { reveal(target); target.focus(); } };
          box.append(note, fill, proceed); box.hidden = false; box.focus(); return;
        }
      }
      const submitter = (event as SubmitEvent).submitter as HTMLButtonElement | null;
      const inlineSave = editable && submitter?.value !== "regenerate";
      if (!inlineSave && dirtyForms.size && !window.confirm("Leaving will discard unsaved edits. Continue? / Несохранённые правки будут потеряны. Продолжить?")) return;
      sending = true;
      box.classList.remove("form-feedback--success"); box.setAttribute("role", "alert");
      box.hidden = true;
      form.setAttribute("aria-busy", "true");
      const buttons = Array.from(form.querySelectorAll<HTMLButtonElement>("button[type='submit']"));
      const sentRevision = editRevision;
      const navigationRevision = pageEditRevision;
      const data = new FormData(form);
      if (submitter?.name) data.set(submitter.name, submitter.value);
      const buttonStates = buttons.map(button => ({button, disabled: button.disabled, text: button.textContent}));
      buttons.forEach(button => { button.disabled = true; if (button.dataset.busyLabel) button.textContent = button.dataset.busyLabel; });
      const navigate = (url: string) => {
        if (dirtyForms.size && pageEditRevision !== navigationRevision) {
          box.textContent = "Action completed, but newer edits are unsaved. Save or copy them before leaving. / Действие выполнено, но новые правки не сохранены.";
          const link = document.createElement("a"); link.href = url; link.textContent = " Open updated status / Открыть обновлённую страницу";
          box.append(link); box.hidden = false;
          return;
        }
        dirtyForms.clear(); window.location.assign(url);
      };
      const controller = new AbortController();
      const timeout = inlineSave ? window.setTimeout(() => controller.abort(), 30000) : undefined;
      try {
        const response = await fetch(form.getAttribute("action") || window.location.href, { method: "POST", body: form.enctype === "multipart/form-data" ? data : new URLSearchParams(Array.from(data.entries(), ([key, value]) => [key, String(value)])), headers: { Accept: "application/json" }, credentials: "same-origin", signal: controller.signal });
        if (response.ok && response.redirected) { navigate(response.url); return; }
        const json = response.headers.get("Content-Type")?.includes("application/json");
        const payload = json ? await response.json() : null;
        if (response.ok && payload?.saved && inlineSave) {
          if (editRevision === sentRevision) { dirtyForms.delete(form); delete form.dataset.unsaved; }
          const status = form.querySelector<HTMLElement>("[data-save-status]");
          if (status) status.textContent = editRevision === sentRevision ? "Saved. Review and rebuild PDF. / Сохранено. Пересоберите PDF." : "Earlier version saved; newer edits are unsaved. / Новые изменения ещё не сохранены.";
          box.classList.add("form-feedback--success"); box.setAttribute("role", "status");
          const card = form.closest("article");
          const origin = card?.querySelector<HTMLElement>("[data-section-origin]");
          if (origin) origin.textContent = "user";
          if (data.get("resolve_missing") === "yes") card?.querySelectorAll<HTMLElement>(".todo").forEach(node => node.hidden = true);
          const pageStatus = document.querySelector<HTMLElement>("[data-document-edit-status]");
          if (pageStatus) { pageStatus.textContent = "Review and rebuild PDF"; pageStatus.className = "pill pill--warn"; }
          const summary = document.querySelector<HTMLElement>("[data-section-summary]");
          if (summary) summary.textContent = "Saved sections have changed. Review all sections and rebuild the PDF on the input page.";
          document.querySelectorAll<HTMLElement>("[data-outline-notice], [data-section-message]").forEach(node => node.hidden = true);
          box.textContent = editRevision === sentRevision
            ? "Saved. Other unsaved edits remain on this page. / Сохранено; остальные правки остались на странице."
            : "Newer edits are NOT saved. Only the submitted version was saved. / Новые правки НЕ сохранены. Сохранена только отправленная версия.";
          if (editRevision !== sentRevision) box.classList.remove("form-feedback--success");
          box.hidden = false;
          return;
        }
        if (response.ok && payload?.url) { navigate(payload.url); return; }
        const message = payload?.error ?? (json ? "" : (await response.text()).slice(0, 1000));
        box.textContent = message || "Не удалось сохранить. Проверьте поля и повторите. / Could not save; check the fields and retry.";
        box.hidden = false; box.focus();
      } catch {
        box.textContent = controller.signal.aborted
          ? "Save response timed out. Your text is kept here; check the saved version before retrying. / Ответ не получен. Текст остался в форме; проверьте сохранённую версию перед повтором."
          : "Ошибка сети. Поля и файлы сохранены в форме; повторите отправку. / Network error; your form is preserved.";
        box.hidden = false; box.focus();
      } finally {
        window.clearTimeout(timeout);
        sending = false;
        form.removeAttribute("aria-busy");
        buttonStates.forEach(({button, disabled, text}) => { button.disabled = disabled; button.textContent = text; });
      }
    });
  }
}

function initConfirmations(): void {
  for (const element of document.querySelectorAll<HTMLElement>("[data-confirm]")) {
    element.addEventListener("click", event => {
      if (!window.confirm(element.dataset.confirm ?? "Are you sure?")) event.preventDefault();
    });
  }
}

export function initForms(): void {
  initProfileGuidance();
  initCreateDefaults();
  initAdaptiveInputs();
  initFileInputs();
  initBusyForms();
  initConfirmations();
}
