package util

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultDirPerm  = 0o755
	defaultFilePerm = 0o644
)

// Exists reports whether the given path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !errors.Is(err, os.ErrNotExist)
}

// FileExists reports whether the given path exists and is a file.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists reports whether the given path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir creates the directory if it does not already exist.
func EnsureDir(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, defaultDirPerm)
}

// EnsureParentDir creates the parent directory of a file path if needed.
func EnsureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return EnsureDir(dir)
}

// ReadBytes reads the file content as bytes.
func ReadBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadString reads the file content as a string.
func ReadString(path string) (string, error) {
	data, err := ReadBytes(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadJSON reads a JSON file into out.
func ReadJSON(path string, out any) error {
	data, err := ReadBytes(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

// ReadYAML reads a YAML file into out.
func ReadYAML(path string, out any) error {
	data, err := ReadBytes(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

// WriteBytes writes bytes to a file and creates parent directories automatically.
func WriteBytes(path string, data []byte) error {
	if err := EnsureParentDir(path); err != nil {
		return err
	}
	return os.WriteFile(path, data, defaultFilePerm)
}

// WriteString writes a string to a file and creates parent directories automatically.
func WriteString(path string, content string) error {
	return WriteBytes(path, []byte(content))
}

// WriteJSON writes value to a JSON file using indented formatting.
func WriteJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return WriteBytes(path, data)
}

// WriteYAML writes value to a YAML file.
func WriteYAML(path string, value any) error {
	data, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	return WriteBytes(path, data)
}
