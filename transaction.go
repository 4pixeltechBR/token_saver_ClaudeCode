package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
)

type Snapshot struct {
	Path       string `json:"path"`
	Kind       string `json:"kind"`
	Before     []byte `json:"before"`
	After      []byte `json:"after"`
	Existed    bool   `json:"existed"`
	Mode       uint32 `json:"mode"`
	BeforeKey  any    `json:"before_key,omitempty"`
	KeyExisted bool   `json:"key_existed"`
	EnvExisted bool   `json:"env_existed"`
}
type Journal struct {
	ID          string     `json:"id"`
	Project     string     `json:"project"`
	Status      string     `json:"status"`
	Files       []Snapshot `json:"files"`
	CreatedDirs []string   `json:"created_dirs"`
}
type Result struct {
	ID      string   `json:"id,omitempty"`
	Status  string   `json:"status"`
	Changes []string `json:"changes"`
	Savings string   `json:"savings"`
}

func (r Result) human() string {
	var b strings.Builder
	fmt.Fprintln(&b, r.Status)
	for _, c := range r.Changes {
		fmt.Fprintln(&b, "- "+c)
	}
	if len(r.Changes) > 0 {
		fmt.Fprintln(&b, "Economia: nao medida. Confirme o resultado em /status e /context.")
	}
	return strings.TrimSpace(b.String())
}

func stateDir(o options) string {
	return filepath.Join(o.Config, "token-saver", "state", hash([]byte(o.Project))[:24])
}
func lock(o options) (func(), error) {
	dir := stateDir(o)
	if err := noLinks(dir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	p := filepath.Join(dir, "operation.lock")
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.New("outra operacao pode estar em andamento. Se houve interrupcao, consulte references/recovery.md antes de remover operation.lock")
	}
	fmt.Fprintln(f, os.Getpid())
	f.Close()
	return func() { os.Remove(p) }, nil
}
func newID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return time.Now().UTC().Format("20060102T150405.000000000") + "-" + hex.EncodeToString(b)
}
func saveJournal(o options, j Journal) error {
	return atomicWrite(filepath.Join(stateDir(o), j.ID+".json"), encode(j), 0600)
}

func latest(o options) (Journal, error) {
	entries, err := os.ReadDir(stateDir(o))
	if errors.Is(err, os.ErrNotExist) {
		return Journal{}, nil
	}
	if err != nil {
		return Journal{}, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() > entries[j].Name() })
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		b, _, _, err := readLimited(filepath.Join(stateDir(o), entry.Name()), maxFile*6)
		if err != nil {
			return Journal{}, err
		}
		var j Journal
		if err = json.Unmarshal(b, &j); err != nil {
			return Journal{}, errors.New("historico invalido; nenhum arquivo foi restaurado")
		}
		if j.ID+".json" != entry.Name() || j.Project != o.Project {
			return Journal{}, errors.New("historico nao corresponde ao projeto")
		}
		if j.Status == "undone" || j.Status == "rolled-back" {
			continue
		}
		if j.Status != "pending" && j.Status != "committed" {
			return Journal{}, errors.New("estado de historico desconhecido")
		}
		if err := validateJournal(o, j); err != nil {
			return Journal{}, err
		}
		return j, nil
	}
	return Journal{}, nil
}

func validateJournal(o options, j Journal) error {
	seen := map[string]bool{}
	for _, s := range j.Files {
		want := ""
		switch s.Kind {
		case "tool-search":
			want = filepath.Join(o.Project, ".claude", "settings.local.json")
		case "concise":
			want = filepath.Join(o.Project, ".claude", "rules", "token-saver.md")
		default:
			return errors.New("tipo de backup desconhecido")
		}
		if s.Path != want || seen[want] {
			return errors.New("caminho de backup invalido")
		}
		seen[want] = true
		if err := noLinks(want); err != nil {
			return err
		}
	}
	for _, d := range j.CreatedDirs {
		if d != filepath.Join(o.Project, ".claude") && d != filepath.Join(o.Project, ".claude", "rules") {
			return errors.New("pasta de backup invalida")
		}
	}
	return nil
}

