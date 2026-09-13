# CCX — Universal Claude Code Compatibility Gateway
## Canonical End-to-End Build Specification

**Status:** Final architecture/specification baseline  
**Audience:** AI coding agents, maintainers, contributors  
**Project:** `ccx`  
**License:** Open-source (license to be selected before repository publication)  
**Runtime model:** Local-only. CCX does not operate a hosted control plane, hosted proxy, telemetry backend, or mandatory cloud service.

---

## 0. Mission

CCX makes an unmodified Claude Code CLI usable with a broad model/provider ecosystem through one local runtime.

The target user experience is:

```text
install ccx
    ↓
ccx
    ↓
choose provider
    ↓
choose model
    ↓
verify compatibility
    ↓
claude
```

After setup, normal usage should require no manually started proxy:

```bash
claude
```

CCX starts/reuses its local runtime, resolves the configured route, and stays out of the user's way. If the user invokes `claude` directly rather than `ccx run`, CCX MUST provide an explicit supported integration mechanism (service, wrapper, or environment launcher); it MUST NOT assume Claude Code can magically discover CCX without that integration.

### Core promise

> **Configure once. Run Claude Code normally. CCX absorbs provider/protocol incompatibility underneath.**

CCX MUST NOT require users to understand protocol translation, SSE, beta headers, capability negotiation, tool schemas, model discovery quirks, or provider-specific workarounds.

Power users MUST be able to inspect all of those details.

---

# 1. Product contract

## 1.1 MUST

CCX MUST:

- be local-first and usable without any CCX-hosted service;
- keep the Claude Code binary unmodified;
- expose a local Anthropic-compatible gateway to Claude Code;
- support cloud providers, aggregators, local runtimes, existing gateways, and arbitrary compatible endpoints;
- isolate provider-specific behavior behind adapters;
- normalize protocol differences through a canonical internal representation;
- preserve streaming and tool-call semantics;
- propagate cancellation;
- prevent duplicate tool execution caused by retries;
- never silently switch providers/models unless explicitly configured;
- never log secrets;
- default to loopback-only listening;
- provide deterministic diagnostics;
- support headless/JSON operation;
- have a compatibility test lab for Claude Code/provider regressions;
- make unsupported capabilities explicit instead of pretending compatibility.

## 1.2 SHOULD

CCX SHOULD:

- auto-start its local runtime;
- detect installed/local model servers;
- discover models where possible;
- cache model metadata safely;
- automatically recover from CCX crashes;
- recover after laptop sleep/wake;
- provide explicit fallback policies;
- provide cost/usage accounting;
- support project-scoped configuration;
- support signed self-updates and rollback;
- support provider/plugin contributions.

## 1.3 MUST NOT

CCX MUST NOT:

- require a CCX account;
- require a CCX cloud backend;
- silently upload prompts, responses, logs, or telemetry;
- silently alter shell startup files;
- expose the local gateway to LAN by default;
- claim a provider capability merely because a marketing/catalog flag says it exists;
- replay a partially completed agent request automatically when doing so could duplicate side effects;
- silently change a model after a failure unless the selected policy permits it.

---

# 2. Current Claude Code compatibility facts

These are compatibility facts to track as versioned evidence, not permanent assumptions.

Claude Code gateway model discovery exists for gateway configurations and can query `/v1/models`; current documentation indicates discovery requires Claude Code v2.1.129 or later. citeturn0search7

Current Claude Code behavior also exposes custom-model configuration mechanisms including `ANTHROPIC_CUSTOM_MODEL_OPTION`, custom headers, and custom capability declarations in newer versions. citeturn0search1

Real Claude Code issues show that gateway discovery can fail under OAuth-passthrough gateway configurations even when the discovery flag is enabled. One 2026 report documents discovery being gated by internal gateway-mode resolution and silently skipped in that setup. citeturn0search4turn0search5

Anthropic's Messages API remains the canonical native protocol for the compatibility ingress: `POST /v1/messages`, with structured messages, tool use, streaming, and thinking/extended-thinking behavior. citeturn0search9turn0search0

Anthropic's streaming API uses SSE and exposes incremental text, tool-use and extended-thinking events. citeturn0search2

**Implementation rule:** Claude Code behavior MUST be continuously re-verified against released versions. Do not hard-code today's quirks as timeless protocol guarantees.

---

# 3. Architecture

```text
                         USER
                          │
                   ccx / claude
                          │
              ┌───────────▼──────────┐
              │   CCX Supervisor     │
              └───────────┬──────────┘
                          │
              ┌───────────▼──────────┐
              │   Session Manager    │
              └───────────┬──────────┘
                          │
                    Local Gateway
                          │
              ┌───────────▼──────────┐
              │ Anthropic Ingress    │
              │ + Compat Layer       │
              └───────────┬──────────┘
                          │
                    Universal IR
                          │
       ┌──────────────────┼──────────────────┐
       │                  │                  │
 Capability Engine   Transform Engine   Resilience
       │                  │               Engine
       └──────────────────┼──────────────────┘
                          │
                       Router
                          │
                    Protocol Layer
                          │
             ┌────────────┼────────────┐
             │            │            │
         Anthropic     OpenAI       Gemini
         /native       Chat/Resp     native
             │            │            │
             └────────────┼────────────┘
                          │
                   Provider Adapters
                          │
       ┌──────────────────┼──────────────────┐
       │                  │                  │
      Cloud             Local             Custom
```

---

# 4. Module boundaries

## 4.1 CLI/TUI

Owns:

- command parsing;
- interactive flows;
- presentation;
- machine-readable output.

Must not own provider/network logic.

## 4.2 Supervisor

Owns:

- runtime lifecycle;
- process/service state;
- health;
- startup;
- shutdown;
- crash recovery.

## 4.3 Gateway

Owns:

- local HTTP server;
- Claude Code ingress;
- request validation;
- response/SSE serialization.

## 4.4 Session Manager

Owns:

- session identity;
- route binding;
- provider/model association;
- config snapshot;
- concurrency isolation.

## 4.5 IR

Owns canonical request/response/event structures.

## 4.6 Protocol adapters

Convert between IR and external wire formats.

## 4.7 Provider adapters

Own:

- base URL;
- auth;
- discovery;
- provider quirks;
- provider limits;
- provider-specific metadata.

Provider logic MUST NOT leak into the universal IR.

## 4.8 Capability engine

Combines:

```text
declared + probed + observed
        ↓
effective capability
```

## 4.9 Transform engine

Uses policy operations:

```text
PASS
CONVERT
REMOVE
EMULATE
DEGRADE
REJECT
```

## 4.10 Resilience engine

Owns:

- timeout;
- retry;
- backoff;
- circuit breaker;
- fallback;
- cancellation;
- reconnect.

## 4.11 Diagnostics

Owns:

- doctor;
- health;
- compatibility reports;
- repair suggestions.

## 4.12 Secrets

Owns credential resolution/storage/redaction.

No adapter should implement its own secret logging/redaction.

---

# 5. Runtime modes

## 5.1 AUTO

Default.

```text
claude
 ↓
detect CCX
 ↓
runtime exists?
 ├─ yes → reuse
 └─ no → start
 ↓
inject environment
 ↓
launch Claude Code
```

## 5.2 DAEMON

Persistent per-user service.

Supported service mechanisms:

- macOS launchd;
- Linux systemd user service;
- Windows per-user service/task mechanism.

No root requirement.

## 5.3 EPHEMERAL

```bash
ccx run --ephemeral
```

Starts only for the command/session and shuts down when safe.

Useful for CI and isolated tests.

---

# 6. Process lifecycle state machine

```text
STOPPED
  ↓
STARTING
  ↓
HEALTHY
  ↓
RUNNING
  ├── DEGRADED
  ├── RECOVERING
  └── STOPPING
```

Failure:

```text
RUNNING
  ↓
CRASHED
  ↓
BACKOFF
  ↓
STARTING
```

Repeated failure MUST trip a startup circuit breaker and provide actionable diagnostics instead of restart-looping forever.

---

# 7. Local gateway

Default bind:

```text
127.0.0.1
```

Port:

- automatically selected;
- persisted only when useful;
- collision-safe.

Gateway MUST support:

```text
POST /v1/messages
```

and the minimum ancillary compatibility endpoints required by supported Claude Code versions.

`/v1/models` support SHOULD exist for gateway discovery and CCX tooling.

Unknown endpoints SHOULD return a useful structured error, not a generic crash.

---

# 8. Claude Code environment integration

CCX MUST avoid permanent shell mutation by default.

Preferred launch:

```text
ccx resolves route
 ↓
creates process environment
 ↓
ANTHROPIC_BASE_URL=local CCX
ANTHROPIC_AUTH_TOKEN=ephemeral local token
other required compatibility variables
 ↓
exec claude
```

CCX should use an ephemeral local token even when upstream credentials are stored elsewhere, so Claude Code gateway behavior is deterministic.

Environment variables MUST be versioned in the compatibility layer because Claude Code changes supported gateway variables over time.

---

# 9. Model discovery

Discovery sources:

1. provider `/v1/models`;
2. provider-native catalog;
3. local runtime API;
4. static provider catalog;
5. manually configured model;
6. CCX compatibility metadata.

Discovery result:

```yaml
id:
display_name:
provider:
protocol:
context_window:
capabilities:
pricing:
metadata:
source:
observed_at:
```

Claude Code's own picker can filter gateway-discovered models. Current evidence shows custom/non-Claude IDs can disappear from its picker, while direct model IDs and custom-model mechanisms remain usable. citeturn0search3turn0search6

Therefore CCX MUST NOT make Claude Code's picker the authoritative model-management UI.

CCX's picker is authoritative.

---

# 10. Model selection strategy

Default:

```text
CCX selects model BEFORE Claude Code starts.
```

CCX sets the resolved model explicitly.

This avoids dependence on Claude Code's internal model-picker filtering.

Supported:

```bash
ccx use <model>
ccx run --model <model>
ccx run --provider <provider> --model <model>
```

The CCX TUI MUST allow searching all models regardless of whether Claude Code's picker would display them. CCX MUST NOT modify undocumented Claude Code cache files as a primary mechanism; any client-file workaround is experimental, version-gated, backed up, and disabled by default.

---

# 11. Universal IR

Minimum structures:

```text
Request
  model
  system
  messages
  tools
  tool_choice
  thinking
  max_tokens
  temperature
  metadata
  cache_control
  stream

Message
  role
  content[]

ContentBlock
  text
  image
  document
  tool_use
  tool_result
  thinking
  redacted_thinking
  unknown

StreamEvent
  message_start
  content_block_start
  content_block_delta
  content_block_stop
  message_delta
  message_stop
  ping
  error

Usage
  input_tokens
  output_tokens
  cache_creation
  cache_read

Error
  category
  provider_code
  retryable
  safe_message
  upstream_request_id
```

IR MUST preserve unknown/extension fields where possible.

Unknown fields MUST NOT be silently discarded if they can be safely retained.

---

# 12. Anthropic compatibility ingress

Support:

- required/optional Anthropic headers;
- API version;
- beta headers;
- custom headers;
- request IDs;
- streaming;
- tool use;
- thinking;
- cache directives;
- usage;
- cancellation;
- errors.

Headers must pass through a policy layer.

Example:

```text
incoming anthropic-beta
        ↓
classify
 ├─ native Anthropic target → preserve
 ├─ translated target → remove/translate
 └─ unknown target → policy decision
```

Never blindly forward Anthropic beta headers to OpenAI/Gemini endpoints.

Anthropic's current API continues to evolve optional headers/features, so header handling must be data-driven and covered by compatibility fixtures. citeturn0search9

---

# 13. Streaming state machine

```text
IDLE
 ↓
MESSAGE_STARTED
 ↓
BLOCK_STARTED
 ↓
DELTA*
 ↓
BLOCK_STOPPED
 ↓
...
 ↓
MESSAGE_DELTA
 ↓
MESSAGE_STOPPED
```

Must handle:

- text;
- reasoning;
- tool input;
- usage;
- ping;
- errors;
- disconnects.

Anthropic documents SSE streaming for text, tool use and extended thinking, including fine-grained tool-input streaming. citeturn0search2

## Keepalive

Translated streams MUST generate compatible keepalive events when the upstream can remain silent for extended periods.

This specifically protects long thinking streams.

---

# 14. Tool execution safety

Tool IDs:

```text
Claude ID ↔ CCX ID ↔ provider ID
```

must be mapped deterministically.

For every tool call:

```text
request fingerprint
session ID
tool ID
call ID
attempt
```

must be tracked.

## Hard invariant

A retry after ambiguous network failure MUST NOT automatically re-execute an already-issued tool call.

If execution status is unknown:

```text
UNKNOWN_EXECUTION_STATE
```

and the system must require a safe continuation strategy.

---

# 15. Tool protocol normalization

Support:

- JSON-schema tools;
- function tools;
- tool choice;
- parallel calls;
- tool results;
- tool errors;
- partial streamed arguments.

Malformed arguments MUST produce controlled errors.

Provider-specific parser quirks belong in provider/model compatibility metadata.

For vLLM-like runtimes, model-family parser differences MUST be represented explicitly rather than assuming every OpenAI-compatible endpoint is behaviorally identical.

---

# 16. Thinking/reasoning

Canonical states:

```text
unsupported
supported
supported_streaming
supported_but_transform_required
unknown
```

CCX MUST NOT synthesize fake reasoning merely to claim support.

If target cannot represent thinking:

```text
policy:
  reject
  degrade
```

Default for required reasoning:

```text
REJECT
```

Default for optional reasoning:

```text
DEGRADE with warning
```

---

# 17. Prompt caching

Represent:

```text
supported
unsupported
gated
inconsistent
unknown
```

Observed cache behavior must be tracked separately from provider documentation.

Cache directives MUST be stripped or transformed when the target cannot support them.

The user-facing diagnostic should distinguish:

```text
protocol unsupported
```

from:

```text
provider/account/model gating
```

---

# 18. Context windows

Model descriptor:

```yaml
context:
  advertised:
  verified:
  safe:
```

Never assume that a custom model has the same context window as a native Claude model.

If Claude Code itself imposes a lower limit for an unrecognized model, CCX MUST detect and report the mismatch.

Policies:

```text
strict
safe
experimental
```

CCX MUST NOT claim to override Claude Code's context handling unless the required supported mechanism exists for the installed Claude Code version.

---

# 19. Capability engine

Each capability has:

```text
DECLARED
PROBED
OBSERVED
EFFECTIVE
```

Effective state:

```text
SUPPORTED
UNSUPPORTED
GATED
UNVERIFIED
BROKEN
```

No arbitrary percentage score.

Example:

```yaml
tools:
  declared: true
  probed: true
  observed: true
  effective: supported
```

---

# 20. Probe levels

## Minimal

- DNS/network;
- endpoint;
- authentication;
- model discovery.

## Standard

- streaming;
- basic text;
- basic tool use.

## Deep

- parallel tools;
- thinking;
- cache;
- multimodal;
- cancellation;
- long-stream behavior.

Probe requests MUST be bounded and visible.

The user should know when a probe may consume provider tokens.

---

# 21. Local provider special cases

## Ollama

Do not blindly forward unsupported Claude-specific probes.

CCX MUST recognize local runtime version/capability differences and avoid repeated unsupported token-count/probe requests.

## LM Studio

Use its supported Anthropic/native interface when available rather than translating through OpenAI unnecessarily.

## vLLM

Treat OpenAI compatibility as wire compatibility, not behavioral equivalence.

Track:

- model family;
- parser;
- tool behavior;
- reasoning;
- streaming quirks.

## llama.cpp / SGLang

Use their actual exposed API and capability metadata.

---

# 22. Protocol adapters

Required architecture:

```text
AnthropicAdapter
OpenAIChatAdapter
OpenAIResponsesAdapter
GeminiAdapter
OllamaAdapter
CustomAnthropicAdapter
CustomOpenAIAdapter
```

Adapters must implement:

```text
validate
serialize_request
parse_response
parse_stream
map_error
discover_models
```

Provider adapters may extend protocol adapters with quirks.

---

# 23. Provider catalog

Initial mainstream coverage SHOULD include:

```text
Anthropic
OpenAI
Google Gemini
OpenRouter
DeepSeek
xAI
Mistral
Groq
Qwen
Z.AI
MiniMax
Amazon Bedrock
Google Vertex AI
Ollama
LM Studio
vLLM
llama.cpp
SGLang
LiteLLM
Claude Code Router
arbitrary custom endpoint
```

Provider support is not considered complete merely because authentication works.

Minimum provider verification:

```text
connectivity
auth
model discovery
streaming
tools
errors
```

---

# 24. Custom endpoint UX

```bash
ccx provider add custom
```

Wizard:

```text
Endpoint URL
Protocol:
  Auto
  Anthropic
  OpenAI Chat
  OpenAI Responses
  Gemini

Authentication:
  API key
  Bearer
  Custom header
  None

Model discovery:
  Automatic
  /v1/models
  Manual
```

Then live test.

Auto protocol detection MUST be conservative. If ambiguous, ask.

---

# 25. Authentication

Credential sources:

```text
OS keychain
environment variable
interactive login
OAuth/device flow
AWS credential chain
GCP ADC
custom command
custom header
plaintext config only when explicitly selected
```

Named identities:

```text
openai:personal
openai:work
openai:backup
```

Credential state:

```text
valid
expired
invalid
cooldown
disabled
unknown
```

Secrets must never appear in:

```text
logs
doctor
status
JSON output
error messages
crash dumps
```

---

# 26. Routing

Routing dimensions:

```text
provider
model
profile
capabilities
cost
latency
health
context
```

V1:

```text
explicit route
profile route
```

V2:

```text
capability-aware route
health-aware route
```

V3:

```text
cost/latency optimization
role-based routing
```

---

# 27. Retry policy

Default:

```text
400/401/403/404 → no retry
429 → bounded backoff
5xx → bounded retry
timeout → bounded retry
network → bounded retry
```

Honor:

```text
Retry-After
```

Use:

```text
exponential backoff
jitter
max attempts
max elapsed time
```

Retries must be aware of request side effects.

---

# 28. Fallback policy

Fallback dimensions:

```text
credential
model
provider
profile
```

Explicit modes:

```text
disabled
ask
automatic
```

