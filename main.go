package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const version = "3.0.0"

type options struct {
	Project       string
	Config        string
	Harness       string
	JSON          bool
	Yes           bool
	Dry           bool
	Concise       bool
	Expect        string
	ClaudeVersion string
	Managed       string
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Token Saver: "+err.Error())
		os.Exit(1)
	}
}

func run(args []string, in io.Reader, out io.Writer) error {
	command := "audit"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		command, args = strings.ToLower(args[0]), args[1:]
	}
	aliases := map[string]string{"auditar": "audit", "status": "audit", "detalhes": "details", "desfazer": "undo", "aplicar": "apply", "setup": "apply", "help": "help", "--help": "help"}
	if alias, ok := aliases[command]; ok {
		command = alias
	}
	if command == "version" {
		fmt.Fprintln(out, version)
		return nil
	}
	f := flag.NewFlagSet("token-saver "+command, flag.ContinueOnError)
	f.SetOutput(out)
	o := options{}
	f.StringVar(&o.Project, "project", ".", "Pasta do projeto")
	f.StringVar(&o.Config, "config-dir", "", "Pasta de configuracao do harness")
	f.StringVar(&o.Harness, "harness", "claude", "Harness: auto, claude, codex, antigravity, opencode, cursor, gemini, copilot, cline ou minimax")
	f.BoolVar(&o.JSON, "json", false, "Saida estruturada, sem valores secretos")
	f.BoolVar(&o.Yes, "yes", false, "Autoriza as mudancas apresentadas")
	f.BoolVar(&o.Dry, "dry-run", false, "Somente previa, sem escrever arquivos")
	f.BoolVar(&o.Concise, "concise", false, "Propor Modo Direto para este projeto")
	f.StringVar(&o.Expect, "expect", "", "Exigir o identificador da previa revisada")
	if err := f.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if f.NArg() != 0 {
		return errors.New("argumento desconhecido; use --help")
	}
	if command == "help" {
		fmt.Fprintln(out, helpText)
		return nil
	}
	valid := map[string]bool{"audit": true, "plan": true, "apply": true, "undo": true, "details": true, "install": true, "detect": true}
	if !valid[command] {
		return errors.New("comando desconhecido. Use audit, plan, apply, undo, details, detect ou install")
	}
	if command == "details" {
		fmt.Fprintln(out, detailsText)
		return nil
	}
	if err := normalize(&o); err != nil {
		return err
	}
	if command == "detect" {
		return outputDetections(out, o.JSON, detectedHarnesses(o.Project))
	}
	if command == "install" {
		return install(o, out)
	}
	if o.Harness != "claude" && command != "audit" {
		if command == "plan" {
			p := Plan{ID: hash([]byte(o.Harness + "\x00" + o.Project)), Notes: []string{harnessInfoText(o.Harness), "Nenhuma configuracao foi alterada."}, Savings: "nao medida"}
			return output(out, o.JSON, p)
		}
		return fmt.Errorf("o harness %s esta em modo somente leitura; use audit ou aguarde um adaptador de escrita", o.Harness)
	}
	if command == "undo" {
		preview, err := undoPreview(o)
		if err != nil {
			return err
		}
		if o.Dry || preview.ID == "" {
			return output(out, o.JSON, preview)
		}
		if !o.Yes {
			if err := output(out, false, preview); err != nil {
				return err
			}
			if !confirm(in, out) {
				fmt.Fprintln(out, "Cancelado. Nenhuma alteracao aplicada.")
				return nil
			}
		}
		result, err := undo(o, preview.ID)
		if err != nil {
			return err
		}
		return output(out, o.JSON, result)
	}
	a, err := auditForHarness(o)
	if err != nil {
		return err
	}
	if command == "audit" {
		return output(out, o.JSON, a)
	}
	p, err := makePlan(o, a)
	if err != nil {
		return err
	}
	if command == "plan" || o.Dry || len(p.Changes) == 0 {
		return output(out, o.JSON, p)
	}
	if o.Expect != "" && o.Expect != p.ID {
		return errors.New("o ambiente mudou desde a previa. Gere e revise um novo plano")
	}
	if !o.Yes {
		if err := output(out, false, p); err != nil {
			return err
		}
		if !confirm(in, out) {
			fmt.Fprintln(out, "Cancelado. Nenhuma alteracao aplicada.")
			return nil
		}
	}
	result, err := apply(o, p)
	if err != nil {
		return err
	}
	return output(out, o.JSON, result)
}

