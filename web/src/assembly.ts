type Palette = {
  accent: string;
  font: string;
  ink: string;
  line: string;
  surface: string;
};

type Fragment = {
  /** Target centre inside the page, in CSS pixels. */
  cx: number;
  cy: number;
  /** Target width of a text line. */
  tw: number;
  th: number;
  /** Orbit parameters. */
  phase: number;
  speed: number;
  /** Small scrap width while the fragment is still flying around. */
  scrap: number;
  glyph: string;
  size: number;
  /** Normalised target inside the page. */
  nx: number;
  ny: number;
  nw: number;
  flyStart: number;
  flying: boolean;
  fromX: number;
  fromY: number;
  fromAngle: number;
};

const lineWidths = [
  0.52, 0.84, 0.66, 0.78, 0.5, 0.82, 0.58, 0.72, 0.4, 0.78, 0.64, 0.75, 0.48, 0.7,
  0.81, 0.56,
];

const glyphs = [
  { glyph: "§", x: 0.05, y: 0.2 },
  { glyph: "A", x: 0.95, y: 0.32 },
  { glyph: "1", x: 0.05, y: 0.62 },
  { glyph: "¶", x: 0.95, y: 0.8 },
  { glyph: "//", x: 0.5, y: 0.96 },
];

/** Turns stage and section progress into one assembly ratio. */
function visualProgress(completed: number, total: number, stage: string, status: string): number {
  if (status === "completed") {
    return 100;
  }

  const sections = total > 0 ? (completed / total) * 100 : 0;
  switch (stage) {
    case "rendering":
      return Math.max(90, Math.round(sections));
    case "generating":
      return Math.max(22, Math.round(sections * 0.85));
    case "selecting_files":
    case "reviewing_evidence":
    case "analyzing":
      return Math.max(10, Math.round(sections * 0.6));
    default:
      return Math.max(4, Math.round(sections));
  }
}

function roundRect(
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  w: number,
  h: number,
  r: number,
): void {
  ctx.beginPath();
  ctx.roundRect(x, y, Math.max(w, 0.1), Math.max(h, 0.1), Math.min(r, h / 2));
}

class Assembly {
  private readonly canvas: HTMLCanvasElement;
  private readonly ctx: CanvasRenderingContext2D;
  private readonly probe: HTMLElement;
  private readonly fragments: Fragment[] = [];
  private readonly reduced: boolean;
  private palette: Palette;
  private page = { x: 0, y: 0, w: 0, h: 0 };
  private width = 0;
  private height = 0;
  private settled = 0;

