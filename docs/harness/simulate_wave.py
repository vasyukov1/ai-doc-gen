#!/usr/bin/env python3
"""User-simulation-driven development: one wave.

Simulated users (personas) critique what you show them; a triage step keeps the product's
directive fixed and turns their findings into ONE small wave of changes. Then you build the
wave, take new screenshots and run this again.

    python simulate_wave.py --directive directive.md --personas personas.yaml \
        --screens a.png b.png --runs 2 --out out/
    #   instead of --screens:  --url http://localhost:8501   (Playwright screenshots it)
    #                          --html mockup.html             (a static mockup)
    #                          --text structure.md            (no screens: maps, situations, specs)
    python simulate_wave.py ... --dry-run    # no API key: print and save the prompts it would send

Every model answer is saved in OUT/responses/<call>.json and reused when it is already there, so
a re-run costs nothing, and you can paste a prompt into Claude or ChatGPT yourself and save the
answer under that name (the "by hand" route). Needs: pip install anthropic pyyaml playwright
"""
import argparse, base64, datetime, json, mimetypes, os, pathlib, re, sys

import yaml

HERE = pathlib.Path(__file__).parent
CRITIC_MODEL, TRIAGE_MODEL = "claude-sonnet-5-5", "claude-opus-5-5"
SEVERITY = {1: "cosmetic", 2: "minor", 3: "major", 4: "blocker"}

# The JSON each step must return. The API enforces it (structured outputs); by hand, the prompt asks.
FINDINGS_SCHEMA = {
    "type": "object", "additionalProperties": False, "required": ["findings"],
    "properties": {"findings": {"type": "array", "items": {
        "type": "object", "additionalProperties": False,
        "required": ["screen", "what_i_tried", "what_happened", "problem", "severity", "heuristic", "suggestion"],
        "properties": {k: {"type": "string"} for k in
                       ["screen", "what_i_tried", "what_happened", "problem", "heuristic", "suggestion"]}
                      | {"severity": {"type": "integer", "enum": [1, 2, 3, 4]}}}}}}
WAVE_SCHEMA = {
    "type": "object", "additionalProperties": False, "required": ["changes", "must_not_change", "rewalk"],
    "properties": {
        "changes": {"type": "array", "items": {
            "type": "object", "additionalProperties": False,
            "required": ["id", "change", "findings", "personas", "decision", "reason"],
            "properties": {"id": {"type": "string"}, "change": {"type": "string"},
                           "findings": {"type": "array", "items": {"type": "string"}},
                           "personas": {"type": "array", "items": {"type": "string"}},
                           "decision": {"type": "string", "enum": ["ACCEPT", "DEFER", "REJECT"]},
                           "reason": {"type": "string"}}}},
        "must_not_change": {"type": "array", "items": {"type": "string"}},
        "rewalk": {"type": "array", "items": {"type": "string"}}}}


def prompt_part(name, **fill):
    """The prompt files are for people too: the model gets the text after the '---8<---' line."""
    text = (HERE / name).read_text().rsplit("---8<---", 1)[-1].strip()
    for key, value in fill.items():
        text = text.replace("{{" + key + "}}", value.strip())
    return text


def screenshot(target, out_dir, paths, size=(390, 844)):
    """--url / --html: open each path on a phone-sized page and save a PNG of what is on screen."""
    from playwright.sync_api import sync_playwright
    shots = []
    with sync_playwright() as p:
        browser = p.chromium.launch(executable_path=os.environ.get("CHROME_PATH"))  # None: Playwright's own
        page = browser.new_page(viewport={"width": size[0], "height": size[1]}, device_scale_factor=2,
                                is_mobile=True, has_touch=True)
        for i, path in enumerate(paths, 1):
            page.goto(target.rstrip("/") + path if target.startswith("http") else target, wait_until="networkidle")
            page.wait_for_timeout(1000)
            shots.append(out_dir / f"screen_{i}.png")
            page.screenshot(path=shots[-1])
        browser.close()
    return shots


def image_block(path):
    data = base64.standard_b64encode(pathlib.Path(path).read_bytes()).decode()
    media = mimetypes.guess_type(str(path))[0] or "image/png"
    return {"type": "image", "source": {"type": "base64", "media_type": media, "data": data}}


