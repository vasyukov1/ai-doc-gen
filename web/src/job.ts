import { setAssemblyState } from "./assembly";

type JobSnapshot = {
  completed_sections: number;
  latest_message: string;
  stage: string;
  status: string;
  total_sections: number;
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
      return "Completed";
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

  const apply = (job: JobSnapshot): void => {
    const pill = document.querySelector<HTMLElement>('[data-field="status"]');
    if (pill) {
      const modifier = statusClass(job.status);
      pill.textContent = statusLabel(job.status);
      pill.className = modifier ? `pill pill--${modifier}` : "pill";
    }

    setText('[data-field="message"]', job.latest_message);

    const done = job.completed_sections;
    const total = job.total_sections;
    const percent = total > 0 ? Math.min(100, Math.round((done / total) * 100)) : 0;

    setText('[data-field="count"]', `${done} of ${total} sections`);
    setText('[data-field="percent"]', `${percent}%`);
    setText('[data-field="assembly-note"]', `${done} of ${total} sections`);

    const bar = document.querySelector<HTMLElement>('[data-field="progress"]');
    if (bar) {
      bar.style.setProperty("--progress", `${percent}%`);
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
    if (!document.hidden) {
      try {
        const response = await fetch(`/api/jobs/${encodeURIComponent(id)}`, {
          credentials: "same-origin",
          headers: { Accept: "application/json" },
        });
        if (!response.ok) {
          throw new Error("job unavailable");
        }

        const job = (await response.json()) as JobSnapshot;
        if (!ACTIVE.has(job.status)) {
          window.location.reload();
          return;
        }

        apply(job);
      } catch {
        // Keep the last rendered state and try again on the next tick.
      }
    }

    timer = window.setTimeout(() => void tick(), 3000);
  };

  document.addEventListener("visibilitychange", () => {
    if (!document.hidden) {
      window.clearTimeout(timer);
      void tick();
    }
  });

  timer = window.setTimeout(() => void tick(), 3000);
}
