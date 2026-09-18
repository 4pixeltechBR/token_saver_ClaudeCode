package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// HarnessInfo describes the smallest contract needed by the portable skill.
// The CLI intentionally starts conservative: only Claude Code has automatic
// settings writes; the other adapters provide discovery and read-only audit.
type HarnessInfo struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	SkillRoots     []string `json:"skill_roots"`
	Invocation     string   `json:"invocation"`
	Capabilities   []string `json:"capabilities"`
	WriteSupported bool     `json:"write_supported"`
}

type HarnessDetection struct {
	HarnessInfo
	Detected bool `json:"detected"`
}

var harnessCatalog = map[string]HarnessInfo{
	"claude": {
		ID: "claude", Name: "Claude Code", SkillRoots: []string{".claude/skills"},
		Invocation: "/token-saver", Capabilities: []string{"detect", "audit", "plan", "apply", "undo", "modo direto"}, WriteSupported: true,
	},
	"codex": {
		ID: "codex", Name: "Codex", SkillRoots: []string{".agents/skills", "~/.agents/skills"},
		Invocation: "$token-saver", Capabilities: []string{"detect", "audit", "details", "install"},
	},
	"antigravity": {
		ID: "antigravity", Name: "Antigravity", SkillRoots: []string{".agents/skills", "~/.gemini/antigravity/skills"},
		Invocation: "mencione token-saver ou peça uma auditoria", Capabilities: []string{"detect", "audit", "details", "install"},
	},
	"opencode": {
		ID: "opencode", Name: "OpenCode", SkillRoots: []string{".opencode/skills", ".agents/skills", ".claude/skills"},
		Invocation: "use a skill token-saver", Capabilities: []string{"detect", "audit", "details", "install"},
	},
	"cursor": {
		ID: "cursor", Name: "Cursor", SkillRoots: []string{".cursor/skills", ".agents/skills", ".claude/skills"},
		Invocation: "/token-saver", Capabilities: []string{"detect", "audit", "details", "install"},
	},
	"gemini": {
		ID: "gemini", Name: "Gemini CLI", SkillRoots: []string{".gemini/skills", ".agents/skills"},
		Invocation: "ative token-saver quando a tarefa mencionar custo ou contexto", Capabilities: []string{"detect", "audit", "details", "install"},
	},
	"copilot": {
		ID: "copilot", Name: "GitHub Copilot", SkillRoots: []string{".github/skills", ".agents/skills", ".claude/skills"},
		Invocation: "/token-saver", Capabilities: []string{"detect", "audit", "details", "install"},
	},
	"cline": {
		ID: "cline", Name: "Cline", SkillRoots: []string{".cline/skills", ".agents/skills"},
		Invocation: "use a skill token-saver", Capabilities: []string{"detect", "audit", "details", "install"},
	},
	"minimax": {
		ID: "minimax", Name: "MiniMax Code", SkillRoots: []string{"skills/token-saver"},
		Invocation: "ative o plugin token-saver", Capabilities: []string{"detect", "audit", "details"},
	},
}

func harnessInfo(id string) (HarnessInfo, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	h, ok := harnessCatalog[id]
	return h, ok
}

func knownHarnesses() []HarnessInfo {
	ids := make([]string, 0, len(harnessCatalog))
	for id := range harnessCatalog {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]HarnessInfo, 0, len(ids))
	for _, id := range ids {
		out = append(out, harnessCatalog[id])
	}
	return out
}

func homePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func executableExists(names ...string) bool {
	for _, name := range names {
		if _, err := exec.LookPath(name); err == nil {
			return true
		}
	}
	return false
}

func detectedHarnesses(project string) []HarnessDetection {
	home := homePath()
	checks := map[string]bool{
		"claude":      executableExists("claude") || pathExists(filepath.Join(home, ".claude")) || pathExists(filepath.Join(project, ".claude")),
		"codex":       executableExists("codex") || pathExists(filepath.Join(home, ".codex")) || pathExists(filepath.Join(project, ".agents")),
		"antigravity": executableExists("agy", "antigravity") || pathExists(filepath.Join(home, ".gemini", "antigravity")) || pathExists(filepath.Join(project, ".agent")),
		"opencode":    executableExists("opencode") || pathExists(filepath.Join(home, ".config", "opencode")) || pathExists(filepath.Join(project, ".opencode")),
		"cursor":      executableExists("cursor") || pathExists(filepath.Join(home, ".cursor")) || pathExists(filepath.Join(project, ".cursor")),
		"gemini":      executableExists("gemini") || pathExists(filepath.Join(home, ".gemini")) || pathExists(filepath.Join(project, ".gemini")),
		"copilot":     executableExists("copilot") || pathExists(filepath.Join(home, ".copilot")) || pathExists(filepath.Join(project, ".github", "skills")),
		"cline":       pathExists(filepath.Join(home, ".cline")) || pathExists(filepath.Join(project, ".cline")),
		"minimax":     executableExists("mcode", "minimax") || pathExists(filepath.Join(home, ".minimax")),
	}
	out := make([]HarnessDetection, 0, len(harnessCatalog))
	for _, h := range knownHarnesses() {
		out = append(out, HarnessDetection{HarnessInfo: h, Detected: checks[h.ID]})
	}
	return out
}

func resolveHarness(id, project string) (string, error) {
	id = strings.ToLower(strings.TrimSpace(id))
	if id == "" {
		id = "claude"
	}
	if id != "auto" {
		if _, ok := harnessInfo(id); !ok {
			return "", fmt.Errorf("harness desconhecido: %s", id)
		}
		return id, nil
	}
	for _, d := range detectedHarnesses(project) {
		if d.Detected && d.ID == "claude" {
			return d.ID, nil
		}
	}
	for _, d := range detectedHarnesses(project) {
		if d.Detected {
			return d.ID, nil
		}
	}
	return "claude", nil
}

func harnessConfigRoot(id, configured string) string {
	if configured != "" {
		return configured
	}
	home := homePath()
	switch id {
	case "claude":
		return filepath.Join(home, ".claude")
	case "codex":
		return filepath.Join(home, ".agents")
	case "antigravity":
		return filepath.Join(home, ".gemini", "antigravity")
	case "opencode":
		return filepath.Join(home, ".config", "opencode")
	case "cursor":
		return filepath.Join(home, ".cursor")
	case "gemini":
		return filepath.Join(home, ".gemini")
	case "copilot":
		return filepath.Join(home, ".copilot")
	case "cline":
		return filepath.Join(home, ".cline")
	case "minimax":
		return filepath.Join(home, ".minimax")
	default:
		return filepath.Join(home, ".token-saver")
	}
}

func harnessSkillDestination(id, configured string) string {
	root := harnessConfigRoot(id, configured)
	if id == "minimax" {
		return filepath.Join(root, "skills", "token-saver")
	}
	return filepath.Join(root, "skills", "token-saver")
}
