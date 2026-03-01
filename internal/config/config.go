package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	HostLang     string            `yaml:"host_lang"`
	TargetLang   string            `yaml:"target_lang"`
	Dictionaries map[string]string `yaml:"dictionaries"`
}

var (
	getenv        = os.Getenv
	userConfigDir = os.UserConfigDir
	readFile      = os.ReadFile
	writeFile     = os.WriteFile
	mkdirAll      = os.MkdirAll
)

func Load() (*Settings, error) {
	settings := &Settings{
		HostLang:     "en",
		TargetLang:   "fr",
		Dictionaries: make(map[string]string),
	}

	configDir, err := userConfigDir()
	if err == nil {
		configPath := filepath.Join(configDir, "voc", "settings.yaml")
		data, err := readFile(configPath)
		if err == nil {
			var fileSettings Settings
			if err := yaml.Unmarshal(data, &fileSettings); err == nil {
				if fileSettings.HostLang != "" {
					settings.HostLang = fileSettings.HostLang
				}
				if fileSettings.TargetLang != "" {
					settings.TargetLang = fileSettings.TargetLang
				}
				if fileSettings.Dictionaries != nil {
					for k, v := range fileSettings.Dictionaries {
						settings.Dictionaries[k] = v
					}
				}
			}
		}
	}

	if lang := getenv("VOC_HOST_LANG"); lang != "" {
		settings.HostLang = lang
	}
	if lang := getenv("VOC_TARGET_LANG"); lang != "" {
		settings.TargetLang = lang
	}

	return settings, nil
}

func (s *Settings) Save() error {
	configDir, err := userConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "voc", "settings.yaml")
	if err := mkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return writeFile(configPath, data, 0644)
}
