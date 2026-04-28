---
name: token-saver
description: >
  Diagnóstica, configura e otimiza o consumo de tokens do Claude Code — cobrindo todas
  as camadas validadas: setup único, variáveis de ambiente, MCPs, CLAUDE.md, sessão,
  subagentes e compressão de output. Ativa SEMPRE que o usuário mencionar tokens, custo,
  contexto cheio, sessão lenta, quota acabando, limites de uso, "economizar", "otimizar",
  "contexto explodindo", "reduzir custo", "/compact", "CLAUDE.md grande", MCP pesado,
  ou qualquer variação de gastar menos com Claude Code. Funciona em CLI/terminal e
  Desktop App (Windows/macOS). Dois modos: AUDIT (lê e reporta) e SETUP (aplica correções).
---

# Token Saver — Guia Definitivo de Economia

Cobre 5 camadas de otimização em ordem decrescente de impacto. Comece sempre pelo AUDIT.

---

## Modo AUDIT — diagnóstico rápido

Rode este bloco inteiro. Ele funciona em bash (Linux/macOS) e PowerShell (Windows).

```bash
# === BASH (Linux / macOS / WSL) ===
echo "=== CLAUDE.md do projeto ==="
[ -f CLAUDE.md ] && wc -w CLAUDE.md | awk '{print $1" palavras (~"int($1*1.3)" tokens)"}' || echo "Não encontrado"

echo "=== CLAUDE.md global ==="
[ -f ~/.claude/CLAUDE.md ] && wc -w ~/.claude/CLAUDE.md | awk '{print $1" palavras"}' || echo "Não encontrado"

echo "=== MCPs configurados ==="
python3 -c "
import json, os
path = os.path.expanduser('~/.claude.json')
if os.path.exists(path):
    d = json.load(open(path))
    mcps = d.get('mcpServers', {})
    print(f'{len(mcps)} MCPs: {list(mcps.keys())}')
else:
    print('~/.claude.json não encontrado')
" 2>/dev/null

echo "=== Skills instaladas ==="
ls ~/.claude/skills/ 2>/dev/null | wc -l | xargs -I{} echo "{} skills (~$(( $(ls ~/.claude/skills/ 2>/dev/null | wc -l) * 100 )) tokens fixos)"

echo "=== .claudeignore ==="
[ -f .claudeignore ] && echo "Existe ($(wc -l < .claudeignore) regras)" || echo "AUSENTE — risco de indexar node_modules/binários"

echo "=== Tool Search ==="
[ "$ENABLE_TOOL_SEARCH" = "true" ] && echo "ATIVO" || echo "INATIVO — definir ENABLE_TOOL_SEARCH=true"

echo "=== Env vars de custo ==="
echo "MAX_THINKING_TOKENS=${MAX_THINKING_TOKENS:-não definido}"
echo "CLAUDE_CODE_SUBAGENT_MODEL=${CLAUDE_CODE_SUBAGENT_MODEL:-não definido (herda modelo principal)}"
```

```powershell
# === POWERSHELL (Windows Desktop App) ===
Write-Host "=== CLAUDE.md do projeto ===" -ForegroundColor Cyan
if (Test-Path CLAUDE.md) {
    $words = (Get-Content CLAUDE.md | Measure-Object -Word).Words
    Write-Host "$words palavras (~$([int]($words * 1.3)) tokens)"
} else { Write-Host "Não encontrado" }

Write-Host "=== MCPs configurados ===" -ForegroundColor Cyan
$claudeJson = "$env:USERPROFILE\.claude.json"
if (Test-Path $claudeJson) {
    $cfg = Get-Content $claudeJson | ConvertFrom-Json
    $mcps = $cfg.mcpServers.PSObject.Properties.Name
    Write-Host "$($mcps.Count) MCPs: $($mcps -join ', ')"
} else { Write-Host "~/.claude.json não encontrado" }

Write-Host "=== Skills instaladas ===" -ForegroundColor Cyan
$skillsPath = "$env:USERPROFILE\.claude\skills"
if (Test-Path $skillsPath) {
    $count = (Get-ChildItem $skillsPath -Directory).Count
    Write-Host "$count skills (~$($count * 100) tokens fixos por sessão)"
}

Write-Host "=== .claudeignore ===" -ForegroundColor Cyan
if (Test-Path .claudeignore) { Write-Host "Existe" } else { Write-Host "AUSENTE" }

Write-Host "=== Tool Search ===" -ForegroundColor Cyan
Write-Host "ENABLE_TOOL_SEARCH=$env:ENABLE_TOOL_SEARCH"
```

