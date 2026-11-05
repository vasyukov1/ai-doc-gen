// Opens the compiled PDF in a dialog anchored to the right of the viewport.
// The trigger stays a plain link: without JavaScript, or when the dialog
// element is unavailable, the browser opens the inline PDF in a new tab.
export function initPdfViewer(): void {
  const viewer = document.querySelector<HTMLDialogElement>("[data-pdf-viewer]");
  const frame = document.querySelector<HTMLIFrameElement>("[data-pdf-frame]");
  if (!viewer || !frame || typeof viewer.showModal !== "function") {
    return;
  }

  const triggers = Array.from(document.querySelectorAll<HTMLAnchorElement>("[data-pdf-open]"));
  if (triggers.length === 0) {
    return;
  }

  let lastTrigger: HTMLAnchorElement | null = null;

  const open = (trigger: HTMLAnchorElement, event: MouseEvent): void => {
    // Modified clicks keep the browser's own new-tab behaviour.
    if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) {
      return;
    }

    const source = trigger.getAttribute("href");
    if (!source) {
      return;
    }

    event.preventDefault();
    lastTrigger = trigger;
    if (frame.getAttribute("src") !== source) {
      frame.setAttribute("src", source);
    }
    if (!viewer.open) {
      viewer.showModal();
    }
  };

  for (const trigger of triggers) {
    trigger.addEventListener("click", (event) => open(trigger, event));
  }

  const close = (): void => {
    if (viewer.open) {
      viewer.close();
    }
  };

  viewer.querySelector("[data-pdf-close]")?.addEventListener("click", close);

  // A click on the backdrop targets the dialog element itself.
  viewer.addEventListener("click", (event) => {
    if (event.target === viewer) {
      close();
    }
  });

  viewer.addEventListener("close", () => {
    // Drop the document so the next open always shows the latest build.
    frame.removeAttribute("src");
    lastTrigger?.focus();
    lastTrigger = null;
  });
}