Default:

```text
ask
```

Automatic fallback MUST check capability compatibility before switching.

Example:

```text
Primary unavailable.

Fallback candidate:
DeepSeek / model-x

Compatibility:
tools ✓
reasoning ✓
context ✓

Use fallback? [Y/n]
```

---

# 29. Circuit breakers

Track per:

```text
provider
credential
model
endpoint
```

State:

```text
CLOSED
OPEN
HALF_OPEN
```

Avoid hammering a failing local or cloud endpoint.

---

# 30. Cancellation

Ctrl+C MUST propagate cancellation upstream whenever the protocol supports it.

On cancellation:

```text
Claude Code
 ↓
CCX cancellation
 ↓
provider cancellation
```

CCX MUST not leave abandoned generations running unnecessarily.

---

# 31. Timeouts

Separate:

```text
connect timeout
request header timeout
idle stream timeout
total request timeout
shutdown grace period
```

Long thinking streams must not be killed merely because no text arrived.

Use event/keepalive activity to reset idle timers.

---

# 32. Crash recovery

If gateway crashes:

```text
detect
 ↓
restart
 ↓
health check
 ↓
restore routing state
```

Never replay uncertain side-effecting tool calls automatically.

Session continuity should preserve metadata/config, not fabricate continuation semantics Claude Code itself did not provide.

---

# 33. Sleep/wake

On wake:

```text
invalidate stale connections
 ↓
reconnect
 ↓
health check
 ↓
resume accepting requests
```

No manual restart should normally be necessary.

---

# 34. Concurrency

CCX MUST support multiple simultaneous Claude Code processes.

Isolation key:

```text
session_id
```

Each session owns:

```text
provider
model
credential
route
config snapshot
request counters
```

A global config change MUST NOT mutate an already-running session unexpectedly.

---

# 35. Multi-hop gateway

Support:

```text
Claude Code
 ↓
CCX
 ↓
LiteLLM
 ↓
OpenRouter
 ↓
Provider
```

and:

```text
Claude Code
 ↓
CCX
 ↓
CCR
 ↓
Provider
```

CCX should display the resolved route.

Never hide multi-hop behavior.

---

# 36. Network support

Support:

```text
HTTP proxy
HTTPS proxy
NO_PROXY
custom CA
IPv4
IPv6
DNS failure
TLS errors
```

Remote local endpoints:

```text
laptop CCX
 ↓
LAN
 ↓
Ollama/vLLM server
```

must require explicit remote endpoint configuration.

Local gateway itself remains loopback by default.

---

# 37. Project configuration

Example:

```yaml
version: 1

profile: coding

requirements:
  tools: required
  reasoning: preferred
  vision: optional
  context: 128000

routing:
  provider: openrouter
  model: qwen-model

fallback:
  mode: ask
```

Project config MUST NOT contain secrets.

---

# 38. Configuration precedence

```text
CLI flags
 ↓
project config
 ↓
user config
 ↓
provider defaults
 ↓
built-in defaults
```

Credentials resolve through the credential subsystem, not ordinary config precedence.

---

# 39. Profiles

Examples:

```text
coding
reasoning
cheap
fast
local
offline
```

A profile is a route policy, not a model alias.

Example:

```yaml
profile:
  name: coding
  requirements:
    tools: required
    reasoning: required
  candidates:
    - openrouter/qwen
    - openai/model
```

---

# 40. Import/export/backup

```bash
ccx export
ccx import <file>
ccx backup
ccx restore
```

Export defaults to non-secret data.

Before migration/update/reset:

```text
snapshot
 ↓
validate
 ↓
change
 ↓
verify
```

---

# 41. Config migrations

Schema versions:

```text
v1 → v2 → v3
```

Every migration must:

- preserve old config;
- validate new schema;
- create backup;
- report changes;
- be reversible when practical.

---

# 42. CLI

## Setup

```bash
ccx
ccx init
ccx setup
```

## Providers

```bash
ccx provider list
ccx provider add
ccx provider remove <name>
ccx provider test <name>
```

## Auth

```bash
ccx auth login <provider>
ccx auth status
ccx auth logout <provider>
```

## Models

```bash
ccx models
ccx models refresh
ccx models test <model>
ccx use <model>
```

## Profiles

```bash
ccx profile list
ccx profile edit <name>
ccx profile use <name>
```

## Execution

```bash
ccx run
ccx run --provider <p> --model <m>
ccx run --profile <name>
ccx run --dry-run
ccx run --debug
ccx run --ephemeral
```

## Diagnostics

```bash
ccx doctor
ccx doctor provider <name>
ccx status
ccx logs
ccx logs --live
```

## Project

```bash
ccx project init
ccx project requirements
```

## Runtime

```bash
ccx daemon install
ccx daemon uninstall
ccx daemon start
ccx daemon stop
ccx daemon status
```

## Config

```bash
ccx config get <key>
ccx config set <key> <value>
ccx config edit
```

## Maintenance

```bash
ccx update
ccx rollback
ccx version
```

Every command SHOULD support:

```text
--json
```

where meaningful.

---

# 43. TUI

## Visual identity

Theme:

```text
background: graphite/near-black
primary: high-contrast white
accent: electric cyan
success: green
warning: amber
danger: red
muted: gray
```

Brand should communicate:

```text
routing
connection
infrastructure
reliability
```

Avoid copying Hermes' yellow visual identity.

## Main dashboard

```text
╭─ CCX ─────────────────────────────────────╮
│ ● RUNNING                                 │
│                                           │
│ Profile       coding                      │
│ Provider      OpenRouter                  │
│ Model         Qwen...                     │
│ Protocol      OpenAI Responses            │
│                                           │
│ Tools         ✓                           │
│ Reasoning     ✓                           │
│ Vision        —                           │
│ Cache         ⚠                           │
│                                           │
│ Sessions      2                           │
│ Fallback      ASK                          │
│ Gateway       127.0.0.1:42173             │
│                                           │
│ [M] Models [P] Profiles [D] Doctor        │
│ [L] Logs   [Q] Quit                       │
╰───────────────────────────────────────────╯
```

---

# 44. First-run wizard

```text
CCX
Universal Claude Code Gateway

How do you want to connect?

❯ Cloud providers
  Local models
  Custom endpoint
  Existing gateway
```

After provider selection:

```text
Connecting...

✓ Endpoint reachable
✓ Authentication valid
✓ Models discovered
```

Then searchable model picker:

```text
Qwen...
DeepSeek...
Gemini...
```

Model capabilities shown:

```text
Tools ✓
Reasoning ✓
Vision —
Context 128k
```

Then compatibility test.

Only after verification:

```text
Launch Claude Code? [Y/n]
```

---

# 45. UX progressive disclosure

Beginner:

```text
Connected ✓
Ready ✓
```

Power user:

```text
protocol
headers
transforms
capabilities
upstream request ID
latency
token usage
```

Do not expose protocol complexity until useful.

---

# 46. Error UX

Every error follows:

```text
WHAT
WHY
IMPACT
WHAT CCX CAN DO
ACTION
```

Example:

```text
✗ Provider rejected the request.

Reason:
  parallel tool calls are unsupported.

CCX can serialize tool calls.

Apply compatibility mode? [Y/n]
```

Never dump raw HTTP errors as the primary UX.

Raw details remain available in debug mode.

---

# 47. Doctor

`ccx doctor` must inspect:

```text
installation
Claude Code version
CCX version
gateway
ports
environment
authentication
provider
model
protocol
streaming
tools
reasoning
context
cache
fallback
network
TLS
local runtimes
config
service
```

Example:

```text
CCX Doctor

Installation       ✓
Claude Code        ✓
Gateway             ✓
Authentication      ✓
Provider            ✓
Model               ✓

Streaming           ✓
Tools               ✓
Reasoning           ✓
Prompt cache        ⚠

Issue:
Provider advertises cache but no cache hit observed.

Resolution:
Cache marked GATED. No automatic cache directives will be sent.
```

`--fix` only performs safe, reversible repairs.

---

# 48. Logging/observability

Every request gets:

```text
ccx_request_id
session_id
provider
model
attempt
```

Logs:

```text
ERROR
WARN
INFO
DEBUG
TRACE
```

Debug wire logs MUST redact secrets centrally.

Body capture is opt-in.

Default body logging:

```text
OFF
```

---

# 49. Privacy

Modes:

```text
normal
private
offline
```

Offline:

- no remote update check;
- no remote model catalog;
- only local endpoints.

Private:

- no prompt/response body logs;
- minimum metadata.

CCX MUST NOT transmit analytics by default.

---

# 50. Usage/cost accounting

Track locally:

```text
requests
input tokens
output tokens
cache tokens
latency
provider
model
session
estimated cost
```

Optional limits:

```yaml
budget:
  session: 1.00
  daily: 5.00
```

When exceeded:

```text
stop
ask
fallback_local
```

Never silently spend beyond configured guardrails.

---

# 51. Resource limits

Configurable:

```text
max_sessions
max_concurrency
max_request_bytes
max_stream_buffer
max_tool_argument_bytes
max_log_bytes
max_retries
max_memory
```

CCX must remain stable under malformed or huge provider responses.

---

# 52. Security threat model

Threats:

```text
malicious provider
malformed response
stolen credential
malicious local process
malicious project config
remote network exposure
TLS interception
log leakage
supply-chain attack
```

Controls:

```text
loopback default
least privilege
OS keychain
central redaction
input validation
signed updates
checksum verification
no root
explicit remote binding
```

---

# 53. Project config security

Project configuration may influence routing and capabilities but MUST NOT automatically:

- execute shell commands;
- install providers;
- change credential storage;
- bind CCX to LAN;
- disable security;
- enable raw logging.

Potentially dangerous project settings require explicit user approval.

---

# 54. Update system

Channels:

```text
stable
beta
nightly
```

Transaction:

```text
download
 ↓
verify signature/checksum
 ↓
stage
 ↓
health test
 ↓
activate
 ↓
rollback if unhealthy
```

`ccx update` must never leave a half-updated installation.

---

# 55. Compatibility laboratory

Core matrix:

```text
Claude Code version
×
CCX version
×
protocol
×
provider
×
model
×
feature
```

Test categories:

```text
basic text
multi-turn
system
tools
parallel tools
tool results
reasoning
reasoning + tools
streaming
silent thinking
vision
documents
cache
context
cancellation
429
5xx
timeouts
disconnects
malformed SSE
model discovery
fallback
```

---

# 56. Golden fixtures

Every regression becomes:

```text
Claude Code version
provider
model
feature
request fixture
response fixture
expected IR
expected output
```

The regression corpus is permanent project memory.

---

# 57. Chaos testing

Inject:

```text
429
500
503
timeout
connection reset
partial SSE
invalid JSON
duplicate event
missing event
credential expiry
model disappearance
slow stream
provider restart
CCX restart
machine sleep
```

Assertions:

```text
CCX survives
Claude Code receives valid protocol
no duplicate tool execution
fallback follows policy
secrets remain hidden
```

---

# 58. Fuzzing

Fuzz:

```text
HTTP headers
JSON
SSE
tool schemas
content blocks
provider errors
model metadata
```

Property:

```text
invalid external input
→ controlled error
→ no panic
```

---

# 59. Long-running tests

Minimum stress scenario:

```text
hours-long session
10k+ stream events
hundreds of tool calls
multiple context compactions
multiple concurrent sessions
```

Measure:

```text
memory
CPU
goroutines/threads
file descriptors
latency
buffer growth
```

No unbounded growth.

---

# 60. Provider verification levels

```text
EXPERIMENTAL
DISCOVERED
CONNECTED
VERIFIED
CERTIFIED
```

Certification requires passing the relevant compatibility suite.

Provider catalog must show verification status.

---

# 61. Provider SDK

Later:

```bash
ccx provider scaffold
```

Generate:

```text
provider definition
adapter
fixtures
tests
docs
```

Third-party providers must pass automated compatibility validation.

Dynamic arbitrary executable plugins are NOT required for 1.0; they add supply-chain/security/versioning complexity.

---

# 62. Packaging

Target:

```text
Linux x64
Linux ARM64
macOS x64
macOS ARM64
Windows x64
Windows ARM64
```

Primary artifact:

```text
single native binary
```

Distribution:

```text
GitHub Releases
Homebrew
winget
```

Additional package managers later.

---

# 63. Terminal compatibility

Test:

```text
bash
zsh
fish
PowerShell
cmd
tmux
screen
SSH
VS Code terminal
JetBrains terminal
Windows Terminal
```

Support:

```text
NO_COLOR
ASCII fallback
non-TTY
high contrast
keyboard-only navigation
```

---

# 64. CI/headless

Everything automation-relevant must work without TUI:

```bash
ccx doctor --json
ccx models --json
ccx status --json
ccx test --json
ccx run --dry-run --json
```

CI must be able to fail on compatibility regressions.

---

# 65. Devcontainers

Test all combinations:

```text
Claude Code host
CCX container
```

```text
Claude Code container
CCX host
```

```text
both container
```

Do not assume `localhost` has identical meaning across namespaces.

---

# 66. SSH/remote development

Support:

```text
SSH terminal
remote filesystem
remote Claude Code
remote CCX
local provider
remote provider
```

The TUI must remain functional over ordinary SSH terminals.

---

# 67. Existing gateways

CCX is compatible with:

```text
LiteLLM
Claude Code Router
OpenRouter
organization gateways
self-hosted gateways
```

It can be:

```text
front proxy
```

or:

```text
client-side compatibility layer
```

or:

```text
multi-hop gateway
```

CCX MUST clearly display the resulting route.

---

# 68. Performance targets

Initial targets:

```text
local gateway overhead:
  <10ms median for non-streaming translation excluding provider latency

streaming:
  bounded buffering
  first downstream event forwarded as soon as safely transformable

memory:
  no request-body accumulation beyond configured limits
```

Benchmarks must measure:

```text
passthrough
Anthropic→OpenAI
Anthropic→Gemini
tool-heavy streams
long reasoning streams
large context
```

---

# 69. Failure-mode master checklist

The initial research identified these concrete failure classes and they MUST become regression tests:

1. unknown `anthropic-beta` forwarded to incompatible upstream;
2. silent long-thinking stream without keepalives;
3. gateway discovery failing under certain auth configurations;
4. non-Claude model IDs filtered by Claude Code's picker;
5. custom-model context-window mismatch;
6. unsupported local token-count probe storms;
7. vLLM tool-parser mismatch;
8. prompt-cache capability being gated/inconsistent;
9. cancellation not reaching upstream;
10. secrets appearing in debug logs;
11. silent provider/model switching.

These are source-derived requirements from the existing CCX research baseline. fileciteturn1file0L12-L26

---

# 70. Nuke-and-corners audit

Before 1.0, explicitly test:

```text
What if Claude Code changes tomorrow?
What if a beta header appears?
What if a beta header disappears?
What if a provider lies about capabilities?
What if provider changes model IDs?
What if /v1/models disappears?
What if discovery silently fails?
What if OAuth is used instead of an API key?
What if auth expires mid-stream?
What if the provider returns malformed SSE?
What if the provider sends duplicate tool events?
What if a tool call was executed but response was lost?
What if the laptop sleeps for 8 hours?
What if the provider is down for 30 minutes?
What if CCX crashes during a tool call?
What if 10 Claude Code sessions run simultaneously?
What if 100 sessions run?
What if a local model server disappears?
What if a local model server restarts?
What if context is 1M but Claude Code assumes 200k?
What if the provider reports false token usage?
What if cache is advertised but unavailable?
What if DNS fails?
What if corporate TLS interception exists?
What if IPv6 works but IPv4 fails?
What if IPv4 works but IPv6 fails?
What if port is occupied?
What if config is corrupted?
What if an update is interrupted?
What if migration fails?
What if a project config is malicious?
What if a provider endpoint is malicious?
What if logs fill the disk?
What if a response is enormous?
What if SSE never terminates?
What if a provider never sends keepalives?
What if a user has no network?
What if a user has only local models?
What if the user has an existing gateway?
What if Claude Code changes its picker?
What if Claude Code changes authentication?
What if Claude Code adds a new content block?
What if a future model uses a new reasoning representation?
What if an upstream protocol adds a required field?
What if CCX receives an unknown event?
```

Every question must map to:

```text
test
implementation invariant
or documented unsupported behavior
```

---

# 71. Hard invariants

These MUST become automated tests:

```text
1. No secret leakage.
2. No duplicate side-effecting tool execution caused by retry.
3. No silent model/provider fallback.
4. No LAN exposure by default.
5. No mandatory CCX cloud.
6. No provider response can crash the gateway.
7. No unsupported capability presented as supported.
8. No unbounded request/stream buffering.
9. Ctrl+C propagates cancellation where supported.
10. Config migrations preserve user data.
11. Failed updates can roll back.
12. One session cannot corrupt another session's route.
13. Unknown protocol fields are preserved where safe.
14. Unknown protocol events fail gracefully.
15. All compatibility workarounds are version/target scoped.
```

---

# 72. 1.0 definition of done

CCX 1.0 is NOT complete until:

```text
[ ] Clean-machine installation works.
[ ] First-run setup works without reading documentation.
[ ] Provider authentication works.
[ ] Model discovery works where supported.
[ ] Arbitrary custom model IDs work.
[ ] Claude Code launches through CCX.
[ ] Normal text works.
[ ] Multi-turn works.
[ ] Streaming works.
[ ] Long thinking streams work.
[ ] Tools work.
[ ] Parallel tools work where target supports them.
[ ] Tool results work.
[ ] Cancellation works.
[ ] 429 recovery works.
[ ] Provider outage behavior is deterministic.
[ ] Multiple sessions are isolated.
[ ] CCX restart is recoverable.
[ ] Local runtimes work.
[ ] OpenAI-compatible endpoints work.
[ ] Custom Anthropic endpoints work.
[ ] Diagnostics explain failures.
[ ] Secrets never appear in logs.
[ ] Gateway is loopback-only by default.
[ ] Headless JSON mode works.
[ ] Update + rollback works.
[ ] Compatibility CI exists.
[ ] Regression corpus exists.
[ ] Chaos tests exist.
[ ] Fuzz tests exist.
[ ] Long-session tests pass.
```

---

# 73. Implementation order

## Phase 0 — foundation

