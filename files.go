package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const maxFile = 8 * 1024 * 1024

func hash(b []byte) string { v := sha256.Sum256(b); return hex.EncodeToString(v[:]) }

// Strict parsing prevents a merge from silently choosing between duplicate keys.
func object(data []byte) (map[string]any, error) {
	data = bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf})
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	v, err := readValue(d)
	if err != nil {
		return nil, err
	}
	if _, err = d.Token(); err != io.EOF {
		return nil, errors.New("conteudo JSON adicional")
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("a raiz JSON deve ser um objeto")
	}
	return m, nil
}

func readValue(d *json.Decoder) (any, error) {
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return t, nil
	}
	switch delim {
	case '{':
		m := map[string]any{}
		for d.More() {
			k, e := d.Token()
			if e != nil {
				return nil, e
			}
			key, ok := k.(string)
			if !ok {
				return nil, errors.New("chave JSON invalida")
			}
			if _, exists := m[key]; exists {
				return nil, errors.New("chave JSON duplicada")
			}
			v, e := readValue(d)
			if e != nil {
				return nil, e
			}
			m[key] = v
		}
		_, err = d.Token()
		return m, err
	case '[':
		v := []any{}
		for d.More() {
			x, e := readValue(d)
			if e != nil {
				return nil, e
			}
			v = append(v, x)
		}
		_, err = d.Token()
		return v, err
	default:
		return nil, errors.New("JSON invalido")
	}
}

func readFile(path string) ([]byte, bool, os.FileMode, error) {
	return readLimited(path, maxFile)
}

func readLimited(path string, limit int64) ([]byte, bool, os.FileMode, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, 0600, nil
	}
	if err != nil {
		return nil, false, 0, err
	}
	if !info.Mode().IsRegular() {
		return nil, false, 0, errors.New("arquivo especial ou link nao suportado: " + path)
	}
	if info.Size() > limit {
		return nil, false, 0, errors.New("arquivo acima do limite suportado: " + path)
	}
	b, err := os.ReadFile(path)
	return b, true, info.Mode().Perm(), err
}

func loadObject(path string) (map[string]any, error) {
	b, exists, _, err := readFile(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return map[string]any{}, nil
	}
	m, err := object(b)
	if err != nil {
		return nil, fmt.Errorf("JSON invalido em %s; original preservado", path)
	}
	return m, nil
}

func noLinks(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for {
		info, e := os.Lstat(path)
		if e == nil {
			redirect, checkErr := isRedirect(path, info)
			if checkErr != nil {
				return checkErr
			}
			if redirect {
				return errors.New("escrita recusada em caminho com link ou redirecionamento: " + path)
			}
		}
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			return e
		}
		parent := filepath.Dir(path)
		if parent == path {
			break
		}
		path = parent
	}
	return nil
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	if err := noLinks(path); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".token-saver-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err = f.Chmod(mode); err == nil {
		_, err = f.Write(data)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func encode(v any) []byte {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return append(b, '\n')
}

func envMap(m map[string]any) (map[string]any, error) {
	v, exists := m["env"]
	if !exists {
		return map[string]any{}, nil
	}
	e, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("env precisa ser um objeto; nenhuma mudanca aplicada")
	}
	return e, nil
}

func stringValue(v any) string { s, _ := v.(string); return s }
