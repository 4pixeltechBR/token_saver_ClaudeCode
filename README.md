# Token Saver

Uma forma simples e reversível de encontrar desperdício de contexto no Claude Code.
O Token Saver preserva seu modelo, seu nível de raciocínio e suas instruções.

[Baixar a versão mais recente](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest) · [Ver o código](https://github.com/4pixeltechBR/token_saver_ClaudeCode)

## Instalar

1. Abra a [release mais recente](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest) e baixe o ZIP correspondente ao seu computador:

   - [Windows Intel/AMD](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-windows-amd64.zip)
   - [Windows ARM](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-windows-arm64.zip)
   - [macOS Apple Silicon](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-darwin-arm64.zip)
   - [macOS Intel](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-darwin-amd64.zip)
   - [Linux Intel/AMD](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-linux-amd64.zip)
   - [Linux ARM64](https://github.com/4pixeltechBR/token_saver_ClaudeCode/releases/latest/download/token-saver-2.0.1-linux-arm64.zip)

2. Extraia o ZIP.
3. No Windows, abra `instalar.cmd`. No macOS/Linux, execute `sh install.sh` na pasta extraída.
4. Abra uma nova sessão do Claude Code e digite `/token-saver`.

O pacote contém o executável e a skill. Não é necessário instalar Python ou Node.
Confira `SHA256SUMS.txt` se quiser validar o download antes de instalar.

## Usar

Na primeira execução, a ferramenta lê o projeto, mostra até três achados e propõe
apenas mudanças compatíveis e verificáveis. Você revisa a prévia antes de aplicar.
Se não houver ganho confirmado, ela informa que nada precisa ser alterado.

| O que você quer | Comando |
|---|---|
| Começar | `/token-saver` |
| Apenas diagnosticar | `/token-saver auditar` |
| Desfazer a última aplicação | `/token-saver desfazer` |
| Ver detalhes e limites | `/token-saver detalhes` |
| Optar por respostas mais objetivas | `/token-saver modo direto` |

O fluxo básico não troca o modelo, não altera o nível de raciocínio, não edita
`CLAUDE.md`, não cria `.claudeignore` e não compacta a conversa automaticamente.
O Modo Direto é opcional e adiciona uma regra curta no projeto.

## O que é medido

O diagnóstico mostra configurações locais, tamanhos de arquivos de instrução,
declarações de MCP e skills detectáveis. Isso não é uma medição de tokens ou custo
da sessão. Confirme o comportamento real no Claude Code com `/status`, `/context`
e `/usage`. Uma assinatura não fica mais barata automaticamente.

## Segurança e desfazer

JSON inválido é preservado e interrompe a aplicação. As escritas são atômicas,
as prévias têm identificador e o histórico fica fora do projeto. O desfazer preserva
edições posteriores em outras chaves e recusa conflitos na mesma chave.

Uma instalação antiga ou personalizada não é sobrescrita. Consulte a [recuperação](skill/references/recovery.md) e os [limites técnicos](skill/references/details.md) quando precisar.

## Para quem desenvolve

O CLI também oferece `audit`, `plan`, `apply`, `undo`, `install` e `details`, com
`--json`, `--dry-run`, `--project`, `--config-dir` e `--expect`. Os testes e o
processo de empacotamento estão em [CONTRIBUTING.md](CONTRIBUTING.md).

## Licença

MIT. [Fontes e validade](skill/references/sources.md) · [English](README.en.md)
