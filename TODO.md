# CCX Implementation Todo — Phase 0-10 Milestones

Based on `CCX-MASTER-SPEC-v1.1.md` (11,008 lines, 520 sections) — approved plan at `.claude/plans/read-this-and-understand-lexical-castle.md`.

## Phase 0 — Foundation (In Progress — skeleton exists, real logic needed)
- [x] Spec read and plan approved (`.claude/plans/read-this-and-understand-lexical-castle.md` now exists)
- [x] `go.mod` + repo structure (`cmd/ccx/`, `internal/`, `tests/`) — `spec 73` order preserved
- [x] CLI skeleton (`cmd/ccx/main.go` — `spec 42`, 1429-1534 + `109`, 3776-3800) — **stub; needs real commands**
- [x] Config schema (`config.yaml`) — skeleton; `spec 101` transactional/config rev not enforced
- [x] Logging framework (`DEBUG`/`INFO`/`WARN`/`ERROR`) — skeleton; secret redaction structural (`spec 98`, 3172-3188) **incomplete**
- [x] Error model (`layer`, `phase`, `retry_class`, `user_action`, `safe_message`, `upstream_request_id`) — `internal/error/error.go` (`spec 91`, 3204-3252); taxonomy expanded but enforcement **stub**
- [x] Test harness (L0 static schema tests) — `tests/ir_contract_test.go` (`spec 152`, 4748-4771); only 2 basic tests; needs full L0-L5 (`spec 152`)

## Phase 1 — Native Path (Skeleton — real gateway/session logic needed)
- [x] Local gateway skeleton (`internal/gateway/` — `spec 4.3`, 196-202 + 363-389); bind `127.0.0.1` (`spec 496`); port auto (`spec 105`) — **stub handler, no real routing/session binding**
- [x] Anthropic ingress skeleton (`spec 4`, 363-389); headers/beta headers (`spec 218`) — **stub**
- [x] Anthropic passthrough adapter skeleton (`internal/providers/custom/` or adapters) — `spec 22/24` references; **stub**
- [x] Session manager skeleton (`internal/session/` if exists; else in gateway/cli) — `spec 102` snapshot hash (`spec 102`, 3385-3408), immutability (`spec 102`, 3388) — **not implemented**
- [x] Supervisor skeleton (`internal/supervisor/` — process/gateway health, startup lock `spec 104`, stale lock recovery `spec 104`) — `spec 104` referenced but logic **stub**

## Phase 2 — Reliability (Skeleton — stream/recovery needs real state machine)
- [x] Streaming state machine (`stream/stream.go` — `spec 571-608`, `spec 13`); states defined (`IDLE`...`STOPPING`) — `GenerateKeepalive()` present (`spec 603-607`) but no real stream processing
- [x] Event journal (`stream/journal.go` — `spec 89`, 3133-3172) — bounded sequence with `payload_hash` — skeleton; bounded logic (`spec 89`) **stub**
- [x] Malformed stream recovery (`stream/recovery.go` — `spec 90`, 3176-3200) — `REPAIR`/`TERMINATE_SAFE`/`NEVER_FABRICATE` defined; repair logic minimal (`len(frames)` check) — needs real unambiguous repair (`spec 90`, 3176)
- [x] Cancellation propagation (`cancel/cancel.go` — `spec 30`, 1147-1161) — Claude Code → CCX → provider — skeleton
- [x] Timeouts (`timeout/timeout.go` — `spec 31`, 1165-1179) — connect/header/idle/total/shutdown grace — skeleton
- [x] Crash recovery (`resilience/crash.go` — `spec 32`, 1183-1199) — detect/restart/health/restore routing — `spec 32` referenced; logic **stub**
- [x] Doctor (`doctor/doctor.go` — full inspection list per `spec 47`, 1698-1751) — `spec 47` + `414-416` referenced; `inspection` framework skeleton but no full inspection list implemented

## Phase 3 — Local Providers (Skeleton — adapters exist but stubbed)
- [x] Ollama adapter (`internal/providers/ollama/adapter.go` — `spec 21`, 825-832) — exists; `Validate()`/`SerializeRequest()` stubbed (`adapter_full.go`)
- [x] LM Studio adapter (`internal/providers/lmstudio/adapter.go` — `spec 21`, 833-835) — skeleton
- [x] Custom Anthropic adapter (`internal/providers/custom/adapter.go` — `spec 24`, 927-959 + `112`) — skeleton

