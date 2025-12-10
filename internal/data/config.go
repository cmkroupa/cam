package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	OllamaModel string `json:"ollama_model"` // e.g. "llama3", "qwen2.5:1.5b"
}

type ConfigStore struct {
	Config Config `json:"config"`
	path   string
	mu     sync.RWMutex
}

func NewConfigStore() *ConfigStore {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := filepath.Join(home, ".config", "cam", "config.json")
	return &ConfigStore{
		path: path,
	}
}

func (cs *ConfigStore) LoadConfig() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	data, err := os.ReadFile(cs.path)
	if os.IsNotExist(err) {
		cs.Config = Config{}
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return json.Unmarshal(data, &cs.Config)
}

func (cs *ConfigStore) SaveConfig() error {
	cs.mu.RLock()
	data, err := json.MarshalIndent(cs.Config, "", "  ")
	cs.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(cs.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(cs.path, data, 0644)
}

func (cs *ConfigStore) SetOllamaModel(model string) error {
	cs.mu.Lock()
	cs.Config.OllamaModel = model
	cs.mu.Unlock()
	return cs.SaveConfig()
}

func (cs *ConfigStore) GetOllamaModel() string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	if cs.Config.OllamaModel == "" {
		return "qwen2.5-coder:7b"
	}
	return cs.Config.OllamaModel
}