func apply(o options, p Plan) (Result, error) {
	r := Result{Changes: []string{}, Savings: "nao medida"}
	release, err := lock(o)
	if err != nil {
		return r, err
	}
	defer release()
	previous, err := latest(o)
	if err != nil {
		return r, err
	}
	if previous.Status == "pending" {
		return r, errors.New("ha uma aplicacao interrompida. Use undo antes de continuar")
	}
	a, err := audit(o)
	if err != nil {
		return r, err
	}
	fresh, err := makePlan(o, a)
	if err != nil {
		return r, err
	}
	if fresh.ID != p.ID {
		return r, errors.New("configuracao mudou depois da previa; gere um novo plano")
	}
	j := Journal{ID: newID(), Project: o.Project, Status: "pending", Files: []Snapshot{}, CreatedDirs: []string{}}
	dirs := map[string]bool{}
	for _, c := range p.Changes {
		j.Files = append(j.Files, Snapshot{c.Path, c.Kind, c.Before, c.After, c.Existed, uint32(c.Mode), c.BeforeKey, c.KeyExisted, c.EnvExisted})
		for dir := filepath.Dir(c.Path); dir != o.Project; dir = filepath.Dir(dir) {
			if _, e := os.Stat(dir); errors.Is(e, os.ErrNotExist) {
				dirs[dir] = true
			} else {
				break
			}
		}
	}
	for d := range dirs {
		j.CreatedDirs = append(j.CreatedDirs, d)
	}
	sort.Slice(j.CreatedDirs, func(i, k int) bool { return len(j.CreatedDirs[i]) > len(j.CreatedDirs[k]) })
	if err = saveJournal(o, j); err != nil {
		return r, err
	}
	for _, s := range j.Files {
		b, exists, _, e := readFile(s.Path)
		if e == nil && (exists != s.Existed || !bytes.Equal(b, s.Before)) {
			e = errors.New("arquivo alterado por outro processo")
		}
		if e == nil {
			e = atomicWrite(s.Path, s.After, os.FileMode(s.Mode))
		}
		if e == nil {
			actual, _, _, readErr := readFile(s.Path)
			if readErr != nil {
				e = readErr
			} else if !bytes.Equal(actual, s.After) {
				e = errors.New("validacao da escrita falhou")
			}
		}
		if e != nil {
			if rollbackErr := restore(o, &j, "rolled-back"); rollbackErr != nil {
				return r, fmt.Errorf("aplicacao interrompida; backup preservado. Execute undo. Causa: %v; recuperacao: %v", e, rollbackErr)
			}
			return r, fmt.Errorf("falha ao aplicar; mudancas revertidas: %w", e)
		}
		r.Changes = append(r.Changes, s.Path)
	}
	j.Status = "committed"
	if err = saveJournal(o, j); err != nil {
		return r, errors.New("arquivos aplicados, mas registro final falhou. Backup pendente preservado; use undo para recuperar")
	}
	r.ID = j.ID
	r.Status = "Alteracoes aplicadas e verificadas. Use /token-saver desfazer para reverter."
	return r, nil
}

type restoration struct {
	path    string
	before  []byte
	existed bool
	data    []byte
	remove  bool
	mode    os.FileMode
}