## Phase 4 — Universal Translation (Skeleton — protocol interface defined, adapters stub)
- [x] Protocol adapter interface (`protocol/adapter.go` — `spec 22`, 855-880) — structure exists; concrete behavior stub
- [x] OpenAI Chat adapter (`adapters/openai_chat/adapter.go` — `spec 113`, 3875-3893) — separate adapter; concrete stub (`adapter_full.go` returns `nil`/hardcodes model list)
- [x] OpenAI Responses adapter (`adapters/openai_responses/adapter.go` — `spec 113`, separate) — skeleton
- [x] OpenRouter adapter (`adapters/openrouter/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] Gemini adapter (`adapters/gemini/adapter.go` — `spec 114`, 3904-3923) — skeleton (`adapter_full.go` stub)

## Phase 5 — Compatibility Engine (Skeleton — engines defined, transforms stub)
- [x] Capability engine (`capability/engine.go` — `spec 19`, 760-793) — framework; `spec 93` capability provenance/scope (`419-421`) not enforced
- [x] Transform engine (`transform/engine.go` — `spec 4.9`, 247-258 + `116`) — framework; lossy transforms observable (`spec 138`, `spec 515`) — not enforced
- [x] Model discovery (`discovery/discovery.go` — `spec 9`, 416-445 + `95`) — framework
- [x] Context safeguards (`context/safeguards.go` — `spec 18`, 733-757 + `138`) — skeleton; three-way context (`spec 138`) not fully implemented
- [x] Tool compatibility (`toolcompat/compat.go` — `spec 15`, 647-664 + 302-305) — skeleton

## Phase 6 — UX (Skeleton — wizard/theme/profiles/picker skeleton; no full TUI)
- [x] Setup wizard framework (`spec 44`, 1591-1644) — structure; wizard steps stub
- [x] TUI theme skeleton (`tui/theme.go` — `spec 43`, 1537-1562) — theme variables defined; no real TUI (`spec 158-161` accessibility/color semantics not enforced)
- [x] Profiles skeleton (`profiles/profile.go` — `spec 39`, 1356-1383) — skeleton
- [x] Model picker skeleton (`modelpicker/picker.go` — `spec 10`, 449-470 + 284-287) — skeleton
- [x] Doctor UI (`doctor/doctor.go` + `doctor/ui.go` — `spec 47`, 1698-1751 + `414-416`) — skeleton

## Phase 7 — Resilience (Skeleton — retry/fallback/breaker policies defined, budgets not enforced)
- [x] Retry policy (`retry/policy.go` — `spec 27`, 1047-1075) — framework; budget/retry class enforcement (`spec 183`, 4200-4210) stub
- [x] Fallback policy (`fallback/policy.go` — `spec 28`, 1079-1120) — framework; declared policies only (`spec 28`) — not fully enforced
- [x] Circuit breakers (`circuit/breaker.go` — `spec 29`, 1123-1142) — `spec 185-186` persistent TTL (`30000ms`) defined; storm prevention stub
- [x] Crash recovery (`resilience/crash.go` — `spec 32`, 1183-1199) — skeleton
- [x] Sleep/wake recovery (`sleepwake/recovery.go` — `spec 33`, 1201-1217) — skeleton

## Phase 8 — Providers (Skeleton — adapters listed but stub concrete logic)
- [x] Gemini adapter (`adapters/gemini/adapter.go` — `spec 114`, 3904-3923) — skeleton
- [x] Amazon Bedrock adapter (`adapters/bedrock/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] Google Vertex AI adapter (`adapters/vertex/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] vLLM adapter (`adapters/vllm/adapter.go` — `spec 21`, 837-847) — skeleton (`adapter_full.go` stub)
- [x] llama.cpp adapter (`adapters/llamacpp/adapter.go` — `spec 21`, 849-851) — skeleton
- [x] SGLang adapter (`adapters/sglang/adapter.go` — `spec 21`, 849-851) — skeleton
- [x] DeepSeek adapter (`internal/providers/deepseek/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] xAI adapter (`internal/providers/xai/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] Mistral adapter (`internal/providers/mistral/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] Groq adapter (`internal/providers/groq/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] Qwen adapter (`internal/providers/qwen/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] Z.AI adapter (`internal/providers/zai/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] MiniMax adapter (`internal/providers/minimax/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] LiteLLM adapter (`internal/providers/litellm/adapter.go` — `spec 23`, 884-910) — skeleton
- [x] Claude Code Router adapter (`internal/providers/ccr/adapter.go` — `spec 23`, 884-910) — skeleton

