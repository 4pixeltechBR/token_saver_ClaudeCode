# Token Saver 3.0

Safe context auditing for Claude Code, Codex, Antigravity, OpenCode, Cursor,
Gemini CLI, GitHub Copilot, Cline, and MiniMax Code.

[Download the release](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest) · [Source](https://github.com/4pixeltechBR/token_saver_ClaudeCode) · [Harness matrix](skill/references/harnesses.md)

## What changed

This release separates the portable skill from harness-specific behavior. The
same package can be discovered by several agents, but it only writes a setting
when the format, version, and rollback behavior are known.

The core is conservative: audit first, show evidence, create an identified plan,
and never turn file size into a made-up token saving estimate.

## Capabilities

| Function | Claude Code | Other harnesses |
|---|---:|---:|
| Detect available environments | Yes | Yes |
| Audit project and instructions | Yes | Yes |
| Detect stack | Yes | Yes |
| Count discoverable skills | Yes | Yes |
| JSON output for automation | Yes | Yes |
| Preview (`plan`) | Yes | Yes, read-only |
| Apply settings | Yes, when eligible | Not yet |
| Undo an application | Yes | Not yet |
| Direct Mode | Yes, opt-in | Manual guidance |
| Telemetry | Never | Never |

### What Claude Code can apply

- Re-enable Tool Search only when an explicit disablement exists and compatibility
  is confirmed.
- Optionally add a short Direct Mode rule to `.claude/rules`.
- Create backups, write atomically, validate, and roll back with conflict and
  symlink protection.

### What every harness can audit

- Long instruction files and possible repetition.
- Stack detected from project manifests.
- Skills found in host-supported directories.
- Whether the skill is discoverable.
- Limits and next steps without inventing usage or savings.

Token Saver does not change models, reduce reasoning effort, edit `CLAUDE.md`,
create `.claudeignore`, compact conversations automatically, or install external
plugins.

## Install

Download the latest system ZIP, extract it, and run the binary. Python and Node
are not required.

```text
token-saver detect --project .
token-saver install --harness claude --project .
token-saver install --harness codex --project .
```

Use `--harness auto` to select the first detected environment. Open a new session
or reload skills after installation.

| Environment | Global skill path | Invocation |
|---|---|---|
| Claude Code | `~/.claude/skills/token-saver` | `/token-saver` |
| Codex | `~/.agents/skills/token-saver` | `$token-saver` |
| Antigravity | `~/.gemini/antigravity/skills/token-saver` | ask for an audit |
| OpenCode | `~/.config/opencode/skills/token-saver` | use the `token-saver` skill |
| Cursor | `~/.cursor/skills/token-saver` | `/token-saver` |
| Gemini CLI | `~/.gemini/skills/token-saver` | `/skills`, then enable it |
| GitHub Copilot | `~/.copilot/skills/token-saver` | `/token-saver` |
| Cline | `~/.cline/skills/token-saver` | use the `token-saver` skill |
| MiniMax Code | plugin `skills/token-saver` | enable the plugin |

## Use

Claude Code’s main flow is:

```text
/token-saver
/token-saver auditar
/token-saver desfazer
/token-saver detalhes
/token-saver modo direto
```

The CLI works for every supported host:

```text
token-saver detect --json
token-saver audit --harness codex --project . --json
token-saver audit --harness opencode --project .
token-saver plan --harness claude --project . --json
token-saver apply --harness claude --project . --yes --expect PLAN_ID --json
token-saver undo --harness claude --project . --json
token-saver details
token-saver version
```

Common options are `--project`, `--harness`, `--config-dir`, `--json`,
`--dry-run`, `--yes`, `--expect`, and `--concise`.

## Honest token reporting

The audit reads local files and configuration. It does not measure input tokens,
output tokens, cache, real cost, or counterfactual savings. A saved setting does
not prove a session improvement. Use the usage or context report provided by your
own harness.

## Safety and recovery

The default mode is read-only. Invalid JSON, managed configuration, unknown
versions, or conflicts stop an application. Claude Code keeps history outside the
project, writes atomically, and preserves later edits to unrelated keys. See
[recovery](skill/references/recovery.md) and [details](skill/references/details.md).

## Develop

```text
go test ./...
go vet ./...
python scripts/package.py --out dist
```

The repository contains the Go core, the portable skill in `skill/`, and discovery
adapters in `harness.go` and `portable_audit.go`.

## License

MIT. See [CHANGELOG](CHANGELOG.md) and [Português](README.md).
