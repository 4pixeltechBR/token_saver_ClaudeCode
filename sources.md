# Referências e fontes — Token Saver

## Fontes primárias (Anthropic)

| Claim | URL |
|-------|-----|
| Tool Search API + 85% redução | https://www.anthropic.com/engineering/advanced-tool-use |
| Best practices oficiais (/compact, /clear, subagents) | https://code.claude.com/docs/en/best-practices |
| Subagents — isolamento de contexto | https://code.claude.com/docs/en/sub-agents |
| /btw command | https://code.claude.com/docs/en/sub-agents |
| Model config / MAX_THINKING_TOKENS / adaptive thinking | https://code.claude.com/docs/en/model-config |
| MCPs — tool definitions em contexto | https://code.claude.com/docs/en/mcp |

## Variáveis de ambiente

| Fonte | URL |
|-------|-----|
| Lista completa reverse-engineered | https://gist.github.com/jedisct1/9627644cda1c3929affe9b1ce8eaf714 |
| Versão v2.1.104 verificada | https://gist.github.com/mculp/e6a573f2a45ef7dbbf30f6a8574c7351 |
| Referência everything-claude-code | https://github.com/affaan-m/everything-claude-code/blob/main/docs/token-optimization.md |

## Benchmarks e casos reais

| Claim | Fonte | Magnitude |
|-------|-------|-----------|
| 143k/200k tokens pré-conversa com MCPs | Scott Spence blog | 72% do contexto consumido antes de digitar |
| .claudeignore elimina 30–40% em projetos Next.js | 32blog.com | Maior ganho único simples |
| /compact a 60% vs 95% — diferença de qualidade | MindStudio blog | Qualitativo + documentado |
| Plan mode + .claudeignore = −50% total | 32blog.com | Combinação validada |
| CLAUDE.md de Boris Cherny = ~2.500 tokens | Cuttlesoft | Benchmark do criador do Claude Code |
| Subagente Haiku vs Opus = −80% | Ratio de preços Anthropic | Haiku $0.25/$1.25 vs Opus $15/$75 per Mtok |
| MAX_THINKING_TOKENS=10000 = −70% thinking | everything-claude-code | Reduz de 31.999 default |
| 51k → 8.5k tokens com Tool Search (167 tools) | Joe Njenga, Medium | −83% tokens de tools |

## Bugs conhecidos (abril 2026)

| Bug | Issue | Workaround |
|-----|-------|-----------|
| ENABLE_TOOL_SEARCH não funciona no Desktop App via env var | #41472 | Usar settings.json |
| ENABLE_TOOL_SEARCH auto mode (10% threshold) não dispara | #18397, #19890, #18298 | Setar `true` explicitamente |
| defer_loading + cache_control incompatíveis | #30920 | `ENABLE_TOOL_SEARCH=false` se necessário |
| CLAUDE_CODE_MAX_OUTPUT_TOKENS não se aplica a subagentes | #25569 | Bug aberto, sem fix |
| Caching bugs de março 2026 causaram 10–20x inflação | #40524 | Resolvido em patch posterior |

## Caveman

| Claim | Fonte | Magnitude |
|-------|-------|-----------|
| 22–87% redução output tokens | README oficial + evals/ | Range real, não média |
| ~14–21% em benchmarks independentes (72 runs) | Kuba Guzik, DEV.to | Metodologia publicada |
| 4–10% de economia total de sessão | Pasquale Pillitteri | Output = ~25% do total |
| /caveman:compress = ~46% redução input por sessão | Documentação Caveman | Input tokens de CLAUDE.md |

## Idioma

| Claim | Fonte | Magnitude |
|-------|-------|-----------|
| Português ~50% mais tokens que inglês | arXiv 2305.15425 | Paper revisado |
| "system..." (EN) = 7 tok vs equivalente PT = 14 tok | Tokenizador OpenAI (platform.openai.com/tokenizer) | Verificável |
