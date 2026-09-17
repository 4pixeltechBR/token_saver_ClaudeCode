package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Finding struct {
	Level   string `json:"level"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
type InstructionSize struct {
	Path  string `json:"path"`
	Bytes int    `json:"bytes"`
	Lines int    `json:"lines"`
}
type Audit struct {
	Version          string            `json:"version"`
	Project          string            `json:"project"`
	ClaudeVersion    string            `json:"claude_version"`
	Stack            []string          `json:"stack"`
	InstructionFiles []InstructionSize `json:"instruction_files"`
	MCPDeclarations  int               `json:"mcp_declarations"`
	SkillFolders     int               `json:"skill_folders"`
	ToolSearch       string            `json:"tool_search"`
	Findings         []Finding         `json:"findings"`
	Savings          string            `json:"savings"`
	Scope            string            `json:"scope"`
	eligible         bool
	fingerprint      string
}

func (a Audit) human() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Token Saver %s | Diagnostico\nProjeto: %s\n", version, a.Project)
	if len(a.Stack) > 0 {
		fmt.Fprintf(&b, "Projeto detectado: %s\n", strings.Join(a.Stack, " + "))
	}
	fmt.Fprintf(&b, "Busca de ferramentas: %s\n", a.ToolSearch)
	for _, f := range a.Findings {
		fmt.Fprintf(&b, "[%s] %s\n", f.Level, f.Message)
	}
	if len(a.Findings) == 0 {
		b.WriteString("Nenhuma mudanca necessaria nas configuracoes verificadas.\n")
	}
	fmt.Fprintf(&b, "Economia: %s. Modelo e raciocinio preservados.\n", a.Savings)
	b.WriteString("Configuracoes locais; confirme a sessao em /status e /context.\n")
	return strings.TrimSpace(b.String())
}

func audit(o options) (Audit, error) {
	a := Audit{Version: version, Project: o.Project, ClaudeVersion: "nao identificado", Stack: []string{}, InstructionFiles: []InstructionSize{}, Findings: []Finding{}, Savings: "nao medida", Scope: "Arquivos locais; nao inclui flags de inicio, politicas remotas, ferramentas carregadas ou consumo da sessao."}
	files := []string{filepath.Join(o.Config, "settings.json"), filepath.Join(o.Project, ".claude", "settings.json"), filepath.Join(o.Project, ".claude", "settings.local.json")}
	combinedEnv := map[string]any{}
	effective := map[string]any{}
	fingerprint := []string{version, o.Project, o.Config}
	for _, p := range files {
		m, err := loadObject(p)
		if err != nil {
			return a, err
		}
		b, _, _, _ := readFile(p)
		fingerprint = append(fingerprint, p, hash(b))
		for k, v := range m {
			effective[k] = v
		}
		e, err := envMap(m)
		if err != nil {
			return a, fmt.Errorf("%s: %w", p, err)
		}
		for k, v := range e {
			combinedEnv[k] = v
		}
	}
	versionText := o.ClaudeVersion
	if versionText == "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if p, err := exec.LookPath("claude"); err == nil {
			b, e := exec.CommandContext(ctx, p, "--version").Output()
			if e == nil {
				versionText = string(b)
			}
		}
	}
	if m := regexp.MustCompile(`\b\d+\.\d+\.\d+\b`).FindString(versionText); m != "" {
		a.ClaudeVersion = m
	}
	fingerprint = append(fingerprint, a.ClaudeVersion)
	setting := stringValue(combinedEnv["ENABLE_TOOL_SEARCH"])
	if v, ok := combinedEnv["ENABLE_TOOL_SEARCH"]; ok {
		if _, ok := v.(string); !ok {
			return a, fmt.Errorf("ENABLE_TOOL_SEARCH deve ser texto nas configuracoes")
		}
	}
	variableNames := []string{"ENABLE_TOOL_SEARCH", "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS", "ANTHROPIC_BASE_URL", "CLAUDE_CODE_USE_BEDROCK", "CLAUDE_CODE_USE_VERTEX", "CLAUDE_CODE_USE_FOUNDRY", "CLAUDE_CODE_PROVIDER_MANAGED_BY_HOST", "ANTHROPIC_MODEL"}
	actualEnv := map[string]string{}
	for _, key := range variableNames {
		actualEnv[key] = stringValue(combinedEnv[key])
		if v, ok := os.LookupEnv(key); ok {
			actualEnv[key] = v
		}
		fingerprint = append(fingerprint, key, actualEnv[key])
	}
	blocked := false
	if a.ClaudeVersion == "nao identificado" || !atLeast(a.ClaudeVersion, 2, 1, 221) {
		blocked = true
		a.Findings = append(a.Findings, Finding{"INFO", "version", "Versao compativel nao confirmada; ajuste automatico de ferramentas indisponivel."})
	}
	for _, path := range managedPaths(o) {
		if _, err := os.Stat(path); err == nil {
			blocked = true
			a.Findings = append(a.Findings, Finding{"INFO", "managed", "Politica gerenciada encontrada; revise as regras da organizacao antes de alterar ferramentas."})
			fingerprint = append(fingerprint, path)
			break
		}
	}
	if os.Getenv("ENABLE_TOOL_SEARCH") != "" {
		blocked = true
		a.Findings = append(a.Findings, Finding{"INFO", "environment", "Busca de ferramentas definida no ambiente desta execucao; origem precisa ser verificada."})
	}
	base := strings.TrimRight(actualEnv["ANTHROPIC_BASE_URL"], "/")
	if base != "" && base != "https://api.anthropic.com" {
		blocked = true
	}
	for _, key := range []string{"CLAUDE_CODE_USE_BEDROCK", "CLAUDE_CODE_USE_VERTEX", "CLAUDE_CODE_USE_FOUNDRY", "CLAUDE_CODE_PROVIDER_MANAGED_BY_HOST", "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS"} {
		v := actualEnv[key]
		if v != "" && v != "0" && v != "false" {
			blocked = true
		}
	}
	model := stringValue(effective["model"])
	if actualEnv["ANTHROPIC_MODEL"] != "" {
		model = actualEnv["ANTHROPIC_MODEL"]
	}
	if !supportedModel(model) {
		blocked = true
	}
	if strings.EqualFold(actualEnv["ENABLE_TOOL_SEARCH"], "false") {
		a.ToolSearch = "desativacao encontrada"
		if !blocked && setting == "false" {
			a.eligible = true
			a.Findings = append(a.Findings, Finding{"ACAO", "tool-search", "A busca sob demanda foi desativada. Podemos reativar apenas neste projeto."})
		} else {
			a.Findings = append(a.Findings, Finding{"REVISAR", "tool-search", "Busca desativada, mas a compatibilidade ou a origem nao esta confirmada. Nenhuma mudanca automatica."})
		}
	} else if setting == "" && actualEnv["ENABLE_TOOL_SEARCH"] == "" {
		a.ToolSearch = "padrao do Claude Code; confirmar na sessao"
	} else {
		a.ToolSearch = "configurada; confirmar na sessao"
	}
	stacks := []struct{ file, label string }{{"package.json", "JavaScript/TypeScript"}, {"pyproject.toml", "Python"}, {"requirements.txt", "Python"}, {"go.mod", "Go"}, {"Cargo.toml", "Rust"}, {"composer.json", "PHP"}, {"pom.xml", "Java"}, {"pubspec.yaml", "Dart/Flutter"}}
	seen := map[string]bool{}
	for _, s := range stacks {
		if _, err := os.Stat(filepath.Join(o.Project, s.file)); err == nil && !seen[s.label] {
			a.Stack = append(a.Stack, s.label)
			seen[s.label] = true
		}
	}
	for _, p := range []string{filepath.Join(o.Config, "CLAUDE.md"), filepath.Join(o.Project, "CLAUDE.md"), filepath.Join(o.Project, ".claude", "CLAUDE.md")} {
		b, exists, _, err := readFile(p)
		if err != nil {
			return a, err
		}
		if exists {
			lines := len(strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"))
			a.InstructionFiles = append(a.InstructionFiles, InstructionSize{p, len(b), lines})
			if lines > 200 {
				a.Findings = append(a.Findings, Finding{"SUGESTAO", "instructions", fmt.Sprintf("%s: %d linhas. Revise repeticoes e mova detalhes para referencias; nao e uma contagem de tokens.", filepath.Base(p), lines)})
			}
		}
	}
	// Count declarations without printing names, commands, environment or credentials.
	for _, p := range []string{filepath.Join(o.Project, ".mcp.json"), filepath.Join(o.Config, "..", ".claude.json")} {
		m, err := loadObject(p)
		if err != nil {
			a.Findings = append(a.Findings, Finding{"REVISAR", "mcp-json", "Uma configuracao de MCP nao pode ser lida; contagem incompleta."})
			continue
		}
		if servers, ok := m["mcpServers"].(map[string]any); ok {
			a.MCPDeclarations += len(servers)
		}
		if projects, ok := m["projects"].(map[string]any); ok {
			for pth, v := range projects {
				if filepath.Clean(pth) == filepath.Clean(o.Project) {
					if cfg, ok := v.(map[string]any); ok {
						if servers, ok := cfg["mcpServers"].(map[string]any); ok {
							a.MCPDeclarations += len(servers)
						}
					}
				}
			}
		}
	}
	for _, p := range []string{filepath.Join(o.Config, "skills"), filepath.Join(o.Project, ".claude", "skills")} {
		entries, _ := os.ReadDir(p)
		for _, entry := range entries {
			if entry.IsDir() {
				if _, err := os.Stat(filepath.Join(p, entry.Name(), "SKILL.md")); err == nil {
					a.SkillFolders++
				}
			}
		}
	}
	// Ancestor settings can affect a session: never guess which project Claude selected.
	userHome, _ := os.UserHomeDir()
	for parent := filepath.Dir(o.Project); parent != filepath.Dir(parent); parent = filepath.Dir(parent) {
		if filepath.Clean(parent) == filepath.Clean(userHome) || filepath.Clean(filepath.Join(parent, ".claude")) == filepath.Clean(o.Config) {
			continue
		}
		for _, name := range []string{"settings.json", "settings.local.json"} {
			if _, err := os.Stat(filepath.Join(parent, ".claude", name)); err == nil {
				a.eligible = false
				a.Findings = append(a.Findings, Finding{"REVISAR", "ancestor", "Configuracao em pasta ancestral: rode o diagnostico na raiz usada pelo Claude Code."})
				break
			}
		}
	}
	sort.Strings(a.Stack)
	a.fingerprint = hash([]byte(strings.Join(fingerprint, "\x00")))
	return a, nil
}

func atLeast(s string, major, minor, patch int) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return false
	}
	want := []int{major, minor, patch}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return false
		}
		if n > want[i] {
			return true
		}
		if n < want[i] {
			return false
		}
	}
	return true
}
func supportedModel(s string) bool {
	m := regexp.MustCompile(`^claude-(?:sonnet|opus|haiku)-(\d+)[.-](\d+)(?:-|$)`).FindStringSubmatch(s)
	if len(m) != 3 {
		return false
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	return major > 4 || (major == 4 && minor >= 5)
}
func managedPaths(o options) []string {
	if o.Managed != "" {
		return []string{o.Managed}
	}
	paths := []string{filepath.Join(o.Config, "managed-settings.json")}
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("ProgramFiles")
		if base == "" {
			base = `C:\Program Files`
		}
		paths = append(paths, filepath.Join(base, "ClaudeCode", "managed-settings.json"))
	case "darwin":
		paths = append(paths, "/Library/Application Support/ClaudeCode/managed-settings.json")
	default:
		paths = append(paths, "/etc/claude-code/managed-settings.json")
	}
	return paths
}
