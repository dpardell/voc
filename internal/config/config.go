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

func Load() (*Settings, error) {
	settings := &Settings{
		HostLang:   "en",
		TargetLang: "fr",
		Dictionaries: map[string]string{
			"fr": "https://kaikki.org/dictionary/French/kaikki.org-dictionary-French.jsonl",
			"sk": "https://kaikki.org/dictionary/Slovak/kaikki.org-dictionary-Slovak.jsonl",
		},
	}

	// 1. Load from environment variables (lowest priority)
	if lang := os.Getenv("VOC_HOST_LANG"); lang != "" {
		settings.HostLang = lang
	}
	if lang := os.Getenv("VOC_TARGET_LANG"); lang != "" {
		settings.TargetLang = lang
	}

	// 2. Load from settings file
	configDir, err := os.UserConfigDir()
	if err == nil {
		configPath := filepath.Join(configDir, "voc", "settings.yaml")
		data, err := os.ReadFile(configPath)
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

	return settings, nil
}

func (s *Settings) Save() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "voc", "settings.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