Depois de rodar, use esta tabela para priorizar:

| Sinal | Risco | Ação imediata |
|-------|-------|---------------|
| MCPs > 5 sem Tool Search | 🔴 | Camada 1-A |
| CLAUDE.md > 3.000 tokens | 🔴 | Camada 2-A |
| Sem .claudeignore | 🔴 | Camada 1-B |
| Skills > 15 | 🟡 | Camada 2-B |
| Subagent model = Opus | 🟡 | Camada 1-C |
| MAX_THINKING_TOKENS não definido | 🟡 | Camada 1-C |

---

## Camada 1 — Setup único (maior ROI)

### 1-A: Tool Search — lazy loading de MCPs

**O que faz:** Em vez de enviar *todas* as definições de tools em cada turno (cada MCP pesado custa 7k–17k tokens), o modelo busca tools sob demanda. Impacto documentado: **−85% nos tokens de tool definitions**.

```bash
# BASH — adicionar ao ~/.bashrc ou ~/.zshrc
echo 'export ENABLE_TOOL_SEARCH=true' >> ~/.bashrc && source ~/.bashrc
```

```powershell
# POWERSHELL (Windows) — persistente para o usuário
[System.Environment]::SetEnvironmentVariable("ENABLE_TOOL_SEARCH", "true", "User")
# Reiniciar o terminal / Desktop App depois
```

Ou via `~/.claude/settings.json` (funciona CLI e Desktop App):
```json
{
  "env": {
    "ENABLE_TOOL_SEARCH": "true"
  }
}
```

