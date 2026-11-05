type Theme = "system" | "light" | "dark";

const STORAGE_KEY = "docgen.theme";

function storedTheme(): Theme {
  try {
    const value = localStorage.getItem(STORAGE_KEY);
    if (value === "light" || value === "dark" || value === "system") {
      return value;
    }
  } catch {
    // Storage can be unavailable in private modes; the system theme still works.
  }
  return "system";
}

function applyTheme(theme: Theme): void {
  const root = document.documentElement;
  if (theme === "system") {
    root.removeAttribute("data-theme");
  } else {
    root.setAttribute("data-theme", theme);
  }
}

export function initTheme(): void {
  let current = storedTheme();
  applyTheme(current);

  const menu = document.querySelector<HTMLElement>("[data-theme-menu]");
  if (!menu) {
    return;
  }

  const trigger = menu.querySelector<HTMLButtonElement>("[data-theme-trigger]");
  const panel = menu.querySelector<HTMLElement>("[data-theme-panel]");
  if (!trigger || !panel) {
    return;
  }

  const choices = Array.from(panel.querySelectorAll<HTMLButtonElement>("[data-theme-choice]"));

  const sync = (): void => {
    for (const choice of choices) {
      choice.setAttribute("aria-checked", String(choice.dataset.themeChoice === current));
    }
  };

  const setOpen = (open: boolean): void => {
    panel.hidden = !open;
    trigger.setAttribute("aria-expanded", String(open));
  };

  sync();
  setOpen(false);

  trigger.addEventListener("click", () => setOpen(panel.hidden));

  for (const choice of choices) {
    choice.addEventListener("click", () => {
      const next = choice.dataset.themeChoice;
      current = next === "light" || next === "dark" ? next : "system";
      try {
        localStorage.setItem(STORAGE_KEY, current);
      } catch {
        // The theme still applies for this page.
      }
      applyTheme(current);
      sync();
      setOpen(false);
      trigger.focus();
    });
  }

  document.addEventListener("click", (event) => {
    if (!menu.contains(event.target as Node)) {
      setOpen(false);
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && !panel.hidden) {
      setOpen(false);
      trigger.focus();
    }
  });
}
