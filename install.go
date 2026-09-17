package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed skill
var skillFiles embed.FS

type installManifest struct {
	Product string            `json:"product"`
	Version string            `json:"version"`
	Files   map[string]string `json:"files"`
}

func install(o options, out io.Writer) error {
	dest := filepath.Join(o.Config, "skills", "token-saver")
	if err := noLinks(dest); err != nil {
		return err
	}
	if _, err := os.Stat(dest); err == nil {
		if err = verifyOwned(dest); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if o.Dry {
		fmt.Fprintf(out, "Previa: instalar a skill em %s. Nenhuma configuracao sera alterada.\n", dest)
		return nil
	}
	parent := filepath.Dir(dest)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return err
	}
	lockPath := filepath.Join(parent, ".token-saver-install.lock")
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New("outra instalacao pode estar em andamento")
	}
	lock.Close()
	defer os.Remove(lockPath)
	stage, err := os.MkdirTemp(parent, ".token-saver-install-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage) // stage is exclusively created by this operation under the verified skills directory.
	manifest := installManifest{Product: "token-saver", Version: version, Files: map[string]string{}}
	err = fs.WalkDir(skillFiles, "skill", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel := strings.TrimPrefix(path, "skill/")
		b, e := skillFiles.ReadFile(path)
		if e != nil {
			return e
		}
		if e = atomicWrite(filepath.Join(stage, filepath.FromSlash(rel)), b, 0600); e != nil {
			return e
		}
		manifest.Files[rel] = hash(b)
		return nil
	})
	if err != nil {
		return err
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	binary, err := os.ReadFile(exe)
	if err != nil {
		return err
	}
	binaryName := "bin/token-saver"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	if err = atomicWrite(filepath.Join(stage, filepath.FromSlash(binaryName)), binary, 0700); err != nil {
		return err
	}
	manifest.Files[binaryName] = hash(binary)
	if err = atomicWrite(filepath.Join(stage, ".token-saver-install.json"), encode(manifest), 0600); err != nil {
		return err
	}
	backup := ""
	if _, err = os.Stat(dest); err == nil {
		if err = verifyOwned(dest); err != nil {
			return err
		}
		backup = filepath.Join(o.Config, "token-saver", "install-backups", newID())
		if err = noLinks(backup); err != nil {
			return err
		}
		if err = os.MkdirAll(filepath.Dir(backup), 0700); err != nil {
			return err
		}
		if err = os.Rename(dest, backup); err != nil {
			return errors.New("nao foi possivel preparar a atualizacao; feche processos usando a skill e tente novamente")
		}
	}
	if err = os.Rename(stage, dest); err != nil {
		if backup != "" {
			if restoreErr := os.Rename(backup, dest); restoreErr != nil {
				return fmt.Errorf("instalacao falhou; versao anterior preservada em %s", backup)
			}
		}
		return err
	}
	if err = verifyOwned(dest); err != nil {
		return err
	}
	fmt.Fprintf(out, "Token Saver %s instalado e verificado.\nAbra uma nova sessao do Claude Code e digite /token-saver.\nPasta: %s\nNenhuma configuracao de modelo ou projeto foi alterada.\n", version, dest)
	if backup != "" {
		fmt.Fprintf(out, "Instalacao anterior preservada em: %s\n", backup)
	}
	return nil
}

func verifyOwned(dir string) error {
	b, _, _, err := readFile(filepath.Join(dir, ".token-saver-install.json"))
	if err != nil {
		return err
	}
	var m installManifest
	if json.Unmarshal(b, &m) != nil || m.Product != "token-saver" || len(m.Files) == 0 {
		return errors.New("ja existe uma skill token-saver sem registro desta versao. Preserve ou renomeie essa pasta antes de instalar")
	}
	seen := 0
	err = filepath.WalkDir(dir, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if err := noLinks(p); err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return errors.New("link inesperado na instalacao")
		}
		if d.IsDir() {
			return nil
		}
		rel, e := filepath.Rel(dir, p)
		if e != nil {
			return e
		}
		rel = filepath.ToSlash(rel)
		if rel == ".token-saver-install.json" {
			return nil
		}
		expected, ok := m.Files[rel]
		if !ok {
			return errors.New("arquivos personalizados encontrados; instalacao existente preservada")
		}
		data, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		if hash(data) != expected {
			return errors.New("a skill foi editada localmente; instalacao existente preservada")
		}
		seen++
		return nil
	})
	if err != nil {
		return err
	}
	if seen != len(m.Files) {
		return errors.New("instalacao incompleta; arquivos existentes preservados")
	}
	return nil
}
