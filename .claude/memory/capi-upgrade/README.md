# CAPx upgrade — this fork's memory (read this first)

This directory is the **only** CAPx content that lives in this fork — its **config, human decision
inputs, notes, and engine-generated records**, in four role-based folders. The engine, skills, and
agents are **not** here; they come from the `spectro-capx-upgrade` plugin (skill `capi-repo-upgrade`).

> **Not onboarded yet?** Create `conf/capi-upgrade.conf` with at least `UPSTREAM_URL`, `PROVIDER`,
> `SPECTRO_SRC`. Everything in `runtime/` is generated for you; `decisions/` you fill in at review time.

## Layout — four folders by role

```
.claude/memory/capi-upgrade/
├── conf/       — repo config you set ONCE (committed)
│   ├── capi-upgrade.conf          UPSTREAM_URL · PROVIDER · SPECTRO_SRC · CORE_CAPI · patterns · IMAGES · CLOUD · CBT_FEATURE · VENDORCRD
│   ├── expected-test-failures.tsv (optional) tests that fail BY DESIGN — used by the UT gate
│   └── env-failure-patterns.txt   (optional) infra/env failure patterns — used by the UT gate
├── decisions/  — Gate-1 inputs you EDIT each cycle; reviewed & approved in the Plan PR (→ Jira)
│   ├── skip-commits.txt           SHAs to SKIP (drop)
│   ├── resolve-commits.txt        SHAs to force-pick as PICK-RESOLVE
│   ├── known-commits.tsv          authoritative  <sha> <category> [decision]
│   └── regwatch-risk.tsv          durable Regression-Watch risk  <sha> <L|M|H> [note]
├── notes/      — human prose (committed), not read by the engine
│   ├── ledger.md                  running upgrade notes for this fork
│   └── parity/<version>-parity.md parity notes for commits skipped-as-upstreamed
└── runtime/    — ENGINE-GENERATED · do NOT hand-edit (committed as the plan artifact)
    ├── summary.json               machine plan (pins S/T/F SHAs + new_branch)
    ├── plan-latest.json           the FROZEN plan Gate-1 approves and execute replays
    └── <ts>-PLAN|APPLY-*.md        per-run records
```

## What goes where (the question everyone asks)

- **You edit:** `conf/` (once) and `decisions/` (each cycle, during Gate-1 review).
- **The engine writes:** everything in `runtime/` — never hand-edit it.
- **Pushed to the Plan PR (linked in Jira):** the whole dir is committed by `capi-plan-pr.sh`; what a
  reviewer actually approves is `decisions/*` together with the frozen `runtime/plan-latest.json`.
- **`notes/`** is free-form context — nothing consumes it.

## Canonical definitions

Terms + full engine behavior live in the `capi-repo-upgrade` skill's `references/glossary.md`,
`SKILL.md`, and `DISCOVERED-INTERFACES.md` (in the `spectro-capx-upgrade` plugin).

## How it runs (quick)

plan (read-only) → **Gate 1** (edit `decisions/`, add `luma/approved`) → execute
(creates `spectro-v<target>-master`) → **Gate 2** (a maintainer merges).
Single fork: `/capx-upgrade <target>`. Whole fleet: Luma's `capi-fleet` coordinator.
