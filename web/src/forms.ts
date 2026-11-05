function reveal(element: HTMLElement): void {
  let parent = element.parentElement;
  while (parent) {
    if (parent instanceof HTMLDetailsElement) parent.open = true;
    parent = parent.parentElement;
  }
}

function initFileInputs(): void {
  for (const input of document.querySelectorAll<HTMLInputElement>("[data-file-input]")) {
    const zone = input.closest<HTMLElement>(".dropzone");
    const title = zone?.querySelector<HTMLElement>("[data-file-name]");
    const empty = title?.textContent ?? "";
    input.addEventListener("change", () => {
      zone?.classList.toggle("is-set", Boolean(input.files?.length));
      if (title) title.textContent = input.files?.[0]?.name ?? empty;
    });
  }
}

function initBusyForms(): void {
  for (const form of document.querySelectorAll<HTMLFormElement>("form[data-busy]")) {
    const create = form.hasAttribute("data-create-document");
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
      const structured = ["purpose", "functions", "constraints", "acceptance"].every(n => field("spec." + n)?.value.trim());
      if (spec) {
        spec.setCustomValidity(academic() && !spec.value.trim() && !file?.files?.length && !structured ? "Добавьте техническое задание: текст или файл .txt/.md. / Provide specification text or a file." : "");
        spec.setAttribute("aria-required", String(academic()));
      }
      const badge = form.querySelector<HTMLElement>("[data-spec-required]");
      if (badge) { badge.textContent = academic() ? "Required · Обязательно" : "Optional · Необязательно"; badge.className = academic() ? "required" : "optional"; }
      if (create) {
        const repo = field("github_url");
        const archive = field("archive") as HTMLInputElement;
        const hasRepo = Boolean(repo?.value.trim());
        const hasArchive = Boolean(archive.files?.length);
        repo?.setCustomValidity(hasRepo === hasArchive ? "Выберите один источник: GitHub или ZIP. / Choose either GitHub or ZIP." : "");
      }
    };
    form.addEventListener("input", () => { approved = false; update(); });
    form.addEventListener("change", () => { approved = false; update(); });
    form.addEventListener("invalid", event => {
      const target = event.target as HTMLElement;
      reveal(target);
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
        const missing = names.filter(n => !field("metadata." + n)?.value.trim());
        if (missing.length) {
          box.replaceChildren();
          const note = document.createElement("p");
          note.textContent = "Не заполнены / Missing: " + missing.join(", ") + ". В PDF будут красные TODO. / The PDF will contain red TODOs.";
          const proceed = document.createElement("button");
          proceed.type = "button"; proceed.className = "button"; proceed.textContent = "Продолжить с TODO / Continue with TODOs";
          proceed.onclick = () => { approved = true; form.requestSubmit(); };
          box.append(note, proceed); box.hidden = false; box.focus(); return;
        }
      }
      sending = true;
      box.hidden = true;
      form.setAttribute("aria-busy", "true");
      const buttons = Array.from(form.querySelectorAll<HTMLButtonElement>("button[type='submit']"));
      const data = new FormData(form);
      const submitter = (event as SubmitEvent).submitter as HTMLButtonElement | null;
      if (submitter?.name) data.set(submitter.name, submitter.value);
      buttons.forEach(button => button.disabled = true);
      try {
        const response = await fetch(form.action, { method: "POST", body: form.enctype === "multipart/form-data" ? data : new URLSearchParams(Array.from(data.entries(), ([key, value]) => [key, String(value)])), headers: { Accept: "application/json" }, credentials: "same-origin" });
        if (response.ok && response.redirected) { window.location.assign(response.url); return; }
        const json = response.headers.get("Content-Type")?.includes("application/json");
        const payload = json ? await response.json() : null;
        if (response.ok && payload?.url) { window.location.assign(payload.url); return; }
        const message = payload?.error ?? (json ? "" : (await response.text()).slice(0, 1000));
        box.textContent = message || "Не удалось сохранить. Проверьте поля и повторите. / Could not save; check the fields and retry.";
        box.hidden = false; box.focus();
      } catch {
        box.textContent = "Ошибка сети. Поля и файлы сохранены в форме; повторите отправку. / Network error; your form is preserved.";
        box.hidden = false; box.focus();
      } finally {
        sending = false;
        form.removeAttribute("aria-busy");
        buttons.forEach(button => button.disabled = false);
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
  initFileInputs();
  initBusyForms();
  initConfirmations();
}
