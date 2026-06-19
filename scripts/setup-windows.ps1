# setup-windows.ps1
# Setup automático de economia de tokens para Claude Code no Windows (refatorado)
# Executar: powershell -ExecutionPolicy Bypass -File setup-windows.ps1

param(
    [switch]$DryRun = $false
)

$ErrorActionPreference = "Stop"

function Write-Step($msg) { Write-Host "`n[+] $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    OK: $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "    AVISO: $msg" -ForegroundColor Yellow }

Write-Host "`n=== Token Saver Setup - Claude Code Windows ===" -ForegroundColor Magenta

# Find python interpreter
function Get-PythonPath {
    $paths = @("python3", "python")
    foreach ($cmd in $paths) {
        $found = Get-Command $cmd -ErrorAction SilentlyContinue
        if ($found) {
            return $found.Source
        }
    }
    return $null
}

$pythonCmd = Get-PythonPath
if ($null -eq $pythonCmd) {
    Write-Warn "Python nao foi encontrado no sistema. Por favor, instale o Python para rodar o setup."
    Exit 1
}

$scriptDir = if ($PSScriptRoot) { $PSScriptRoot } else { Get-Location }
$mergeScript = Join-Path $scriptDir "merge_json.py"

if (-not (Test-Path $mergeScript)) {
    Write-Warn "Script de mesclagem nao encontrado em: $mergeScript"
    Exit 1
}

# CORRIGIDO: O arquivo correto é ~/.claude.json, não ~/.claude/settings.json
$settingsPath = "$env:USERPROFILE\.claude.json"

# 1. ~/.claude.json
Write-Step "Configurando ~/.claude.json"

if (-not $DryRun) {
    & $pythonCmd $mergeScript --file $settingsPath --config-type settings --model "claude-sonnet-4-6" --update-env ENABLE_TOOL_SEARCH=true CLAUDE_CODE_SUBAGENT_MODEL=claude-haiku-4-5-20251001 MAX_THINKING_TOKENS=10000 MAX_MCP_OUTPUT_TOKENS=10000 DISABLE_NON_ESSENTIAL_MODEL_CALLS=1
    Write-Ok "Configuracoes aplicadas com sucesso em: $settingsPath"
} else {
    Write-Host "    [DRY RUN] Iria mesclar as seguintes configuracoes em ${settingsPath}:" -ForegroundColor Gray
    & $pythonCmd $mergeScript --file $settingsPath --config-type settings --model "claude-sonnet-4-6" --update-env ENABLE_TOOL_SEARCH=true CLAUDE_CODE_SUBAGENT_MODEL=claude-haiku-4-5-20251001 MAX_THINKING_TOKENS=10000 MAX_MCP_OUTPUT_TOKENS=10000 DISABLE_NON_ESSENTIAL_MODEL_CALLS=1 --dry-run
}

# 2. .claudeignore no diretório atual
Write-Step "Criando .claudeignore no projeto atual"
$ignoreContent = @"
# Gerado por token-saver - economiza de 30-40% de contexto
# Padrões comuns / originais
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

# Extensões e Mídias
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

# Lockfiles comuns
package-lock.json
yarn.lock
pnpm-lock.yaml

# Cache e logs de ferramentas
.cache/
tmp/
temp/
.turbo/
.vercel/
.netlify/
coverage/
.nyc_output/

# --- Cobertura Robusta de Stacks ---
# 1. JS/TS Stack (Adicional)
.yarn/cache/
.yarn/unplugged/
.pnpm-store/
bower_components/
jspm_packages/
yarn-error.log

# 2. Python Stack (Adicional)
.tox/
.nox/
.ipynb_checkpoints/
pip-log.txt
poetry.lock
pipfile.lock

# 3. Java / JVM Stack
target/
.gradle/
.m2/
*.class
*.jar
*.war
*.ear

# 4. Dart / Flutter Stack
.dart_tool/
.flutter-plugins
.flutter-plugins-dependencies
.pub-cache/
.pub/
build/flutter/

# 5. PHP Stack
vendor/
composer.lock
.phpunit.result.cache
"@

$ignorePath = Join-Path (Get-Location) ".claudeignore"
if (Test-Path $ignorePath) {
    Write-Warn ".claudeignore ja existe no diretorio atual - nao sobrescrevendo. Revise se necessario."
} else {
    if (-not $DryRun) {
        $ignoreContent | Set-Content $ignorePath -Encoding UTF8
        Write-Ok ".claudeignore criado no diretorio atual"
    } else {
        Write-Host "    [DRY RUN] Criaria .claudeignore no diretorio atual" -ForegroundColor Gray
    }
}

# 3. Diagnóstico de MCPs
Write-Step "Verificando MCPs configurados"
if (Test-Path $settingsPath) {
    try {
        $countOutput = & $pythonCmd -c @"
import json
try:
    d = json.load(open(r'$settingsPath'))
    print(len(d.get('mcpServers', {}).keys()))
except Exception:
    print('0')
"@
        $count = [int]$countOutput.Trim()
        if ($count -gt 5) {
            Write-Warn "$count MCPs configurados. Risco de overhead alto (cada MCP pesado = 7k-17k tokens/turno)."
            Write-Host "    -> Desabilitar MCPs nao usados em: Claude Code -> Settings -> Connectors" -ForegroundColor Yellow
        } else {
            Write-Ok "$count MCPs configurados. OK."
        }
    } catch {
        Write-Warn "Nao foi possivel parsear ~/.claude.json"
    }
} else {
    Write-Ok "~/.claude.json nao encontrado (sessao limpa)"
}

# 4. Verificar CLAUDE.md
Write-Step "Auditando CLAUDE.md"
if (Test-Path "CLAUDE.md") {
    $words = (Get-Content CLAUDE.md | Measure-Object -Word).Words
    $tokens = [int]($words * 1.3)
    if ($tokens -gt 3000) {
        Write-Warn "CLAUDE.md = ~$tokens tokens. Acima do ideal (2.500). Revisar e compactar."
        Write-Host "    -> Dica: peca ao Claude para auditar cada secao e identificar o que remover." -ForegroundColor Yellow
    } else {
        Write-Ok "CLAUDE.md = ~$tokens tokens. Dentro do limite ideal (<= 2.500)."
    }
} else {
    Write-Ok "CLAUDE.md nao encontrado no diretorio atual."
}

# 5. Skills
Write-Step "Contando skills instaladas"
$skillsPath = "$env:USERPROFILE\.claude\skills"
if (Test-Path $skillsPath) {
    $count = @(Get-ChildItem $skillsPath -Directory -ErrorAction SilentlyContinue).Count
    $fixedTokens = $count * 100
    if ($count -gt 15) {
        Write-Warn "$count skills = ~$fixedTokens tokens fixos por sessao. Considerar desinstalar as menos usadas."
    } else {
        Write-Ok "$count skills = ~$fixedTokens tokens fixos. OK."
    }
} else {
    Write-Ok "Nenhuma skill instalada ainda."
}

Write-Host "`n=== Setup concluido ===" -ForegroundColor Magenta
Write-Host "Proximos passos:" -ForegroundColor White
Write-Host "  1. Reiniciar o Claude Code" -ForegroundColor White
Write-Host "  2. Rodar /context para verificar o uso atual" -ForegroundColor White
Write-Host "  3. Usar /compact quando contexto atingir 60%" -ForegroundColor White
Write-Host "  4. Usar /clear entre tasks nao relacionadas" -ForegroundColor White
Write-Host "  5. Shift+Tab antes de tasks grandes (Plan mode)" -ForegroundColor White
