interface ReadingPosition {
  path: string;
  section: string;
  offset: number;
  scroll: number;
}

const POSITION_KEY = "docgenSectionsReading";

function restorePosition(): void {
  const position = history.state?.[POSITION_KEY] as ReadingPosition | undefined;
  if (!position || position.path !== location.pathname) return;

  const state = { ...history.state };
  delete state[POSITION_KEY];
  try {
    history.replaceState(state, "");
  } catch {
    // A restricted browser may prevent clearing the reload checkpoint.
  }

  const restore = (): void => {
    const heading = document.getElementById(position.section);
    const section = heading?.closest<HTMLElement>("[data-reading-section]");
    const top = section ? window.scrollY + section.getBoundingClientRect().top - position.offset : position.scroll;
    window.scrollTo({ top, behavior: "instant" });
  };
  if (document.readyState === "complete") requestAnimationFrame(restore);
  else window.addEventListener("load", () => requestAnimationFrame(restore), { once: true });
}

export function initSectionsReading(): void {
  const page = document.querySelector<HTMLElement>("[data-sections-reading]");
  if (!page) return;

  restorePosition();
  const status = page.querySelector<HTMLElement>("[data-reading-status]");

  for (const link of page.querySelectorAll<HTMLAnchorElement>("[data-sections-refresh]")) {
    link.addEventListener("click", event => {
      if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
      event.preventDefault();

      if (page.querySelector("form[aria-busy='true']")) {
        if (status) status.textContent = "An action is still in progress. Wait for its result before refreshing.";
        return;
      }
      const unsaved = page.querySelector<HTMLFormElement>("form[data-unsaved]");
      if (unsaved) {
        const message = "Your unsaved edits are kept here. Save all edited sections before refreshing.";
        if (status) status.textContent = message;
        const saveStatus = unsaved.querySelector<HTMLElement>("[data-save-status]");
        if (saveStatus) saveStatus.textContent = message;
        unsaved.querySelector<HTMLTextAreaElement>("textarea")?.focus({ preventScroll: true });
        return;
      }

      const sections = Array.from(page.querySelectorAll<HTMLElement>("[data-reading-section]"));
      const section = sections.find(item => item.getBoundingClientRect().bottom > 0);
      const position: ReadingPosition = {
        path: location.pathname,
        section: section?.querySelector("h2")?.id ?? "",
        offset: section?.getBoundingClientRect().top ?? 0,
        scroll: window.scrollY,
      };

      try {
        history.replaceState({ ...history.state, [POSITION_KEY]: position }, "");
      } catch {
        // Native reload still works if history state cannot be written.
      }
      location.reload();
    });
  }
}
