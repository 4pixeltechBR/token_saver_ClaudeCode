#!/bin/bash
# token-saver-setup.sh
# Setup automático de economia de tokens para Claude Code (bash)
# Uso: bash token-saver-setup.sh [--dry-run]

set -euo pipefail
DRY_RUN=false
[[ "${1:-}" == "--dry-run" ]] && DRY_RUN=true

step() { echo -e "\n\033[36m[+] $1\033[0m"; }
ok()   { echo -e "    \033[32mOK: $1\033[0m"; }
warn() { echo -e "    \033[33mAVISO: $1\033[0m"; }

echo -e "\n\033[35m=== Token Saver Setup — Claude Code ===\033[0m"

CLAUDE_DIR="$HOME/.claude"
SETTINGS="$CLAUDE_DIR/settings.json"

# 1. settings.json
step "Configurando ~/.claude/settings.json"
mkdir -p "$CLAUDE_DIR"

NEW_SETTINGS=$(cat << 'EOF'
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
EOF
)

if [[ -f "$SETTINGS" ]]; then
    warn "settings.json já existe. Verificar e mesclar manualmente se necessário."
    echo "    Conteúdo atual:"
    cat "$SETTINGS" | head -30
    echo ""
    echo "    Mesclar com:"
    echo "$NEW_SETTINGS"
else
    if [[ "$DRY_RUN" == "false" ]]; then
        echo "$NEW_SETTINGS" > "$SETTINGS"
        ok "settings.json criado: $SETTINGS"
    else
        echo "    [DRY RUN] Criaria $SETTINGS com:"
        echo "$NEW_SETTINGS"
    fi
fi

# 2. .claudeignore
step "Criando .claudeignore no projeto atual"
IGNORE_PATH="$(pwd)/.claudeignore"

IGNORE_CONTENT="# Gerado por token-saver — economiza 30-40% de contexto
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
.nyc_output/"

if [[ -f "$IGNORE_PATH" ]]; then
    warn ".claudeignore já existe — não sobrescrevendo."
else
    if [[ "$DRY_RUN" == "false" ]]; then
        echo "$IGNORE_CONTENT" > "$IGNORE_PATH"
        ok ".claudeignore criado"
    else
        echo "    [DRY RUN] Criaria .claudeignore"
    fi
fi

# 3. Diagnóstico MCPs
step "Verificando MCPs configurados"
CLAUDE_JSON="$HOME/.claude.json"
if [[ -f "$CLAUDE_JSON" ]]; then
    COUNT=$(python3 -c "
import json
d = json.load(open('$CLAUDE_JSON'))
mcps = list(d.get('mcpServers', {}).keys())
print(len(mcps))
print(','.join(mcps))
" 2>/dev/null | head -1)
    if [[ -n "$COUNT" && "$COUNT" -gt 5 ]]; then
        warn "$COUNT MCPs configurados. Considerar desabilitar os não usados."
    elif [[ -n "$COUNT" ]]; then
        ok "$COUNT MCPs. OK."
    fi
else
    ok "~/.claude.json não encontrado (sessão limpa)"
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
SKILLS_PATH="$CLAUDE_DIR/skills"
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
