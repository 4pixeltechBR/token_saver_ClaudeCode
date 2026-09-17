# Token Saver 2.0

**Encontre desperdicios no Claude Code, entenda as mudancas e desfaça quando precisar.**

Sem trocar seu modelo. Sem exigir Python ou Node. Sem promessa de economia inventada.

## Comece em tres passos

1. Baixe o pacote da versao publicada para seu sistema e extraia o ZIP.
2. **Windows:** abra `instalar.cmd`. **macOS/Linux:** rode `sh install.sh` na pasta extraida.
3. Abra uma nova sessao do Claude Code e digite **`/token-saver`**.

O pacote local entregue ja pode ser usado. O download publico so estara disponivel
depois que o mantenedor publicar a release; nao ha comando remoto ficticio neste guia.

| Seu computador | Pacote |
|---|---|
| Windows Intel/AMD | `token-saver-2.0.0-windows-amd64.zip` |
| Windows ARM | `token-saver-2.0.0-windows-arm64.zip` |
| Mac Apple Silicon | `token-saver-2.0.0-darwin-arm64.zip` |
| Mac Intel | `token-saver-2.0.0-darwin-amd64.zip` |
| Linux Intel/AMD | `token-saver-2.0.0-linux-amd64.zip` |
| Linux ARM64 | `token-saver-2.0.0-linux-arm64.zip` |

## Como funciona

`/token-saver` verifica o projeto, apresenta ate tres achados relevantes e mostra
o que pode melhorar. Quando voce autoriza, aplica e valida as mudancas.
Se tudo ja estiver adequado, informa isso e encerra.

| Voce quer | Comando |
|---|---|
| Comecar | `/token-saver` |
| Somente verificar | `/token-saver auditar` |
| Voltar a ultima aplicacao | `/token-saver desfazer` |
| Entender as recomendacoes | `/token-saver detalhes` |
| Respostas mais objetivas, opcional | `/token-saver modo direto` |

**O que muda no basico:** reativacao local de Tool Search somente quando ha uma
desativacao explicita e compatibilidade confirmavel. Nas versoes atuais, muitas
instalacoes ja usam o padrao adequado: nenhuma mudanca sera necessaria.

**O que permanece:** seu modelo, esforco de raciocinio, permissoes e CLAUDE.md.
Modo Direto e opcional e adiciona uma regra curta ao projeto. Nao fazemos bloqueio
de arquivos, compactacao invisivel nem alteracoes automaticas em instrucoes.

## Resultado honesto

Mostramos configuracoes verificadas, tamanhos de instrucoes e mudancas aplicadas.
O consumo real da sessao deve ser verificado em `/context` e `/usage`.
Economia financeira depende da forma de cobranca; uma assinatura nao fica mais
barata automaticamente. [Limites e opcoes](skill/references/details.md).

## Se algo der errado

Use `/token-saver desfazer`. Outras edicoes posteriores sao preservadas; conflitos
na mesma configuracao sao reportados sem sobrescrever o usuario.
JSON invalido e interrompido, sem tentativa de "conserto" por regex.
[Recuperacao e desinstalacao](skill/references/recovery.md).

Uma skill token-saver antiga ou personalizada nao sera sobrescrita. Preserve ou
renomeie a pasta antiga antes de instalar esta versao. O executavel pode ser usado
diretamente do pacote para auditar antes de migrar.

## Para desenvolvedores

CLI com `audit`, `plan`, `apply`, `undo`, `install`, `details`; `--json` para automacao,
`--dry-run` para previa e `--expect` para aplicar apenas o plano revisado.
AUDIT/SETUP continuam reconhecidos. Todos os ajustes sao locais; nenhuma telemetria.

[Arquitetura, testes e publicacao](CONTRIBUTING.md) · [English](README.en.md) · [Fontes](skill/references/sources.md)

Licenca MIT. Evolucao do trabalho de Victor Machado Mendonça / 4pixeltechBR.
