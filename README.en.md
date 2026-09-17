# Token Saver 2.0

Inspect Claude Code configuration, review useful changes, apply and undo them.
Preserves model choice and reasoning effort. No Python or Node runtime required.

## Install

Download and extract the release package for your operating system and CPU.
On Windows open `instalar.cmd` (or run `install.ps1`). On macOS/Linux run
`sh install.sh` in the extracted directory. Open a new Claude Code session and
run `/token-saver`. The skill responds in your language.

Public download commands become available when the maintainer publishes a release.
Local packages work without publication or a network connection during installation.

## Use

- `/token-saver`: guided diagnosis and review.
- `/token-saver audit`: read-only diagnosis.
- `/token-saver undo`: revert the last application.
- `/token-saver details`: evidence and limitations.
- Ask for “modo direto” to opt into a concise-response project rule.

The basic adjustment only re-enables Tool Search in local project settings when
an explicit disabling setting and compatible environment can be identified.
Missing environment variables do not imply a problem. Most modern installations
may need no automatic changes. Models, thinking, permissions and CLAUDE.md are
preserved. No `.claudeignore` generation or background compaction.

Configuration improvements are not measurements of token savings. Verify session
behavior through `/status`, `/context` and `/usage`. The CLI does not measure a
counterfactual baseline, actual session tokens, or subscription savings.

Invalid JSON stops the operation. Backups remain local outside the project;
they may contain sensitive configuration. Undo preserves unrelated later edits
and refuses same-field conflicts. Existing customized skill installations are
preserved; review and rename them before migrating.

The CLI offers `--project`, `--config-dir`, `--json`, `--dry-run`, `--yes` and
`--expect`. Sources and binaries are MIT licensed. See CONTRIBUTING.md for builds
and native-platform validation requirements.