```text
repo
Go module
CLI skeleton
config schema
logging
error model
test harness
```

## Phase 1 — native path

```text
local gateway
Anthropic ingress
Anthropic passthrough
session manager
supervisor
```

## Phase 2 — reliability

```text
streaming state machine
keepalive
cancellation
timeouts
error normalization
doctor
```

## Phase 3 — local providers

```text
Ollama
LM Studio
custom Anthropic
```

## Phase 4 — universal translation

```text
IR
OpenAI Chat
OpenAI Responses
OpenRouter
```

## Phase 5 — compatibility engine

```text
capability engine
transforms
model discovery
context safeguards
tool compatibility
```

## Phase 6 — UX

```text
setup wizard
TUI
profiles
model picker
doctor UI
```

## Phase 7 — resilience

```text
retry
fallback
circuit breakers
crash recovery
sleep/wake
```

## Phase 8 — providers

```text
Gemini
Bedrock
Vertex
vLLM
llama.cpp
SGLang
additional mainstream providers
```

## Phase 9 — release engineering

```text
installers
signing
update
rollback
migration
packaging
```

## Phase 10 — compatibility laboratory

```text
Claude Code release matrix
golden fixtures
chaos
fuzzing
stress
nightly regression
```

Only after these should CCX 1.0 be declared.

---

# 74. V2

```text
capability-aware routing
cost-aware routing
latency-aware routing
role-based model routing
advanced fallback
provider scoring
automatic compatibility remediation
provider SDK
community certification
benchmark suite
```

---

# 75. V3

```text
multiple agent protocols
OpenCode
Codex
Aider
other coding agents

shared universal agent IR
agent-specific ingress adapters
agent-specific compatibility profiles
```

---

# 76. V4 / research direction

Potential:

```text
adaptive model routing
task classification
model tournament
parallel candidate generation
judge-based routing
automatic cost optimization
local/cloud hybrid orchestration
agent performance learning
```

These MUST remain separate from the core compatibility mission.

---

# 77. What CCX is NOT

CCX is not:

```text
a hosted LLM service
a model provider
a replacement for Claude Code
a mandatory cloud account
a generic AI chatbot
a promise that every model supports every feature
```

CCX is:

```text
a local compatibility/runtime layer between agent and model ecosystem
```

---

# 78. Strategic differentiation

Existing gateways already demonstrate that basic protocol translation, provider routing, fallbacks, and broad provider access are valuable.

Therefore CCX must not compete solely on:

```text
"supports OpenAI"
"supports 50 providers"
"has a proxy"
```

The differentiator is:

```text
Claude Code compatibility
+
automatic lifecycle
+
capability verification
+
tool/stream correctness
+
excellent terminal UX
+
diagnostics
+
compatibility regression laboratory
```

The strongest long-term asset is the compatibility corpus:

```text
Claude Code version
×
provider
×
model
×
feature
×
observed behavior
```

Every real-world failure becomes a permanent regression test.

---

# 79. AI-agent implementation rules

An AI coding agent implementing CCX MUST:

1. Read this specification before modifying architecture.
2. Prefer existing modules over duplicate abstractions.
3. Never bypass the IR for translated protocols unless a documented fast path exists.
4. Keep protocol and provider logic separate.
5. Add a regression test for every compatibility bug fixed.
6. Add a fixture before changing compatibility behavior.
7. Never introduce silent fallback.
8. Never log credentials.
9. Never assume OpenAI compatibility implies behavioral compatibility.
10. Never treat provider metadata as verified truth.
11. Preserve unknown fields where safe.
12. Reject unsafe ambiguity instead of guessing.
13. Keep platform-specific service logic behind interfaces.
14. Keep TUI presentation separate from runtime logic.
15. Keep CLI commands usable without TUI.
16. Ensure every public feature has automated tests.
17. Update the compatibility matrix when adding provider support.
18. Do not add V2/V3 functionality to the V1 critical path unless the abstraction is required by V1.
19. Prefer small, testable state machines over hidden global state.
20. When Claude Code behavior is undocumented or version-dependent, implement it behind a versioned compatibility rule and add evidence/fixture coverage.

---

# 80. Final product shape

The final experience should feel like this:

```text
FIRST DAY

$ ccx

╭──────────────────────────────────────╮
│ CCX                                  │
│ Universal Claude Code Gateway        │
╰──────────────────────────────────────╯

How do you want to connect?

❯ Cloud
  Local
  Custom
  Existing gateway

        ↓

✓ Provider connected
✓ Models discovered
✓ Model verified
✓ Tools ✓
✓ Reasoning ✓
✓ Streaming ✓

Launch Claude Code? Y

        ↓

$ claude
```

Then every normal day:

```text
$ claude

● CCX
  coding / provider / model
  tools ✓  reasoning ✓  streaming ✓
```

When everything works, CCX is almost invisible.

When something breaks:

```text
✗ Provider rejected parallel tools.

CCX detected:
  model does not support parallel tool calls.

Available:
  [1] Serialize automatically
  [2] Choose another model
  [3] Abort
```

When something truly goes wrong:

```text
$ ccx doctor

Gateway       ✓
Auth          ✓
Provider      ✓
Model         ✓
Streaming     ✓
Tools         ✗

Root cause:
  provider returned malformed tool arguments.

Evidence:
  fixture #CCX-00421

Recommended:
  switch to compatibility mode "strict"
```

That is the standard.

**CCX should feel less like a proxy and more like a local compatibility operating layer for Claude Code.**

---

## Appendix A — Research baseline

The original CCX research plan established the initial architecture, command surface, wizard UX, capability model, 11 concrete failure modes, and risk-ordered roadmap. Its strongest requirements are retained here rather than discarded. fileciteturn1file0L12-L44

The source architecture explicitly separates Claude Code from a loopback CCX gateway, then places the compatibility shim, model-picker bridge, context handling, IR, capability engine, transform pipeline, protocol adapters, resilience engines, credential store and redacted observability inside CCX. fileciteturn1file2L93-L144

The source command surface establishes the intended terminal-first operating model, including `ccx run`, `doctor`, `status`, `logs`, profiles, project configuration, and JSON-friendly operation. fileciteturn1file4L252-L287

The source UX specifically calls for probing the endpoint before model selection, capability verification after model selection, and a compatibility-oriented model picker. fileciteturn1file3L175-L223

The current external evidence reinforces that model discovery and custom-model selection are still moving targets in Claude Code, including documented/custom environment mechanisms and real gateway-auth discovery failures. citeturn0search1turn0search5

---

# END OF CCX MASTER SPECIFICATION


---

# 81. Deep self-review — Revision 1.1

This section is deliberately adversarial. The previous specification was strong architecturally, but it was **not yet implementation-complete enough to hand blindly to an AI coding agent**. The following corrections are mandatory.

## 81.1 Verdict

**Architecture:** strong  
**UX direction:** strong  
**Provider breadth:** strong  
**Reliability thinking:** strong  
**Protocol completeness:** incomplete  
**Security model:** incomplete  
**Lifecycle semantics:** incomplete  
**Compatibility strategy:** incomplete  
**Implementation contract:** incomplete

The biggest risk is not that CCX lacks features. The biggest risk is **building a beautiful proxy that is 90% compatible and fails on the remaining 10% that Claude Code actually exercises**.

The research baseline already correctly identified beta-header drift, keepalives, model filtering, context mismatch, probe storms, vLLM tool-parser failures, gated caching, cancellation, secret leakage, and silent fallback as real failure classes. fileciteturn2file1L51-L65

The current Anthropic API surface demonstrates why the original IR was too small: current Messages content includes images, documents, search results, thinking/redacted thinking, client tools, server tool use, and web-search results, while token counting is a first-class endpoint. citeturn1search0turn1search10

### Critical conclusion

**Do not call the original document “final” without this revision.**

---

# 82. P0 correction — define the compatibility boundary precisely

CCX has two different jobs and they MUST remain separate:

```text
A. Claude Code client compatibility
B. Model/provider protocol translation
```

The architecture must therefore have:

```text
ClaudeCodeCompatLayer
        ↓
AnthropicWireIR
        ↓
UniversalSemanticIR
        ↓
TargetProtocolAdapter
```

Do NOT collapse all three into one generic IR.

Why:

- Anthropic wire details are version-sensitive.
- Claude Code sends client-specific behavior that is not necessarily an Anthropic API invariant.
- Provider translation should not know Claude Code quirks.
- Future agent clients can reuse the semantic IR without pretending their wire contracts are identical.

---

# 83. P0 correction — wire compatibility must be versioned

Create:

```text
compat/claude-code/
    versions/
        v2.1.x/
        v2.2.x/
    headers/
    endpoints/
    env/
    discovery/
    model-resolution/
    fixtures/
```

Every compatibility rule needs:

```yaml
id:
introduced:
last_verified:
applies_to:
evidence:
severity:
action:
```

Example:

```yaml
id: gateway-discovery-requires-auth-env
applies_to: ">=2.1.129"
action: ensure_local_gateway_token
```

Never write:

```text
if claude_code then do X
```

without version scope.

---

# 84. P0 correction — the wire contract is larger than `/v1/messages`

The gateway compatibility surface MUST explicitly model:

```text
POST /v1/messages
POST /v1/messages/count_tokens
GET  /v1/models
```

plus any additional endpoints Claude Code versions actually exercise.

Anthropic's current API documents `/v1/messages/count_tokens` as a real endpoint that counts tokens including tools, images and documents. citeturn1search0

CCX MUST decide per endpoint:

```text
PASS
EMULATE
TRANSLATE
SHORT-CIRCUIT
REJECT
```

For example:

```text
Claude Code → count_tokens
        ↓
CCX
 ├─ provider has compatible counter → forward
 ├─ local tokenizer available        → calculate locally
 ├─ approximation allowed            → estimate + mark approximate
 └─ impossible                       → return explicit unsupported response
```

Never forward a known-invalid probe repeatedly.

---

# 85. P0 correction — expand the content-block model

The IR MUST NOT stop at:

```text
text
image
document
tool_use
tool_result
thinking
redacted_thinking
```

It must support an extensible registry for:

```text
server_tool_use
web_search_tool_result
search_result
citations
code_execution
bash_code_execution
computer-use related blocks
MCP-related blocks
future Anthropic blocks
provider extensions
```

Current Anthropic documentation already exposes server tool use, web-search results, code execution output, images, PDFs/documents, thinking signatures and additional block types. citeturn1search0turn1search10

Unknown block behavior:

```text
known + supported       → transform
known + unsupported     → policy
unknown + safely opaque → preserve
unknown + unsafe        → reject
```

---

# 86. P0 correction — thinking is not just “reasoning text”

Thinking state needs:

```text
type
signature
redacted data
budget
display policy
encryption/opaque payload
ordering
round-trip preservation
```

Current Anthropic documentation explicitly requires thinking blocks/signatures to be passed back unmodified and in original order in applicable flows. citeturn1search0

Therefore:

**CCX MUST preserve opaque thinking artifacts byte-for-byte where the target protocol requires them.**

Never parse, summarize, regenerate, or normalize signatures.

---

# 87. P0 correction — tool retry semantics need a real state machine

The old “don't duplicate tool calls” rule is correct but underspecified.

Implement:

```text
MODEL_REQUEST
  ↓
UPSTREAM_ACCEPTED
  ↓
STREAMING
  ↓
TOOL_CALL_EMITTED
  ↓
DOWNSTREAM_RECEIVED
```

Ambiguous failure:

```text
UPSTREAM_ACCEPTED
      ↓
connection lost
      ↓
UNKNOWN
```

Retry policy:

```text
No tool call observed
  → retry may be allowed

Tool call observed
  → default NO RETRY

Tool call observed + downstream ACK unavailable
  → NO RETRY

Pure text response + connection lost before terminal event
  → retry only under explicit safe policy
```

CCX is not the tool executor. The important risk is **replaying the model request and causing Claude Code to receive a second logically identical tool request**.

---

# 88. P0 correction — response IDs and continuation identity

Track separately:

```text
request_id
upstream_request_id
response_id
message_id
session_id
tool_call_id
attempt_id
```

Never reuse one ID for all layers.

Every translation must maintain a mapping table.

Example:

```text
Claude toolu_x
     ↕
CCX call_ccx_y
     ↕
OpenAI call_z
```

Mapping must survive streaming and error handling.

---

# 89. P0 correction — stream translation needs a formal event journal

Do not implement streaming as ad-hoc callbacks.

Use:

```text
WireEvent
 ↓
EventNormalizer
 ↓
SemanticEvent
 ↓
Transform
 ↓
TargetEvent
 ↓
Serializer
```

Maintain a bounded event journal:

```text
sequence
timestamp
event_type
block_id
call_id
payload_hash
```

This enables:

- duplicate detection;
- missing-event detection;
- debugging;
- fixture generation;
- stream reconstruction.

Never store full payloads by default.

---

# 90. P0 correction — malformed stream recovery

The gateway MUST define behavior for:

```text
partial JSON
partial UTF-8
duplicate event
out-of-order event
missing content_block_stop
missing message_stop
provider disconnect
provider sends terminal error after content
provider sends invalid tool JSON
```

Rules:

```text
recoverable framing issue → repair if unambiguous
semantic corruption       → terminate safely
ambiguous tool state       → never fabricate tool result
```

A translated stream MUST always end in a protocol-valid terminal state whenever possible.

---

# 91. P0 correction — error taxonomy

The error model needs more than:

```text
category
provider_code
retryable
safe_message
request_id
```

Add:

```text
layer:
  client
  ccx
  protocol
  provider
  network
  auth
  policy

phase:
  discovery
  auth
  request
  stream
  transform
  cancellation
  shutdown

retry_class:
  never
  safe
  conditional
  unknown

user_action:
  reauth
  retry
  choose_model
  run_doctor
  disable_feature
  none
```

Anthropic's current API defines distinct 400/401/402/403/404/409/413/429/500/504/529 classes and includes request IDs for diagnosis. citeturn1search1

CCX should preserve upstream error meaning while adding CCX context.

---

# 92. P0 correction — rate limits are multidimensional

Do not model rate limits as one integer.

Track:

```text
RPM
input TPM
output TPM
concurrency
provider-specific quota
account spend cap
acceleration limit
```

Anthropic currently documents RPM, input-token/minute and output-token/minute limits, plus acceleration behavior and `Retry-After` handling. citeturn1search5

CCX's retry engine must distinguish:

```text
temporary capacity
hard spend cap
invalid credentials
model unavailable
provider overloaded
```

A spend cap should not enter an endless retry loop.

---

# 93. P0 correction — capability system needs provenance

Every capability result must include:

```yaml
state: supported|unsupported|gated|unverified|broken
source:
  declared:
  probed:
  observed:
verified_at:
confidence:
evidence_id:
scope:
  provider:
  model:
  credential:
  protocol:
```

The important part is `scope`.

Example:

```text
cache:
  provider: OpenRouter
  model: X
  credential: account-A
  state: gated
```

Do not globally mark model X as cache-supported forever.

---

# 94. P0 correction — probes need safety classes

Every probe must declare:

```text
network_cost
token_cost
side_effect_risk
privacy_risk
duration
```

Probe classes:

```text
SAFE
BILLABLE
DESTRUCTIVE
NEVER_AUTOMATIC
```

Default automatic:

```text
SAFE only
```

Standard:

```text
SAFE + low-cost BILLABLE
```

Deep:

```text
explicit user approval
```

Never use a real external tool execution as a generic capability probe.

---

# 95. P0 correction — model discovery needs source precedence

Define:

```text
manual explicit model
        >
project alias
        >
provider discovery
        >
provider catalog
        >
CCX registry
        >
stale cache
```

But stale cache MUST be visibly marked.

Each model record needs:

```text
canonical_id
provider_id
display_name
aliases
protocol
context
capabilities
pricing
region
availability
source
observed_at
expires_at
```

A stale model MUST never be silently presented as currently available.

---

# 96. P0 correction — model identity must be canonical

Define:

```text
provider/model
provider/region/model
gateway/provider/model
```

as separate identifiers.

Internally:

```text
ModelRef {
  provider
  canonical_id
  requested_id
  public_alias
}
```

This avoids destructive string-prefix hacks.

Prefix stripping can be a compatibility transform, but it must not become the model identity system.

---

# 97. P0 correction — authentication architecture

Authentication needs:

```text
CredentialProvider
CredentialResolver
CredentialStore
CredentialLease
CredentialRedactor
```

Credential sources:

```text
env
OS keychain
OAuth
AWS chain
GCP ADC
Azure identity
custom command
manual secret
```

The resolver must return a temporary in-memory credential object.

Adapters must never directly read arbitrary config files.

---

# 98. P0 correction — secret handling needs memory discipline

For secrets:

```text
never log
never serialize into ordinary config
never include in crash reports
never include in telemetry
never show in `doctor`
```

Where language/runtime permits:

```text
minimize lifetime
zero buffers where practical
avoid immutable copies
```

Do not promise perfect memory zeroization if the runtime cannot guarantee it; document the limitation honestly.

---

# 99. P0 correction — local gateway authentication

Even loopback services need authentication.

Threat:

```text
malicious local process
```

Therefore:

```text
127.0.0.1
+
random high-entropy session token
+
constant-time validation
```

Token must be:

- generated per runtime;
- passed only to child Claude Code process;
- never written to normal logs;
- invalidated when runtime exits.

---

# 100. P0 correction — remote exposure requires an explicit security profile

If the user deliberately exposes CCX over LAN:

```text
bind: 0.0.0.0
```

that MUST require:

```text
--allow-remote
```

and a security checklist:

```text
strong auth
TLS or trusted tunnel
allowed networks
rate limit
audit log
no anonymous access
```

Default remains loopback.

---

# 101. P0 correction — config is a state machine

Config changes should be:

```text
READ
VALIDATE
PLAN
SNAPSHOT
APPLY
VERIFY
COMMIT
```

Failure:

```text
ROLLBACK
```

Never mutate live state halfway through a failed provider edit.

---

# 102. P0 correction — active sessions must be immutable

This is a major correction.

A running session owns:

```text
route snapshot
credential reference
model
protocol
capability policy
transform policy
retry policy
```

Changing global config affects:

```text
future sessions
```

not:

```text
existing sessions
```

