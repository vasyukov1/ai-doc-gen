import { setAssemblyState } from "./assembly";

type JobSnapshot = {
  completed_sections: number;
  latest_message: string;
  stage: string;
  status: string;
  total_sections: number;
  model_sections: number;
  recovery_pending: boolean;
};

const ACTIVE = new Set(["queued", "running", "retry_wait"]);

// These three helpers mirror statusLabel, statusClass and stageIndex in web/pages.go.
function statusLabel(status: string): string {
  switch (status) {
    case "queued":
      return "Queued";
    case "running":
      return "Running";
    case "input_ready":
      return "Input ready";
    case "needs_input":
      return "Needs input";
    case "retry_wait":
      return "Waiting to retry";
    case "ready_for_review":
      return "Ready for review";
    case "completed":
      return "PDF compiled";
    case "failed":
      return "Failed";
    case "cancelled":
      return "Cancelled";
    default:
      return status;
  }
}

function statusClass(status: string): string {
  switch (status) {
    case "completed":
      return "ok";
    case "failed":
    case "cancelled":
      return "error";
    case "input_ready":
    case "needs_input":
    case "ready_for_review":
      return "warn";
    case "queued":
    case "running":
    case "retry_wait":
      return "running";
    default:
      return "";
  }
}

function stageIndex(stage: string, status: string): number {
  if (status === "completed") {
    return 5;
  }

  switch (stage) {
    case "selecting_files":
    case "reviewing_evidence":
    case "analyzing":
      return 2;
    case "generating":
      return 3;
    case "rendering":
      return 4;
    default:
      return 1;
  }
}

function setText(selector: string, value: string): void {
  const element = document.querySelector<HTMLElement>(selector);
  if (element) {
    element.textContent = value;
  }
}

export function initJobPolling(): void {
  const root = document.querySelector<HTMLElement>("[data-job-poll]");
  const id = root?.dataset.jobPoll;
  if (!root || !id) {
    return;
  }

  let timer = 0;
  let pending = false;
  const notice = document.createElement("div");
  notice.className = "callout callout--warn"; notice.hidden = true; notice.setAttribute("role", "status");
  const message = document.createElement("p");
  const retry = document.createElement("button");
  retry.type = "button"; retry.className = "button button--quiet"; retry.textContent = "Refresh status / Обновить статус";
  notice.append(message, retry); root.prepend(notice);

  const apply = (job: JobSnapshot): void => {
    const pill = document.querySelector<HTMLElement>('[data-field="status"]');
    if (pill) {
      const modifier = job.recovery_pending ? "warn" : statusClass(job.status);
      pill.textContent = job.recovery_pending ? "Recovering" : statusLabel(job.status);
      pill.className = modifier ? `pill pill--${modifier}` : "pill";
    }

    setText('[data-field="message"]', job.latest_message);

    const done = job.completed_sections;
    const total = job.total_sections;
    const percent = total > 0 ? Math.min(100, Math.round((done / total) * 100)) : 0;

    setText('[data-field="count"]', `${done} of ${total} sections available`);
    setText('[data-field="percent"]', total > 0 ? `${percent}%` : "Planning sections");
    setText('[data-field="model-count"]', `${job.model_sections ?? 0} improved by a model.`);
    setText('[data-field="assembly-note"]', `${done} of ${total} sections available`);

    const bar = document.querySelector<HTMLElement>('[data-field="progress"]');
    if (bar) {
      bar.querySelector<HTMLElement>(".progress__fill")?.style.setProperty("--progress", `${percent}%`);
      bar.setAttribute("aria-valuenow", String(percent));
    }

    const index = stageIndex(job.stage, job.status);
    for (const step of document.querySelectorAll<HTMLElement>("[data-step]")) {
      const number = Number(step.dataset.step);
      step.classList.toggle("step--done", number < index);
      step.classList.toggle("step--current", number === index);
    }

    setAssemblyState(done, total, job.stage, job.status);
  };

  const tick = async (): Promise<void> => {
    window.clearTimeout(timer);
    if (pending || document.hidden) return;
    pending = true; retry.disabled = true;
    const controller = new AbortController();
    const timeout = window.setTimeout(() => controller.abort(), 10000);
    let reschedule = true;
    try {
      const response = await fetch(`/api/jobs/${encodeURIComponent(id)}`, {
        credentials: "same-origin", headers: { Accept: "application/json" }, signal: controller.signal,
      });
      if (!response.ok) {
        reschedule = response.status !== 404;
        throw new Error(response.status === 404 ? "This document is no longer accessible. It may have expired or this browser session has changed." : "Status update failed. The last known state is shown; processing may still be running.");
      }
      const job = (await response.json()) as JobSnapshot;
      notice.hidden = true;
      if (!ACTIVE.has(job.status)) { reschedule = false; window.location.reload(); return; }
      apply(job);
    } catch (error) {
      message.textContent = error instanceof Error && error.name !== "AbortError" && error.message !== "Failed to fetch"
        ? error.message : "Connection interrupted. The last known state is shown; processing may still be running. / Связь прервана; показан последний известный статус.";
      notice.hidden = false;
    } finally {
      window.clearTimeout(timeout); pending = false; retry.disabled = false;
      if (reschedule) timer = window.setTimeout(() => void tick(), 3000);
    }
  };

  retry.addEventListener("click", () => void tick());
  document.addEventListener("visibilitychange", () => {
    if (!document.hidden) void tick();
    else window.clearTimeout(timer);
  });
  timer = window.setTimeout(() => void tick(), 3000);
}