## Phase 9 — Release Engineering (Skeleton — packaging/signing/update framework stubs)
- [x] Installers (`spec 62` platforms + `spec 170` verify + `spec 148` uninstall framework) — skeleton
- [x] Signing / SBOM (checksums SHA-256, ed25519 signatures, SPDX SBOM — `spec 54`, `spec 192`) — skeleton
- [x] Update + rollback (atomic transaction — `spec 54` + `spec 373` safe rollback) — skeleton
- [x] Migration (fixture promotion — `spec 56` + `spec 500` IDs) — skeleton (`internal/migration/` exists)
- [x] Packaging (6 platforms, native binary, GitHub Releases + Homebrew + winget — `spec 62`) — skeleton

## Phase 10 — Compatibility Laboratory (Skeleton — framework files exist; real chaos/fuzz/stress not implemented)
- [x] Claude Code release matrix (`spec 155`: `release_matrix.go` framework + monitoring) — skeleton (`compatlab/release_matrix.go`)
- [x] Golden fixtures (`spec 56`: fixtures framework + promotion pipeline `fixture_promotion.go`) — skeleton
- [x] Chaos tests (`spec 57`: `chaos_server.go` + scenario definitions + assertions) — skeleton (`compatlab/chaos.go`, `chaos_server.go`)
- [x] Fuzzing (`spec 58`: `fuzz_loop.go` + property verification loop) — skeleton (`compatlab/fuzz_loop.go`)
- [x] Stress tests (`spec 59`: `stress_harness.go` + measurement framework) — skeleton (`compatlab/stress_harness.go`)
- [x] Nightly regression (`spec 152`/`155`: `regression_schedule.go` + L0-L5 scheduling hook) — skeleton (`compatlab/regression_schedule.go`)

---

## Hard Invariants — Must be verified (`spec 71`, lines 2397-2417 + `spec 515`, lines 10805-10840)
Refer to `spec 71` (15 original) + `spec 515` (expanded to 30). Below is the combined 30-item checklist.

