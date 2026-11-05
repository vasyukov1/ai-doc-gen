function resize(textarea: HTMLTextAreaElement): void {
  textarea.style.height = "auto";
  textarea.style.height = `${textarea.scrollHeight + 2}px`;
}

export function initTextareas(): void {
  const textareas = document.querySelectorAll<HTMLTextAreaElement>("textarea[data-autogrow]");

  for (const textarea of textareas) {
    resize(textarea);
    textarea.addEventListener("input", () => resize(textarea));
  }
}