func normalize(o *options) error {
	var err error
	o.Project, err = filepath.Abs(o.Project)
	if err != nil {
		return err
	}
	o.Project, err = filepath.EvalSymlinks(o.Project)
	if err != nil {
		return errors.New("pasta do projeto nao encontrada")
	}
	info, err := os.Stat(o.Project)
	if err != nil || !info.IsDir() {
		return errors.New("o projeto precisa ser uma pasta")
	}
	if o.Harness == "" {
		o.Harness = "claude"
	}
	o.Harness, err = resolveHarness(o.Harness, o.Project)
	if err != nil {
		return err
	}
	if o.Config == "" && o.Harness == "claude" {
		o.Config = os.Getenv("CLAUDE_CONFIG_DIR")
	}
	o.Config = harnessConfigRoot(o.Harness, o.Config)
	o.Config, err = filepath.Abs(o.Config)
	return err
}

func confirm(in io.Reader, out io.Writer) bool {
	fmt.Fprint(out, "Aplicar estas mudancas? [s/N] ")
	line, _ := bufio.NewReader(in).ReadString('\n')
	return strings.EqualFold(strings.TrimSpace(line), "s") || strings.EqualFold(strings.TrimSpace(line), "y")
}

type printable interface{ human() string }

func output(out io.Writer, structured bool, value printable) error {
	if structured {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(value)
	}
	_, err := fmt.Fprintln(out, value.human())
	return err
}

const helpText = `Token Saver 3.0 Multi-Harness — configuracoes claras, mudancas reversiveis.
  token-saver audit                 Diagnosticar sem alterar nada
  token-saver plan                  Ver mudancas propostas
  token-saver apply                 Revisar e aplicar
  token-saver apply --concise       Incluir Modo Direto (opcional)
  token-saver undo                  Desfazer a ultima aplicacao
  token-saver install               Instalar a skill no harness escolhido
  token-saver detect                Detectar harnesses disponiveis
  token-saver details               Entender os limites

Opcoes: --project PASTA, --harness ID, --config-dir PASTA, --json, --dry-run, --yes, --expect ID.
AUDIT e SETUP antigos continuam aceitos. Nenhuma telemetria ou chamada de IA.`

const detailsText = `Token Saver 3.0 Multi-Harness
O diagnostico le configuracoes locais e mede bytes/linhas; nao inventa tokens.
Quantidade de MCPs ou skills nao prova desperdicio. O CLI nao le dados da sessao
nem mede tokens, cache, custo real ou economia contrafactual.

No Claude Code, o ajuste automatico de Tool Search so e proposto para uma desativacao explicita,
com Claude Code >= 2.1.221, modelo 4.5+ identificavel e acesso direto a Anthropic.
Overrides de ambiente, configuracao gerenciada e provedores desconhecidos
impedem esse ajuste. Flags de inicializacao e politicas remotas nao sao visiveis:
confirme /status e /context depois. Ausencia de variavel nao significa desligado.

O Modo Direto e opcional e adiciona uma regra curta ao projeto, sem reduzir
raciocinio, trocar modelos ou reescrever CLAUDE.md. A regra tambem ocupa contexto.
Nao cria .claudeignore, nao bloqueia arquivos, nao compacta conversas e nao
promete alterar mensalidades. As regras de modelo, thinking e permissoes ficam
preservadas. O CLI nunca imprime valores de env, credenciais ou JSON completo.

Backups locais podem conter configuracoes sensiveis. Ficam em
o diretorio de configuracao do harness/token-saver/state, fora do projeto, com acesso restrito onde
suportado pelo sistema. Undo preserva edicoes posteriores em outras chaves e
recusa conflitos. Execute fora de uma gravacao concorrente das configuracoes.
Uma interrupcao deixa journal recuperavel por undo; nao apague backups pendentes.`
