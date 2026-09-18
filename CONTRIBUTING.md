# Manutencao

## Construir e testar

Requer Go 1.24+ apenas para desenvolver. O usuario final recebe binarios nativos.

```sh
go test ./...
go vet ./...
go build -trimpath -buildvcs=false -o token-saver .
```

Os testes usam projetos e perfis temporarios. Nunca execute os instaladores de
setup antigos: eles foram substituidos pelo comando `install`, que nao otimiza
configuracoes durante a instalacao.

`python scripts/package.py --out dist` compila Windows/macOS/Linux, amd64 e arm64,
e cria ZIPs com instaladores e checksums, usando somente a biblioteca padrao Python.
Python e Go sao dependencias de desenvolvimento, nao de uso da skill.

## Arquitetura

- `skill/`: instrucoes e referencias embutidas no executavel.
- `audit.go`: leitura de configuracoes e evidencias locais, sem gravacao.
- `plan.go`: lista fechada de mudancas permitidas e identificador de previa.
- `transaction.go`: journal antes da primeira escrita, validacao e undo seletivo.
- `install.go`: instalacao da skill com manifesto de propriedade e backup de update.
- `harness.go`: catalogo, deteccao e destinos dos harnesses suportados.
- `portable_audit.go`: auditoria somente leitura para hosts sem adaptador de escrita.
- `files.go`: parser JSON estrito, limite de tamanho e gravacao temporaria + rename.

O CLI nao usa chamadas de IA, rede ou telemetria. O adaptador Claude consulta
`claude --version` com timeout durante o diagnostico. Nao presume que os arquivos representem flags
de inicio ou politicas remotas. A verificacao em uma sessao real continua necessaria.

Os locks coordenam instancias desta ferramenta. Outro editor ou o proprio Claude
pode escrever simultaneamente: hashes sao conferidos antes de cada escrita, mas
nao existe transacao distribuida com esses processos. Evite aplicar durante uma
gravacao concorrente. Renomeacao e durabilidade dependem do sistema de arquivos;
prefira disco local. Interrupcoes preservam journal para recuperacao.

O parser rejeita comentarios, virgulas finais, chaves duplicadas e raizes nao objeto.
Aceita BOM UTF-8. Nunca imprima valores de configuracao ou backups em logs de CI.

## Publicar

1. Execute o workflow `Test` e resolva qualquer falha em Windows, macOS e Linux.
2. Valide instalacao e reconhecimento em cada harness listado. Valide aplicacao e
   undo em Claude Code real. Teste iniciantes e seniores; nao substitua isso por
   contagem de testes.
3. Atualize `version` em main.go e `VERSION` em scripts/package.py juntos.
4. Gere os pacotes; confira SHA256SUMS.txt e publique uma release com a tag v3.0.0.
5. O workflow `Packages` gera artefatos para download, mas nao publica release.

Checksums detectam corrupcao; nao equivalem a assinatura de identidade. Binarios
nao estao assinados nem notarizados. Windows/macOS podem exibir verificacoes do
sistema. Assinatura de codigo e notarizacao sao etapas de distribuicao do mantenedor.

O pacote de Agent Skill/Plugin e um formato comunitario; nao e um plugin oficial
de Anthropic, OpenAI, Google, Cursor, GitHub, Cline ou MiniMax. O companion CLI
continua sendo distribuido nos ZIPs de plataforma.

## Aceitacao

Testar: JSON invalido intacto; URLs e chaves preservadas; preview sem escrita;
reexecucao sem duplicacao; modelo inalterado; undo exato ou seletivo; conflito
sem sobrescrita; recuperacao pendente; instalacao sem Python/Node; paths com espacos.
Benchmarks de economia devem usar tarefas equivalentes, repeticoes, criterios de
qualidade e contabilizacao de cache e retrabalho. Nao publique percentuais sem isso.
