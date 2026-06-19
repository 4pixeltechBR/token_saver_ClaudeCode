#!/bin/bash
# setup.sh
# Setup automático de economia de tokens para Claude Code (bash refatorado)
# Uso: bash setup.sh [--dry-run]

set -euo pipefail
DRY_RUN=false
[[ "${1:-}" == "--dry-run" ]] && DRY_RUN=true

step() { echo -e "\n\033[36m[+] $1\033[0m"; }
ok()   { echo -e "    \033[32mOK: $1\033[0m"; }
warn() { echo -e "    \033[33mAVISO: $1\033[0m"; }

echo -e "\n\033[35m=== Token Saver Setup — Claude Code ===\033[0m"

# Find python interpreter
get_python() {
    if command -v python3 &>/dev/null; then
        echo "python3"
    elif command -v python &>/dev/null; then
        echo "python"
    else
        echo ""
    fi
}

PYTHON_CMD=$(get_python)
if [[ -z "$PYTHON_CMD" ]]; then
    warn "Python não foi encontrado no sistema. Por favor, instale o Python para rodar o setup."
    exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MERGE_SCRIPT="$SCRIPT_DIR/merge_json.py"

if [[ ! -f "$MERGE_SCRIPT" ]]; then
    warn "Script de mesclagem não encontrado em: $MERGE_SCRIPT"
    exit 1
fi

# CORRIGIDO: O arquivo correto é ~/.claude.json, não ~/.claude/settings.json
SETTINGS="$HOME/.claude.json"

# 1. ~/.claude.json
step "Configurando ~/.claude.json"

if [[ "$DRY_RUN" == "false" ]]; then
    "$PYTHON_CMD" "$MERGE_SCRIPT" \
        --file "$SETTINGS" \
        --config-type settings \
        --model "claude-sonnet-4-6" \
        --update-env \
            ENABLE_TOOL_SEARCH=true \
            CLAUDE_CODE_SUBAGENT_MODEL=claude-haiku-4-5-20251001 \
            MAX_THINKING_TOKENS=10000 \
            MAX_MCP_OUTPUT_TOKENS=10000 \
            DISABLE_NON_ESSENTIAL_MODEL_CALLS=1
    ok "Configurações aplicadas com sucesso em: $SETTINGS"
else
    echo "    [DRY RUN] Iria mesclar as seguintes configurações em $SETTINGS:"
    "$PYTHON_CMD" "$MERGE_SCRIPT" \
        --file "$SETTINGS" \
        --config-type settings \
        --model "claude-sonnet-4-6" \
        --update-env \
            ENABLE_TOOL_SEARCH=true \
            CLAUDE_CODE_SUBAGENT_MODEL=claude-haiku-4-5-20251001 \
            MAX_THINKING_TOKENS=10000 \
            MAX_MCP_OUTPUT_TOKENS=10000 \
            DISABLE_NON_ESSENTIAL_MODEL_CALLS=1 \
        --dry-run
fi

# 2. .claudeignore
step "Criando .claudeignore no projeto atual"
IGNORE_PATH="$(pwd)/.claudeignore"

IGNORE_CONTENT="# Gerado por token-saver — economiza de 30-40% de contexto
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
.phpunit.result.cache"

if [[ -f "$IGNORE_PATH" ]]; then
    warn ".claudeignore já existe no diretório atual — não sobrescrevendo. Revise se necessário."
else
    if [[ "$DRY_RUN" == "false" ]]; then
        echo "$IGNORE_CONTENT" > "$IGNORE_PATH"
        ok ".claudeignore criado no diretório atual"
    else
        echo "    [DRY RUN] Criaria .claudeignore no diretório atual"
    fi
fi

# 3. Diagnóstico MCPs
step "Verificando MCPs configurados"
if [[ -f "$SETTINGS" ]]; then
    COUNT=$("$PYTHON_CMD" -c "
import json
try:
    d = json.load(open('$SETTINGS'))
    mcps = list(d.get('mcpServers', {}).keys())
    print(len(mcps))
except Exception:
    print('0')
" 2>/dev/null || echo "0")
    
    if [[ "$COUNT" -gt 5 ]]; then
        warn "$COUNT MCPs configurados. Risco de overhead alto (cada MCP pesado = 7k-17k tokens/turno)."
    else
        ok "$COUNT MCPs configurados. OK."
    fi
else
    ok "~/.claude.json não encontrado"
fi

# 4. CLAUDE.md
step "Auditando CLAUDE.md"
if [[ -f "CLAUDE.md" ]]; then
    WORDS=$(wc -w < CLAUDE.md)
    TOKENS=$(( WORDS * 13 / 10 ))
    if [[ $TOKENS -gt 3000 ]]; then
        warn "CLAUDE.md = ~$TOKENS tokens. Acima do ideal (2.500)."
    else
        ok "CLAUDE.md = ~$TOKENS tokens. OK."
    fi
else
    ok "CLAUDE.md não encontrado no diretório atual."
fi

# 5. Skills
step "Contando skills instaladas"
SKILLS_PATH="$HOME/.claude/skills"
if [[ -d "$SKILLS_PATH" ]]; then
    COUNT=$(ls "$SKILLS_PATH" | wc -l)
    TOKENS=$(( COUNT * 100 ))
    if [[ $COUNT -gt 15 ]]; then
        warn "$COUNT skills = ~$TOKENS tokens fixos/sessão. Considerar limpeza."
    else
        ok "$COUNT skills = ~$TOKENS tokens fixos. OK."
    fi
else
    ok "Nenhuma skill instalada ainda."
fi

echo -e "\n\033[35m=== Setup concluído ===\033[0m"
echo "Próximos passos:"
echo "  1. Reiniciar Claude Code"
echo "  2. Rodar /context para verificar uso atual"
echo "  3. Usar /compact quando contexto atingir 60%"
echo "  4. Usar /clear entre tasks não relacionadas"
echo "  5. Shift+Tab antes de tasks grandes (Plan mode)"