def ask(args, name, model, system, content, schema, effort):
    """One model call, saved under responses/<name>.json; a saved answer is reused (re-runs are free)."""
    saved = args.out / "responses" / f"{name}.json"
    if saved.exists():
        return json.loads(saved.read_text())
    texts = [b["text"] for b in content if b["type"] == "text"]
    images = [b["path"] for b in content if b["type"] == "image_path"]
    prompt = f"# {name} → {model}\n\n## System\n\n{system}\n\n## User\n\n" + \
             "".join(f"[image: {i}]\n" for i in images) + "\n" + "\n\n".join(texts)
    (args.out / "prompts").mkdir(parents=True, exist_ok=True)
    (args.out / "prompts" / f"{name}.md").write_text(prompt)
    if args.dry_run:
        print(prompt + "\n" + "=" * 100)
        return None
    import anthropic
    blocks = [image_block(b["path"]) if b["type"] == "image_path" else b for b in content]
    try:
        reply = anthropic.Anthropic().beta.messages.create(
            model=model, max_tokens=16000, system=system, messages=[{"role": "user", "content": blocks}],
            output_config={"effort": effort, "format": {"type": "json_schema", "schema": schema}},
            betas=["server-side-fallback-2026-07-01"], fallbacks="default")  # a refusal re-runs on a fallback
    except (anthropic.AuthenticationError, TypeError) as error:  # TypeError: no key found at all
        sys.exit(f"No API key ({error}). Use --dry-run and answer prompts/{name}.md by hand (see ux_critic.md).")
    if reply.stop_reason == "refusal":
        sys.exit(f"{name}: the model declined ({reply.stop_details}).")
    answer = json.loads(next(b.text for b in reply.content if b.type == "text"))
    saved.parent.mkdir(parents=True, exist_ok=True)
    saved.write_text(json.dumps(answer, indent=2, ensure_ascii=False))
    print(f"  {name}: {model}, {reply.usage.input_tokens} in / {reply.usage.output_tokens} out")
    return answer


def screen_key(screen):
    """The first file name or situation id a finding names, so two runs' wordings can be compared."""
    m = re.search(r"[\w-]+\.(?:png|jpe?g)|US\.C\.\d+\.\d+|TS\.C\.\d+\.\d+|section \d+", screen, re.I)
    return (m.group(0) if m else screen).lower()


def spread(findings, persona, runs):
    """How much two runs of the same persona agree: the share of (screen, heuristic) pairs seen in both."""
    keys = [{(screen_key(f["screen"]), f["heuristic"].lower()) for f in findings
             if f["persona"] == persona and f["run"] == r} for r in range(1, runs + 1)]
    union, both = set().union(*keys), set.intersection(*keys) if keys else set()
    return {"per_run": [len([f for f in findings if f["persona"] == persona and f["run"] == r])
                        for r in range(1, runs + 1)],
            "agreement": round(len(both) / len(union), 2) if union else None}