unless an explicit hot-reload operation says otherwise.

---

# 103. P0 correction — supervisor and gateway need separate health

Health states:

```text
PROCESS_HEALTH
GATEWAY_HEALTH
UPSTREAM_HEALTH
SESSION_HEALTH
```

Example:

```text
CCX process ✓
gateway ✓
OpenRouter ✗
active session DEGRADED
```

Do not represent the whole system as one green/red state.

---

# 104. P0 correction — startup locking

Multiple concurrent:

```bash
ccx
ccx
ccx
```

must not launch three gateways accidentally.

Use:

```text
runtime lock
PID
health endpoint
owner metadata
```

Handle stale locks safely.

---

# 105. P0 correction — port ownership

Runtime metadata must include:

```yaml
pid:
port:
bind:
started_at:
version:
runtime_token_id:
```

On startup:

```text
lock exists
 ↓
process alive?
 ├─ yes → reuse
 └─ no → stale → recover
```

Never kill an unrelated process just because it occupies the expected port.

---

# 106. P0 correction — updates must be compatibility-aware

Before update:

```text
backup config
record active CCX version
record Claude Code version
record provider state
```

After update:

```text
start isolated health check
run protocol smoke test
run local fixture suite
```

If failure:

```text
rollback CCX
```

Do not automatically roll back Claude Code itself.

---

# 107. P1 correction — no undocumented Claude Code filesystem surgery by default

The earlier architecture proposed pre-seeding:

```text
~/.claude/cache/gateway-models.json
```

That is too risky as a default architecture.

Treat client filesystem mutation as:

```text
EXPERIMENTAL
VERSION-GATED
BACKED-UP
OPT-IN
```

Prefer:

```text
ANTHROPIC_MODEL
custom model environment mechanisms
documented gateway discovery
CCX launch-time configuration
```

Current real-world evidence shows gateway discovery and model resolution behavior changes across Claude Code releases, including model-list-dependent subagent resolution. citeturn0search1

---

# 108. P1 correction — define the “Claude Code launcher”

The phrase “run `claude` normally and CCX automatically appears” is too magical.

Choose one explicit supported mechanism:

### Option A — shell shim

```text
ccx install-shell
```

### Option B — Claude Code settings/environment integration

CCX writes only documented configuration.

### Option C — OS service + environment launcher

### Option D — `ccx claude`

```bash
ccx claude
```

For 1.0:

**Recommended: `ccx claude` + optional shell integration.**

This is deterministic and does not hijack arbitrary `claude` executions.

---

# 109. P1 correction — CLI semantics need one canonical entry point

Current spec has:

```text
ccx run
ccx
claude
```

Define:

```text
ccx
    → interactive launcher/setup

ccx claude [args...]
    → guaranteed CCX-wrapped Claude Code

ccx run
    → lower-level runtime/agent command
```

This removes ambiguity.

---

# 110. P1 correction — `ccx switch` should not mutate active sessions

Add:

```bash
ccx use <model>
ccx use --profile <profile>
```

These set the default for future sessions.

Do NOT promise hot switching of a running Claude Code conversation unless Claude Code itself supports it.

---

# 111. P1 correction — provider registry vs protocol registry

Separate:

```text
Protocol:
  anthropic
  openai-chat
  openai-responses
  gemini

Provider:
  openai
  openrouter
  groq
  local-ollama
  custom-x
```

A provider may expose multiple protocols.

A protocol may be used by hundreds of providers.

This separation is mandatory for maintainability.

---

# 112. P1 correction — custom providers need a declarative path

Most custom providers should require:

```yaml
provider:
  id:
  protocol:
  base_url:
  auth:
  models:
```

No code/plugin required.

Code plugins are for:

```text
custom signing
custom OAuth
custom stream dialect
custom model discovery
custom routing
```

This dramatically lowers the barrier for community support.

---

# 113. P1 correction — OpenAI Responses must not be treated as Chat++

Create separate adapter implementations.

They share reusable primitives:

```text
tool mapper
content mapper
usage mapper
error mapper
stream event mapper
```

but remain distinct protocol modules.

---

# 114. P1 correction — Gemini needs a dedicated semantic mapper

Do not implement:

```text
Gemini = OpenAI with renamed fields
```

Use:

```text
Anthropic
  ↓
Semantic IR
  ↓
Gemini contents/tools/parts
```

with explicit support for:

```text
function calling
tool results
system instruction
multimodal parts
thinking
streaming
finish reasons
usage
```

Unsupported semantics must be surfaced.

---

# 115. P1 correction — “any model” needs a precise promise

The marketing/product claim should be:

> Any model that can be represented through a supported protocol and whose required Claude Code capabilities can be preserved or safely transformed.

NOT:

> Literally every AI model.

A model without tool calling cannot magically become a reliable Claude Code agent model.

This distinction protects the project's credibility.

---

# 116. P1 correction — transform policies need deterministic contracts

Every transform:

```yaml
id:
input:
output:
loss:
conditions:
reversible:
risk:
default:
```

Example:

```text
parallel_tools → serialized_tools

loss:
  concurrency

risk:
  latency increase

default:
  allowed only if policy=degrade
```

Never hide semantic loss.

---

# 117. P1 correction — compatibility report must show loss

Use:

```text
✓ preserved
≈ transformed
⚠ degraded
✗ unavailable
? unverified
```

Example:

```text
Tools          ✓ preserved
Parallel tools ≈ serialized
Thinking       ⚠ transformed
Vision         ✗ unavailable
Cache          ? unverified
```

This is better than fake “9/10” scoring.

---

# 118. P1 correction — remove score language

The previous examples use:

```text
Core 10/10
Agentic 7/8
```

Remove this.

It creates false precision.

Use a capability matrix instead.

---

# 119. P1 correction — observability needs OpenTelemetry-style concepts without requiring OpenTelemetry

Internally define:

```text
trace_id
span_id
request_id
session_id
provider_request_id
```

But CCX does not need to depend on an external telemetry backend.

Export formats can be added later.

---

# 120. P1 correction — metrics need concrete names

Define metrics such as:

```text
ccx_requests_total
ccx_requests_failed_total
ccx_request_duration_seconds
ccx_stream_first_event_seconds
ccx_stream_duration_seconds
ccx_upstream_retries_total
ccx_fallbacks_total
ccx_tool_calls_observed_total
ccx_capability_probe_total
ccx_provider_health
ccx_active_sessions
ccx_bytes_in_total
ccx_bytes_out_total
```

All local.

---

# 121. P1 correction — diagnostics must produce a support bundle

Add:

```bash
ccx doctor --bundle
```

Bundle:

```text
CCX version
Claude Code version
OS
architecture
sanitized config
provider/protocol identifiers
capability results
recent error codes
compatibility fixture IDs
```

Never include:

```text
API keys
OAuth tokens
prompt bodies
tool arguments
file contents
```

unless the user explicitly exports a separately labeled sensitive bundle.

---

# 122. P1 correction — provider support matrix must be generated

Do not manually maintain a huge README table.

Generate it from machine-readable provider metadata:

```yaml
provider:
protocols:
auth:
discovery:
streaming:
tools:
reasoning:
vision:
cache:
verification:
```

Then generate:

```text
docs/providers.md
TUI provider cards
doctor output
```

from the same source.

---

# 123. P1 correction — provider versions matter

Provider compatibility should be:

```text
provider
provider API revision
model
model revision/family
CCX adapter version
```

Local runtimes additionally need:

```text
runtime version
server flags
tool parser
```

---

# 124. P1 correction — local model servers need endpoint identity

For local providers:

```text
host
port
runtime
version
model
server configuration hash
```

A vLLM model with a different server command can behave differently even with the same model ID.

---

# 125. P1 correction — security needs SSRF defense

Because CCX accepts arbitrary custom endpoints, protect against:

```text
localhost targets
cloud metadata endpoints
private network scanning
IPv6 loopback
DNS rebinding
redirect-to-private-IP
```

Custom endpoint probing must have explicit network policy.

Example:

```text
allow_public
allow_private
allow_loopback
```

Default for custom remote providers:

```text
explicit user approval for private-network targets
```

---

# 126. P1 correction — redirects

Never blindly follow arbitrary HTTP redirects for credential-bearing requests.

Policy:

```text
same-origin redirect → allowed
cross-origin redirect → re-evaluate auth + explicit policy
private/public boundary → reject by default
```

---

# 127. P1 correction — TLS

Support:

```text
system CA
custom CA
mTLS
TLS minimum version
proxy CA
```

But never offer:

```text
--insecure
```

as a silent workaround.

If insecure TLS exists for development:

```text
explicit flag
loud warning
never persisted silently
```

---

# 128. P1 correction — request size and decompression bombs

Limits must apply to:

```text
compressed body
decompressed body
JSON nesting
string lengths
array counts
SSE event size
tool argument size
header size
```

Never let gzip/brotli decompression bypass resource limits.

---

# 129. P1 correction — model catalog poisoning

Provider-supplied model metadata is untrusted.

Never allow provider metadata to directly cause:

```text
shell execution
plugin installation
filesystem writes
credential access
network access
```

Catalog is data.

---

# 130. P1 correction — project config trust boundary

A repository can contain:

```text
.ccx/config.yaml
```

Therefore opening an untrusted repository must NOT automatically activate dangerous routing/auth/network settings.

Add:

```text
trusted project
untrusted project
```

First encounter:

```text
This project contains CCX configuration.

Trust this project? [y/N]
```

Safe fields may be read without trust.

---

# 131. P1 correction — environment poisoning

Environment variables can override provider endpoints.

Doctor must display:

```text
effective source
```

for important values:

```text
model → project
base URL → user
credential → keychain
```

Do not print secret values.

---

# 132. P1 correction — shell injection

Never construct shell commands by concatenating:

```text
provider
model
URL
credential
project path
```

Use argv arrays/process APIs.

This applies especially to:

```text
custom credential commands
ccx claude
shell integration
update installer
```

---

# 133. P1 correction — credential command sandbox

If supporting:

```yaml
credential:
  command: ...
```

make it:

```text
opt-in
explicit
documented
never project-controlled by default
```

Project config must never be allowed to choose arbitrary credential commands automatically.

---

# 134. P1 correction — offline mode semantics

Define exact behavior:

```text
offline:
  provider endpoints allowed: local/private explicitly configured
  update checks: disabled
  catalog refresh: disabled
  telemetry: disabled
```

Do not call `offline` a guarantee unless CCX itself can enforce that no network request occurs except explicitly permitted targets.

---

# 135. P1 correction — privacy mode must include request bodies

Privacy mode:

```text
no prompt logs
no response logs
no tool argument logs
no model input snapshots
```

Hashes/metadata only.

---

# 136. P1 correction — usage accounting must distinguish estimated vs authoritative

Every token/cost number needs:

```text
source:
  provider
  CCX estimate
  local tokenizer

accuracy:
  authoritative
  estimated
  unavailable
```

Never present an estimate as an invoice.

---

# 137. P1 correction — pricing cannot be a hardcoded truth

Pricing metadata should include:

```text
effective_from
effective_until
currency
input_price
output_price
cached_input_price
source
```

Cost is:

```text
estimate
```

unless directly supplied by provider billing data.

---

# 138. P1 correction — context management needs a three-way distinction

Separate:

```text
provider_context_limit
claude_code_context_limit
ccx_safe_context_limit
```

The safe limit is:

```text
min(provider, client, policy)
```

unless a documented client override exists.

Never claim CCX can magically increase Claude Code's internal limit.

---

# 139. P1 correction — compaction must be treated as client behavior

CCX can expose correct metadata where supported.

It must NOT implement its own hidden conversation compaction and pretend Claude Code performed it.

Otherwise conversation semantics diverge.

---

# 140. P1 correction — MCP/server tools

Because modern Claude Code workflows can involve MCP and server-side tools, the compatibility model must explicitly distinguish:

```text
client-executed tool
server-executed tool
provider-native tool
CCX-emulated tool
```

CCX MUST NOT blindly translate a server-side tool into a client-side tool.

If a target provider cannot execute a server-native tool:

```text
reject
or explicit emulation
```

Never silently fabricate it.

---

# 141. P1 correction — citations and documents

Document/image translation must define:

```text
mime type
base64
URL
file reference
size limit
provider support
```

Do not silently convert a document to text unless policy allows semantic loss.

---

# 142. P1 correction — streaming backpressure

Define:

```text
upstream reader
bounded channel
transform queue
downstream writer
```

If downstream is slow:

```text
bounded buffering
backpressure
cancellation
```

Never use an unbounded in-memory stream queue.

---

# 143. P1 correction — graceful shutdown

Shutdown sequence:

```text
STOP ACCEPTING
 ↓
signal active sessions
 ↓
cancel upstream requests
 ↓
wait grace period
 ↓
close listeners
 ↓
flush safe logs/metrics
 ↓
release lock
 ↓
exit
```

Forced termination:

```text
hard stop
```

must not attempt unsafe cleanup.

---

# 144. P1 correction — signal handling

Support:

```text
SIGINT
SIGTERM
SIGHUP where meaningful
Windows console events
```

Semantics must be documented.

`SIGHUP` should not accidentally destroy a user's active session unless explicitly configured.

---

# 145. P1 correction — database/storage choice

The spec currently lacks a concrete persistence strategy.

Recommended V1:

```text
config.yaml
credentials → OS keychain
runtime state → small local state DB/file
logs → rotating files
compatibility cache → versioned JSON/SQLite
```

Do not introduce a database merely for complexity.

Use SQLite only if atomic querying/history becomes necessary.

---

# 146. P1 correction — filesystem layout

Define one canonical layout:

```text
~/.ccx/
  config.yaml
  state/
  cache/
  logs/
  backups/
  runtime/
  compat/
```

Project:

```text
<repo>/.ccx/config.yaml
```

No random writes into `~/.claude/` except explicitly version-gated compatibility integrations.

---

# 147. P1 correction — lock/config permissions

On Unix:

```text
config: 0600 when secrets could ever be present
credentials: OS keychain
state: user-only
logs: user-only
runtime token: user-only
```

On Windows use equivalent user ACL restrictions.

---

# 148. P1 correction — uninstall must be first-class

Add:

```bash
ccx uninstall
```

Must show:

```text
binary
shell integration
services
config
cache
logs
credentials
```

with independent choices.

Never delete credentials by default.

---

# 149. P1 correction — reset must be safe

Add:

```bash
ccx reset
```

Modes:

```text
reset runtime
reset cache
reset config
reset everything
```

Always preview before destructive operations.

---

# 150. P1 correction — migration testing

For every config schema:

```text
vN fixture
 ↓
migration
 ↓
vN+1 validation
 ↓
semantic equivalence
```

Keep migration fixtures permanently.

---

# 151. P1 correction — deterministic test recording

Add:

```bash
ccx capture
ccx replay <fixture>
```

Purpose:

```text
capture a sanitized wire interaction
replay it locally
```

This becomes extremely valuable for debugging provider regressions without repeatedly spending API money.

Sensitive capture must be opt-in.

---

# 152. P1 correction — compatibility test levels

Define:

```text
L0 static schema
L1 wire fixture
L2 mock provider
L3 local runtime
L4 live provider
L5 real Claude Code
```

CI should run:

```text
L0-L3 every PR
L4 nightly/controlled
L5 scheduled compatibility matrix
```

This makes the test program financially realistic.

---

# 153. P1 correction — Claude Code binary acquisition

The spec does not define how CCX finds Claude Code.

Implement:

```text
PATH discovery
explicit path
platform-specific standard locations
version detection
```

Do not download Claude Code automatically.

CCX must work with the user's existing installation.

---

# 154. P1 correction — version compatibility policy

Maintain:

```text
supported
best-effort
known-broken
unknown
```

Example:

```text
Claude Code 2.1.245 → VERIFIED
Claude Code 2.1.246 → UNKNOWN
Claude Code 2.0.x   → BEST-EFFORT
```

If known-broken:

```text
CCX refuses launch by default
--force allows explicit override
```

---

# 155. P1 correction — release compatibility watch

The existing roadmap correctly proposes tracking Claude Code releases and automatically flagging behavior/header changes. fileciteturn2file0L17-L31

Make this concrete:

```text
nightly:
  fetch release metadata
  compare documented API/gateway changes
  run fixture suite
  classify:
    unchanged
    additive
    behavior change
    breaking
```

Open a compatibility issue automatically in the repository when maintainers choose to enable that workflow.

---

# 156. P1 correction — provider certification must expire

A provider certificate should have:

```text
verified_at
expires_at
tested_version
tested_model
tested_features
```

A certification can become:

```text
STALE
```

without becoming false.

---

# 157. P1 correction — UX needs a persistent “route pill”

Every TUI page should show:

```text
CCX ●  coding  ›  OpenRouter  ›  Qwen...
```

This prevents users from forgetting what is actually running.

---

# 158. P1 correction — TUI needs navigation rules

Keyboard:

```text
↑↓        navigate
Enter     select
Esc       back
/         search
?         help
r         refresh
d         doctor
l         logs
p         providers
m         models
q         quit
```

No mouse dependency.

Every screen must have:

```text
title
breadcrumb
primary action
back action
status
```

---

# 159. P1 correction — TUI must work at tiny terminal sizes

Minimum supported:

```text
80×24
```

At smaller dimensions:

```text
compact mode
```

Never render overflowing borders/text.

---

# 160. P1 correction — color is not semantic

Every state must have:

```text
icon/symbol
text label
color
```

So:

```text
✓ VERIFIED
⚠ DEGRADED
✗ FAILED
? UNVERIFIED
```

works without color.

---

# 161. P1 correction — accessibility

Support:

```text
NO_COLOR
high contrast
screen-reader-friendly text mode
keyboard-only navigation
ASCII fallback
```

TUI should degrade to a readable textual interface.

---

# 162. P1 correction — zero-config path

The ideal user should be able to:

```bash
ccx claude
```

and receive:

```text
No route configured.

Choose:
  [1] Anthropic
  [2] OpenRouter
  [3] Local
  [4] Custom
```

After successful setup:

```bash
ccx claude
```

should not ask again.

---

# 163. P1 correction — shell integration must be reversible

If installed:

```bash
ccx install-shell
```

also provide:

```bash
ccx uninstall-shell
```

Show exactly what files/lines are modified.

