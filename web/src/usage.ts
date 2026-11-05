type Call = { input_tokens: number; output_tokens: number; provider: string; model: string; error_code: string };

export function initUsage(): void {
  const root = document.querySelector<HTMLElement>("[data-usage-report]");
  if (!root) return;
  const summary = root.querySelector<HTMLElement>("[data-usage-summary]")!;
  const refresh = root.querySelector<HTMLButtonElement>("[data-usage-refresh]")!;
  refresh.hidden = false;
  let pending = false;
  const load = async () => {
    if (pending) return;
    pending = true;
    refresh.disabled = true;
    const controller = new AbortController();
    const timeout = window.setTimeout(() => controller.abort(), 10000);
    try {
      const response = await fetch(`/api/jobs/${root.dataset.usageReport}/generation-report`, { signal: controller.signal, cache: "no-store" });
      if (!response.ok) throw Error("Unavailable report");
      const report = await response.json() as { calls: Call[] };
      if (!Array.isArray(report.calls)) throw Error("Invalid report");
      const calls = report.calls;
      const input = calls.reduce((sum, call) => sum + call.input_tokens, 0);
      const output = calls.reduce((sum, call) => sum + call.output_tokens, 0);
      const models = [...new Set(calls.map(call => `${call.provider} / ${call.model}`))];
      summary.textContent = calls.length
        ? `${input.toLocaleString()} input tokens · ${output.toLocaleString()} output tokens · ${calls.length} recorded calls (including discovery and failures). ${models.join(", ")}. Snapshot refreshed now.${calls.length === 100 ? " Showing the first 100 calls only; totals may be incomplete." : ""}`
        : "No provider calls recorded yet. No model-generated prose has been measured in this report.";
    } catch {
      summary.textContent = "Usage could not be refreshed; totals are unavailable. Retry or open the generation report.";
    } finally {
      clearTimeout(timeout);
      pending = false;
      refresh.disabled = false;
    }
  };
  refresh.addEventListener("click", () => void load());
  void load();
}
