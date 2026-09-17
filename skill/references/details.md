# Diagnostico e opcoes

## O que e medido

O CLI identifica stack por manifests, bytes e linhas de CLAUDE.md, pastas de skills
com SKILL.md e declaracoes locais de MCP. Nao mede tokens da sessao, ferramentas
carregadas, chamadas, cache, custo real ou economia contrafactual. MCPs repetidos
em escopos diferentes contam como declaracoes, nao como servidores efetivos.
Plugins, politicas remotas e arquivos de instrucoes aninhados nao entram na contagem.

O limiar de 200 linhas gera uma sugestao editorial, nao diagnostico de desperdicio.
Antes de reduzir instrucoes, identifique repeticoes e preserve decisoes, comandos
de verificacao e restricoes. Uma proposta de edicao precisa de revisao especifica.

## Compatibilidade conservadora

O CLI le settings do usuario, compartilhado do projeto e local do projeto nessa
ordem. Respeita CLAUDE_CONFIG_DIR. Nunca escreve ~/.claude.json. Esse diagnostico
nao substitui /status: flags, politicas de servidor, providers e sessoes incorporadas
podem alterar o comportamento efetivo.

Tool Search so e proposto quando ha `ENABLE_TOOL_SEARCH: "false"` explicito nos
settings, Claude Code >= 2.1.221 identificado, modelo Sonnet/Opus/Haiku 4.5+ com
identificador reconhecivel e nenhum override ou provider incompatível detectado.
Aliases de modelo sem versao ficam para revisao. Ausencia da variavel preserva o
padrao nativo. Nao force `true` em um proxy apenas para obter um placar melhor.

O ajuste vai para `.claude/settings.local.json`, apenas no projeto. Se criado
manualmente, esse arquivo pode aparecer como nao rastreado no Git: revise seus
ignores antes de fazer commit. A ferramenta nao altera .gitignore automaticamente.

## Modo Direto

Opcional, instalado em `.claude/rules/token-saver.md`. Corta cortesias repetitivas,
sem ordenar menos investigacao, testes ou codigo incompleto. A propria regra ocupa
contexto; use quando a verbosidade estiver atrapalhando. Nao sobrescreve regra
existente personalizada. E compartilhavel por Git, diferentemente de settings.local.

## Melhorias que precisam de avaliacao especifica

- Modelo menor em trabalho auxiliar: comparar sucesso, tempo e retrabalho.
- Menor esforco de raciocinio: apenas para tarefas apropriadas e modelo compativel.
- Instrucao de compactacao: preservar objetivo, decisoes e estado; usar o recurso
  nativo quando necessario, sem compactar obrigatoriamente em 60%.
- Saidas extensas: filtrar na origem mantendo erros, codigo de saida e acesso ao log
  completo. Cortar `npm test` com `head` pode esconder a falha real.
- Navegacao: buscas direcionadas e recursos de inteligencia de codigo existentes.

Nenhuma dessas opcoes e aplicada automaticamente pelo CLI.

## Comparacao antes/depois

Use tarefas equivalentes, mesmo modelo e criterios de aceitacao. Registre tokens
de entrada, saida, cache read/write, duracao, repeticoes e testes. Compare varias
execucoes, separando cache frio e quente. Considere o custo do proprio diagnostico.
Em assinatura, consumo menor nao reduz a mensalidade automaticamente. Sem uma
linha de base comparavel, relate a mudanca observada, nao "economia comprovada".