Never silently overwrite shell configuration.

---

# 164. P1 correction — command passthrough

`ccx claude` MUST preserve:

```text
all Claude Code CLI arguments
stdin
stdout
stderr
exit code
signals
TTY
working directory
environment
```

This is a critical acceptance test.

---

# 165. P1 correction — exit codes

Define CCX exit codes:

```text
0 success
1 general failure
2 invalid CLI usage
3 configuration error
4 authentication error
5 provider unavailable
6 compatibility failure
7 interrupted
8 security policy rejection
9 update/migration failure
```

When wrapping Claude Code, preserve Claude Code's exit code unless CCX itself fails before launch.

---

# 166. P1 correction — no accidental TUI contamination

When Claude Code runs:

```text
CCX TUI MUST NOT capture or rewrite Claude Code's interactive terminal UI.
```

The TUI is for setup/status/diagnostics.

`ccx claude` should hand the terminal to Claude Code cleanly.

This is essential.

---

# 167. P1 correction — stdin/stdout discipline

CCX internal diagnostics must go to:

```text
stderr
```

unless the user explicitly requests JSON/output mode.

Never corrupt Claude Code's stdout stream.

---

# 168. P1 correction — JSON schema for automation

Every JSON command should have a version:

```json
{
  "schema_version": 1,
  "data": {}
}
```

Never make scripts depend on human-readable strings.

---

# 169. P1 correction — machine-readable errors

JSON error:

```json
{
  "schema_version": 1,
  "error": {
    "code": "PROVIDER_AUTH_EXPIRED",
    "message": "...",
    "retryable": false,
    "action": "reauth"
  }
}
```

Stable `code`, unstable human `message`.

---

# 170. P1 correction — test the installer, not just the binary

Clean-machine tests:

```text
Linux
macOS
Windows
```

Verify:

```text
install
upgrade
rollback
uninstall
shell integration
daemon
permissions
PATH
first run
```

---

# 171. P1 correction — architecture language should stop overpromising “daemon”

The runtime should be:

```text
supervised local service
```

A persistent daemon is optional.

This better matches:

```text
AUTO
EPHEMERAL
PERSISTENT
```

and avoids forcing service installation on beginners.

---

# 172. P1 correction — local-only does not mean “no remote traffic”

Clarify:

> CCX has no hosted backend. It may send user requests directly to the user's selected provider endpoint.

This distinction must appear prominently in docs.

---

# 173. P1 correction — privacy disclosure

First-run should state:

```text
CCX does not host your requests.

Your prompts are sent directly to the provider you configure,
through your local CCX process.
```

For multi-hop:

```text
Claude Code → CCX → Gateway → Provider
```

show every hop.

---

# 174. P1 correction — provider terms are outside CCX

CCX cannot guarantee:

```text
provider privacy
provider retention
provider billing
provider availability
```

Doctor/docs should identify:

```text
CCX-controlled
provider-controlled
Claude-Code-controlled
```

---

# 175. P1 correction — route explainability

Add:

```bash
ccx route explain
```

Output:

```text
Claude Code
 ↓
CCX local gateway
 ↓
OpenRouter
 ↓
Qwen model

Transforms:
  Anthropic tools → OpenAI tools
  thinking        → provider reasoning
  parallel tools → preserved

Policy:
  fallback = ask
  retries = 2
```

This is a killer debugging feature.

---

# 176. P1 correction — dry-run should be extremely useful

`ccx run --dry-run` should show:

```text
Claude Code version
effective profile
effective provider
effective model
protocol path
capability matrix
transforms
losses
retry policy
fallback policy
credential source
network target
```

No request is sent.

---

# 177. P1 correction — preflight

Add:

```bash
ccx preflight
```

Fast checks:

```text
config
provider
auth
model
compatibility
```

No expensive capability probes unless requested.

---

# 178. P1 correction — doctor levels

```bash
ccx doctor
ccx doctor --standard
ccx doctor --deep
ccx doctor --network
ccx doctor --provider <p>
```

Default must be cheap.

---

# 179. P1 correction — provider probe caching

Probe result cache:

```text
TTL
provider/model scoped
credential scoped
CCX version scoped
```

A new CCX version invalidates compatibility-sensitive probes.

---

# 180. P1 correction — no automatic deep probing

This is essential because provider APIs can charge.

The research baseline already calls out configurable probe aggressiveness because probes can consume paid tokens. fileciteturn2file7L423-L431

Therefore:

```text
default = minimal
```

and deep tests require explicit action.

---

# 181. P1 correction — fallback must be capability-aware

Before fallback:

```text
candidate model
 ↓
requirements
 ↓
effective capability check
 ↓
cost/policy check
 ↓
user policy
```

Never:

```text
provider failed → random other model
```

---

# 182. P1 correction — fallback must preserve user intent

If the active model supports:

```text
vision
```

and fallback does not:

```text
DO NOT fallback silently.
```

If the user explicitly permits degraded fallback:

```text
show degradation
```

---

# 183. P1 correction — retry and fallback are separate

Retry:

```text
same route
```

Fallback:

```text
different route
```

They must have separate budgets.

Example:

```text
retry_budget: 2
fallback_budget: 1
```

---

# 184. P1 correction — no retry after semantic partial output by default

If provider has already produced meaningful content:

```text
do not automatically replay
```

unless the request is proven safe.

This prevents duplicated model actions/tool proposals.

---

# 185. P1 correction — circuit breaker state must be persistent enough

Do not persist transient failure state forever.

Use:

```text
memory
short TTL
```

with optional persisted health history.

---

# 186. P1 correction — provider outage should not break CCX

A dead provider is:

```text
UPSTREAM_UNAVAILABLE
```

not:

```text
CCX_BROKEN
```

The UX must clearly separate them.

---

# 187. P1 correction — upstream request IDs

Always surface:

```text
CCX request ID
upstream request ID
```

when available.

Anthropic documents a `request-id` response header and corresponding error-body request ID. citeturn1search1

---

# 188. P1 correction — compatibility fixture IDs become user-facing

When a known failure happens:

```text
Compatibility issue CCX-COMP-00421
```

Doctor can say:

```text
Known regression.
```

This makes support dramatically easier.

---

# 189. P1 correction — issue templates

Repository should have:

```text
provider bug
Claude Code compatibility bug
streaming bug
tool-call bug
auth bug
TUI bug
security report
```

with automated diagnostic bundle attachment support.

---

# 190. P1 correction — security disclosure

Create:

```text
SECURITY.md
```

with:

```text
supported versions
reporting channel
credential leakage procedure
response expectations
```

---

# 191. P1 correction — dependency security

CI must run:

```text
dependency audit
SBOM generation
license scan
static analysis
secret scan
```

Release artifacts should have provenance metadata.

---

# 192. P1 correction — supply-chain hardening

Release:

```text
checksums
signatures
SBOM
build provenance
```

Installer must verify what it downloads.

---

# 193. P1 correction — plugin security

If plugins are introduced later:

```text
signed plugins
versioned API
permissions
sandbox where possible
explicit install
```

No arbitrary plugin execution from project config.

---

# 194. P1 correction — “provider logo” is cosmetic

Provider branding should be optional metadata:

```text
logo
display_name
website
```

The TUI must still work perfectly in:

```text
ASCII
NO_COLOR
offline
```

Do not make branding architecture-critical.

---

# 195. P1 correction — documentation structure

Repository docs should become:

```text
README.md
QUICKSTART.md
ARCHITECTURE.md
PROTOCOLS.md
PROVIDERS.md
COMPATIBILITY.md
SECURITY.md
PRIVACY.md
TROUBLESHOOTING.md
DEVELOPING.md
RELEASING.md
```

The huge master specification remains the engineering source of truth.

---

# 196. P1 correction — generated compatibility docs

Generate:

```text
Claude Code compatibility matrix
provider matrix
protocol feature matrix
```

from machine-readable metadata.

No duplicate truth.

---

# 197. P1 correction — benchmark suite

Benchmark:

```text
passthrough latency
translation latency
first-byte latency
stream throughput
tool-call translation
large-context memory
100 concurrent sessions
```

Report:

```text
p50
p95
p99
CPU
memory
allocations
```

---

# 198. P1 correction — realistic concurrency targets

Do not promise “100 sessions” as a normal laptop workload.

Define tiers:

```text
small:
  1–5

medium:
  6–20

stress:
  21–100
```

Pass/fail separately.

---

# 199. P1 correction — memory limits need behavior

If memory pressure occurs:

```text
reject new request
```

rather than:

```text
OOM crash
```

Expose:

```text
RESOURCE_LIMIT
```

with a useful action.

---

# 200. P1 correction — file descriptor limits

Doctor should detect:

```text
ulimit
Windows handle availability
```

when concurrency is high.

---

# 201. P1 correction — test terminal resizing

During TUI:

```text
resize 80×24
resize 120×40
resize 200×60
```

No panic or corrupted state.

---

# 202. P1 correction — TTY passthrough tests

`ccx claude` acceptance test:

```text
stdin bytes preserved
stdout bytes preserved
stderr semantics preserved
exit code preserved
Ctrl+C works
terminal resize works
colors work
raw mode works
```

This deserves its own test suite.

---

# 203. P1 correction — Windows needs process-tree handling

On Windows, stopping CCX must not leave child Claude Code processes orphaned.

Test:

```text
CCX
 └─ Claude Code
     └─ subprocesses
```

and ensure process-group/job-object semantics are correct.

---

# 204. P1 correction — macOS/Linux child process behavior

Likewise test:

```text
signals
process groups
terminal ownership
```

for Unix systems.

---

# 205. P1 correction — auto-start must not create zombie processes

Repeated:

```bash
ccx claude
```

should reuse one runtime.

After all sessions exit:

```text
auto mode may idle/stop
```

according to policy.

---

# 206. P1 correction — runtime idle policy

Add:

```yaml
runtime:
  mode: auto
  idle_timeout: 10m
```

Options:

```text
never
5m
10m
30m
1h
```

Persistent daemon ignores idle shutdown.

---

# 207. P1 correction — startup latency target

Measure:

```text
cold CCX startup
warm CCX startup
cold provider connection
```

Do not confuse:

```text
CCX startup
```

with:

```text
model first token
```

---

# 208. P1 correction — cache invalidation

Caches require explicit keys:

```text
CCX version
Claude Code version
provider
model
protocol
credential identity hash
feature set
```

Never reuse a capability result after a materially relevant change.

---

# 209. P1 correction — clock skew

TTL systems should handle:

```text
clock moved backwards
sleep/wake
NTP correction
```

Use monotonic clocks for durations.

---

# 210. P1 correction — timezone should not affect protocol state

Persist timestamps as:

```text
UTC
```

Display in local time.

---

# 211. P1 correction — request timestamps

Record:

```text
accepted_at
upstream_started_at
first_event_at
last_event_at
completed_at
```

This makes latency diagnosis precise.

---

# 212. P1 correction — request body hashing

For privacy-safe correlation:

```text
body_hash
```

can identify repeated requests without storing content.

Use a keyed hash if hashes could expose sensitive equality relationships.

---

# 213. P1 correction — config interpolation

If supporting:

```text
${ENV_VAR}
```

define:

```text
allowed variables
escaping
missing variable behavior
secret redaction
```

Do not allow arbitrary command substitution.

---

# 214. P1 correction — YAML hazards

Config parser must reject or safely handle:

```text
anchors
aliases
duplicate keys
unexpected types
huge nesting
```

Avoid YAML features that create surprising execution semantics.

---

# 215. P1 correction — provider URL normalization

Normalize:

```text
scheme
host
port
path
trailing slash
```

but do not destructively rewrite user intent.

For example:

```text
https://host/api
```

must not become:

```text
https://host/
```

and lose `/api`.

---

# 216. P1 correction — endpoint joining

Never use naive string concatenation:

```text
base + "/v1/messages"
```

Implement URL resolution that understands:

```text
https://host
https://host/
https://host/api
https://host/api/
```

This is a classic custom-provider bug.

---

# 217. P1 correction — auth header collisions

If a user configures:

```text
Authorization
x-api-key
```

and the adapter also generates auth:

```text
explicit policy
```

must decide which wins.

Never send duplicate contradictory auth headers.

---

# 218. P1 correction — header allowlist/denylist

Define:

```text
safe pass-through headers
provider-generated headers
hop-by-hop headers
sensitive headers
```

Never blindly proxy:

```text
Connection
Keep-Alive
Transfer-Encoding
Upgrade
Proxy-Authorization
```

or arbitrary credential headers.

---

# 219. P1 correction — response header policy

Do not blindly forward provider headers to Claude Code.

Define:

```text
preserve
translate
strip
```

especially for:

```text
content-length
transfer-encoding
connection
set-cookie
server
```

---

# 220. P1 correction — SSE framing correctness

Test:

```text
multi-line data
empty data
comments
UTF-8 boundaries
CRLF
LF
chunk boundaries
event ordering
```

Do not assume one TCP chunk equals one SSE event.

---

# 221. P1 correction — compression

Handle:

```text
gzip
br
identity
```

without buffering the entire response.

---

# 222. P1 correction — HTTP/2 and HTTP/1.1

Provider clients should support:

```text
HTTP/1.1
HTTP/2
```

where the runtime library safely supports them.

Do not assume HTTP/2 availability.

---

# 223. P1 correction — proxy authentication

Support common:

```text
HTTP_PROXY
HTTPS_PROXY
NO_PROXY
```

and authenticated corporate proxies where the runtime permits.

Doctor must show:

```text
proxy detected
```

without exposing credentials.

---

# 224. P1 correction — DNS rebinding defense

For custom endpoints:

```text
resolve
connect
verify destination policy
```

Do not resolve once and assume the IP remains safe.

---

# 225. P1 correction — SSRF redirect defense

Re-check destination after redirects.

Do not let:

```text
public.example
```

redirect to:

```text
169.254.169.254
```

or local admin endpoints.

---

# 226. P1 correction — provider metadata refresh should be atomic

```text
download
validate
write temp
fsync where appropriate
rename
```

Never leave a half-written catalog.

---

# 227. P1 correction — logs need rotation

Define:

```text
max file size
max files
max total disk
retention
```

When full:

```text
drop oldest
```

never:

```text
crash CCX
```

---

# 228. P1 correction — log levels must be per subsystem

Allow:

```text
gateway=debug
provider=openai=trace
tui=info
```

without turning everything into TRACE.

---

# 229. P1 correction — redaction must be structural

Redact:

```text
known secret fields
known secret headers
API key patterns
OAuth token patterns
JWT-like values
```

But do not rely solely on regex.

The strongest protection is:

```text
never serialize credential objects into logs
```

---

# 230. P1 correction — debug mode must have a safety banner

```text
DEBUG LOGGING ENABLED

Request metadata will be logged.
Prompt/response bodies remain disabled by default.
```

If body capture is enabled:

```text
SENSITIVE MODE
```

must be explicit.

---

# 231. P1 correction — crash reports

Crash report should contain:

```text
stack
version
platform
sanitized state
```

not:

```text
prompt
response
credential
tool arguments
environment dump
```

---

# 232. P1 correction — model aliases need collision handling

If two providers expose:

```text
sonnet
```

CCX must never resolve ambiguously.

Use:

```text
provider/model
```

internally.

Friendly aliases are presentation-only.

---

# 233. P1 correction — provider deletion safety

Cannot remove a provider if:

```text
active session
```

uses it.

Options:

```text
block
or mark pending deletion
```

Never invalidate an active session silently.

---

# 234. P1 correction — credential deletion safety

Similarly:

```text
auth logout
```

should warn:

```text
2 profiles depend on this credential.
```

Active sessions may continue using an in-memory credential lease until termination.

---

# 235. P1 correction — profile inheritance

Profiles should support:

```text
extends: base
```

only if merge semantics are formally defined.

Otherwise avoid inheritance in V1.

Simple is safer.

---

# 236. P1 correction — configuration explainability

Add:

```bash
ccx config explain <key>
```

Example:

```text
model = qwen...
source = project
overridden by = CLI? no
```

This will save enormous debugging time.

---

# 237. P1 correction — route resolution must be deterministic

Given identical:

```text
config
environment
provider catalog
```

resolution must produce identical:

```text
route
```

unless a policy explicitly permits dynamic selection.

---

# 238. P1 correction — dynamic routing needs reproducibility

For auto-routing, record:

```text
candidate set
scores
decision
reason
timestamp
```

so a user can reproduce why CCX chose a model.

---

# 239. P1 correction — no “AI magic” in compatibility core

Do not use an LLM to decide:

```text
whether a protocol field should be dropped
```

Compatibility must be deterministic.

AI-assisted routing belongs in later research scope.

---

# 240. P1 correction — benchmark correctness before performance

For translation:

```text
semantic equivalence
```

beats:

```text
microsecond optimization
```

No optimization should bypass correctness fixtures.

---

# 241. P1 correction — native fast path needs equivalence tests

Native passthrough should not bypass all testing.

Test:

```text
native request
header policy
stream handling
errors
cancellation
```

The source plan's native-provider fast path is sound, but the fast path still needs full contract coverage. fileciteturn2file2L93-L144

---

# 242. P1 correction — OpenAI-compatible does not mean one behavior

Provider metadata should record:

```text
protocol=openai-chat
behavior_profile=vllm-qwen
```

This preserves the source plan's important vLLM parser-awareness requirement. fileciteturn2file7L431-L441

---

# 243. P1 correction — provider behavior profiles

Create:

```text
behavior profiles
```

for:

```text
tool streaming
reasoning
parallel tools
JSON mode
vision
context
```

Provider/model adapters select the profile.

---

# 244. P1 correction — compatibility profiles need precedence

```text
exact model
>
model family
>
provider behavior
>
protocol defaults
```

This makes quirks maintainable.

---

# 245. P1 correction — generated test matrix

From provider metadata, generate:

```text
provider × protocol × feature
```

and automatically create test cases.

Avoid manually forgetting a provider when a new feature is added.

---

# 246. P1 correction — regression corpus should distinguish bug classes

Tags:

```text
headers
discovery
tools
streaming
reasoning
context
auth
network
fallback
security
TUI
```

Then CI can run targeted suites.

---

# 247. P1 correction — live-provider tests must be opt-in

Never put paid API tests in ordinary pull-request CI.

