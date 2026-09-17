# Token Saver

A simple, reversible way to find wasted context in Claude Code. Token Saver
preserves your model, reasoning effort and instructions.

[Download the latest release](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest) · [View the source](https://github.com/4pixeltechBR/token_saver_ClaudeCode)

## Install

1. Open the [latest release](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest) and download the ZIP for your computer:

   - [Windows Intel/AMD](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-windows-amd64.zip)
   - [Windows ARM](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-windows-arm64.zip)
   - [macOS Apple Silicon](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-darwin-arm64.zip)
   - [macOS Intel](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-darwin-amd64.zip)
   - [Linux Intel/AMD](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-linux-amd64.zip)
   - [Linux ARM64](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-linux-arm64.zip)

2. Extract the ZIP.
3. On Windows open `instalar.cmd`; on macOS/Linux run `sh install.sh` in the extracted directory.
4. Open a new Claude Code session and run `/token-saver`.

The package includes the executable and skill. Python and Node are not required.
Use `SHA256SUMS.txt` to verify the download before installing if desired.

## Use

The first run reads the project, shows up to three findings and proposes only
compatible, verifiable changes. You review the preview before applying anything.
If no confirmed gain exists, it says that nothing needs changing.

| Goal | Command |
|---|---|
| Start | `/token-saver` |
| Diagnose only | `/token-saver audit` |
| Undo the last application | `/token-saver undo` |
| See details and limits | `/token-saver details` |
| Opt into concise responses | Ask for “modo direto” |

The basic flow does not change the model or reasoning effort, edit `CLAUDE.md`,
create `.claudeignore` or compact conversations automatically. Modo Direto is
optional and adds one short project rule.

## What is measured

The diagnostic reports local settings, instruction-file sizes, detectable MCP
declarations and skills. It does not measure session tokens or cost. Confirm real
behavior in Claude Code with `/status`, `/context` and `/usage`. A subscription
does not become cheaper automatically.

## Safety and undo

Invalid JSON is preserved and stops application. Writes are atomic, previews have
an identifier and history stays outside the project. Undo preserves later edits in
other keys and refuses same-key conflicts.

Existing customized installations are not overwritten. See [recovery](skill/references/recovery.md) and [technical limits](skill/references/details.md).

## For contributors

The CLI also provides `audit`, `plan`, `apply`, `undo`, `install` and `details`,
with `--json`, `--dry-run`, `--project`, `--config-dir` and `--expect`. See
[CONTRIBUTING.md](CONTRIBUTING.md) for tests and packaging.

## License

MIT. [Sources and validity](skill/references/sources.md) · [Português](README.md)
