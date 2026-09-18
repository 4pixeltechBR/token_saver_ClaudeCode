package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// auditForHarness keeps the proven Claude adapter intact and routes every
// other supported host through a conservative, read-only adapter.
func auditForHarness(o options) (Audit, error) {
	if o.Harness == "claude" || o.Harness == "" {
		return audit(o)
	}
	return auditPortable(o), nil
}

func harnessInstructionNames(id string) []string {
	common := []string{"AGENTS.md"}
	switch id {
	case "claude":
		return []string{"CLAUDE.md"}
	case "codex":
		return common
	case "antigravity", "gemini":
		return []string{"GEMINI.md", "AGENTS.md"}
	case "opencode", "cursor", "copilot", "cline", "minimax":
		return common
	default:
		return common
	}
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = filepath.Clean(value)
		if value == "." || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func portableSkillRoots(o options) []string {
	home := homePath()
	roots := []string{}
	if h, ok := harnessInfo(o.Harness); ok {
		for _, root := range h.SkillRoots {
			if strings.HasPrefix(root, "~/") {
				roots = append(roots, filepath.Join(home, strings.TrimPrefix(root, "~/")))
			} else {
				roots = append(roots, filepath.Join(o.Project, filepath.FromSlash(root)))
			}
		}
	}
	return uniqueStrings(roots)
}

func countSkillFolders(roots []string) int {
	count := 0
	seen := map[string]bool{}
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			path := filepath.Join(root, entry.Name())
			if seen[path] {
				continue
			}
			if _, err := os.Stat(filepath.Join(path, "SKILL.md")); err == nil {
				seen[path] = true
				count++
			}
		}
	}
	return count
}

func readInstructionSizes(project string, names []string) ([]InstructionSize, []Finding) {
	paths := []string{}
	for _, name := range names {
		paths = append(paths, filepath.Join(project, name))
	}
	files := []InstructionSize{}
	findings := []Finding{}
	for _, path := range uniqueStrings(paths) {
		b, exists, _, err := readFile(path)
		if err != nil || !exists {
			continue
		}
		lines := len(strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"))
		files = append(files, InstructionSize{Path: path, Bytes: len(b), Lines: lines})
		if lines > 200 {
			findings = append(findings, Finding{Level: "SUGESTAO", Code: "instructions", Message: fmt.Sprintf("%s: %d linhas. Revise repeticoes e mova detalhes para referencias; nao e uma contagem de tokens.", filepath.Base(path), lines)})
		}
	}
	return files, findings
}

func auditPortable(o options) Audit {
	h, _ := harnessInfo(o.Harness)
	a := Audit{Version: version, Project: o.Project, Harness: o.Harness, ClaudeVersion: "nao aplicavel", Stack: []string{}, InstructionFiles: []InstructionSize{}, Findings: []Finding{}, Savings: "nao medida", Scope: "Arquivos de instrucao do projeto e diretorios de skills locais/globais; este adaptador nao altera configuracoes.", SkillFolders: countSkillFolders(portableSkillRoots(o)), ToolSearch: "nao aplicavel neste harness"}
	a.Stack = detectStack(o.Project)
	a.InstructionFiles, a.Findings = readInstructionSizes(o.Project, harnessInstructionNames(o.Harness))
	a.Findings = append(a.Findings, Finding{Level: "INFO", Code: "readonly", Message: fmt.Sprintf("%s: auditoria disponivel; aplicacao automatica ainda nao habilitada.", h.Name)})
	a.Findings = append(a.Findings, Finding{Level: "INFO", Code: "verification", Message: "Confirme o comportamento no painel ou comando de contexto do seu harness."})
	a.fingerprint = hash([]byte(o.Harness + "\x00" + o.Project + "\x00" + fmt.Sprint(a.SkillFolders)))
	return a
}

func detectStack(project string) []string {
	stacks := []struct{ file, label string }{{"package.json", "JavaScript/TypeScript"}, {"pyproject.toml", "Python"}, {"requirements.txt", "Python"}, {"go.mod", "Go"}, {"Cargo.toml", "Rust"}, {"composer.json", "PHP"}, {"pom.xml", "Java"}, {"pubspec.yaml", "Dart/Flutter"}}
	seen := map[string]bool{}
	out := []string{}
	for _, s := range stacks {
		if _, err := os.Stat(filepath.Join(project, s.file)); err == nil && !seen[s.label] {
			seen[s.label] = true
			out = append(out, s.label)
		}
	}
	sort.Strings(out)
	return out
}

func (a Audit) portableHuman() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Token Saver %s | Diagnostico %s\nProjeto: %s\n", version, a.Harness, a.Project)
	if len(a.Stack) > 0 {
		fmt.Fprintf(&b, "Projeto detectado: %s\n", strings.Join(a.Stack, " + "))
	}
	fmt.Fprintf(&b, "Skills detectaveis: %d\n", a.SkillFolders)
	for _, f := range a.Findings {
		fmt.Fprintf(&b, "[%s] %s\n", f.Level, f.Message)
	}
	fmt.Fprintf(&b, "Economia: %s. Nenhuma configuracao foi alterada.\n", a.Savings)
	return strings.TrimSpace(b.String())
}

func harnessInfoText(id string) string {
	h, ok := harnessInfo(id)
	if !ok {
		return "Harness desconhecido; nenhuma configuracao foi alterada."
	}
	return fmt.Sprintf("%s esta em modo somente leitura nesta versao. O diagnostico pode ser executado, mas nao altera configuracoes.", h.Name)
}

func outputDetections(out io.Writer, structured bool, detections []HarnessDetection) error {
	if structured {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(detections)
	}
	for _, d := range detections {
		status := "nao detectado"
		if d.Detected {
			status = "detectado"
		}
		fmt.Fprintf(out, "%s: %s — instalar em %s — %s\n", d.Name, status, strings.Join(d.SkillRoots, " ou "), d.Invocation)
	}
	return nil
}
