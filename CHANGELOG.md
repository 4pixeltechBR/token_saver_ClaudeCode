# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
