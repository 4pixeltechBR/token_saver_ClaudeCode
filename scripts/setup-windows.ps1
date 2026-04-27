# token-saver-setup.ps1
# Setup automatico de economia de tokens para Claude Code no Windows
# Executar: powershell -ExecutionPolicy Bypass -File token-saver-setup.ps1

param(
    [switch]$DryRun = $false
)

$ErrorActionPreference = "Stop"
$claudeDir = "$env:USERPROFILE\.claude"
$settingsPath = "$claudeDir\settings.json"

function Write-Step($msg) { Write-Host "`n[+] $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    OK: $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "    AVISO: $msg" -ForegroundColor Yellow }

Write-Host "`n=== Token Saver Setup — Claude Code Windows ===" -ForegroundColor Magenta

# 1. settings.json
Write-Step "Configurando ~/.claude/settings.json"
if (-not (Test-Path $claudeDir)) { New-Item -ItemType Directory -Path $claudeDir -Force | Out-Null }

$newEnv = @{
    ENABLE_TOOL_SEARCH               = "true"
    CLAUDE_CODE_SUBAGENT_MODEL       = "claude-haiku-4-5-20251001"
    MAX_THINKING_TOKENS              = "10000"
    MAX_MCP_OUTPUT_TOKENS            = "10000"
    DISABLE_NON_ESSENTIAL_MODEL_CALLS = "1"
}

if (Test-Path $settingsPath) {
    $existing = Get-Content $settingsPath | ConvertFrom-Json
    Write-Warn "settings.json existente encontrado. Mesclando variáveis de ambiente."
    if (-not $existing.env) { $existing | Add-Member -NotePropertyName env -NotePropertyValue @{} }
    foreach ($key in $newEnv.Keys) {
        $existing.env | Add-Member -NotePropertyName $key -NotePropertyValue $newEnv[$key] -Force
    }
    $config = $existing
} else {
    $config = @{
        model = "claude-sonnet-4-6"
        env   = $newEnv
    }
}

if (-not $DryRun) {
    $config | ConvertTo-Json -Depth 5 | Set-Content $settingsPath -Encoding UTF8
    Write-Ok "settings.json atualizado: $settingsPath"
} else {
    Write-Host "    [DRY RUN] Escreveria em: $settingsPath" -ForegroundColor Gray
    $config | ConvertTo-Json -Depth 5
}

# 2. .claudeignore no diretório atual
Write-Step "Criando .claudeignore no projeto atual"
$ignoreContent = @"
# Gerado por token-saver — economiza 30-40% de contexto
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
*.png
*.jpg
*.jpeg
*.gif
*.webp
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
.cache/
tmp/
temp/
.turbo/
.vercel/
.netlify/
coverage/
.nyc_output/
"@

$ignorePath = Join-Path (Get-Location) ".claudeignore"
if (Test-Path $ignorePath) {
    Write-Warn ".claudeignore já existe — não sobrescrevendo. Revisar manualmente."
} else {
    if (-not $DryRun) {
        $ignoreContent | Set-Content $ignorePath -Encoding UTF8
        Write-Ok ".claudeignore criado em: $ignorePath"
    } else {
        Write-Host "    [DRY RUN] Criaria: $ignorePath" -ForegroundColor Gray
    }
}

# 3. Diagnóstico de MCPs
Write-Step "Verificando MCPs configurados"
$claudeJson = "$env:USERPROFILE\.claude.json"
if (Test-Path $claudeJson) {
    try {
        $cfg = Get-Content $claudeJson | ConvertFrom-Json
        $mcps = $cfg.mcpServers.PSObject.Properties.Name
        $count = $mcps.Count
        if ($count -gt 5) {
            Write-Warn "$count MCPs configurados. Risco de overhead alto (cada MCP pesado = 7k-17k tokens/turno)."
            Write-Host "    MCPs: $($mcps -join ', ')" -ForegroundColor Gray
            Write-Host "    → Desabilitar MCPs não usados em: Claude Code → Settings → Connectors" -ForegroundColor Yellow
        } else {
            Write-Ok "$count MCPs configurados. Dentro do limite razoável."
        }
    } catch {
        Write-Warn "Não foi possível parsear ~/.claude.json"
    }
} else {
    Write-Ok "~/.claude.json não encontrado (sessão limpa)"
}

# 4. Verificar CLAUDE.md
Write-Step "Auditando CLAUDE.md"
if (Test-Path "CLAUDE.md") {
    $words = (Get-Content CLAUDE.md | Measure-Object -Word).Words
    $tokens = [int]($words * 1.3)
    if ($tokens -gt 3000) {
        Write-Warn "CLAUDE.md = ~$tokens tokens. Acima do ideal (2.500). Revisar e compactar."
        Write-Host "    → Dica: peça ao Claude para auditar cada seção e identificar o que remover." -ForegroundColor Yellow
    } else {
        Write-Ok "CLAUDE.md = ~$tokens tokens. Dentro do limite ideal (<= 2.500)."
    }
} else {
    Write-Ok "CLAUDE.md não encontrado no diretório atual."
}

# 5. Skills
Write-Step "Contando skills instaladas"
$skillsPath = "$env:USERPROFILE\.claude\skills"
if (Test-Path $skillsPath) {
    $count = (Get-ChildItem $skillsPath -Directory -ErrorAction SilentlyContinue).Count
    $fixedTokens = $count * 100
    if ($count -gt 15) {
        Write-Warn "$count skills = ~$fixedTokens tokens fixos por sessão. Considerar desinstalar as menos usadas."
    } else {
        Write-Ok "$count skills = ~$fixedTokens tokens fixos. OK."
    }
} else {
    Write-Ok "Nenhuma skill instalada ainda."
}

Write-Host "`n=== Setup concluido ===" -ForegroundColor Magenta
Write-Host "Proximos passos:" -ForegroundColor White
Write-Host "  1. Reiniciar o Claude Code Desktop App" -ForegroundColor White
Write-Host "  2. Rodar /context para verificar o uso atual" -ForegroundColor White
Write-Host "  3. Usar /compact quando contexto atingir 60%" -ForegroundColor White
Write-Host "  4. Usar /clear entre tasks nao relacionadas" -ForegroundColor White
Write-Host "  5. Shift+Tab antes de tasks grandes (Plan mode)" -ForegroundColor White