> **Nota Desktop App (Windows):** `ENABLE_TOOL_SEARCH` via variável de sistema tem bug ativo (issue #41472). O método mais confiável é via `settings.json`. Verificar se está funcionando com `/context` — se tool definitions ainda aparecem grandes, usar o fallback de desabilitar MCPs não usados manualmente.

### 1-B: .claudeignore — cortar indexação desnecessária

Cria o arquivo na raiz do projeto. Impacto: até **−40% de contexto** em projetos JS/TS antes do primeiro prompt.

```bash
# BASH — criar .claudeignore no projeto atual
cat > .claudeignore << 'EOF'
# Dependências e builds
node_modules/
.next/
.nuxt/
dist/
build/
out/
.output/
__pycache__/
.venv/
venv/
*.egg-info/
.pytest_cache/

# Binários e mídia
*.png
*.jpg
*.jpeg
*.gif
*.webp
*.svg
*.ico
*.woff
*.woff2
*.ttf
*.mp4
*.mp3
*.zip
*.tar.gz
*.exe
*.dll
*.so

# Dados e logs
*.log
logs/
*.sqlite
*.sqlite3
*.db
*.csv
*.parquet
*.lock
package-lock.json
yarn.lock
pnpm-lock.yaml

# Cache e temporários
.cache/
tmp/
temp/
.turbo/
.vercel/
.netlify/
coverage/
.nyc_output/
EOF
echo "✓ .claudeignore criado"
```

```powershell
# POWERSHELL (Windows)
@"
node_modules/
.next/
dist/
build/
__pycache__/
*.png
*.jpg
*.gif
*.exe
*.dll
*.zip
*.log
*.sqlite
*.sqlite3
package-lock.json
yarn.lock
pnpm-lock.yaml
coverage/
.cache/
"@ | Set-Content .claudeignore -Encoding UTF8
Write-Host "✓ .claudeignore criado"
```

### 1-C: settings.json — variáveis de ambiente de custo

Arquivo de configuração central. Localização:
- **Linux/macOS:** `~/.claude/settings.json`
- **Windows:** `%USERPROFILE%\.claude\settings.json`

Configuração recomendada para economizar tokens:

```json
{
  "model": "claude-sonnet-4-6",
  "env": {
    "ENABLE_TOOL_SEARCH": "true",
    "CLAUDE_CODE_SUBAGENT_MODEL": "claude-haiku-4-5-20251001",
    "MAX_THINKING_TOKENS": "10000",
    "MAX_MCP_OUTPUT_TOKENS": "10000",
    "DISABLE_NON_ESSENTIAL_MODEL_CALLS": "1"
  }
}
```

O que cada variável faz:

| Variável | Padrão | Impacto |
|----------|--------|---------|
| `CLAUDE_CODE_SUBAGENT_MODEL=haiku` | herda principal | Subagentes ~80% mais baratos |
| `MAX_THINKING_TOKENS=10000` | 31.999 | −70% no custo de thinking |
| `MAX_MCP_OUTPUT_TOKENS=10000` | 25.000 | Limita respostas gigantes de MCP |
| `DISABLE_NON_ESSENTIAL_MODEL_CALLS=1` | desligado | Elimina chamadas de background |
| `ENABLE_TOOL_SEARCH=true` | false | −85% em tool definitions |

> Para Opus 4.7: `MAX_THINKING_TOKENS` não se aplica (usa adaptive reasoning). Usar `/effort low` ou `/effort medium` para controlar custo de raciocínio.

---

## Camada 2 — CLAUDE.md e Skills

### 2-A: CLAUDE.md eficiente

**Regra do criador do Claude Code (Boris Cherny):** CLAUDE.md ideal tem ~2.500 tokens (~1.900 palavras). Se o seu está maior, instruções começam a ser ignoradas.

Teste de cada linha: *"Se eu remover isso, Claude vai errar?"* Se não → remova.

Estrutura ótima para cache máximo:
```markdown
<!-- PARTE ESTÁTICA (nunca muda — maximiza cache hit) -->
## Project
[1-2 frases sobre o projeto]

## Stack
[tecnologias em bullet points — sem parágrafos]

## Conventions
[regras técnicas curtas — inglês preferido aqui]

## Tool preference
Prefer CLI over MCP when task can be done with a shell command.

## Compact policy
When compacting: preserve modified files list, architectural decisions, error solutions.
Omit: failed attempts, exploratory output, file contents already committed.

<!-- PARTE DINÂMICA (muda por task — vai no fim) -->
## Current task
[estado atual — atualizar por sessão]
```

Mover para arquivos separados (zero custo até Claude precisar):
- Documentação detalhada → `docs/ARCHITECTURE.md`
- Bugs históricos → `.claude/COMMON_MISTAKES.md`
- Comandos do projeto → `.claude/QUICK_START.md`

### 2-B: Skills — manter ≤ 12 ativas

Cada skill injeta ~100 tokens de frontmatter em toda sessão. Com 35 skills: 3.500 tokens fixos antes do primeiro prompt.

```bash
# Ver skills instaladas com tamanho
for dir in ~/.claude/skills/*/; do
    name=$(basename "$dir")
    lines=$(wc -l < "$dir/SKILL.md" 2>/dev/null || echo 0)
    printf "%4d linhas  %s\n" "$lines" "$name"
done | sort -rn
```

```powershell
# PowerShell
Get-ChildItem "$env:USERPROFILE\.claude\skills" -Directory | ForEach-Object {
    $skillMd = Join-Path $_.FullName "SKILL.md"
    $lines = if (Test-Path $skillMd) { (Get-Content $skillMd).Count } else { 0 }
    [PSCustomObject]@{ Lines = $lines; Name = $_.Name }
} | Sort-Object Lines -Descending | Format-Table
```

---

## Camada 3 — Gestão de sessão (hábitos diários)

### Comandos essenciais

| Comando | Quando usar | O que faz |
|---------|-------------|-----------|
| `/context` | A qualquer momento | Breakdown exato do uso por categoria |
| `/cost` | Fim de sessão / auditoria | Custo total da sessão atual |
| `/compact [instrução]` | A ~60% de contexto | Comprime preservando o que você especificar |
| `/clear` | Entre tasks não relacionadas | Reset completo (mantém CLAUDE.md) |
| `/btw pergunta` | Dúvidas rápidas | Resposta efêmera — não entra no histórico |
| `Shift+Tab` | Antes de tasks grandes | Ativa Plan mode — vê o plano antes de executar |

### Timing de sessão (usuários Pro)

A janela de uso do Claude Pro dura ~5 horas a partir da primeira mensagem da sessão.

**Estratégia:** Envie um "oi" ou comando simples 2–3 horas antes do seu bloco de trabalho intenso. Quando você estiver no pico de produtividade, a janela reinicia — entregando carga renovada no momento certo.

**Horários de menor limitação (horário de Brasília):**

| Horário | Recomendação | Motivo |
|---------|-------------|--------|
| 6h–9h | ✅ Ideal | Antes da sobreposição com EUA |
| 9h–15h | ❌ Evitar | Pico simultâneo Brasil + EUA |
| Após 20h | ✅ Bom | Fora do horário comercial americano |
| Fins de semana | ✅ Ideal | Demanda global reduzida |

**Hierarquia de decisão para gerenciar contexto:**

```
Contexto < 60%  → continue normalmente
Contexto 60–80% → /compact "preserve [X crítico]"
Contexto > 80%  → /compact imediatamente (antes de degradar)
Task nova, diferente → /clear + handoff note
Mesmo loop de erro (2x) → /clear + prompt melhorado
```

**Compact com instrução sempre que possível:**
```
/compact Preserve: list of modified files, decisions on auth flow, error solutions.
         Omit: all failed attempts, exploratory reads, terminal output.
```

**Handoff note antes de /clear:**
```
Before we clear: write a 200-word handoff note. Include:
what we built, key decisions and why, current state, next steps, gotchas.
```

### Plan mode — elimina trial-and-error

`Shift+Tab` antes de qualquer task com 3+ arquivos. Claude mostra o plano sem executar nada. Você corta o que não precisa. Só então executa. Elimina a maior fonte de desperdício: iterações por erro.

### /btw — perguntas efêmeras

```
/btw What does the calculateMetrics() function return?
/btw Is there a test for the auth flow?
```

A resposta aparece em overlay, vê o contexto completo, mas **não entra no histórico**. Para dúvidas pontuais vale muito mais que fazer uma pergunta normal.

---

## Camada 4 — Subagentes e MCP discipline

### Subagentes para isolamento de contexto

Regra prática: **tarefa abrange 3+ arquivos ou gera output grande → delegar a subagente**.

```
# Exploração de codebase (não polui o contexto principal)
Use a subagent to investigate how authentication handles token refresh.
Report back: files involved, current flow, any existing OAuth utilities.

# Execução de testes (logs ficam isolados)
Use a subagent to run the full test suite.
Report back: only failing tests with their error messages.

# Pesquisa paralela
Use separate subagents to investigate the auth, database, and API modules in parallel.
```

Overhead de subagente só vale a pena para tasks > 3 arquivos. Para 1–2 arquivos, faça inline.

### MCP discipline: CLI first

```
# Em vez de MCP pesado
"Use the AWS MCP to describe instance i-0123456789"

# CLI equivalente (zero overhead de MCP)
"Run: aws ec2 describe-instances --instance-ids i-0123456789"
```

Quando o MCP não tem alternativa CLI, mantém. Quando tem, prefere CLI.

Para ver tokens por MCP na sessão atual, dentro do Claude Code:
```
/mcp
```

### 4-C: Code Review Graph — mapear codebase sem ler cada arquivo

**Problema:** Em projetos grandes, Claude lê a base de código repetidamente, esgotando tokens rapidamente. Um projeto médio de 50+ arquivos pode queimar 60–70% do contexto só em exploração.

**Solução:** Gerar um mapa estruturado e hierárquico do projeto e fornecê-lo ao Claude. Impacto: **−60 a −70% de tokens** em projetos com muitos arquivos.

```bash
# Instalar a ferramenta (requer Node.js)
npm install -g code-review-graph

# Gerar mapa do projeto em texto estruturado
crg --format text > project-map.txt

# Alternativa: formato JSON para projetos grandes
crg --format json > project-map.json
```

Uso no Claude:
```
"Here is the complete project map: [colar conteúdo do project-map.txt]
Now, without reading any files, tell me where the authentication logic lives
and navigate directly to those files only."
```

Em vez de deixar Claude explorar arquivo por arquivo, forneça o mapa e instrua navegação cirúrgica. Use subagentes para leitura — o contexto principal recebe apenas os achados.

---

## Camada 5 — Compressão de output (Caveman)

Instalar:
```bash
# CLI
claude plugin marketplace add JuliusBrussee/caveman
claude plugin install caveman@caveman
# Reiniciar Claude Code
```

No Desktop App: botão `+` → Plugins → Add plugin → buscar "caveman".

Modos disponíveis:
```
/caveman lite   # remove filler, mantém gramática
/caveman full   # fragmentos — recomendado para debug
/caveman ultra  # máxima compressão — para log/terminal pesado
```

Comprimir CLAUDE.md existente (~46% redução de input por sessão):
```
/caveman:compress CLAUDE.md
```

Impacto real: 14–87% de redução em output discursivo. Médias reais em benchmarks independentes: ~14–45%. A maior diferença aparece em sessões de debug onde Claude explica muito. Código gerado não é afetado.

---

## Camada 6 — Otimização de Input

### 6-A: PDFs — nunca suba o arquivo direto

PDFs contêm metadados, layouts, cabeçalhos, rodapés e formatação invisível que consomem **70–90% de tokens desnecessários** antes de qualquer análise real.

**Fluxo recomendado:**

1. Abra o PDF em outra ferramenta gratuita (ChatGPT Free, Google Gemini, Claude.ai web)
2. Use este prompt de extração:
```
Read this document. Remove all repetition, headers, footers, page numbers
and formatting artifacts. Return only the essential points and key data
in plain text, organized logically.
```
3. Cole o texto limpo no Claude Code

**Economia:** 70–90% de redução vs. upload direto do arquivo.

Aplicável também a:
- Páginas web extensas → extrair só o artigo (sem nav, footer, ads)
- Relatórios corporativos → remover boilerplate legal e formatação
- Transcrições longas → resumir antes de enviar

### 6-B: Contexto cirúrgico — só envie o que Claude precisa

Regra: **quanto menos contexto irrelevante, mais precisa a resposta.**

```
# Em vez de colar 500 linhas de log
"Here are the last 20 lines around the error: [apenas o trecho relevante]"

# Em vez de colar todo o arquivo de config
"Here is the relevant config section: [apenas o bloco específico]"

# Em vez de colar o schema inteiro do banco
"Here are the 3 tables involved in this query: [apenas essas tabelas]"
```

---

## Seleção de modelo por complexidade

```bash
# settings.json — modelo padrão sempre Sonnet
{ "model": "claude-sonnet-4-6" }

# Durante sessão — mudar conforme a task
/model sonnet    # default, 80% das tasks
/model opus      # apenas para: arquitetura complexa, refactor profundo, análise crítica
/model haiku     # busca simples, formatação, tarefas mecânicas

# Alias híbrido (plan=Opus, execute=Sonnet)
/model opusplan
```

Distribuição recomendada de uso:

| Modelo | % ideal | Quando usar |
|--------|---------|-------------|
| **Sonnet** | 80% | Codificação diária, análise de dados, relatórios, resumos |
| **Opus** | 15% | Arquitetura complexa, bugs raros, escrita com alta nuance |
| **Haiku** | 5% | Classificações rápidas, automações simples, buscas |

Regra prática: **comece sempre com Sonnet. Só mude para Opus se a resposta não for satisfatória.** Na dúvida, Sonnet resolve.

Custo relativo aproximado: Haiku = 1×, Sonnet = 6×, Opus = 30×.

---

## Referência: impacto validado por camada

| Otimização | Redução típica | Fonte | Onde funciona |
|-----------|---------------|-------|---------------|
| Tool Search (`ENABLE_TOOL_SEARCH`) | −85% tool tokens | Anthropic Engineering | CLI ✓ / Desktop (via settings.json) |
| `.claudeignore` (node_modules+) | −30 a −40% por sessão | Benchmarks community | CLI ✓ / Desktop ✓ |
| Subagente model = Haiku | −80% custo subagentes | Ratio de preços Anthropic | CLI ✓ / Desktop ✓ |
| `MAX_THINKING_TOKENS=10000` | −70% thinking tokens | Docs oficiais | CLI ✓ / Desktop ✓ |
| `/compact` a 60% (vs 95%) | Previne degradação | Cuttlesoft / MindStudio | CLI ✓ / Desktop ✓ |
| CLAUDE.md < 2.500 tokens | Base mais leve toda sessão | Boris Cherny pattern | CLI ✓ / Desktop ✓ |
| Subagentes p/ explorações | Isola output volumoso | Anthropic blog | CLI ✓ / Desktop ✓ |
| Caveman output | −14 a −87% output | Benchmarks independentes | CLI ✓ / Desktop ✓ |
| Plan mode antes de executar | Elimina trial-and-error | 32blog (50% redução total) | CLI ✓ / Desktop ✓ |
| `/btw` para dúvidas pontuais | Resposta sem histórico | Docs oficiais | CLI ✓ / Desktop ✓ |
| Inglês em CLAUDE.md/instruções | −30 a −50% input instruções | arXiv 2305.15425 | CLI ✓ / Desktop ✓ |
| CLI em vez de MCP quando possível | −100% overhead do MCP | Community best practices | CLI ✓ / Desktop ✓ |
| Code Review Graph | −60 a −70% tokens de codebase | Community benchmarks | CLI ✓ / Desktop ✓ |
| PDFs pré-processados | −70 a −90% tokens de input | Community best practices | Todos |
| Contexto cirúrgico (trechos vs. arquivos) | −40 a −60% input por tarefa | Community best practices | Todos |
| Timing de sessão (janela 5h Pro) | Carga renovada no pico | Community tip | Claude Pro |
| Horários de baixa demanda | Menos limitações de servidor | Community tip | Todos |

---

## Quick start — 10 minutos, máximo impacto

Para quem quer aplicar agora sem ler tudo:

```bash
# 1. settings.json — coloca isso em ~/.claude/settings.json
# (Windows: %USERPROFILE%\.claude\settings.json)
{
  "model": "claude-sonnet-4-6",
  "env": {
    "ENABLE_TOOL_SEARCH": "true",
    "CLAUDE_CODE_SUBAGENT_MODEL": "claude-haiku-4-5-20251001",
    "MAX_THINKING_TOKENS": "10000",
    "MAX_MCP_OUTPUT_TOKENS": "10000",
    "DISABLE_NON_ESSENTIAL_MODEL_CALLS": "1"
  }
}

# 2. .claudeignore — cria na raiz do projeto (veja Camada 1-B)

# 3. CLAUDE.md — auditoria rápida
# Perguntar ao Claude: "Audit my CLAUDE.md. For each section, tell me if removing
# it would cause you to make mistakes. List what to cut."

# 4. Hábito de sessão
# - /context quando começar
# - /compact a ~60%
# - /clear entre tasks diferentes
# - Shift+Tab antes de tasks grandes
```

Referência detalhada de cada otimização: ver `references/`