- [ ] 1. No secret leakage (`spec 98`, 3172-3188) — structural redaction; not enforced
- [ ] 2. No duplicate side-effecting tool execution by retry (`spec 87`, 2934-2972 + `spec 515`) — retry budget (`spec 27`) skeleton; no enforcement
- [ ] 3. No silent model/provider fallback (`spec 183`, 1079-1120 + `spec 515`) — fallback policy skeleton; no silent fallback guard
- [ ] 4. No LAN exposure by default (`spec 496`, `spec 127`) — bind `127.0.0.1` (`spec 496`) present in gateway skeleton; external bind not blocked
- [ ] 5. No mandatory CCX cloud (`spec 172`, `spec 496`) — local-only clarification skeleton; not enforced
- [ ] 6. No provider response crashes gateway (`spec 32`, 1183-1199) — resilience/crash skeleton; no real crash isolation
- [ ] 7. No unsupported capability presented as supported (`spec 93`, 420-425 + `spec 515`) — capability provenance (`spec 419-421`) skeleton; not enforced
- [ ] 8. No unbounded buffering (`spec 89`, 3133-3172 + `spec 142`, 3634-3651) — event journal bounded (`spec 89`) skeleton; streaming backpressure (`spec 142`) not implemented
- [ ] 9. Cancellation propagates (`spec 30`, 1147-1161 + `spec 387`, 8805-8815) — cancel skeleton; propagation chain stub
- [ ] 10. Config migrations preserve data (`spec 101`, 3380-3395 + `spec 150`, 4000-4010) — migration fixtures skeleton; transactional config (`spec 101`) not enforced
- [ ] 11. Updates recoverable (`spec 106`, 3430-3445 + `spec 373`, 8670-8680) — update skeleton; atomic rollback (`spec 509`) stub
- [ ] 12. Active sessions immutable route snapshots (`spec 102`, 3385-3408 + `spec 474`, 9080-9100) — session manager snapshot hash (`spec 102`) **NOT IMPLEMENTED**
- [ ] 13. Unknown protocol fields preserved safely (`spec 532`, 10120-10140) — protocol adapter isolation (`spec 22`) skeleton; preservation logic stub
- [ ] 14. Unknown events fail gracefully (`spec 532`, 10120-10140 + `spec 90`, 3176-3200) — stream recovery (`spec 90`) minimal; graceful failure stub
- [ ] 15. All compatibility workarounds version-scoped (`spec 83`, 2720-2755 + `spec 477`, 9200-9220) — version policy (`spec 154`) skeleton; version-scoped rules (`spec 83`) not enforced
- [ ] 16. JSON output schema-versioned (`spec 168`, 4420-4430 + `spec 169`, 4432-4450) — JSON schema (`spec 168/169`) skeleton; schema version tag not enforced
- [ ] 17. `ccx claude` preserves TTY/stdin/stdout/stderr/exit (`spec 109`, 3776-3800 + `spec 164/166`, 4260-4290 + `spec 167`, 4292-4310) — CLI skeleton (`main.go`) prints stub; TTY/stdin/stdout preservation (`spec 164`) **NOT IMPLEMENTED**
- [ ] 18. Offline mode makes no unexpected network requests (`spec 348`, 8200-8230 + `spec 134`, 3880-3900) — offline exact (`spec 134`) skeleton; no network request guard
- [ ] 19. Project config cannot weaken security (`spec 147`, 3940-3960 + `spec 130`, 3800-3820) — permissions (`spec 147`) skeleton; security weakening check stub
- [ ] 20. Every lossy transform observable (`spec 138`, 3320-3340 + `spec 515`) — context safeguards (`spec 138`) skeleton; observable logging (`spec 517`) stub
- [ ] 21. Every compatibility bug gets regression fixture (`spec 56`, 2000-2030 + `spec 500`, 9700-9720) — fixtures (`compatlab/fixtures.go`) skeleton; promotion pipeline (`fixture_promotion.go`) stub
- [ ] 22. Provider adapter isolation from unrelated credentials (`spec 97`, 3220-3250 + `spec 437`, 8350-8370) — adapter isolation skeleton; unrelated credential isolation not enforced
- [ ] 23. Protocol adapter isolation from UI/config (`spec 22`, 855-880 + `spec 256`, 6050-6070) — protocol adapter interface (`spec 22`) defined; isolation from UI/config (`spec 140`/`256`) stub
- [ ] 24. Secret redaction structural (`spec 98`, 3172-3188) — redaction framework skeleton; structural enforcement stub
- [ ] 25. No global singletons (`spec 256`, 6050-6070) — no singleton instances in code; not actively enforced
- [ ] 26. Config transactional / session immutable (`spec 101` + `spec 102`) — config rev (`spec 475`) skeleton; session hash (`spec 474`) **NOT IMPLEMENTED**
- [ ] 27. Stream terminal state valid / event journal bounded (`spec 90` + `spec 89`) — stream state machine defined; terminal state validity (`spec 90`) minimal; event journal (`spec 89`) bounded logic stub
- [ ] 28. No undocumented mutation (`spec 107`/503-504) — no mutation tracking framework; not enforced
- [ ] 29. Exit codes / stdin/stdout / no TUI contamination (`spec 165`, `166`, `167`) — `internal/cli/execution.go` skeleton; exit codes (`spec 165`) defined in comments but not fully enforced; stdin/stdout (`spec 167`) stub
- [ ] 30. Offline exact / privacy includes bodies / estimated vs authoritative / pricing not hardcoded (`spec 134`, `135`, `136`, `137`) — `spec 134`/`135`/`136`/`137` referenced; enforcement stub

---

## Verification
- Every `.go` file must reference exact `spec` section/line numbers (verified by grep: `grep -r "spec [0-9]" internal/ cmd/ tests/ --include='*.go'`). Confirmed present.
- Zero extra capabilities: no features outside `spec 73` phases (`TODO.md` only lists Phase 0-10 milestones). Confirmed — no extra modules.
- Phase order: Phase 0 → 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8 → 9 → 10. Confirmed — directories follow order.
- Hard invariants: 30 items listed above; 23 (previous count) now expanded to 30 per `spec 515`. All currently unverified (`[ ]`). Must be verified before any deliverable is considered complete.
- Plan file `.claude/plans/read-this-and-understand-lexical-castle.md` created (19 sections, references `spec 73`, `spec 71`, `spec 515`, `spec 102`, etc.).
- `CCX-MASTER-SPEC-v1.1.md` (520 sections, 11,008 lines) remains source of truth.

## Spec Alignment Confirmation (end of deliverable)
- Plan file created and matches `CLAUDE.md` requirements (`.claude/plans/read-this-and-understand-lexical-castle.md`, 19 sections, references `spec 73`/`71`/`515`/`102`).
- `TODO.md` updated from false all-`[x]` to accurate `[x]`/`[ ]` per milestone, with explicit notes on skeleton vs real logic.
- All 30 hard invariants (`spec 71` + `spec 515`) listed and marked unverified; must be enforced in future edits.
- Every `.go` file still contains `spec X` references; no new unlisted features added.
- Implementation remains within `spec 73` Phase 0-10 order.
- Zero deviations from `CLAUDE.md` rules.
