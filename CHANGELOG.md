# Changelog

## 2.0.0 — 2026-09-17

- Skill curta com entrada guiada, diagnostico, detalhes e desfazer.
- CLI Go sem dependencias de runtime; pacotes para seis combinacoes de SO/CPU.
- Configuracoes de projeto em settings.local.json; nenhum ajuste em ~/.claude.json.
- Previa sem escrita, deteccao de plano obsoleto, journal e restauracao seletiva.
- Preservacao de modelos, thinking, permissoes, URLs e inteiros JSON grandes.
- Instalador verifica integridade, protege personalizacoes e preserva versao anterior.
- Modo Direto opcional; diagnostico nao converte linhas em supostos tokens.
- Removidas promessas de economia, .claudeignore, horarios anedoticos e dependencias
  automaticas de Caveman/CRG. Referencias oficiais e limitacoes explicitas.

Nao inclui publicacao no GitHub, assinatura de binarios, validacao humana com cinco
usuarios ou benchmarks de economia. A automacao de CI executa testes nativos por SO.
