<!-- oss-factory:v1
{"id":"20261019-crush-project-local-promotion-queue","kind":"state","project":"crush","from":"omp:campaign-parent","to":["factory"],"status":"working","dedupe_key":"promotion-queue:charmbracelet/crush","civ":{"city_id":"oss:crush"}}
-->

# Crush — single canonical promotion queue

Origin: https://github.com/charmbracelet/crush. Coordination: https://github.com/kvnloo/crush. Registry identity: `crush` / `oss:crush`. Audit: 2026-10-19.

This Docs-only branch and its coordination PR are the single project-local queue because fork Issues are disabled. It replaces the obsolete factory-local #164 exception; source PRs, review threads and experiments remain evidence, not duplicate queues. Do not merge this coordination branch into the product default. Factory #14 is a directory only. No label or check grants publication, legal-signature or merge authority.

## Ready for upstream

None independently verified publishable as a new PR. Existing origin carriers must be repaired or reviewed, not duplicated.

## Existing upstream PRs needing action

All six owner-origin PRs were read for current comments, reviews, checks and heads. Local/native qualification is separate from hosted checks and maintainer acceptance.

| Existing origin / owner branch | Exact head | Observed state and next action |
| --- | --- | --- |
| [4022](https://github.com/charmbracelet/crush/pull/4022) — `fix/context-samefile-4018`; owner kvnloo | `ef2b4ce37fb8c66bb8bd30c089a21a7a8a16dbd9` | REPAIR / BLOCKED: CLA check FAILURE; unresolved CodeQL path findings require semantic triage, not suppression. Preserve original source and owner review; do not forge CLA or hide AI provenance. |
| [3690](https://github.com/charmbracelet/crush/pull/3690) — `fix/3645-lsp-diagnostics-windows-path`; owner kvnloo | `10b95f8f947f0210f3240836ca7a502374e4d99e` | VERIFY: observed current-head hosted checks SUCCESS, but Windows/native acceptance and present-main applicability are not newly qualified here. Retain original #3645/PR history. |
| [3759](https://github.com/charmbracelet/crush/pull/3759) — `fix/sourcegraph-truncate-output`; owner kvnloo | `fc6e0315027f680fa78ef895d69202fa527b7b69` | ORIGIN_LIVE / LIMITED: native actual-tool HTTP and full Go suites pass; current-head hosted checks SUCCESS. Present upstream `7dadde0b` still returns the unbounded formatted result. Local pinned lint fails on Go export-data compatibility; real Sourcegraph/account behavior unrun. Existing carrier only. |
| [3757](https://github.com/charmbracelet/crush/pull/3757) — `per-464-custom-model-id-override`; owner kvnloo | `de62417b4ca3197cb45352f08f1de44672ce4eba` | REPAIR / CI: recovered macOS job110292126579 fails inherited TestClientServerSpawnRace: nested build killed at300.14s; internal/config PASS23.078s. Source uses a five-minute context and is unchanged from PR base. Resource contention is inference. Authenticated origin READ; rerun not attempted. Maintainer-controlled CI gate, no guessed source repair. |
| [3760](https://github.com/charmbracelet/crush/pull/3760) — `fix/whitespace-assistant-skip`; owner kvnloo | `7bbc3bbe2e7f720504fd261a55c05b5ee856a082` | OWNER_DIRECTION / SUPERSEDED CARRIER, not verified bug LANDED: maintainer prefers #3806. Native preparePrompt/ToAIMessage controls show this candidate's intended filtering, not actual provider acceptance. #3806 merge393e55788032f40efabbc39028aa94ce42bff482 has empty accepted go.mod/go.sum net diff; do not infer effective provider fix. Preserve head/history and establish current dependency behavior before further action. |
| [3694](https://github.com/charmbracelet/crush/pull/3694) — `fix/3671-headless-unresolved-large`; owner kvnloo | `83b1da62a9a1beb9171d27e56822e3f6742f92ca` | VERIFY: hosted checks SUCCESS. Hahaknight withdrew the earlier config opt-out suggestion; do not invent another escape-hatch feature from stale feedback. Read later same-thread fixes and preserve existing -m behavior. |

## Exact native qualification — 2026-10-19

| Target and exact control | Observed actual consumer path | Boundary |
| --- | --- | --- |
| #3759 control `76cc5c574e15072b15aaed0f4f843a5711fae0d9` versus published head above | Actual `NewSourcegraphTool.Run` HTTP POST: control inline640122bytes without truncation/spill; candidate136844bytes, valid UTF-8, truncation marker and full omitted payload recovered from spill | Real SDK/tool and loopback HTTP; injected network transport, not a Sourcegraph service/account. Output limit is display width, not UTF-8 bytes. |
| #3760 control `7f9a8e4ca8c0365ac981d082470999900a5d10ce` versus published head above | Actual private preparePrompt via temporary package bridge, then real ToAIMessage: old whitespace emitted one assistant with zero payload parts; candidate drops empty/whitespace while retaining text, reasoning and tool-call payloads | Executable consumer smoke, not a provider, llama.cpp or real HTTP400 reproduction. Maintainer direction remains separate from behavior. |

Both target full `go test -failfast ./...`, focused/adjacent tests, selected-package vet, product build/`--version`, changed-file gofmt and capitalization checks PASS. Go1.27.2/Linux ARM64, CGO disabled; no local race or multi-OS proof. Pinned golangci-lint2.13.1 FAILS because its analyzer supports export-data version4 but compiler emits5; no source suppression, not green lint. Executable revision metadata pins the requested targets with `+dirty` from initially missing checkout indexes; these are not clean release artifacts. Final tracked trees match HEAD after index recovery. Temporary bridge/harness files were archived and removed.

Native execution preceded the public migration claim while no owning queue existed; no earlier public lease is invented. Original code/authors were unchanged; no commits, source pushes, origin comments, reruns or merges by the native worker. Durable raw commands, paired logs, full-suite/lint output and CI diagnosis are stored separately from production carriers under this coordination branch's `docs/oss-factory/evidence/native-20261019/`.

## Requires verification — downstream carriers

| Original carrier | Branch / exact head | Next action |
| --- | --- | --- |
| [15: fix(prompt): one context inclusion per file identity](https://github.com/kvnloo/crush/pull/15) | `fix/context-samefile-4018@ef2b4ce37fb8c66bb8bd30c089a21a7a8a16dbd9` | VERIFY: requirement, current origin overlap and base, author provenance, native checks and policy must be independently reconciled before promotion. Keep original carrier; no automatic duplicate publication. |
| [14: [promotion] fix(ui): prewarm restored session render cache](https://github.com/kvnloo/crush/pull/14) | `promote/3556-session-load-prewarm@7f306231131ee4ad9557309f606eaacb361f1beb` | VERIFY: requirement, current origin overlap and base, author provenance, native checks and policy must be independently reconciled before promotion. Keep original carrier; no automatic duplicate publication. |
| [13: [promotion] fix(release): keep Windows zip contents at archive root](https://github.com/kvnloo/crush/pull/13) | `promote/3732-windows-flat-zip@f80dc6ab639bf6d5f9ebd34936326fc776a77f96` | VERIFY: requirement, current origin overlap and base, author provenance, native checks and policy must be independently reconciled before promotion. Keep original carrier; no automatic duplicate publication. |
| [12: [promotion] fix(mcp): keep stdio child stderr off the TUI](https://github.com/kvnloo/crush/pull/12) | `promote/3714-mcp-stderr@46f6109f778f69cb20caf5b6def78ac7e6552e9b` | VERIFY: requirement, current origin overlap and base, author provenance, native checks and policy must be independently reconciled before promotion. Keep original carrier; no automatic duplicate publication. |
| [11: [promotion] fix(copilot): avoid early device-flow polling](https://github.com/kvnloo/crush/pull/11) | `promote/3399-copilot-poll-spacing@2683cf607e07b6b88e1dbcf0abe9806414039f4f` | VERIFY: requirement, current origin overlap and base, author provenance, native checks and policy must be independently reconciled before promotion. Keep original carrier; no automatic duplicate publication. |
| [10: maintain(verify-crush): align skill/map with real harness](https://github.com/kvnloo/crush/pull/10) | `lauren/maintain-verify-crush-20260910@032b53650adba68612cc2f7abef0f70b84bd7216` | VERIFY: requirement, current origin overlap and base, author provenance, native checks and policy must be independently reconciled before promotion. Keep original carrier; no automatic duplicate publication. |

## Blocked or awaiting authorization

- #4022 CLA failure requires the actual contributor/legal-account owner; this worker cannot personally sign for anyone. Current policy permits accountable agent assistance, not fabricated human review or sign-offs.
- #3757 CI failure is observed. Origin rerun permissions are distinct from permission to push the owned source branch; read-only origin access does not grant Actions write.
- Platform/provider/Sourcegraph-account behavior is unrun unless a result explicitly names the observed runtime. Existing green metadata does not qualify an unexecuted local path.
- Current contribution/AI/CLA requirements, explicit human origin-action scope and exact-head claim/duplicate checks apply to every publication. Linear is retired; old references below are historical only.

## Already landed

No new source head in this queue is claimed landed by this audit. Upstream #3806 has a verified merge identity, but that does not prove its intended dependency fix survived the accepted tree or that #3760 landed. Fork #1/#2/#3/#4 are historical fork merges, not automatically origin acceptance.

## Superseded and historical carriers

- Fork #6/#7/#8/#9 are CLOSED. Their later families #13/#12/#11/#14 remain distinct OPEN evidence carriers; no source history or branch was deleted by this migration.
- Maintainer preference for #3806 supersedes promotion of #3760’s direction; bug-resolution proof remains separate.
- Factory164 becomes a historical redirect only after this document/coordination PR is live. Its full original body is retained below; its comments and history remain at the original URL.

## Priority order and exact next action

1. Correct actual same-PR maintainer objections and required policy/CI failures; never sign a CLA/DCO or claim human review on a contributor’s behalf.
2. Use the completed #3759 exact-head native proof for the existing review carrier; keep local lint-tool incompatibility and real Sourcegraph-account limits explicit.
3. #3757 actual failure is diagnosed; only an actor with origin Actions write can rerun the existing failed macOS job. No source regression or resource root cause is asserted from config-test success alone.
4. Reconcile #3760 with maintainer direction and actual current dependency behavior, preserving the branch and authorship.
5. Qualify the retained downstream families against present origin before any new submission.

## Attribution and concurrency

Preserve all original commit authors, code/design provenance and review links. The repository owner is not a separate third-party contributor. No code is reimplemented, rewritten or recredited by this queue migration. Worker claims/results belong on the coordination PR using the factory envelope and bounded leases; the global six-worker cap remains. Do not notify people from downstream tracking.

## Historical factory164 body — superseded, not current policy

The following unchanged source is retained for lossless migration. Its old issues-disabled exception, open-count snapshot, readiness statements and Linear language are obsolete and grant no authority.

---

<!-- oss-factory:v1
{"id":"20261016-crush-promotion-fallback-001","kind":"state","project":"crush","from":"chatgpt","to":["factory"],"status":"working","dedupe_key":"promotion-queue:kvnloo/crush","civ":{"city_id":"oss:crush"}}
-->

# Crush — canonical upstream-promotion queue

**Issues-disabled fork fallback:** `kvnloo/crush` has GitHub Issues disabled (create_issue returned HTTP 410). This is its **one logical per-repo promotion queue**, hosted in `kvnloo/oss-factory` as an explicit routing exception. Implementation source stays in the fork PR/branch or origin author thread; no second project registry is created. [Registered project routing](https://github.com/kvnloo/oss-factory/blob/main/factory/projects.json). [Factory-wide dashboard #14](https://github.com/kvnloo/oss-factory/issues/14).

## Reconciliation (2026-10-16)
Historical duplicate family: #14/#9 session prewarm; #13/#6 release zip flattening; #12/#7 MCP stderr isolation; #11/#8 Copilot device-flow polling. Select one author-coordinated final patch per requirement.

**Observed open fork PRs:** 10, from a live linked-account search. This is a candidate source inventory, not a ready-to-publish count.

| Original owner/source | Work | Current classification |
|---|---|---|
| [#15](https://github.com/kvnloo/crush/pull/15) | fix(prompt): one context inclusion per file identity | CHECK / NEEDS_RECHECK |
| [#14](https://github.com/kvnloo/crush/pull/14) | [promotion] fix(ui): prewarm restored session render cache | CHECK / NEEDS_RECHECK |
| [#13](https://github.com/kvnloo/crush/pull/13) | [promotion] fix(release): keep Windows zip contents at archive root | CHECK / NEEDS_RECHECK |
| [#12](https://github.com/kvnloo/crush/pull/12) | [promotion] fix(mcp): keep stdio child stderr off the TUI | CHECK / NEEDS_RECHECK |
| [#11](https://github.com/kvnloo/crush/pull/11) | [promotion] fix(copilot): avoid early device-flow polling | CHECK / NEEDS_RECHECK |
| [#10](https://github.com/kvnloo/crush/pull/10) | maintain(verify-crush): align skill/map with real harness | CHECK / NEEDS_RECHECK |
| [#9](https://github.com/kvnloo/crush/pull/9) | fix: start prewarm on SetMessages to avoid 30s hang on session load | CHECK / NEEDS_RECHECK |
| [#8](https://github.com/kvnloo/crush/pull/8) | fix: prevent early polling in Copilot device-flow to avoid slow_down escalation | CHECK / NEEDS_RECHECK |
| [#7](https://github.com/kvnloo/crush/pull/7) | fix: discard stderr from stdio MCP children to prevent TUI clobbering | CHECK / NEEDS_RECHECK |
| [#6](https://github.com/kvnloo/crush/pull/6) | fix: configure Windows archives with flat layout at zip root | CHECK / NEEDS_RECHECK |

## Promotion acceptance
Record original human/agent author and credit; current origin requirement and supported impact; exact base/head code SHAs; reproducible native/fixture tests and failed controls; newest maintainer feedback and CI; duplicate/owner conflicts; contribution/licensing policy; Linear authorization. Promote only the smallest genuinely ready nonduplicative change. All other work remains CHECK/LIMITED/HOLD/LOCAL_ONLY or has already landed.

No upstream GitHub issue/PR/comment, worker claim or default-branch merge is authorized here. Keep this one issue and [global #14](https://github.com/kvnloo/oss-factory/issues/14) current; link the exact origin URL when authorized publication actually occurs.
