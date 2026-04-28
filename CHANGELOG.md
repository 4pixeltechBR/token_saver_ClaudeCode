# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-04-28

### Added
- **Camada 6 — Input Optimization**: Nova camada com foco em otimização de input
  - 6-A: PDFs pré-processados (70–90% de economia vs. upload direto)
  - 6-B: Contexto cirúrgico (enviar apenas o trecho relevante)
- **Camada 4-C: Code Review Graph**: Mapear codebase com `crg` para navegação cirúrgica (−60 a −70% em projetos grandes)
- **Timing de sessão (Pro)**: Estratégia da janela de 5 horas para maximizar carga disponível
- **Tabela de horários**: Melhores e piores horários de uso no fuso horário de Brasília
- **Percentuais de uso por modelo**: Distribuição recomendada 80% Sonnet / 15% Opus / 5% Haiku
- **Tabela de referência**: 4 novos itens adicionados com impactos validados

### Changed
- Seção "Seleção de modelo" expandida com tabela de distribuição ideal de uso

---

## [1.0.0] - 2026-04-27

### Added
- Initial release of token-saver skill
- 5-layer token optimization framework
- Automatic setup scripts for Windows (PowerShell), Linux, macOS, and WSL
- AUDIT mode for quick token consumption diagnosis
- Comprehensive documentation in SKILL.md
- Integration with Claude Code settings.json
- .claudeignore generation for context optimization
- MCP configuration validation
- CLAUDE.md audit and optimization
- Support for MAX_THINKING_TOKENS and other optimization variables

### Features
- Layer 1: Setup único (Tool Search, .claudeignore, settings.json)
- Layer 2: CLAUDE.md e Skills management
- Layer 3: Session management commands (/context, /compact, /clear, /btw)
- Layer 4: Subagent e MCP discipline
- Layer 5: Output compression (Caveman mode)

### Documentation
- Complete SKILL.md guide (504 lines)
- Validated sources.md with Anthropic official references
- Setup scripts with automatic configuration
- Impact benchmarks for each optimization layer

---

For more information, see [SKILL.md](SKILL.md) and [sources.md](sources.md).
