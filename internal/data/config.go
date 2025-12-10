package data

import (
	"cam/internal/crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	GeminiAPIKey string `json:"gemini_api_key"`
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

func (cs *ConfigStore) SetAPIKey(key string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	keysDir := filepath.Join(home, ".config", "cam", ".keys")
	if err := crypto.EnsureKeysExists(keysDir); err != nil {
		return fmt.Errorf("failed to ensure keys exist: %w", err)
	}

	pubKeyPath := filepath.Join(keysDir, "public_key.pem")
	encryptedBytes, err := crypto.Encrypt([]byte(key), pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to encrypt API key: %w", err)
	}

	encryptedKey := base64.StdEncoding.EncodeToString(encryptedBytes)

	cs.mu.Lock()
	cs.Config.GeminiAPIKey = encryptedKey
	cs.mu.Unlock()
	return cs.SaveConfig()
}

func (cs *ConfigStore) GetAPIKey() string {
	cs.mu.RLock()
	val := cs.Config.GeminiAPIKey
	cs.mu.RUnlock()

	if val == "" {
		return ""
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return val
	}
	privKeyPath := filepath.Join(home, ".config", "cam", ".keys", "private_key.pem")

	ciphertext, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return val
	}

	plaintext, err := crypto.Decrypt(ciphertext, privKeyPath)
	if err != nil {
		return val
	}

	return string(plaintext)
}