  constructor(
    canvas: HTMLCanvasElement,
    probe: HTMLElement,
    initial: { completed: number; total: number; stage: string; status: string },
  ) {
    this.canvas = canvas;
    this.probe = probe;

    const ctx = canvas.getContext("2d");
    if (!ctx) {
      throw new Error("canvas is unavailable");
    }
    this.ctx = ctx;
    this.reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    this.palette = this.readPalette();
    this.build();
    this.update(initial.completed, initial.total, initial.stage, initial.status);
    this.resize();
    this.start();

    new ResizeObserver(() => this.resize()).observe(canvas);
    window
      .matchMedia("(prefers-color-scheme: dark)")
      .addEventListener("change", () => this.refreshPalette());
    new MutationObserver(() => this.refreshPalette()).observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["data-theme"],
    });
  }

  /** Called by the status poller while a job is still running. */
  update(completed: number, total: number, stage: string, status: string): void {
    const settled = Math.round((visualProgress(completed, total, stage, status) / 100) * this.fragments.length);

    if (settled > this.settled) {
      const now = performance.now();
      for (let i = this.settled; i < settled; i += 1) {
        const fragment = this.fragments[i];
        if (!fragment) {
          continue;
        }
        fragment.flyStart = now + (i - this.settled) * 70;
        fragment.flying = false;
      }
    }

    if (settled < this.settled) {
      for (let i = settled; i < this.settled; i += 1) {
        const fragment = this.fragments[i];
        if (fragment) {
          fragment.flying = false;
          fragment.flyStart = 0;
        }
      }
    }

    this.settled = settled;
    if (this.reduced) {
      this.draw(0);
    }
  }

  private build(): void {
    for (let i = 0; i < lineWidths.length; i += 1) {
      this.fragments.push({
        cx: 0,
        cy: 0,
        tw: 0,
        th: i === 0 ? 4 : 2.4,
        phase: (i * 2.399) % (Math.PI * 2),
        speed: 0.08 + ((i * 7) % 5) * 0.022,
        scrap: 8 + (i % 4) * 2.5,
        glyph: "",
        size: 0,
        nx: 0.09,
        ny: i === 0 ? 0.14 : 0.26 + (i - 1) * 0.046,
        nw: lineWidths[i] ?? 0.5,
        flyStart: 0,
        flying: false,
        fromX: 0,
        fromY: 0,
        fromAngle: 0,
      });
    }

    for (const [i, glyph] of glyphs.entries()) {
      this.fragments.push({
        cx: 0,
        cy: 0,
        tw: 0,
        th: 0,
        phase: (i * 1.7 + 0.9) % (Math.PI * 2),
        speed: 0.1 + i * 0.018,
        scrap: 0,
        glyph: glyph.glyph,
        size: 9 + (i % 3) * 1.6,
        nx: glyph.x,
        ny: glyph.y,
        nw: 0,
        flyStart: 0,
        flying: false,
        fromX: 0,
        fromY: 0,
        fromAngle: 0,
      });
    }
  }

  private resize(): void {
    const rect = this.canvas.getBoundingClientRect();
    if (rect.width === 0 || rect.height === 0) {
      return;
    }

    const dpr = Math.min(window.devicePixelRatio || 1, 2);
    this.width = rect.width;
    this.height = rect.height;
    this.canvas.width = Math.round(rect.width * dpr);
    this.canvas.height = Math.round(rect.height * dpr);
    this.ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    const pageWidth = Math.min(rect.width * 0.56, 176);
    const pageHeight = Math.min(rect.height * 0.78, 150);
    this.page = {
      x: (rect.width - pageWidth) / 2,
      y: (rect.height - pageHeight) / 2,
      w: pageWidth,
      h: pageHeight,
    };

    for (const fragment of this.fragments) {
      fragment.tw = fragment.nw * pageWidth;
      fragment.cx = this.page.x + fragment.nx * pageWidth + fragment.tw / 2;
      fragment.cy = this.page.y + fragment.ny * pageHeight;
    }

    if (this.reduced) {
      this.draw(0);
    }
  }

  private start(): void {
    if (this.reduced) {
      this.draw(0);
      return;
    }

    requestAnimationFrame(this.loop);
  }

  private loop = (now: number): void => {
    this.draw(now);
    requestAnimationFrame(this.loop);
  };

  private refreshPalette(): void {
    this.palette = this.readPalette();
    if (this.reduced) {
      this.draw(0);
    }
  }

  private readPalette(): Palette {
    const styles = getComputedStyle(this.probe);
    return {
      accent: styles.textDecorationColor || styles.color,
      font: styles.fontFamily,
      ink: styles.color,
      line: styles.borderTopColor,
      surface: styles.backgroundColor,
    };
  }

  private orbit(fragment: Fragment, now: number): { x: number; y: number; angle: number } {
    const t = now / 1000;
    const angle = fragment.phase + t * fragment.speed;
    const cx = this.width / 2;
    const cy = this.height / 2;
    const rx = this.width * 0.4;
    const ry = this.height * 0.37;

    return {
      x: cx + Math.cos(angle) * rx * (0.86 + 0.14 * Math.sin(t * 0.7 + fragment.phase)),
      y: cy + Math.sin(angle) * ry,
      angle: Math.sin(t * 0.45 + fragment.phase) * 0.55,
    };
  }

  private draw(now: number): void {
    const { ctx, width, height } = this;
    if (width === 0 || height === 0) {
      return;
    }

    ctx.clearRect(0, 0, width, height);
    this.drawPage();

    for (let i = 0; i < this.fragments.length; i += 1) {
      const fragment = this.fragments[i];
      if (!fragment) {
        continue;
      }

      if (i >= this.settled) {
        const point = this.orbit(fragment, now);
        this.drawFragment(fragment, point.x, point.y, point.angle, fragment.scrap, 0.32);
        continue;
      }

      if (this.reduced) {
        this.drawSettled(fragment);
        continue;
      }

      if (!fragment.flying) {
        if (now < fragment.flyStart) {
          const point = this.orbit(fragment, now);
          this.drawFragment(fragment, point.x, point.y, point.angle, fragment.scrap, 0.32);
          continue;
        }

        const point = this.orbit(fragment, now);
        fragment.fromX = point.x;
        fragment.fromY = point.y;
        fragment.fromAngle = point.angle;
        fragment.flying = true;
      }

      const t = Math.min(1, (now - fragment.flyStart) / 700);
      const eased = 1 - Math.pow(1 - t, 3);
      const widthNow = fragment.glyph
        ? 0
        : fragment.scrap + (fragment.tw - fragment.scrap) * eased;

      this.drawFragment(
        fragment,
        fragment.fromX + (fragment.cx - fragment.fromX) * eased,
        fragment.fromY + (fragment.cy - fragment.fromY) * eased,
        fragment.fromAngle * (1 - eased),
        widthNow,
        0.35 + 0.27 * eased,
      );
    }
  }

  private drawPage(): void {
    const { ctx, page, palette } = this;

    roundRect(ctx, page.x, page.y, page.w, page.h, 7);
    ctx.fillStyle = palette.surface;
    ctx.fill();
    ctx.strokeStyle = palette.line;
    ctx.lineWidth = 1;
    ctx.stroke();

    ctx.save();
    ctx.globalAlpha = 0.55;
    ctx.fillStyle = palette.line;
    for (const fragment of this.fragments) {
      if (fragment.glyph) {
        continue;
      }
      roundRect(
        ctx,
        fragment.cx - fragment.tw / 2,
        fragment.cy - fragment.th / 2,
        fragment.tw,
        fragment.th,
        fragment.th / 2,
      );
      ctx.fill();
    }
    ctx.restore();
  }

  private drawSettled(fragment: Fragment): void {
    const { ctx, palette } = this;

    ctx.save();
    if (fragment.glyph) {
      ctx.globalAlpha = 0.3;
      ctx.font = `italic ${fragment.size}px ${palette.font}`;
      ctx.fillStyle = palette.accent;
      ctx.textAlign = "center";
      ctx.fillText(fragment.glyph, fragment.cx, fragment.cy + fragment.size * 0.35);
    } else {
      ctx.globalAlpha = 0.62;
      ctx.fillStyle = palette.ink;
      roundRect(
        ctx,
        fragment.cx - fragment.tw / 2,
        fragment.cy - fragment.th / 2,
        fragment.tw,
        fragment.th,
        fragment.th / 2,
      );
      ctx.fill();
    }
    ctx.restore();
  }

  private drawFragment(
    fragment: Fragment,
    x: number,
    y: number,
    angle: number,
    width: number,
    alpha: number,
  ): void {
    const { ctx, palette } = this;

    ctx.save();
    ctx.translate(x, y);
    ctx.rotate(angle);
    ctx.globalAlpha = alpha;

    if (fragment.glyph) {
      ctx.font = `italic ${fragment.size}px ${palette.font}`;
      ctx.fillStyle = palette.accent;
      ctx.textAlign = "center";
      ctx.fillText(fragment.glyph, 0, fragment.size * 0.35);
    } else {
      ctx.fillStyle = palette.ink;
      roundRect(ctx, -width / 2, -fragment.th / 2, width, fragment.th, fragment.th / 2);
      ctx.fill();
    }

    ctx.restore();
  }
}

let controller: Assembly | null = null;

export function initAssembly(): void {
  const canvas = document.querySelector<HTMLCanvasElement>("[data-assembly]");
  if (!canvas) {
    return;
  }

  const probe = canvas.closest(".assembly-card")?.querySelector<HTMLElement>(".assembly-palette");
  if (!probe) {
    return;
  }

  controller = new Assembly(canvas, probe, {
    completed: Number(canvas.dataset.completed ?? "0"),
    total: Number(canvas.dataset.total ?? "0"),
    stage: canvas.dataset.stage ?? "",
    status: canvas.dataset.status ?? "",
  });
}

export function setAssemblyState(
  completed: number,
  total: number,
  stage: string,
  status: string,
): void {
  controller?.update(completed, total, stage, status);
}