Use:

```text
CCX_LIVE_TESTS=1
```

and provider credentials from CI secrets.

---

# 248. P1 correction — test credentials need isolation

Use dedicated:

```text
test provider accounts
```

where possible.

Never run compatibility tests against a user's personal account in automation.

---

# 249. P1 correction — golden fixtures must be legally safe

Do not commit:

```text
private prompts
provider secrets
copyright-sensitive huge payloads
personal data
```

Use synthetic fixtures.

---

# 250. P1 correction — fuzz corpus persistence

When fuzzing finds a crash:

```text
promote minimized input
→ regression fixture
```

This closes the loop automatically.

---

# 251. P1 correction — release gates

Release MUST fail if:

```text
hard invariant regression
security regression
known Claude Code breakage
tool semantic corruption
stream parser regression
secret leakage
```

Documentation typo alone should not block release.

---

# 252. P1 correction — semantic versioning

Define:

```text
MAJOR
MINOR
PATCH
```

and compatibility implications.

A new provider should not require a major version.

Changing config schema incompatibly should.

---

# 253. P1 correction — backward compatibility

V1 config should remain readable by future CCX versions through migration.

Future CCX must not silently reinterpret old route semantics.

---

# 254. P1 correction — repository layout

Recommended:

```text
cmd/ccx/
internal/
  cli/
  tui/
  supervisor/
  runtime/
  gateway/
  session/
  compat/
  ir/
  protocol/
  provider/
  capability/
  transform/
  resilience/
  auth/
  config/
  storage/
  diagnostics/
  telemetry/
  security/
  update/
pkg/
  api/
  provider-sdk/
tests/
  contract/
  fixtures/
  integration/
  chaos/
  fuzz/
  e2e/
compat/
providers/
docs/
```

The exact language is still an implementation choice; the module boundaries are not.

---

# 255. P1 correction — dependency direction

Enforce:

```text
CLI/TUI
  ↓
application
  ↓
domain
  ↓
ports/interfaces
  ↓
adapters
```

Never:

```text
provider adapter → TUI
```

or:

```text
TUI → HTTP implementation
```

---

# 256. P1 correction — no global singleton state

Avoid global:

```text
current provider
current model
current session
```

Use explicit runtime/session contexts.

This is essential for concurrent sessions and tests.

---

# 257. P1 correction — cancellation context propagation

Every layer gets a cancellation context:

```text
CLI
 ↓
session
 ↓
request
 ↓
provider
 ↓
HTTP
```

Cancellation must be testable at each boundary.

---

# 258. P1 correction — context propagation should carry deadlines

Use:

```text
deadline
request_id
session_id
trace_id
```

but never put secrets into context metadata.

---

# 259. P1 correction — provider adapters need lifecycle hooks

Interface should support:

```text
Init
Validate
Discover
Probe
Send
Stream
Cancel
Close
```

Not every provider needs every method; capabilities declare support.

---

# 260. P1 correction — connection pooling

Provider clients should reuse connections where safe.

But pool keys must include:

```text
scheme
host
port
proxy
TLS config
auth identity where necessary
```

Never share a credential-bearing transport across incompatible identities.

---

# 261. P1 correction — per-provider concurrency control

Provider configuration should support:

```text
max_concurrency
queue_timeout
request_queue_limit
```

Queue behavior:

```text
reject
wait
```

must be explicit.

---

# 262. P1 correction — queue fairness

If multiple sessions share one provider:

```text
per-session fairness
```

should prevent one agent from monopolizing the entire provider connection pool.

---

# 263. P1 correction — starvation detection

If a request waits too long:

```text
QUEUE_TIMEOUT
```

with actionable diagnostics.

---

# 264. P1 correction — retry budget must be global enough to prevent storms

Per request:

```text
max attempts
max elapsed
```

Per provider:

```text
retry rate budget
```

Otherwise 100 sessions can each retry 3 times and create a 300-request storm.

---

# 265. P1 correction — retry jitter must be deterministic in tests

Production:

```text
random jitter
```

Tests:

```text
seeded jitter
```

so failure behavior is reproducible.

---

# 266. P1 correction — fallback storms

If primary fails:

```text
100 sessions → 100 fallback switches
```

could overload the fallback.

Use provider-level admission control.

---

# 267. P1 correction — overload protection

When all candidates are unhealthy:

```text
fail fast
```

rather than endlessly cycling providers.

---

# 268. P1 correction — health probes should be cheap

Provider health should not generate full model generations.

Use:

```text
TCP/TLS
HTTP endpoint
models
```

and only feature probes when requested.

---

# 269. P1 correction — health != capability

A provider can be:

```text
healthy
```

while a model is:

```text
tool-broken
```

Keep those states separate.

---

# 270. P1 correction — model availability can be regional

For cloud providers:

```text
provider
region
model
account
```

all matter.

---

# 271. P1 correction — credential-specific capabilities

Two credentials for the same provider/model can differ.

Capability key:

```text
provider + model + credential_scope + feature
```

when necessary.

---

# 272. P1 correction — account gating

The source plan correctly identifies caching as account/model gated. fileciteturn2file1L61-L65

Generalize this to:

```text
feature gating
```

not just caching.

Possible:

```text
tools
reasoning
context
cache
vision
server tools
```

---

# 273. P1 correction — capability expiration

Some capabilities should expire quickly:

```text
availability
quota
health
```

Others can live longer:

```text
protocol support
```

Each capability class gets a TTL policy.

---

# 274. P1 correction — “observed” must never mean “guaranteed forever”

Observed:

```text
worked at time T
```

not:

```text
always works
```

Doctor should show verification age.

---

# 275. P1 correction — compatibility UI

Model card should show:

```text
Verified 3m ago
```

or:

```text
Last verified 4d ago
```

This is more honest than a permanent green check.

---

# 276. P1 correction — provider setup should be resumable

If user exits halfway:

```text
ccx
```

next launch:

```text
Resume setup?
```

or cleanly start over.

---

# 277. P1 correction — onboarding should not force account creation

This is a core local-first property.

---

# 278. P1 correction — setup should support “skip verification”

For exotic custom endpoints:

```text
manual configuration
```

but show:

```text
UNVERIFIED
```

and never call it certified.

---

# 279. P1 correction — setup must have backtracking

At every step:

```text
Back
Cancel
```

must work without corrupting partial config.

---

# 280. P1 correction — provider wizard should infer protocol carefully

Auto-detection:

```text
probe known endpoints
inspect response
```

If ambiguous:

```text
ask user
```

Never guess based solely on domain name.

---

# 281. P1 correction — custom endpoint security warning

For:

```text
http://
```

show:

```text
Unencrypted connection.
Credentials and prompts may be exposed.
```

Require explicit confirmation for non-loopback HTTP.

---

# 282. P1 correction — local HTTP is acceptable

For:

```text
127.0.0.1
```

HTTP is fine because TLS is unnecessary overhead in the default local process boundary, but authentication should still exist.

---

# 283. P1 correction — endpoint trust model

Represent:

```text
local trusted
private trusted
public TLS
public HTTP
unknown
```

and use this in auth policy.

---

# 284. P1 correction — model picker must expose exact IDs

Friendly name:

```text
Qwen Coder
```

but detail view:

```text
provider: openrouter
model: qwen/...
```

Copy command:

```text
ccx use openrouter/qwen/...
```

---

# 285. P1 correction — favorites and recents

Add:

```text
favorites
recently used
verified
```

to TUI.

This dramatically improves daily UX.

---

# 286. P1 correction — model filtering

Filters:

```text
/tools
/reasoning
/vision
/local
/cheap
/fast
/verified
```

Search should support fuzzy matching.

---

# 287. P1 correction — “recommended” must be explainable

If CCX says:

```text
Recommended
```

show why:

```text
Matches coding profile:
tools ✓
reasoning ✓
context ✓
verified recently ✓
```

No opaque AI ranking in V1.

---

# 288. P1 correction — profiles should be boring

Avoid dozens of smart settings.

V1 profile:

```text
requirements
candidate routes
fallback policy
probe policy
```

That's enough.

---

# 289. P1 correction — policy inheritance should be avoided initially

No complex policy DSL in V1.

YAML should remain understandable.

---

# 290. P1 correction — environment import

Support:

```bash
ccx provider import-env
```

to detect common environment variables.

But:

```text
never upload
never expose
```

and show exactly what was detected.

---

# 291. P1 correction — migration from existing gateways

Add importers for common config styles:

```text
OpenRouter
LiteLLM
Claude Code Router
custom Anthropic gateway
```

But import should create:

```text
CCX config
```

not attempt to take over the old gateway.

---

# 292. P1 correction — coexistence

Users may have:

```text
LiteLLM
CCR
CCX
```

running simultaneously.

Doctor should detect likely port conflicts and route loops.

---

# 293. P1 correction — loop detection

Multi-hop routes can accidentally create:

```text
CCX → CCR → CCX → CCR
```

Add hop marker:

```text
X-CCX-Hop
```

or equivalent local metadata, and stop loops.

Do not trust external copies of the header blindly.

---

# 294. P1 correction — route depth limit

Default:

```text
max_hops = 8
```

Reject deeper chains.

---

# 295. P1 correction — provider endpoint loopback detection

Detect if custom endpoint points to:

```text
CCX itself
```

and reject unless explicitly intended.

---

# 296. P1 correction — no transparent credential forwarding

In multi-hop mode:

```text
Claude credential
```

must not automatically become:

```text
upstream credential
```

unless policy explicitly maps it.

---

# 297. P1 correction — response size limits must be streaming-safe

Don't reject legitimate long responses simply because they are long.

Use:

```text
event size limit
aggregate safety limit
```

with streaming output.

---

# 298. P1 correction — large context should avoid duplicate buffers

For very large requests:

```text
read → parse → transform → send
```

should minimize copies.

Where streaming request bodies are possible, use them.

---

# 299. P1 correction — JSON parser safety

Use bounded parsers.

Reject:

```text
extreme nesting
```

and enforce:

```text
max string
max object fields
max array length
```

---

# 300. P1 correction — Unicode correctness

Test:

```text
emoji
CJK
combining characters
RTL
invalid UTF-8
surrogate edge cases
```

especially in streamed tool arguments.

---

# 301. P1 correction — tool JSON correctness

Test:

```text
escaped quotes
unicode
large arrays
nested objects
partial chunks
strings containing braces
```

Never parse tool arguments by concatenating strings with naive delimiters.

---

# 302. P1 correction — fine-grained tool streaming

Anthropic currently supports fine-grained tool-input streaming where fragments can be partial or invalid JSON until accumulated. citeturn0search0

Therefore CCX must distinguish:

```text
streamed partial tool input
```

from:

```text
invalid final tool input
```

Do not reject every partial fragment.

---

# 303. P1 correction — tool argument buffering policy

For fine-grained tools:

```text
per-call bounded accumulator
```

with:

```text
max_tool_input_bytes
```

If exceeded:

```text
cancel/terminate safely
```

---

# 304. P1 correction — tool call ordering

Preserve:

```text
tool call order
content block order
message order
```

unless the target explicitly supports equivalent parallel semantics.

---

# 305. P1 correction — parallel tool degradation must be explicit

If target cannot parallelize:

```text
parallel → serialized
```

must be marked:

```text
DEGRADED
```

not:

```text
SUPPORTED
```

---

# 306. P1 correction — tool choice semantics

Map:

```text
auto
any
none
specific tool
disable_parallel_tool_use
```

where the target supports them.

If target cannot represent a specific choice:

```text
reject or explicit degrade
```

---

# 307. P1 correction — stop sequences

Normalize:

```text
stop
stop_sequences
finish_reason
```

without inventing semantics.

---

# 308. P1 correction — max tokens semantics

Providers may interpret output limits differently.

Canonical IR needs:

```text
max_output_tokens
```

and adapter mapping.

Never assume `max_tokens` means exactly the same thing across protocols.

---

# 309. P1 correction — temperature/top-p/etc.

Unsupported sampling parameters must follow policy:

```text
drop
warn
reject
```

Default:

```text
drop only when semantically harmless
```

Otherwise warn/reject.

---

# 310. P1 correction — structured output

Add future-proof IR for:

```text
JSON schema
response format
structured output
```

Do not force it into tool calling.

---

# 311. P1 correction — stop reasons need canonical enum + raw value

```text
canonical:
  end_turn
  tool_use
  max_tokens
  stop_sequence
  error
  unknown

raw:
  provider-specific
```

---

# 312. P1 correction — usage fields need extensibility

Canonical:

```text
input
output
cache_read
cache_write
reasoning
```

plus:

```text
raw_usage
```

for provider-specific fields.

---

# 313. P1 correction — usage reconciliation

At completion:

```text
stream usage
final usage
```

may differ.

Define authoritative source:

```text
final response usage
```

unless provider documents otherwise.

---

# 314. P1 correction — request body preservation

For unknown fields:

```text
raw extensions map
```

preserve them through compatible paths.

But do not blindly forward unknown fields to incompatible providers.

---

# 315. P1 correction — beta header registry

Instead of code:

```text
if header == X
```

use:

```yaml
beta:
  id:
  first_seen:
  last_seen:
  semantic_area:
  native_safe:
  openai_action:
  gemini_action:
  probe:
```

This directly operationalizes the source plan's beta-header drift strategy. fileciteturn2file6L398-L413

---

# 316. P1 correction — unknown beta header behavior

For a newly observed beta header:

```text
native Anthropic target → preserve
translated target → strip by default
unknown semantic importance → log compatibility warning
```

If the header appears required for a feature:

```text
capability becomes unverified
```

Do not pretend stripping is always safe.

---

# 317. P1 correction — compatibility policy should be data-driven

Example:

```yaml
rules:
  - when:
      claude_code: ">=2.1.240"
      target_protocol: openai-chat
      feature: beta_header_unknown
    action: strip
    severity: warning
```

---

# 318. P1 correction — Claude Code release watcher must detect more than headers

Watch:

```text
environment variables
gateway endpoints
model discovery
request fields
response expectations
stream events
tool schemas
error behavior
auth
```

---

# 319. P1 correction — official docs vs observed behavior

Every compatibility fact should have:

```text
source_type:
  official
  issue
  observed
  community
```

Priority:

```text
official contract
>
reproducible observation
>
issue report
>
community anecdote
```

When they conflict:

```text
mark discrepancy
```

rather than silently choosing one.

---

# 320. P1 correction — current facts need timestamps

Every volatile fact:

```text
verified_at
```

The spec itself should not pretend to be timeless.

---

# 321. P1 correction — no hardcoded model catalogs in binary

Ship a minimal catalog if useful, but make it:

```text
versioned
refreshable
optional
```

Never make provider availability depend on binary release cadence.

---

# 322. P1 correction — provider catalog update must work offline

Use cached catalog:

```text
stale
```

and clearly label it.

---

# 323. P1 correction — local model discovery should be zero-cost

For local runtimes:

```text
list models
```

should be preferred over generation probes.

---

# 324. P1 correction — local runtime detection

Detect:

```text
Ollama
LM Studio
vLLM
llama.cpp
SGLang
```

through:

```text
process
known port
health endpoint
```

but do not automatically start them.

---

# 325. P1 correction — don't manage third-party servers in V1

CCX can detect local runtimes.

It should NOT automatically:

```text
install model
download multi-GB weights
start GPU server
```

That belongs outside the compatibility gateway.

---

# 326. P1 correction — GPU/resource awareness

For local providers, display:

```text
runtime reachable
model loaded?
VRAM not reliably knowable
```

Do not claim hardware support from HTTP availability alone.

---

# 327. P1 correction — provider health UI should distinguish “configured” from “reachable”

States:

```text
NOT_CONFIGURED
CONFIGURED
REACHABLE
AUTHENTICATED
MODEL_AVAILABLE
VERIFIED
```

---

# 328. P1 correction — onboarding should end with a reproducible command

After setup:

```text
Your route is saved.

Run:
  ccx claude
```

Optionally:

```text
Install shell integration? [y/N]
```

---

# 329. P1 correction — “normal `claude`” is an opt-in convenience

The product should never force a shell hijack.

Recommended:

```text
ccx claude
```

as the guaranteed path.

Optional:

```text
ccx install-shell
```

for users who want bare `claude`.

---

# 330. P1 correction — bare `claude` integration must be detectable

Doctor:

```text
Bare `claude` integration:
  NOT INSTALLED
```

and:

```text
Install now? [y/N]
```

---

# 331. P1 correction — Claude Code updates should trigger compatibility warning

If Claude Code changes:

```text
2.1.245 → 2.1.246
```

and CCX has not verified it:

```text
⚠ New Claude Code version detected.
Compatibility status: UNVERIFIED.
```

User can continue.

---

# 332. P1 correction — known-broken should be loud

```text
✗ Claude Code 2.1.xxx is known incompatible with CCX 1.x.

Reason:
  gateway streaming contract changed.

Upgrade CCX or use a verified Claude Code version.
```

---

# 333. P1 correction — safe rollback of Claude Code compatibility

Do not downgrade Claude Code automatically.

Offer:

```text
show compatible versions
```

and leave installation control to user.

---

# 334. P1 correction — support a compatibility lock

Optional project setting:

```yaml
compatibility:
  claude_code:
    max_verified: "2.1.245"
```

But it should warn, not forcibly block, unless the project explicitly selects strict mode.

---

# 335. P1 correction — strict mode

Add:

```text
compatibility_mode:
  strict
  balanced
  permissive
```

Strict:

```text
unknown capability → reject
known-broken → reject
lossy transform → reject
```

Balanced:

```text
safe degradation
```

Permissive:

```text
best effort
```

Default:

```text
balanced
```

---

# 336. P1 correction — transform audit

Every request should internally produce:

```text
compatibility decision record
```

Example:

```text
removed:
  anthropic-beta:X

converted:
  tool_choice:any → required

degraded:
  parallel_tools → serial

preserved:
  images
```

This is invaluable for `doctor --explain`.

---

# 337. P1 correction — add `ccx explain`

```bash
ccx explain last
ccx explain request <id>
ccx explain route
```

Useful without exposing request bodies.

---

# 338. P1 correction — provider debug should be reproducible

Add:

```bash
ccx test provider <name> --case streaming-tools
```

Cases:

