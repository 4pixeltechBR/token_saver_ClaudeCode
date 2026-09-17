# Recuperacao

`undo --dry-run` mostra o que seria restaurado sem criar backup, pasta ou lock.
`undo --yes` desfaz a ultima aplicacao pendente ou concluida do projeto.

Os backups ficam em `CLAUDE_CONFIG_DIR/token-saver/state/<id-do-projeto>/`.
Sem CLAUDE_CONFIG_DIR, a base e `~/.claude`. Cada JSON contem o estado anterior
dos arquivos; pode conter segredos. Nao envie backups a issues ou ao modelo.
No Windows o acesso segue as permissoes herdadas da pasta pessoal; em Unix os
novos arquivos e pastas usam permissoes restritas.

## Se uma operacao for interrompida

1. Confirme que nenhum processo token-saver ainda esta rodando.
2. Se houver `operation.lock` obsoleto, remova apenas esse arquivo na pasta do
   projeto correspondente em `state`. Nao remova journals nem backups.
3. Rode `undo --dry-run`, depois `undo --yes`.

Uma aplicacao pendente bloqueia novas aplicacoes ate ser recuperada. Undo pode
ser retomado se a restauracao tambem for interrompida. Nao ha reparo de JSON por
regex nem descarte silencioso de chaves.

## Conflitos

Edicoes posteriores em outras chaves de settings sao preservadas. Se o usuario
alterou a mesma chave ENABLE_TOOL_SEARCH, ou editou a regra do Modo Direto, undo
recusa o conflito antes da primeira restauracao. O usuario deve revisar a chave
ou regra com seu backup; nao restaure o arquivo inteiro sobre novas edicoes.

## Instalacao

Reinstalar reconhece apenas instalacoes registradas e intactas. Skills antigas ou
personalizadas sao preservadas: renomeie a pasta depois de revisa-la, ou instale
com `--config-dir` em outro local para experimentar. A versao anterior fica em
`token-saver/install-backups`. Atualize a partir do pacote baixado, nao do proprio
executavel instalado (o Windows pode bloquear executaveis em uso).

Para desinstalar, primeiro desfaça as alteracoes nos projetos em que aplicou.
Depois remova somente a pasta `skills/token-saver`. Preserve os backups ate
confirmar a recuperacao. Remover a skill sozinho nao desfaz ajustes de projeto.
