package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fixture(t *testing.T) options {
	t.Helper()
	for _, key := range []string{"ENABLE_TOOL_SEARCH", "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS", "ANTHROPIC_BASE_URL", "CLAUDE_CODE_USE_BEDROCK", "CLAUDE_CODE_USE_VERTEX", "CLAUDE_CODE_USE_FOUNDRY", "CLAUDE_CODE_PROVIDER_MANAGED_BY_HOST", "ANTHROPIC_MODEL"} {
		key := key
		old, had := os.LookupEnv(key)
		os.Unsetenv(key)
		t.Cleanup(func() {
			if had {
				os.Setenv(key, old)
			} else {
				os.Unsetenv(key)
			}
		})
	}
	root := t.TempDir()
	o := options{Project: filepath.Join(root, "project"), Config: filepath.Join(root, "profile"), ClaudeVersion: "2.1.221", Managed: filepath.Join(root, "no-managed-policy")}
	if err := os.MkdirAll(o.Project, 0700); err != nil {
		t.Fatal(err)
	}
	return o
}
func put(t *testing.T, path string, b []byte) {
	t.Helper()
	if err := atomicWrite(path, b, 0600); err != nil {
		t.Fatal(err)
	}
}
func settings(o options) string { return filepath.Join(o.Project, ".claude", "settings.local.json") }
func seed(t *testing.T, o options) {
	t.Helper()
	put(t, filepath.Join(o.Config, "settings.json"), []byte(`{"model":"claude-sonnet-4-6","env":{"ENABLE_TOOL_SEARCH":"false"}}`))
}
func planFor(t *testing.T, o options) Plan {
	t.Helper()
	a, err := audit(o)
	if err != nil {
		t.Fatal(err)
	}
	p, err := makePlan(o, a)
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func applyFor(t *testing.T, o options) Result {
	t.Helper()
	p := planFor(t, o)
	if len(p.Changes) == 0 {
		t.Fatal("expected useful changes")
	}
	r, err := apply(o, p)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func tree(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
		if os.IsNotExist(e) {
			return nil
		}
		if e != nil {
			return e
		}
		if d.IsDir() {
			m[p] = "dir"
		} else {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			m[p] = hash(b)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestStrictJSON(t *testing.T) {
	for _, data := range []string{`{"a":1,}`, `{"a":1,"a":2}`, `{"env":{"X":"a","X":"b"}}`, `[]`, `null`, `{"a":1} {"a":2}`, `{"a":1}//comment`, ""} {
		if _, err := object([]byte(data)); err == nil {
			t.Errorf("accepted invalid JSON: %s", data)
		}
	}
	for _, data := range []string{`{"url":"https://example.com/a//b","number":9007199254740993}`, "\ufeff{\"env\":{}}"} {
		if _, err := object([]byte(data)); err != nil {
			t.Error(err)
		}
	}
}
func TestDryRunsDoNotWrite(t *testing.T) {
	o := fixture(t)
	before := tree(t, filepath.Dir(o.Project))
	for _, cmd := range []string{"audit", "plan", "apply", "undo", "install"} {
		var b bytes.Buffer
		err := run([]string{cmd, "--project", o.Project, "--config-dir", o.Config, "--dry-run", "--json"}, strings.NewReader(""), &b)
		if err != nil {
			t.Fatal(cmd, err)
		}
	}
	after := tree(t, filepath.Dir(o.Project))
	if !reflect.DeepEqual(before, after) {
		t.Fatal("read-only command wrote files")
	}
}
func TestInvalidSettingsRemainUntouched(t *testing.T) {
	o := fixture(t)
	bad := []byte(`{"env":{"URL":"https://example.com"},"permissions":{"deny":["Read(.env)"]},}`)
	put(t, settings(o), bad)
	before := tree(t, filepath.Dir(o.Project))
	_, err := audit(o)
	if err == nil {
		t.Fatal("malformed JSON accepted")
	}
	after := tree(t, filepath.Dir(o.Project))
	if !reflect.DeepEqual(before, after) {
		t.Fatal("invalid original changed")
	}
}
func TestApplyPreservesConfigurationAndExactUndo(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	before := []byte("{\n \"permissions\": {\"deny\": [\"Read(.env)\"]},\n \"model\":\"claude-opus-4-6\",\n \"env\":{\"URL\":\"https://example.com/mcp\",\"SECRET\":\"never-print-this\"},\n \"largeInteger\":9007199254740993\n}\n")
	put(t, settings(o), before)
	r := applyFor(t, o)
	m, err := loadObject(settings(o))
	if err != nil {
		t.Fatal(err)
	}
	if m["model"] != "claude-opus-4-6" || m["largeInteger"] != json.Number("9007199254740993") {
		t.Fatal("unrelated values changed")
	}
	e, _ := envMap(m)
	if e["URL"] != "https://example.com/mcp" || e["SECRET"] != "never-print-this" || e["ENABLE_TOOL_SEARCH"] != "true" {
		t.Fatal("env merge failed")
	}
	if _, ok := m["permissions"]; !ok {
		t.Fatal("permissions lost")
	}
	if _, err = undo(o, r.ID); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(settings(o))
	if !bytes.Equal(b, before) {
		t.Fatal("undo did not preserve original bytes")
	}
}
func TestUndoPreservesLaterEdits(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	r := applyFor(t, o)
	m, _ := loadObject(settings(o))
	m["model"] = "claude-opus-4-6"
	e, _ := envMap(m)
	e["LATER"] = "keep"
	put(t, settings(o), encode(m))
	if _, err := undo(o, r.ID); err != nil {
		t.Fatal(err)
	}
	m, _ = loadObject(settings(o))
	e, _ = envMap(m)
	if m["model"] != "claude-opus-4-6" || e["LATER"] != "keep" {
		t.Fatal("later changes lost")
	}
	if _, ok := e["ENABLE_TOOL_SEARCH"]; ok {
		t.Fatal("owned key remains")
	}
}
func TestUndoConflictsBeforeAnyWrite(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	o.Concise = true
	r := applyFor(t, o)
	rule := filepath.Join(o.Project, ".claude", "rules", "token-saver.md")
	put(t, rule, []byte("my new rule"))
	before := tree(t, filepath.Dir(o.Project))
	if _, err := undo(o, r.ID); err == nil {
		t.Fatal("expected conflict")
	}
	after := tree(t, filepath.Dir(o.Project))
	if !reflect.DeepEqual(before, after) {
		t.Fatal("conflict caused partial restore")
	}
}
func TestSameKeyConflict(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	r := applyFor(t, o)
	m, _ := loadObject(settings(o))
	e, _ := envMap(m)
	e["ENABLE_TOOL_SEARCH"] = "auto:5"
	put(t, settings(o), encode(m))
	if _, err := undo(o, r.ID); err == nil {
		t.Fatal("same-key edit was overwritten")
	}
}
func TestDefaultNoopAndIdempotence(t *testing.T) {
	o := fixture(t)
	p := planFor(t, o)
	if len(p.Changes) != 0 {
		t.Fatal("missing variable is not an optimization")
	}
	seed(t, o)
	applyFor(t, o)
	p = planFor(t, o)
	if len(p.Changes) != 0 {
		t.Fatal("second execution should be noop")
	}
}
func TestStalePlanRefused(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	p := planFor(t, o)
	put(t, settings(o), []byte(`{"theme":"dark"}`))
	if _, err := apply(o, p); err == nil {
		t.Fatal("stale plan accepted")
	}
	m, _ := loadObject(settings(o))
	if len(m) != 1 {
		t.Fatal("stale plan mutated file")
	}
}
func TestConservativeCompatibility(t *testing.T) {
	cases := []struct {
		name   string
		change func(*testing.T, options)
	}{
		{"proxy", func(t *testing.T, o options) { t.Setenv("ANTHROPIC_BASE_URL", "https://example.com") }},
		{"environment", func(t *testing.T, o options) { t.Setenv("ENABLE_TOOL_SEARCH", "false") }},
		{"provider", func(t *testing.T, o options) { t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1") }},
		{"managed", func(t *testing.T, o options) { put(t, o.Managed, []byte(`{}`)) }},
		{"unknown-model", func(t *testing.T, o options) { put(t, settings(o), []byte(`{"model":"sonnet"}`)) }},
		{"betas", func(t *testing.T, o options) { t.Setenv("CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS", "1") }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := fixture(t)
			seed(t, o)
			c.change(t, o)
			if len(planFor(t, o).Changes) != 0 {
				t.Fatal("unsafe compatibility assumption")
			}
		})
	}
}
func TestOldVersionRefused(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	o.ClaudeVersion = "2.1.100"
	if len(planFor(t, o).Changes) != 0 {
		t.Fatal("old version accepted")
	}
}
func TestEffectiveSettingsPrecedence(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	put(t, filepath.Join(o.Project, ".claude", "settings.json"), []byte(`{"env":{"ENABLE_TOOL_SEARCH":"true"}}`))
	if len(planFor(t, o).Changes) != 0 {
		t.Fatal("project override ignored")
	}
}
func TestOutputDoesNotLeakSecrets(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	put(t, settings(o), []byte(`{"env":{"API_KEY":"dont-leak-me","URL":"https://private.example/token"}}`))
	a, err := audit(o)
	if err != nil {
		t.Fatal(err)
	}
	p := planFor(t, o)
	for _, v := range []printable{a, p} {
		for _, structured := range []bool{true, false} {
			var b bytes.Buffer
			output(&b, structured, v)
			if strings.Contains(b.String(), "dont-leak-me") || strings.Contains(b.String(), "private.example") {
				t.Fatal("secret in output")
			}
		}
	}
}
func TestConciseOptInAndUndo(t *testing.T) {
	o := fixture(t)
	o.Concise = true
	r := applyFor(t, o)
	p := planFor(t, o)
	if len(p.Changes) != 0 {
		t.Fatal("duplicate rule")
	}
	if _, err := undo(o, r.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(o.Project, ".claude")); !os.IsNotExist(err) {
		t.Fatal("new empty project folders remain")
	}
}
func TestExistingRulePreserved(t *testing.T) {
	o := fixture(t)
	o.Concise = true
	put(t, filepath.Join(o.Project, ".claude", "rules", "token-saver.md"), []byte("custom"))
	a, _ := audit(o)
	if _, err := makePlan(o, a); err == nil {
		t.Fatal("custom rule overwritten")
	}
}
func TestPendingJournalRecovery(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	o.Concise = true
	p := planFor(t, o)
	j := Journal{ID: newID(), Project: o.Project, Status: "pending"}
	for _, c := range p.Changes {
		j.Files = append(j.Files, Snapshot{c.Path, c.Kind, c.Before, c.After, c.Existed, uint32(c.Mode), c.BeforeKey, c.KeyExisted, c.EnvExisted})
	}
	if err := saveJournal(o, j); err != nil {
		t.Fatal(err)
	}
	put(t, p.Changes[0].Path, p.Changes[0].After)
	if _, err := apply(o, planFor(t, o)); err == nil {
		t.Fatal("pending operation not blocked")
	}
	if _, err := undo(o, j.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(settings(o)); !os.IsNotExist(err) {
		t.Fatal("partial write not undone")
	}
}
func TestTamperedJournalCannotWriteOutsideProject(t *testing.T) {
	o := fixture(t)
	j := Journal{ID: newID(), Project: o.Project, Status: "pending", Files: []Snapshot{{Kind: "tool-search", Path: filepath.Join(o.Config, "settings.json")}}}
	saveJournal(o, j)
	if _, err := undoPreview(o); err == nil {
		t.Fatal("out-of-scope restore accepted")
	}
}
func TestInstallUpdateAndCustomProtection(t *testing.T) {
	o := fixture(t)
	var b bytes.Buffer
	if err := install(o, &b); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(o.Config, "skills", "token-saver")
	if err := verifyOwned(dest); err != nil {
		t.Fatal(err)
	}
	if err := install(o, &b); err != nil {
		t.Fatal("reinstall", err)
	}
	put(t, filepath.Join(dest, "SKILL.md"), []byte("custom instructions"))
	before := tree(t, dest)
	if err := install(o, &b); err == nil {
		t.Fatal("custom installation overwritten")
	}
	if !reflect.DeepEqual(before, tree(t, dest)) {
		t.Fatal("custom files changed")
	}
}
func TestUnknownSkillNotOverwritten(t *testing.T) {
	o := fixture(t)
	put(t, filepath.Join(o.Config, "skills", "token-saver", "SKILL.md"), []byte("old skill"))
	if err := install(o, &bytes.Buffer{}); err == nil {
		t.Fatal("old skill overwritten")
	}
}
func TestLockPreventsConcurrentWriter(t *testing.T) {
	o := fixture(t)
	seed(t, o)
	release, err := lock(o)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if _, err := apply(o, planFor(t, o)); err == nil {
		t.Fatal("lock ignored")
	}
}
func TestEmptyAnswerCancels(t *testing.T) {
	if confirm(strings.NewReader("\n"), &bytes.Buffer{}) {
		t.Fatal("empty answer authorizes mutation")
	}
}
func TestWrongEnvTypeStops(t *testing.T) {
	o := fixture(t)
	put(t, settings(o), []byte(`{"env":[]}`))
	if _, err := audit(o); err == nil {
		t.Fatal("invalid env accepted")
	}
}
func TestStackDetectionUsesManifests(t *testing.T) {
	o := fixture(t)
	put(t, filepath.Join(o.Project, "package.json"), []byte(`{}`))
	put(t, filepath.Join(o.Project, "pyproject.toml"), []byte("[project]"))
	a, err := audit(o)
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Stack) != 2 {
		t.Fatal(a.Stack)
	}
}
func TestSymlinkWriteRefused(t *testing.T) {
	o := fixture(t)
	outside := filepath.Join(filepath.Dir(o.Project), "outside")
	os.MkdirAll(outside, 0700)
	if err := os.Symlink(outside, filepath.Join(o.Project, ".claude")); err != nil {
		t.Skip("symlinks require permission on this OS")
	}
	o.Concise = true
	a, _ := audit(o)
	if _, err := makePlan(o, a); err == nil {
		t.Fatal("symlink allowed")
	}
}

func TestPortableHarnessAuditIsReadOnly(t *testing.T) {
	o := fixture(t)
	o.Harness = "codex"
	put(t, filepath.Join(o.Project, "AGENTS.md"), []byte("# project\n"))
	before := tree(t, filepath.Dir(o.Project))
	a, err := auditForHarness(o)
	if err != nil {
		t.Fatal(err)
	}
	if a.Harness != "codex" || a.Scope == "" || len(a.InstructionFiles) != 1 {
		t.Fatalf("unexpected portable audit: %+v", a)
	}
	if !reflect.DeepEqual(before, tree(t, filepath.Dir(o.Project))) {
		t.Fatal("portable audit wrote files")
	}
}

func TestPortableHarnessPlanDoesNotWrite(t *testing.T) {
	o := fixture(t)
	var out bytes.Buffer
	if err := run([]string{"plan", "--harness", "opencode", "--project", o.Project, "--config-dir", o.Config, "--json"}, strings.NewReader(""), &out); err != nil {
		t.Fatal(err)
	}
	var p Plan
	if err := json.Unmarshal(out.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	if len(p.Changes) != 0 || len(p.Notes) == 0 {
		t.Fatalf("portable plan should be read-only: %+v", p)
	}
}

func TestHarnessInstallDestination(t *testing.T) {
	o := fixture(t)
	o.Harness = "codex"
	if err := install(o, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(o.Config, "skills", "token-saver")
	if err := verifyOwned(want); err != nil {
		t.Fatal(err)
	}
}
