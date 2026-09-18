# Token Saver 3.0

Auditoria segura de contexto para Claude Code, Codex, Antigravity, OpenCode,
Cursor, Gemini CLI, GitHub Copilot, Cline e MiniMax Code.

[Baixar a release](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest) · [Ver o código](https://github.com/4pixeltechBR/token_saver_ClaudeCode) · [Matriz de harnesses](skill/references/harnesses.md)

## O que mudou

Esta versão separa a skill portátil do comportamento específico de cada harness.
O mesmo pacote pode ser descoberto por vários agentes, mas só altera uma
configuração quando conhece o formato, a versão e o rollback daquele ambiente.

O núcleo é conservador: diagnostica primeiro, mostra evidências, cria uma prévia
com identificador e nunca transforma tamanho de arquivo em uma falsa economia de
tokens.

## Capacidades

| Função | Claude Code | Outros harnesses |
|---|---:|---:|
| Detectar ambientes disponíveis | Sim | Sim |
| Auditar projeto e instruções | Sim | Sim |
| Detectar stack | Sim | Sim |
| Contar skills detectáveis | Sim | Sim |
| Saída JSON para automação | Sim | Sim |
| Prévia (`plan`) | Sim | Sim, sem escrita |
| Aplicar configurações | Sim, quando elegível | Ainda não |
| Desfazer aplicação | Sim | Ainda não |
| Modo Direto | Sim, opt-in | Orientação manual |
| Telemetria | Nunca | Nunca |

### O que o Claude Code pode aplicar

- Reativar Tool Search somente quando existe uma desativação explícita e a
  compatibilidade é confirmada.
- Adicionar opcionalmente uma regra curta de Modo Direto em `.claude/rules`.
- Fazer backup, escrita atômica, validação e rollback com proteção contra
  conflitos e links simbólicos.

### O que todos os harnesses podem auditar

- Arquivos de instrução longos e possíveis repetições.
- Stack detectada por manifests do projeto.
- Skills encontradas nos diretórios suportados pelo host.
- Estado de descoberta da própria skill.
- Limites e próximos passos sem inventar consumo ou economia.

O Token Saver não troca modelos, não reduz thinking, não edita `CLAUDE.md`, não
cria `.claudeignore`, não compacta conversas automaticamente e não instala
plugins externos.

## Instalar

Baixe o ZIP da [release mais recente](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest)
para seu sistema, extraia e execute o binário. O pacote não exige Python ou Node.

```text
token-saver detect --project .
token-saver install --harness claude --project .
token-saver install --harness codex --project .
```

Use `--harness auto` para escolher o primeiro ambiente detectado. Para uma
instalação em um projeto específico, prefira uma pasta de skill versionada pelo
projeto; para uso global, deixe o instalador usar o diretório pessoal.

Depois de instalar, abra uma nova sessão ou recarregue as skills do harness.
As formas usuais são:

| Ambiente | Instalação global | Como chamar |
|---|---|---|
| Claude Code | `~/.claude/skills/token-saver` | `/token-saver` |
| Codex | `~/.agents/skills/token-saver` | `$token-saver` |
| Antigravity | `~/.gemini/antigravity/skills/token-saver` | peça uma auditoria |
| OpenCode | `~/.config/opencode/skills/token-saver` | use a skill `token-saver` |
| Cursor | `~/.cursor/skills/token-saver` | `/token-saver` |
| Gemini CLI | `~/.gemini/skills/token-saver` | `/skills` e ative |
| GitHub Copilot | `~/.copilot/skills/token-saver` | `/token-saver` |
| Cline | `~/.cline/skills/token-saver` | use a skill `token-saver` |
| MiniMax Code | plugin `skills/token-saver` | ative o plugin |

## Usar

No Claude Code, o fluxo principal é:

```text
/token-saver
/token-saver auditar
/token-saver desfazer
/token-saver detalhes
/token-saver modo direto
```

No terminal, todos os hosts podem usar o CLI:

```text
token-saver detect --json
token-saver audit --harness codex --project . --json
token-saver audit --harness opencode --project .
token-saver plan --harness claude --project . --json
token-saver apply --harness claude --project . --yes --expect ID_DA_PREVIA --json
token-saver undo --harness claude --project . --json
token-saver details
token-saver version
```

Opções comuns:

```text
--project PASTA       raiz do projeto analisado
--harness ID          auto, claude, codex, antigravity, opencode, cursor,
                      gemini, copilot, cline ou minimax
--config-dir PASTA    perfil do harness; útil para testes isolados
--json                saída estruturada, sem segredos
--dry-run             somente prévia
--yes                 autoriza uma prévia já revisada
--expect ID           recusa aplicar se o ambiente mudou
--concise             inclui o Modo Direto no plano do Claude Code
```

## Honestidade sobre tokens

O diagnóstico lê arquivos e configurações locais. Ele não mede tokens de entrada,
saída, cache, custo real ou a economia contrafactual. Uma alteração gravada não
prova ganho na sessão. Verifique o painel, comando de contexto ou relatório de
uso oferecido pelo seu próprio harness.

## Segurança e recuperação

O modo padrão é leitura. JSON inválido, configuração gerenciada, versão não
confirmada ou conflito interrompem a aplicação. No Claude Code, o histórico fica
fora do projeto, as escritas são atômicas e o rollback preserva edições posteriores
em chaves não relacionadas. Consulte [recovery](skill/references/recovery.md) e
[details](skill/references/details.md) antes de remover qualquer backup pendente.

## Desenvolver

```text
go test ./...
go vet ./...
python scripts/package.py --out dist
```

O repositório mantém um núcleo em Go, a skill portável em `skill/` e adaptadores
de descoberta por harness em `harness.go` e `portable_audit.go`. A matriz completa
fica em [skill/references/harnesses.md](skill/references/harnesses.md).

## Licença

MIT. Veja [CHANGELOG](CHANGELOG.md), [VALIDATION](VALIDATION.md) e a versão em
[English](README.en.md).