```text
basic
stream
tools
parallel-tools
thinking
vision
cache
cancel
errors
```

---

# 339. P1 correction — provider certification CLI

```bash
ccx provider certify <name>
```

should run the relevant test suite and produce:

```text
certificate.json
```

---

# 340. P1 correction — certification must never imply provider endorsement

It means:

```text
CCX compatibility verified under listed conditions
```

not:

```text
CCX recommends this provider commercially
```

---

# 341. P1 correction — community provider submissions

Submission must include:

```text
provider.yaml
fixtures
tests
auth docs
capability declaration
```

CI validates.

---

# 342. P1 correction — provider plugins should not run in-process if avoidable

If custom executable plugins arrive in V2:

```text
separate process
RPC
permission model
timeout
```

rather than arbitrary in-process code.

---

# 343. P1 correction — license decision is actually blocking

The existing file says:

```text
license to be selected
```

That is not final.

Before public repository release choose a license deliberately.

For a compatibility infrastructure project, permissive licensing may maximize adoption, but this is a project-owner decision, not an implementation detail.

---

# 344. P1 correction — trademark/branding

CCX documentation must avoid implying endorsement by Anthropic.

Use wording:

```text
CCX is an independent open-source project.
Claude Code and Anthropic are trademarks of their respective owner.
```

---

# 345. P1 correction — provider terms

Each provider integration should link to its official docs/terms in documentation, not embed legal claims into code.

---

# 346. P1 correction — no provider key brokerage

CCX should not become:

```text
central credential marketplace
```

Everything remains user's local credentials.

---

# 347. P1 correction — no hosted telemetry

The project should have a clear policy:

```text
telemetry default = OFF
```

If optional telemetry is ever added:

```text
explicit opt-in
inspectable payload
disable command
```

---

# 348. P1 correction — “offline” should be testable

Add a test mode that rejects all unexpected outbound connections.

For example:

```text
network policy = deny-by-default
```

and integration tests verify no hidden calls.

---

# 349. P1 correction — network destination audit

`ccx doctor --network` should show:

```text
CCX listener
provider targets
proxy
DNS
TLS
recent destinations
```

without payloads.

---

# 350. P1 correction — local network exposure scanner

Doctor should detect:

```text
gateway bound to 0.0.0.0
IPv6 wildcard
```

and flag it loudly.

---

# 351. P1 correction — IPv6 loopback

Do not assume:

```text
127.0.0.1
```

is enough.

If IPv6 is enabled:

```text
::1
```

must be handled explicitly.

---

# 352. P1 correction — dual-stack binding

Never accidentally bind:

```text
[::]:port
```

when the user intended loopback-only.

Verify actual socket exposure, not only configured address.

---

# 353. P1 correction — DNS names for local gateway

Claude Code should preferably use:

```text
127.0.0.1
```

or:

```text
::1
```

not:

```text
localhost
```

when deterministic address family matters.

---

# 354. P1 correction — browser/OAuth flows

OAuth may require a local callback.

Define:

```text
callback bind
random state
PKCE
timeout
CSRF protection
callback cleanup
```

Do not leave an OAuth callback server permanently listening.

---

# 355. P1 correction — OAuth tokens

Store refresh/access tokens in OS keychain where supported.

Access tokens should remain in memory as much as possible.

---

# 356. P1 correction — OAuth provider identity

Credential records:

```text
provider
account
scopes
expires_at
refreshable
```

No token values in config.

---

# 357. P1 correction — clock skew for OAuth

Use provider expiry plus a safety window.

If clock is obviously wrong:

```text
doctor
```

should explain.

---

# 358. P1 correction — provider auth refresh during stream

Do not attempt to replace credentials halfway through an already-open HTTP stream unless the provider protocol supports it.

Instead:

```text
finish/cancel
refresh
new request under safe policy
```

and do not replay tool-producing requests automatically.

---

# 359. P1 correction — API key rotation

Support multiple credentials per provider:

```text
primary
backup
```

but key rotation is not the same as fallback.

Rotation may occur:

```text
expired/invalid credential
```

without changing model/provider.

---

# 360. P1 correction — credential fallback must respect billing boundaries

If two credentials belong to different organizations/accounts:

```text
do not automatically switch
```

unless the user explicitly configured that relationship.

---

# 361. P1 correction — provider aliases

Separate:

```text
provider ID
display name
endpoint
```

so changing branding doesn't break config.

---

# 362. P1 correction — import/export must be redaction-aware

`ccx export` should default to:

```text
safe export
```

and require explicit:

```text
--include-secrets
```

which should be strongly discouraged and clearly labeled.

Better:

```text
credentials never exported
```

in V1.

---

# 363. P1 correction — backup encryption

If backup contains anything sensitive:

```text
encrypted
```

using a local key/password.

But safest V1 backup:

```text
no credentials
```

---

# 364. P1 correction — provider URL secrecy

Some endpoints contain tokens in query strings.

Doctor/logging must redact:

```text
query
fragment
userinfo
```

when URLs are logged.

---

# 365. P1 correction — request header privacy

Some providers put account IDs or tenant IDs in headers.

Classify:

```text
public metadata
sensitive metadata
secret
```

and redact appropriately.

---

# 366. P1 correction — content privacy in error messages

Never echo entire user input into an error.

Use:

```text
field path
length
hash
```

where useful.

---

# 367. P1 correction — file/document privacy

CCX should not persist uploaded document/image bytes merely to translate them.

Streaming/in-memory where possible.

---

# 368. P1 correction — temp files

If temp files are unavoidable:

```text
user-only permissions
random names
atomic cleanup
no secret in filename
```

---

# 369. P1 correction — disk exhaustion

Doctor should check:

```text
free disk
log usage
cache usage
```

and fail safely before disk exhaustion.

---

# 370. P1 correction — update disk space

Before updating:

```text
enough free space
```

or fail before modifying the active binary.

---

# 371. P1 correction — atomic binary replacement

Use platform-safe strategy:

```text
new binary staged
verify
replace
```

Windows may require helper process if executable is locked.

---

# 372. P1 correction — update rollback retention

Keep:

```text
current
previous
```

at minimum.

Clean older versions according to policy.

---

# 373. P1 correction — interrupted update recovery

On startup detect:

```text
update transaction incomplete
```

then:

```text
rollback/repair
```

rather than starting half-updated code.

---

# 374. P1 correction — self-update trust root

The update verifier must have a pinned trust root/key mechanism.

Do not download a “new public key” from the same untrusted channel and trust it automatically.

---

# 375. P1 correction — repository release provenance

Document:

```text
build environment
version
commit
artifact hash
SBOM
signature
```

---

# 376. P1 correction — build reproducibility

Aim for reproducible builds where practical.

At minimum:

```text
same source + pinned dependencies
```

should produce auditable artifacts.

---

# 377. P1 correction — dependency policy

Avoid dependencies that:

```text
phone home
download code at runtime
```

without explicit reason.

---

# 378. P1 correction — no runtime package manager

CCX should ship as a self-contained executable where practical.

Do not make first launch depend on:

```text
npm install
pip install
go install
```

---

# 379. P1 correction — architecture should be language-neutral

The spec should not force Go merely because the sample repository layout uses `internal/`.

The implementation can choose:

```text
Go
Rust
TypeScript + native runtime
```

but must satisfy the architecture and runtime constraints.

For a terminal/system utility, Go or Rust are strong candidates.

---

# 380. P1 correction — TUI framework choice is implementation detail

Do not bake a specific TUI library into the product contract.

---

# 381. P1 correction — API compatibility fixtures need schema snapshots

Store:

```text
request headers
request body
response headers
response body
SSE event sequence
expected semantic IR
```

with sensitive fields removed.

---

# 382. P1 correction — snapshot tests need normalization

Normalize:

```text
timestamps
random IDs
request IDs
provider latency
```

before snapshot comparison.

---

# 383. P1 correction — differential testing

Where possible:

```text
CCX → provider
```

and:

```text
reference client → provider
```

compare:

```text
accepted request
capabilities
usage
stream semantics
```

This catches translation bugs.

---

# 384. P1 correction — metamorphic tests

Examples:

```text
stream=true vs stream=false
```

should produce semantically equivalent final content.

Likewise:

```text
tool call split into N stream chunks
```

must equal:

```text
same tool call in one chunk
```

where protocol semantics allow.

---

# 385. P1 correction — property tests

Properties:

```text
serialize(parse(x)) ≈ x
```

where exact equality is expected.

And:

```text
parse(serialize(IR)) = IR
```

for supported fields.

---

# 386. P1 correction — protocol downgrade tests

For unsupported feature:

```text
strict → reject
balanced → documented degrade
permissive → best effort
```

must be deterministic.

---

# 387. P1 correction — test cancellation race conditions

Test cancellation at:

```text
before connect
after connect
before headers
after headers
during stream
during tool args
after tool call
during retry
during fallback
```

This is one of the most important reliability suites.

---

# 388. P1 correction — test shutdown race conditions

Test:

```text
SIGINT during startup
SIGINT during provider connection
SIGINT during stream
SIGTERM with 20 sessions
shutdown during update
```

---

# 389. P1 correction — test config race conditions

Concurrent:

```text
ccx config set
ccx provider test
ccx claude
```

must not corrupt config or session snapshots.

---

# 390. P1 correction — test provider catalog races

Concurrent:

```text
refresh
use model
remove provider
```

must resolve consistently.

---

# 391. P1 correction — stale model selection

If a selected model disappears between discovery and request:

```text
MODEL_UNAVAILABLE
```

then fallback policy applies.

Do not silently replace the model.

---

# 392. P1 correction — provider 404 ambiguity

404 may mean:

```text
wrong path
unknown model
unsupported endpoint
```

Doctor should inspect context before giving generic advice.

---

# 393. P1 correction — authentication ambiguity

401/403 can mean:

```text
expired credential
wrong credential type
missing scope
model permission
region permission
```

Expose the provider's safe error detail.

---

# 394. P1 correction — billing errors

Do not retry:

```text
402
spend cap
```

unless provider explicitly says transient.

---

# 395. P1 correction — 429 semantics

Use:

```text
Retry-After
```

when present, but recognize that some 429s represent spend caps and should not be retried indefinitely. Anthropic's current documentation explicitly distinguishes rate-limit behavior from spend-limit scenarios. citeturn1search1turn1search5

---

# 396. P1 correction — overloaded vs unavailable

Distinguish:

```text
429 rate limited
529 overloaded
503 unavailable
504 timeout
```

because remediation differs.

---

# 397. P1 correction — error translation

If target emits:

```text
provider-specific error
```

CCX should map:

```text
canonical error type
```

while retaining:

```text
raw provider code
```

for diagnosis.

---

# 398. P1 correction — don't manufacture Anthropic semantics

When a provider doesn't have a real equivalent:

```text
raw extension / degraded state
```

is better than fake Anthropic behavior.

---

# 399. P1 correction — stream terminal error semantics

If upstream fails after downstream has received content:

```text
cannot send normal HTTP error
```

so CCX must use a valid stream-level error/termination strategy supported by Claude Code.

This needs golden fixtures.

---

# 400. P1 correction — connection-close behavior

Test what Claude Code actually does when:

```text
SSE connection closes without message_stop
```

and build compatibility behavior around observed client semantics.

Do not invent behavior solely from the API docs.

---

# 401. P1 correction — request body timeout

Large uploads should not be killed merely because the request body takes time to arrive.

Use:

```text
header timeout
body progress timeout
total timeout
```

separately.

---

# 402. P1 correction — slow provider

A slow provider should produce:

```text
status = waiting/streaming
```

not:

```text
dead
```

if bytes/events are still arriving.

---

# 403. P1 correction — keepalive generation must be protocol-aware

The source plan correctly identifies keepalives as a distinct failure class. fileciteturn2file6L410-L413

But CCX should not blindly invent arbitrary events.

Define:

```text
wire-native keepalive
translated keepalive
transport-level heartbeat
```

and verify which event Claude Code accepts for each compatibility version.

---

# 404. P1 correction — keepalive interval must be configurable

Example:

```yaml
stream:
  keepalive:
    enabled: true
    interval: 15s
```

Default should be safely below common intermediary idle timeouts without creating excessive traffic.

---

# 405. P1 correction — upstream keepalive pass-through

If upstream sends valid keepalive events:

```text
preserve
```

unless target protocol requires translation.

---

# 406. P1 correction — stream clock

Use monotonic timing for:

```text
idle timeout
keepalive
first event latency
```

---

# 407. P1 correction — content block lifecycle validator

Every stream must pass:

```text
event sequence validator
```

before serialization when possible.

This catches:

```text
block_stop without block_start
message_stop with open blocks
duplicate message_start
```

---

# 408. P1 correction — “repair” must be conservative

Only auto-repair:

```text
missing terminal close
```

when the intended close is unambiguous.

Never auto-repair:

```text
tool ID
tool arguments
thinking signatures
```

---

# 409. P1 correction — tool ID collision prevention

Namespace IDs by session/request.

Never allow two concurrent sessions to share generated IDs.

---

# 410. P1 correction — request fingerprinting

Fingerprint:

```text
normalized request
route
session
```

but never use the fingerprint alone as proof that a request is safe to replay.

---

# 411. P1 correction — idempotency support

If provider supports idempotency keys:

```text
use them where semantically correct
```

But don't assume all providers implement them.

---

# 412. P1 correction — request replay policy

Every request class:

```text
read-only generation
tool-producing generation
```

has a replay policy.

Default:

```text
tool-producing = no automatic replay after upstream acceptance
```

---

# 413. P1 correction — provider health circuit breaker must not block explicit user override

If user explicitly says:

```text
ccx test provider X
```

the circuit breaker may be bypassed for the diagnostic operation.

Normal production traffic still respects the breaker.

---

# 414. P1 correction — doctor should not mutate state by default

`doctor`:

```text
read-only
```

`doctor --fix`:

```text
explicit mutations
```

with a preview.

---

# 415. P1 correction — fix plans should be transactional

```text
doctor --fix
 ↓
show plan
 ↓
confirm
 ↓
snapshot
 ↓
apply
 ↓
verify
```

Noninteractive:

```bash
ccx doctor --fix --yes
```

only for safe classified fixes.

---

# 416. P1 correction — fix safety classes

```text
SAFE
CAUTION
DESTRUCTIVE
```

`--yes` only applies SAFE automatically.

---

# 417. P1 correction — repair logs

Every mutation records:

```text
what changed
before hash
after hash
reason
timestamp
```

No secrets.

---

# 418. P1 correction — compatibility override escape hatch

Add:

```bash
ccx compat list
ccx compat explain <rule>
ccx compat disable <rule>
```

but disabling safety rules should be explicit and visible.

---

# 419. P1 correction — user overrides should be scoped

Override can target:

```text
provider
model
protocol
Claude Code version
```

Never globally disable a safety rule when a narrow override works.

---

# 420. P1 correction — configuration provenance

Every effective config value should have:

```text
source
```

This is essential for debugging.

---

# 421. P1 correction — project vs user policy precedence

Do not let an untrusted project override user security settings.

Security precedence:

```text
hardcoded safety
>
user security policy
>
project policy
>
provider defaults
```

Route preference can be project-scoped.

Security controls cannot be weakened by project config.

---

# 422. P1 correction — command-line security flags

A dangerous flag should be:

```text
explicit
```

and ideally:

```text
non-persistent
```

unless saved intentionally.

---

# 423. P1 correction — environment variables must not bypass security policy silently

If:

```text
CCX_ALLOW_REMOTE=1
```

changes a dangerous setting, doctor must show it.

Better: dangerous network exposure requires a CLI/config policy rather than a hidden environment override.

---

# 424. P1 correction — shell integration should only set CCX variables

It must not redefine unrelated user variables.

---

# 425. P1 correction — Claude Code environment cleanup

When CCX launches Claude Code:

```text
set only required overrides
```

Do not overwrite unrelated Anthropic/provider environment variables unless necessary.

---

# 426. P1 correction — nested CCX invocation

If user runs:

```bash
ccx claude
```

inside an existing CCX session:

```text
detect nested runtime
```

and avoid infinite nesting.

---

# 427. P1 correction — CI recursion

CI might already set:

```text
ANTHROPIC_BASE_URL
```

CCX should show:

```text
existing gateway detected
```

and ask whether to wrap it.

---

# 428. P1 correction — explicit gateway chaining

Support:

```bash
ccx run --upstream http://127.0.0.1:xxxx
```

but detect loops.

---

# 429. P1 correction — gateway health endpoint

Minimum:

```text
GET /healthz
```

Response should expose:

```text
status
version
uptime
```

not secrets.

---

# 430. P1 correction — readiness endpoint

Separate:

```text
/healthz
/readyz
```

where:

```text
health = process works
ready = gateway can accept requests
```

Upstream provider does not need to be healthy for CCX itself to be ready.

---

# 431. P1 correction — metrics endpoint should remain local

If exposing metrics:

```text
127.0.0.1 only
```

by default.

---

# 432. P1 correction — admin API

Avoid building a broad local admin API in V1.

CLI/TUI should control local runtime through a narrow authenticated control socket/API.

This reduces attack surface.

---

# 433. P1 correction — control plane vs data plane

Separate:

```text
data plane:
  /v1/messages
  /v1/models

control plane:
  status
  reload
  diagnostics
```

Use different auth scopes if both are network-accessible.

---

# 434. P1 correction — Unix socket / named pipe

For control operations:

```text
Unix domain socket
Windows named pipe
```

is preferable to exposing admin HTTP on LAN.

---

# 435. P1 correction — local token scopes

If using a local gateway token:

```text
data-plane token
control-plane token
```

should be distinct where practical.

---

# 436. P1 correction — runtime token rotation

On runtime restart:

```text
new token
```

Old token invalid.

---

# 437. P1 correction — session authentication

If multiple local users could access a machine, do not assume loopback means same user.

Bind/control access should respect OS user boundaries.

---

# 438. P1 correction — Windows ACLs

Windows named pipes/control files need user ACLs.

---

# 439. P1 correction — Unix socket permissions

Use:

```text
0600
```

or equivalent user-only permissions.

---

# 440. P1 correction — state file corruption recovery

If state is corrupt:

```text
move to .corrupt timestamped backup
rebuild minimal state
```

Do not crash permanently.