func prepareRestore(j Journal) ([]restoration, error) {
	changes := []restoration{}
	for _, s := range j.Files {
		current, exists, mode, err := readFile(s.Path)
		if err != nil {
			return nil, err
		}
		if exists == s.Existed && bytes.Equal(current, s.Before) {
			continue
		}
		if !exists {
			return nil, fmt.Errorf("conflito: %s foi removido depois da aplicacao", s.Path)
		}
		r := restoration{path: s.Path, before: current, existed: exists, mode: mode}
		if bytes.Equal(current, s.After) {
			r.data = s.Before
			r.remove = !s.Existed
			if s.Existed {
				r.mode = os.FileMode(s.Mode)
			}
			changes = append(changes, r)
			continue
		}
		if s.Kind != "tool-search" {
			return nil, fmt.Errorf("conflito: %s foi editado depois da aplicacao; original e backup preservados", s.Path)
		}
		m, err := object(current)
		if err != nil {
			return nil, errors.New("configuracao atual invalida; nada foi restaurado")
		}
		e, err := envMap(m)
		if err != nil {
			return nil, err
		}
		value, has := e["ENABLE_TOOL_SEARCH"]
		if has == s.KeyExisted && reflect.DeepEqual(value, s.BeforeKey) {
			continue
		}
		if value != "true" {
			return nil, errors.New("conflito: ENABLE_TOOL_SEARCH foi editado depois da aplicacao; nenhuma mudanca restaurada")
		}
		if s.KeyExisted {
			e["ENABLE_TOOL_SEARCH"] = s.BeforeKey
		} else {
			delete(e, "ENABLE_TOOL_SEARCH")
		}
		if len(e) == 0 && !s.EnvExisted {
			delete(m, "env")
		} else {
			m["env"] = e
		}
		r.data = encode(m)
		r.remove = !s.Existed && len(m) == 0
		changes = append(changes, r)
	}
	return changes, nil
}

func restore(o options, j *Journal, status string) error {
	changes, err := prepareRestore(*j)
	if err != nil {
		return err
	}
	// Validate every conflict before the first write. Recheck each file at its write.
	for _, r := range changes {
		if err := noLinks(r.path); err != nil {
			return err
		}
		b, exists, _, err := readFile(r.path)
		if err != nil {
			return err
		}
		if exists != r.existed || !bytes.Equal(b, r.before) {
			return errors.New("arquivo mudou durante a restauracao; execute undo novamente")
		}
		if r.remove {
			err = os.Remove(r.path)
		} else {
			err = atomicWrite(r.path, r.data, r.mode)
		}
		if err != nil {
			return err
		}
	}
	for _, d := range j.CreatedDirs {
		os.Remove(d)
	} // Only succeeds for empty directories created by this transaction.
	j.Status = status
	return saveJournal(o, *j)
}

func undoPreview(o options) (Result, error) {
	r := Result{Status: "Nenhuma aplicacao para desfazer.", Changes: []string{}, Savings: "nao medida"}
	j, err := latest(o)
	if err != nil {
		return r, err
	}
	if j.ID == "" {
		return r, nil
	}
	changes, err := prepareRestore(j)
	if err != nil {
		return r, err
	}
	r.ID = j.ID
	r.Status = "Previa: desfazer a ultima aplicacao, preservando outras edicoes."
	for _, c := range changes {
		r.Changes = append(r.Changes, c.path)
	}
	if len(changes) == 0 {
		r.Status = "Os arquivos ja correspondem ao estado anterior. Use undo --yes para fechar o registro."
	}
	return r, nil
}

func undo(o options, expected string) (Result, error) {
	r := Result{Status: "Alteracoes desfeitas.", Changes: []string{}, Savings: "nao medida"}
	release, err := lock(o)
	if err != nil {
		return r, err
	}
	defer release()
	j, err := latest(o)
	if err != nil {
		return r, err
	}
	if j.ID == "" {
		r.Status = "Nenhuma aplicacao para desfazer."
		return r, nil
	}
	if j.ID != expected {
		return r, errors.New("historico mudou; revise a nova previa")
	}
	changes, err := prepareRestore(j)
	if err != nil {
		return r, err
	}
	for _, c := range changes {
		r.Changes = append(r.Changes, c.path)
	}
	if err = restore(o, &j, "undone"); err != nil {
		return r, err
	}
	r.ID = j.ID
	return r, nil
}