def write_wave(args, wave, findings, personas, inputs):
    """wave.md: the decisions, the invariants, what to walk again, and how the runs spread."""
    by = {d: [c for c in wave["changes"] if c["decision"] == d] for d in ("ACCEPT", "DEFER", "REJECT")}
    names = {p["id"]: p["name"] for p in personas}
    row = lambda c: (f"| {c['id']} | {c['change']} | {', '.join(names.get(p, p) for p in c['personas'])} "
                     f"| {', '.join(c['findings'])} | {c['reason']} |")
    lines = [f"# Wave — {args.out.name}", "",
             f"{len(findings)} findings from {len(personas)} personas × {args.runs} runs on: "
             + ", ".join(pathlib.Path(i).name for i in inputs) + ".", "",
             f"Critique: `{args.critic_model}` · triage: `{args.triage_model}` · "
             f"{datetime.datetime.now():%Y-%m-%d %H:%M}", ""]
    for d, title in (("ACCEPT", "Accepted into this wave"), ("DEFER", "Deferred"), ("REJECT", "Rejected")):
        lines += [f"## {title} ({len(by[d])})", "", "| # | Change | Who asked | Findings | Why |",
                  "|---|---|---|---|---|"] + [row(c) for c in by[d]] + [""]
    lines += ["## What must not change", ""] + [f"- {m}" for m in wave["must_not_change"]] + [""]
    lines += ["## Walk again after the wave", ""] + [f"{i}. {s}" for i, s in enumerate(wave["rewalk"], 1)] + [""]
    lines += ["## Spread between runs", "", "| Persona | Findings per run | Agreement (screen + heuristic) |",
              "|---|---|---|"]
    for p in personas:
        s = spread(findings, p["id"], args.runs)
        lines.append(f"| {p['name']} | {s['per_run']} | {s['agreement']} |")
    (args.out / "wave.md").write_text("\n".join(lines) + "\n")
    return "\n".join(lines)


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--directive", required=True, help="the product's fixed directive (markdown)")
    ap.add_argument("--personas", required=True, help="personas.yaml")
    what = ap.add_mutually_exclusive_group(required=True)
    what.add_argument("--screens", nargs="+", help="screenshots (png/jpg)")
    what.add_argument("--url", help="a running app; screenshotted at each of --paths")
    what.add_argument("--html", help="a static mockup file")
    what.add_argument("--text", help="no screens: a text description of the structure")
    ap.add_argument("--paths", nargs="+", default=["/"], help="with --url: the pages to screenshot")
    ap.add_argument("--runs", type=int, default=2, help="runs per persona (shows the spread)")
    ap.add_argument("--max-changes", type=int, default=5, help="the cap on ACCEPT in one wave")
    ap.add_argument("--out", type=pathlib.Path, default=pathlib.Path("out"))
    ap.add_argument("--critic-model", default=CRITIC_MODEL)
    ap.add_argument("--triage-model", default=TRIAGE_MODEL)
    ap.add_argument("--dry-run", action="store_true", help="print the prompts; call nothing")
    args = ap.parse_args()
    args.out.mkdir(parents=True, exist_ok=True)

    directive = pathlib.Path(args.directive).read_text().strip()
    cast = yaml.safe_load(pathlib.Path(args.personas).read_text())
    personas, product = cast["personas"], cast.get("product", "")

    # 1. What the personas look at: screens (taken now, or given) or the structure as text.
    if args.url or args.html:
        target = args.url or pathlib.Path(args.html).resolve().as_uri()
        inputs = screenshot(target, args.out, args.paths)          # local; --dry-run still takes them
    else:
        inputs = args.screens or [args.text]
    if args.text:
        shown = [{"type": "text", "text": "What you are shown (no screens yet — the product's maps, "
                  "situations and specs):\n\n" + pathlib.Path(args.text).read_text()}]
    else:
        shown = [{"type": "image_path", "path": str(s)} for s in inputs] + [{"type": "text", "text":
                 "The screens, in order: " + ", ".join(pathlib.Path(s).name for s in inputs) + "."}]

    # 2. Critique: every persona, every run, on its own (no persona sees another's answer).
    findings, answered = [], 0
    for p in personas:
        persona = yaml.safe_dump(p, allow_unicode=True, sort_keys=False, width=1000)
        system = prompt_part("ux_critic.md", PRODUCT=product, PERSONA=persona)
        for run in range(1, args.runs + 1):
            answer = ask(args, f"critique_{p['id']}_run{run}", args.critic_model, system,
                         shown + [{"type": "text", "text": f"Run {run}. Walk it as {p['name']} and report."}],
                         FINDINGS_SCHEMA, effort="medium")
            answered += answer is not None
            for f in (answer or {}).get("findings", []):
                findings.append({"id": f"F{len(findings) + 1:02d}", "persona": p["id"], "run": run, **f})
    (args.out / "findings.json").write_text(json.dumps(
        {"inputs": [str(i) for i in inputs], "findings": findings,
         "spread": {p["id"]: spread(findings, p["id"], args.runs) for p in personas}}, indent=2, ensure_ascii=False))
    expected = len(personas) * args.runs
    print(f"{len(findings)} findings from {answered} of {expected} critique calls → {args.out / 'findings.json'}")

    # 3. Triage against the directive: deduplicate, then ACCEPT (≤ cap) / DEFER / REJECT.
    if answered < expected:
        findings_text = "(the findings from step 2 go here)"
    else:
        findings_text = json.dumps(findings, indent=1, ensure_ascii=False)
    system = prompt_part("wave_triage.md", DIRECTIVE=directive, MAX=str(args.max_changes))
    wave = ask(args, "triage", args.triage_model, system,
               [{"type": "text", "text": "Findings:\n\n" + findings_text}], WAVE_SCHEMA, effort="high")
    if wave is None:
        return
    accepted = [c for c in wave["changes"] if c["decision"] == "ACCEPT"]
    for c in accepted[args.max_changes:]:              # the cap is a rule, so the code enforces it too
        c["decision"], c["reason"] = "DEFER", "over the wave's cap — " + c["reason"]
    print(write_wave(args, wave, findings, personas, inputs))


if __name__ == "__main__":
    main()
