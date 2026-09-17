package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const directRule = "# Token Saver — Modo Direto\n\nResponda de forma objetiva, sem cortesias repetitivas ou repetir a pergunta.\nPreserve explicacoes necessarias, riscos relevantes, evidencias e codigo completo.\nNao reduza a investigacao, os testes ou o raciocinio necessario para concluir a tarefa.\nSiga o idioma e o nivel de detalhe pedidos pelo usuario.\n"

type Change struct {
	Path        string      `json:"path"`
	Kind        string      `json:"kind"`
	Description string      `json:"description"`
	Before      []byte      `json:"-"`
	After       []byte      `json:"-"`
	Existed     bool        `json:"-"`
	Mode        os.FileMode `json:"-"`
	BeforeKey   any         `json:"-"`
	KeyExisted  bool        `json:"-"`
	EnvExisted  bool        `json:"-"`
}
type Plan struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
	Notes   []string `json:"notes"`
	Savings string   `json:"savings"`
}

func (p Plan) human() string {
	if len(p.Changes) == 0 {
		return "Nenhuma mudanca automatica necessaria ou compativel.\nUse audit para revisar os achados. Economia: nao medida."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Previa %s\n", p.ID[:12])
	for _, c := range p.Changes {
		fmt.Fprintf(&b, "- %s\n  Arquivo: %s\n", c.Description, c.Path)
	}
	for _, n := range p.Notes {
		fmt.Fprintln(&b, n)
	}
	fmt.Fprintln(&b, "Economia: nao medida. As mudancas podem ser desfeitas.")
	return strings.TrimSpace(b.String())
}

func makePlan(o options, a Audit) (Plan, error) {
	p := Plan{Changes: []Change{}, Notes: []string{"Modelo, raciocinio, permissoes e CLAUDE.md serao preservados."}, Savings: "nao medida"}
	if a.eligible {
		path := filepath.Join(o.Project, ".claude", "settings.local.json")
		before, exists, mode, err := readFile(path)
		if err != nil {
			return p, err
		}
		m := map[string]any{}
		if exists {
			m, err = object(before)
			if err != nil {
				return p, err
			}
		}
		e, err := envMap(m)
		if err != nil {
			return p, err
		}
		old, had := e["ENABLE_TOOL_SEARCH"]
		_, envHad := m["env"]
		e["ENABLE_TOOL_SEARCH"] = "true"
		m["env"] = e
		p.Changes = append(p.Changes, Change{Path: path, Kind: "tool-search", Description: "Carregar definicoes de ferramentas sob demanda neste projeto.", Before: before, After: encode(m), Existed: exists, Mode: mode, BeforeKey: old, KeyExisted: had, EnvExisted: envHad})
	}
	if o.Concise {
		path := filepath.Join(o.Project, ".claude", "rules", "token-saver.md")
		before, exists, mode, err := readFile(path)
		if err != nil {
			return p, err
		}
		if exists && string(before) != directRule {
			return p, fmt.Errorf("ja existe uma regra personalizada em %s; original preservado", path)
		}
		if !exists {
			p.Changes = append(p.Changes, Change{Path: path, Kind: "concise", Description: "Modo Direto: respostas objetivas sem cortar codigo, evidencias ou testes.", Before: before, After: []byte(directRule), Existed: exists, Mode: mode})
			p.Notes = append(p.Notes, "Modo Direto e opcional: a regra ocupa contexto e nao garante economia.")
		}
	}
	identity := []string{a.fingerprint}
	for _, c := range p.Changes {
		if err := noLinks(c.Path); err != nil {
			return p, err
		}
		identity = append(identity, c.Path, c.Kind, hash(c.Before), hash(c.After))
	}
	p.ID = hash([]byte(strings.Join(identity, "\x00")))
	return p, nil
}
