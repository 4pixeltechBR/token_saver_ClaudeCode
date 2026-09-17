# Validacao da entrega 2.0.1

Data: 17/09/2026. Base analisada: repositorio 4pixeltechBR/token_saver_ClaudeCode,
commit `088b4775840854ee64da667fdbbc59772b5575cf`.

## Executado neste ambiente

- Compilacao Go 1.27.1 para Windows, macOS e Linux, amd64 e arm64.
- Suite automatizada no Windows amd64: 30 testes/subtestes aprovados, um teste
  de symlink Unix ignorado por falta de privilegio para criar symlinks no Windows.
- Teste adicional nativo de junction do Windows aprovado: escrita recusada e
  destino redirecionado preservado.
- `go vet ./...` aprovado.
- Validacao estrutural de SKILL.md aprovada.
- Sintaxe do instalador PowerShell aprovada pelo parser; instalador shell aprovado
  por `bash -n`.
- Fluxo do pacote Windows em perfil e projeto isolados, inclusive caminhos com
  espacos: instalar, auditar, gerar plano, aplicar Modo Direto e desfazer.
- Integridade dos seis ZIPs e dos checksums dos binarios conferida.

## Comportamentos cobertos

JSON invalido e chaves duplicadas; URLs e numeros grandes; preservacao de modelo,
permissoes e variaveis; ausencia de segredos na saida; preview sem escrita;
configuracao efetiva local; compatibilidade conservadora; idempotencia; plano
obsoleto; lock; conflito na mesma chave; preservacao de edicoes posteriores;
undo exato; recuperacao de journal pendente; caminhos de backup fora do escopo;
instalacao, atualizacao e protecao de skills personalizadas.

## Limites desta validacao

macOS, Linux e Windows ARM foram compilados, mas nao executados nativamente aqui.
O workflow Test executara testes nativos nos tres sistemas quando enviado ao GitHub.
A CLI local do Claude Code respondeu a consulta de versao, mas nao foi iniciada
uma chamada paga de IA para comprovar a invocacao da skill em uma sessao real.

Não foram feitos benchmark de economia nem testes com cinco usuários. Binários
nao estao assinados/notarizados. Nenhuma release foi publicada e nenhuma configuracao
real do usuario foi alterada; instalacao e mutacoes foram testadas em perfis isolados.
