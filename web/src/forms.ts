function initFileInputs(): void {
  const inputs = document.querySelectorAll<HTMLInputElement>("[data-file-input]");

  for (const input of inputs) {
    const zone = input.closest<HTMLElement>(".dropzone");
    if (!zone) {
      continue;
    }

    const title = zone.querySelector<HTMLElement>("[data-file-name]");
    const empty = title?.textContent ?? "";

    input.addEventListener("change", () => {
      const file = input.files?.[0];
      zone.classList.toggle("is-set", Boolean(file));
      if (title) {
        title.textContent = file ? file.name : empty;
      }
    });
  }
}

function initBusyForms(): void {
  const forms = document.querySelectorAll<HTMLFormElement>("form[data-busy]");

  for (const form of forms) {
    form.addEventListener("submit", () => {
      if (!form.checkValidity()) {
        return;
      }

      form.setAttribute("aria-busy", "true");
      const button = form.querySelector<HTMLButtonElement>("button[type='submit']");
      if (button) {
        if (button.dataset.busyLabel) {
          button.textContent = button.dataset.busyLabel;
        }
        button.disabled = true;
      }
    });
  }
}

function initConfirmations(): void {
  const elements = document.querySelectorAll<HTMLElement>("[data-confirm]");

  for (const element of elements) {
    element.addEventListener("click", (event) => {
      const message = element.dataset.confirm ?? "Are you sure?";
      if (!window.confirm(message)) {
        event.preventDefault();
      }
    });
  }
}

export function initForms(): void {
  initFileInputs();
  initBusyForms();
  initConfirmations();
}
