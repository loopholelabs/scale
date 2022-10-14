package scalefile

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Build struct {
	Language     string       `json:"language"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

type ScaleFile struct {
	Name  string `json:"name"`
	Build Build  `json:"build"`
	File  string `json:"file"`
}

func Read(path string) (ScaleFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return ScaleFile{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	return decode(file)
}

func decode(data io.Reader) (ScaleFile, error) {
	decoder := yaml.NewDecoder(data)
	manifest := ScaleFile{}
	err := decoder.Decode(&manifest)
	if err != nil {
		return ScaleFile{}, err
	}

	return manifest, nil
}

func Write(path string, scalefile ScaleFile) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return encode(file, scalefile)
}

func encode(data io.Writer, scalefile ScaleFile) error {
	encoder := yaml.NewEncoder(data)
	encoder.SetIndent(2)
	return encoder.Encode(scalefile)
}