---

# 441. P1 correction — cache corruption recovery

Caches should always be disposable:

```bash
ccx cache clear
```

and automatic rebuild should work.

---

# 442. P1 correction — provider model catalog failures

If refresh fails:

```text
retain previous valid catalog
mark stale
```

Do not replace a valid catalog with an empty catalog.

---

# 443. P1 correction — empty discovery is not always zero models

Distinguish:

```text
SUCCESS_EMPTY
ERROR
UNSUPPORTED
AUTH_FAILED
```

An empty catalog can mean provider supports discovery but has no models available to that credential.

---

# 444. P1 correction — discovery pagination

Model endpoints may paginate.

Support:

```text
next cursor
page token
```

and enforce a maximum.

---

# 445. P1 correction — duplicate model IDs

Provider discovery can return duplicates.

Normalize:

```text
canonical identity
```

and preserve provider raw metadata.

---

# 446. P1 correction — model deprecation

Track:

```text
deprecated
retired
replacement
```

but do not automatically replace active models without policy.

---

# 447. P1 correction — aliases can be user-created

Support:

```bash
ccx alias model coder openrouter/qwen/...
```

but aliases must resolve to immutable canonical IDs at session start.

---

# 448. P1 correction — project reproducibility

Project config should be able to pin:

```text
provider
model
capability requirements
```

so another contributor can reproduce the route without sharing credentials.

---

# 449. P1 correction — environment-specific credentials

The project must not contain:

```text
credential ID that only exists on one machine
```

unless marked optional.

---

# 450. P1 correction — team config vs personal config

Future-friendly structure:

```text
.ccx/config.yaml        # shareable
~/.ccx/config.yaml      # personal
```

Credentials remain outside both.

---

# 451. P1 correction — project trust metadata

Trust should be stored outside the repository:

```text
~/.ccx/trust/
```

so a repo cannot declare itself trusted.

---

# 452. P1 correction — malicious model metadata

A provider's model description can contain prompt-injection-like text.

Never feed model catalog descriptions into an LLM decision process in V1.

---

# 453. P1 correction — provider docs are not executable

Provider metadata should never be interpreted as:

```text
shell
code
template
```

unless a signed plugin explicitly declares executable behavior.

---

# 454. P1 correction — update metadata is untrusted until verified

Do not parse update instructions from an unsigned remote document and execute them.

---

# 455. P1 correction — crash-loop protection

If CCX crashes N times within a window:

```text
disable auto-restart
```

and tell the user:

```text
CCX entered recovery mode.
Run ccx doctor.
```

---

# 456. P1 correction — safe mode

Add:

```bash
ccx --safe-mode
```

which disables:

```text
plugins
experimental compatibility hacks
automatic fallback
```

and uses minimal functionality.

This is extremely useful for recovery.

---

# 457. P1 correction — recovery command

Add:

```bash
ccx repair
```

which is distinct from:

```text
doctor
```

Doctor diagnoses.

Repair applies a known safe recovery procedure.

---

# 458. P1 correction — reset-to-known-good

Add:

```bash
ccx repair --known-good
```

which can restore the previous known-good configuration snapshot.

---

# 459. P1 correction — last-known-good route

Store:

```text
last successful route
```

but never automatically switch to it without policy.

It can be offered as a recovery suggestion.

---

# 460. P1 correction — “working” means end-to-end

Provider verification must not stop at:

```text
HTTP 200
```

A provider is `VERIFIED` only when the required feature tests pass.

The source research explicitly emphasizes live probing before model selection and capability verification after selection. fileciteturn2file3L175-L223

---

# 461. P1 correction — acceptance tests must include Claude Code itself

Mock tests are insufficient.

Have at least:

```text
real Claude Code binary
real CCX
mock provider
```

for deterministic end-to-end tests.

Then:

```text
real Claude Code
real CCX
real provider
```

for controlled live tests.

---

# 462. P1 correction — Claude Code test harness

Automate:

```text
launch Claude Code
send controlled prompt
observe gateway request
return deterministic response
assert Claude Code behavior
```

This becomes the core compatibility lab.

---

# 463. P1 correction — test subagents

Because gateway model discovery can affect subagent model resolution in current Claude Code behavior, test:

```text
main model
subagent model
agent() model selection
```

not just top-level generation. citeturn0search1

---

# 464. P1 correction — test compaction

Long context tests must verify:

```text
status line
context usage
auto-compaction
post-compaction continuation
```

This is where model context metadata becomes user-visible.

---

# 465. P1 correction — test permissions

Claude Code may have permission/tool execution behavior independent of CCX.

CCX should verify it does not corrupt:

```text
tool names
tool inputs
tool results
```

but must not attempt to replace Claude Code's permission system.

---

# 466. P1 correction — test hooks/plugins interaction

CCX should test Claude Code installations with:

```text
hooks
plugins
MCP
custom settings
```

to detect environmental conflicts.

CCX should not attempt to manage those systems.

---

# 467. P1 correction — environment snapshot in compatibility tests

Record sanitized:

```text
Claude Code version
CCX version
OS
terminal
provider
model
relevant env keys
```

This makes bugs reproducible.

---

# 468. P1 correction — compatibility bug severity

Use:

```text
P0 data/semantic corruption
P1 major feature break
P2 degraded feature
P3 cosmetic/diagnostic
```

Release policy maps severity to action.

---

# 469. P1 correction — semantic corruption is highest priority

Examples:

```text
wrong tool arguments
wrong tool ID
wrong model
wrong context
duplicate tool proposal
lost user content
```

must outrank:

```text
TUI color issue
```

---

# 470. P1 correction — “silent data loss” invariant

Add hard invariant:

> CCX MUST NOT silently discard user-visible content, tool arguments, tool results, citations, thinking artifacts, or provider response metadata required for continuation.

If loss is unavoidable:

```text
explicit policy
visible diagnostic
```

---

# 471. P1 correction — “wrong model” invariant

Add:

> CCX MUST never route a request to a different canonical model than the resolved route without recording and surfacing an explicit fallback decision.

---

# 472. P1 correction — “wrong provider” invariant

Same for provider.

---

# 473. P1 correction — “wrong credential” invariant

A request must not accidentally cross credential identities because of connection pooling.

---

# 474. P1 correction — session route snapshot hash

Every session should have:

```text
route_hash
```

so logs/doctor can prove which configuration was active.

---

# 475. P1 correction — config revision

Global config:

```text
revision number
```

Session records:

```text
config_revision
```

This makes config races observable.

---

# 476. P1 correction — runtime revision

Runtime should expose:

```text
runtime_id
```

for correlating sessions across restarts.

---

# 477. P1 correction — compatibility rule revision

Every request can record:

```text
compat_ruleset_version
```

This makes old bugs reproducible after compatibility rules evolve.

---

# 478. P1 correction — deterministic route artifact

`ccx run --dry-run --json` should emit a portable route artifact:

```json
{
  "schema_version": 1,
  "route": {},
  "capabilities": {},
  "transforms": [],
  "policies": {}
}
```

Secrets excluded.

---

# 479. P1 correction — route artifact replay

Future test command:

```bash
ccx replay-route route.json
```

for deterministic local debugging.

---

# 480. P1 correction — compatibility snapshots

Every certified provider/model can produce:

```text
compatibility snapshot
```

that records:

```text
protocol
features
rules
verification
```

---

# 481. P1 correction — provider onboarding should have a “why failed” tree

Instead of:

```text
Connection failed
```

show:

```text
DNS ✓
TCP ✓
TLS ✓
HTTP ✓
Auth ✗
```

This is much more useful.

---

# 482. P1 correction — model verification should have stage output

```text
Request construction ✓
Upstream accepted ✓
First stream event ✓
Final event ✓
Tool call ✓
Cancellation ✓
```

---

# 483. P1 correction — verification should support partial success

A model can be:

```text
text ✓
stream ✓
tools ✗
reasoning ?
```

Do not collapse to one boolean.

---

# 484. P1 correction — capability requirements syntax

Define:

```yaml
requirements:
  tools: required
  parallel_tools: preferred
  reasoning: required
  vision: optional
  context:
    min: 200000
```

States:

```text
required
preferred
optional
forbidden
```

This is better than only required/optional.

---

# 485. P1 correction — candidate selection

Candidate route is valid only if:

```text
all required capabilities satisfied
```

Preferred capabilities improve ranking.

Forbidden capabilities can exclude candidates.

---

# 486. P1 correction — routing ranking should be deterministic

V1 ranking:

```text
required capability satisfaction
>
explicit user order
>
verification freshness
>
health
>
cost
>
latency
```

Do not use opaque model scoring.

---

# 487. P1 correction — cost ranking needs user budget

A cheap model isn't better if it lacks tools.

---

# 488. P1 correction — latency ranking should use observed measurements

Not marketing claims.

---

# 489. P1 correction — provider health scores need decay

Old good health must decay over time.

---

# 490. P1 correction — routing decisions need reason codes

Examples:

```text
SELECT_EXPLICIT
SELECT_PROFILE_FIRST
FALLBACK_CAPABILITY_MATCH
FALLBACK_HEALTH
```

---

# 491. P1 correction — user should be able to pin a route

```bash
ccx pin openrouter/qwen/...
```

This disables dynamic selection for that profile.

---

# 492. P1 correction — route locking

Session should record:

```text
locked = true
```

if user explicitly pinned it.

---

# 493. P1 correction — fallback visibility

If automatic fallback is enabled:

```text
╭──────────────────────────────────╮
│ FALLBACK                          │
│ Primary unavailable               │
│ OpenRouter/Qwen → OpenAI/GPT     │
│ Reason: upstream 503             │
│ [Continue] [Abort]               │
╰──────────────────────────────────╯
```

For noninteractive/headless:

```text
structured event/log
```

---

# 494. P1 correction — fallback in CI

CI should default:

```text
fallback = disabled
```

because reproducibility matters.

---

# 495. P1 correction — cost guardrails in CI

Allow:

```text
max_cost
max_requests
```

and fail fast.

---

# 496. P1 correction — local-only CI

Support:

```text
mock provider
```

so users can test CCX without API credentials.

---

# 497. P1 correction — test provider

Ship an in-process/local deterministic provider:

```text
ccx provider test-server
```

for development.

This should simulate:

```text
streaming
tools
thinking
errors
timeouts
malformed events
```

---

# 498. P1 correction — chaos test server

A deterministic local chaos server dramatically lowers testing cost.

---

# 499. P1 correction — protocol playground

Add developer command:

```bash
ccx dev replay fixture.json
```

for fast iteration.

---

# 500. P1 correction — fixture minimization

When a provider bug is found:

```text
capture
sanitize
minimize
promote to fixture
```

---

# 501. P1 correction — documentation must explain what CCX cannot fix

Examples:

```text
model genuinely lacks tools
provider refuses long context
Claude Code hardcodes behavior with no supported override
provider account lacks access
```

Honesty is part of the product.

---

# 502. P1 correction — compatibility ceiling

Define:

```text
CCX cannot override Claude Code behavior unless:
1. documented configuration exists, or
2. a version-gated supported integration exists.
```

No undocumented binary patching in V1.

---

# 503. P1 correction — do not patch Claude Code binary

Hard invariant:

> CCX MUST NOT modify, inject into, patch, or reverse-engineer the installed Claude Code binary as a runtime dependency.

Compatibility comes from the boundary.

---

# 504. P1 correction — no filesystem monkey-patching by default

Same principle for internal Claude Code files.

---

# 505. P1 correction — compatibility shim is a product subsystem

It needs:

```text
version registry
rule registry
fixtures
tests
release watch
```

not scattered `if` statements.

---

# 506. P1 correction — compatibility changes require changelog entries

Every compatibility rule change:

```text
CHANGELOG
compatibility notes
fixture
```

---

# 507. P1 correction — provider changes require compatibility notes

Same.

---

# 508. P1 correction — release process

Release checklist:

```text
build
test
compatibility
security
SBOM
sign
package
smoke
publish
```

---

# 509. P1 correction — post-release smoke

After publishing:

```text
download released artifact
verify signature
run clean-machine smoke test
```

---

# 510. P1 correction — rollback release

Maintain ability to:

```text
mark release bad
```

and direct users to previous verified version.

---

# 511. P1 correction — project roadmap separation

### V1 must have

```text
Claude Code ingress
local runtime
native Anthropic path
OpenAI Chat
OpenAI Responses
custom endpoints
model discovery
capability engine
streaming
tools
thinking handling
doctor
TUI
profiles
auth
security
reliability
compatibility lab
packaging
```

### V1 should NOT require

```text
LLM-based routing
multi-agent orchestration
web dashboard
hosted control plane
dynamic executable plugins
automatic model downloads
provider billing APIs
```

This keeps the project shippable.

---

# 512. V2 priorities

```text
capability-aware routing
cost/latency optimization
provider SDK
community certification
advanced compatibility remediation
role-based model routing
benchmark portal
```

---

# 513. V3 priorities

```text
OpenCode
Codex
Aider
other agent ingress adapters
shared agent semantic IR
```

---

# 514. V4 research

```text
adaptive routing
learned routing
multi-model deliberation
automatic cost/performance optimization
```

Keep experimental work out of the reliability core.

---

# 515. Revised 1.0 hard-invariant set

The previous invariant list should be expanded to:

```text
1. No secret leakage.
2. No silent model change.
3. No silent provider change.
4. No duplicate side-effecting tool proposal caused by CCX retry.
5. No silent user-visible content loss.
6. No LAN exposure by default.
7. No mandatory CCX cloud.
8. No undocumented Claude Code binary patching.
9. No undocumented Claude Code filesystem mutation by default.
10. No provider response can crash CCX.
11. No unsupported capability presented as supported.
12. No unbounded buffering.
13. Cancellation propagates where supported.
14. Active sessions have immutable route snapshots.
15. Config changes are transactional.
16. Updates are atomic and recoverable.
17. Credentials never enter ordinary config/logs/crash reports.
18. Unknown protocol fields/events fail safely.
19. Tool IDs remain correctly mapped and isolated.
20. Protocol adapters cannot directly mutate UI/config.
21. Provider adapters cannot access unrelated credentials.
22. Project config cannot weaken user security policy.
23. Retry cannot exceed configured global/per-request budgets.
24. Fallback cannot occur unless policy allows it.
25. Every lossy transform is observable.
26. Every compatibility workaround is version-scoped.
27. Every compatibility bug gets a regression fixture.
28. JSON automation output is schema-versioned.
29. `ccx claude` preserves TTY/stdin/stdout/stderr/exit semantics.
30. Offline mode makes no unexpected network requests.
```

---

# 516. Revised 1.0 acceptance matrix

| Area | Pass condition |
|---|---|
| Install | clean machine succeeds |
| First run | setup completes without editing YAML |
| Launch | `ccx claude` starts real Claude Code |
| TTY | interactive Claude Code remains intact |
| Auth | credential source is explicit and secret-safe |
| Discovery | models are discoverable or clearly manual |
| Text | multi-turn generation works |
| Streaming | valid SSE semantics preserved |
| Thinking | supported artifacts survive round trip |
| Tools | IDs/input/results survive round trip |
| Fine-grained tools | partial JSON is handled correctly |
| Images/docs | supported paths preserve semantics |
| Count tokens | endpoint behavior is deterministic |
| Models | stale/error/empty discovery states are distinct |
| Context | safe context limit is correctly reported |
| Retry | no unsafe replay |
| Fallback | explicit policy only |
| Cancellation | upstream request is cancelled where possible |
| Sleep/wake | runtime recovers |
| Crash | runtime restarts safely |
| Concurrency | sessions remain isolated |
| Security | no secret/log/SSRF regressions |
| Privacy | no hidden telemetry |
| Diagnostics | failures produce actionable reason codes |
| TUI | 80×24 + NO_COLOR works |
| Headless | JSON commands are scriptable |
| Update | verify → activate → rollback works |
| Compatibility | tested Claude Code versions are green |
| Provider | certified providers pass required suite |
| Fuzzing | no known crashes |
| Stress | bounded memory/resources |
| Docs | install/troubleshooting/security documented |

---

# 517. Final self-review verdict

After this audit:

### What was already genuinely good

The original plan correctly focused on the hard stuff rather than merely “supporting many APIs”:

- local-only architecture;
- native-provider fast paths;
- explicit capability verification;
- model-picker problem awareness;
- beta-header drift;
- keepalive-safe streaming;
- tool-parser-aware vLLM support;
- cancellation;
- secret redaction;
- explicit fallback;
- doctor-first thinking. fileciteturn2file6L398-L413

### What was dangerously underspecified

The biggest holes were:

```text
wire surface
content-block completeness
thinking signatures
count_tokens
tool replay state
request/response identity
error taxonomy
auth lifecycle
local gateway auth
SSRF
project trust
TTY passthrough
Claude Code launcher semantics
config transactions
update recovery
route immutability
stream event journaling
MCP/server tools
fine-grained tool streaming
compatibility versioning
real Claude Code E2E testing
```

Those are now explicitly covered.

### The most important architectural principle

**CCX should not become a giant pile of provider-specific `if/else` statements.**

The real architecture is:

```text
Claude Code
    ↓
ClaudeCodeCompat
    ↓
Anthropic Wire Model
    ↓
Semantic IR
    ↓
Capability / Policy
    ↓
Transform
    ↓
Target Protocol
    ↓
Provider Adapter
```

with:

```text
Compatibility Rules
Capability Evidence
Provider Behavior Profiles
Regression Fixtures
```

surrounding it.

That is what makes CCX maintainable when Claude Code and providers inevitably move.

---

# 518. Final product standard

The bar for CCX 1.0 should be:

> **If a user can install CCX, configure a provider once, run Claude Code through it for hours, use tools and streaming normally, survive ordinary provider/network failures, understand exactly what model/provider is active, and diagnose the rare failures without manually studying five API specifications — CCX is doing its job.**

If it merely converts:

```text
Anthropic JSON → OpenAI JSON
```

then honestly, **don't build it**.

There are already projects doing that.

The defensible product is:

```text
Universal protocol compatibility
+
Claude Code-specific compatibility
+
verified capability matrix
+
safe degradation
+
automatic lifecycle
+
professional terminal UX
+
diagnostics
+
compatibility laboratory
```

That is the version worth building.
