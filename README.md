# token-saver: Otimização Profissional de Tokens para Claude Code

<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-blue.svg)](CHANGELOG.md)
[![Claude Code Skill](https://img.shields.io/badge/Claude_Code-Skill-purple.svg)](https://claude.com/claude-code)

**English version below** ⬇️

</div>

---

## 🇧🇷 Português

### O que é token-saver?

**token-saver** é uma ferramenta profissional de diagnóstico e otimização de consumo de tokens no Claude Code. Reduz o consumo de tokens em até **85%** através de 5 camadas de otimizações validadas.

### ⚡ Quick Start (3 passos)

1. **Instale o skill:**
   - Copie `token-saver.skill` para `~/.claude/skills/`
   - Reinicie o Claude Code

2. **Execute o diagnóstico:**
   ```bash
   /token-saver AUDIT
   ```

3. **Aplique as otimizações:**
   ```bash
   /token-saver SETUP
   ```

### 🎯 5 Camadas de Otimização

| Camada | Otimização | Redução | Impacto |
|--------|-----------|---------|---------|
| **1** | Setup único (Tool Search, .claudeignore, settings.json) | −85% | 🟢 Alto |
| **2** | CLAUDE.md e Skills management | −30-40% | 🟢 Alto |
| **3** | Session commands (/compact, /context, /clear) | −60% | 🟡 Médio |
| **4** | Subagent e MCP discipline | −80% | 🟢 Alto |
| **5** | Output compression (Caveman mode) | −40% | 🟡 Médio |

### 📚 Documentação Completa

- **[SKILL.md](SKILL.md)** — Guia completo com todas as otimizações
- **[sources.md](sources.md)** — Referências validadas e benchmarks
- **[CHANGELOG.md](CHANGELOG.md)** — Histórico de versões

### 🛠 Scripts de Setup Automático

- **Linux/macOS/WSL:** `scripts/setup.sh`
- **Windows (PowerShell):** `scripts/setup-windows.ps1`

Execute automaticamente para:
- Configurar `~/.claude/settings.json`
- Criar `.claudeignore` no projeto
- Diagnosticar MCPs instalados
- Auditar CLAUDE.md

### 🔍 Modo AUDIT vs SETUP

**AUDIT** (Diagnóstico)
```bash
/token-saver AUDIT
```
- Relatório rápido do consumo atual
- Identifica oportunidades de otimização
- Sem fazer mudanças

**SETUP** (Aplicar Otimizações)
```bash
/token-saver SETUP
```
- Aplica as 5 camadas automaticamente
- Configura settings.json
- Cria .claudeignore
- Pronto para uso imediato

### 📊 Impactos Validados

- **Tool Search:** −85% de tokens de definições de ferramentas
- **.claudeignore:** −30 a −40% por sessão
- **Subagent Haiku:** −80% custo de subagentes
- **MAX_THINKING=10000:** −70% thinking tokens
- **/compact a 60%:** Previne degradação de contexto

### 🐛 Bugs Conhecidos (Abril 2026)

| Bug | Workaround |
|-----|-----------|
| ENABLE_TOOL_SEARCH não funciona no Desktop via env var | Usar settings.json |
| Auto mode (10% threshold) não dispara | Setar `true` explicitamente |
| defer_loading + cache_control incompatíveis | Desabilitar ENABLE_TOOL_SEARCH se necessário |

### 🤝 Contribua

Encontrou um problema ou tem uma sugestão? Abra uma [Issue](https://github.com/4pixeltechBR/token_saver_ClaudeCode/issues) ou uma [Discussion](https://github.com/4pixeltechBR/token_saver_ClaudeCode/discussions).

### 📄 Licença

MIT License — veja [LICENSE](LICENSE) para detalhes.

---

## 🇬🇧 English

### What is token-saver?

**token-saver** is a professional tool for diagnosing and optimizing token consumption in Claude Code. Reduces token usage by up to **85%** through 5 validated optimization layers.

### ⚡ Quick Start (3 steps)

1. **Install the skill:**
   - Copy `token-saver.skill` to `~/.claude/skills/`
   - Restart Claude Code

2. **Run diagnosis:**
   ```bash
   /token-saver AUDIT
   ```

3. **Apply optimizations:**
   ```bash
   /token-saver SETUP
   ```

### 🎯 5 Optimization Layers

| Layer | Optimization | Reduction | Impact |
|-------|--------------|-----------|---------|
| **1** | Single setup (Tool Search, .claudeignore, settings.json) | −85% | 🟢 High |
| **2** | CLAUDE.md and Skills management | −30-40% | 🟢 High |
| **3** | Session commands (/compact, /context, /clear) | −60% | 🟡 Medium |
| **4** | Subagent and MCP discipline | −80% | 🟢 High |
| **5** | Output compression (Caveman mode) | −40% | 🟡 Medium |

### 📚 Full Documentation

- **[SKILL.md](SKILL.md)** — Complete guide with all optimizations
- **[sources.md](sources.md)** — Validated references and benchmarks
- **[CHANGELOG.md](CHANGELOG.md)** — Version history

### 🛠 Automatic Setup Scripts

- **Linux/macOS/WSL:** `scripts/setup.sh`
- **Windows (PowerShell):** `scripts/setup-windows.ps1`

Execute automatically to:
- Configure `~/.claude/settings.json`
- Create `.claudeignore` in your project
- Diagnose installed MCPs
- Audit CLAUDE.md

### 🔍 AUDIT vs SETUP Mode

**AUDIT** (Diagnosis)
```bash
/token-saver AUDIT
```
- Quick report of current consumption
- Identifies optimization opportunities
- No changes made

**SETUP** (Apply Optimizations)
```bash
/token-saver SETUP
```
- Applies all 5 layers automatically
- Configures settings.json
- Creates .claudeignore
- Ready for immediate use

### 📊 Validated Impacts

- **Tool Search:** −85% tool definition tokens
- **.claudeignore:** −30 to −40% per session
- **Subagent Haiku:** −80% subagent cost
- **MAX_THINKING=10000:** −70% thinking tokens
- **/compact to 60%:** Prevents context degradation

### 🐛 Known Bugs (April 2026)

| Bug | Workaround |
|-----|-----------|
| ENABLE_TOOL_SEARCH doesn't work on Desktop via env var | Use settings.json |
| Auto mode (10% threshold) doesn't trigger | Set `true` explicitly |
| defer_loading + cache_control incompatible | Disable ENABLE_TOOL_SEARCH if needed |

### 🤝 Contribute

Found an issue or have a suggestion? Open an [Issue](https://github.com/4pixeltechBR/token_saver_ClaudeCode/issues) or [Discussion](https://github.com/4pixeltechBR/token_saver_ClaudeCode/discussions).

### 📄 License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

Made with ❤️ for Claude Code developers

[⬆ back to top](#)

</div>
