# Validação da entrega 3.0.0

Data: 18/09/2026. Base: release multi-harness do repositório
`4pixeltechBR/token_saver_ClaudeCode`.

## Executado

- `go test ./...` aprovado no Windows amd64, incluindo auditoria portátil,
  plano somente leitura e destinos de instalação por harness.
- `go vet ./...` aprovado.
- Compilação cruzada para Windows, macOS e Linux, amd64 e arm64.
- `token-saver version` retorna `3.0.0`.
- `token-saver detect --json` lista os nove ambientes suportados pelo catálogo.
- Auditoria `codex` executada em modo somente leitura e sem escrever no projeto.
- Seis ZIPs de plataforma e pacote `agent-plugin` gerados.
- `SHA256SUMS.txt` gerado para todos os sete ZIPs.
- O pacote de plugin contém `plugin.json`, `.claude-plugin/plugin.json`,
  `skills/SKILL.md` e `skills/token-saver/SKILL.md`, sem binário nativo.

## Capacidades verificadas

O Claude Code mantém aplicação, backup, validação e undo completos. Codex,
Antigravity, OpenCode, Cursor, Gemini CLI, GitHub Copilot, Cline e MiniMax Code
foram adicionados ao catálogo de descoberta, instalação e auditoria conservadora.
Eles permanecem somente leitura até que exista um adaptador de escrita testado
para o formato e as métricas de cada host.

## Limites

Os binários macOS, Linux e ARM foram compilados, mas não executados nativamente
neste Windows. O plugin sem binário exige o companion CLI para auditoria JSON
determinística, aplicação e rollback. Não foram feitos benchmarks de economia,
testes com usuários ou assinatura/notarização de binários.
