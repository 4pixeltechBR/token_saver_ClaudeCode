# Matriz de harnesses

## Capacidades da versão 3.0

| Harness | Skill local | Invocação usual | Auditoria | Escrita/undo |
|---|---|---|---|---|
| Claude Code | `~/.claude/skills` | `/token-saver` | Sim | Sim |
| Codex | `~/.agents/skills` | `$token-saver` | Sim | Ainda não |
| Antigravity | `~/.gemini/antigravity/skills` | Por intenção | Sim | Ainda não |
| OpenCode | `~/.config/opencode/skills` | Skill do agente | Sim | Ainda não |
| Cursor | `~/.cursor/skills` | `/token-saver` | Sim | Ainda não |
| Gemini CLI | `~/.gemini/skills` | Skill ativada | Sim | Ainda não |
| GitHub Copilot | `~/.copilot/skills` | `/token-saver` | Sim | Ainda não |
| Cline | `~/.cline/skills` | Skill do agente | Sim | Ainda não |
| MiniMax Code | Plugin `skills/token-saver` | Skill do plugin | Sim | Ainda não |

As pastas `.agents/skills` e equivalentes de cada host são caminhos de
descoberta. Elas não tornam configurações, métricas ou comandos de sessão
intercambiáveis. A ferramenta permanece em modo somente leitura quando não
consegue verificar a semântica do host.
