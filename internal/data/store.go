package data

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"cam/internal/crypto"
)

type Command struct {
	Cmd       string   `json:"cmd,omitempty"`
	Encrypted string   `json:"encrypted,omitempty"`
	IsPrivate bool     `json:"is_private"`
	Tags      []string `json:"tags"`
	Timestamp string   `json:"timestamp"`
}

type DataStore struct {
	Stacks map[string][]Command `json:"stacks"`
	path   string
	mu     sync.RWMutex
}

func NewDataStore() *DataStore {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "CRITICAL: Could not find user home directory: %v\n", err)
		os.Exit(1)
	}

	path := filepath.Join(home, ".config", "cam", "data.json")

	return &DataStore{
		Stacks: make(map[string][]Command),
		path:   path,
	}
}

func (ds *DataStore) LoadData(decryptPrivate bool) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.path == "" {
		return fmt.Errorf("could not determine config path")
	}

	data, err := os.ReadFile(ds.path)
	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}

	err = json.Unmarshal(data, &ds.Stacks)
	if err != nil {
		return fmt.Errorf("failed to parse data file: %w", err)
	}

	if ds.Stacks == nil {
		ds.Stacks = make(map[string][]Command)
	}

	if !decryptPrivate {
		for stackName, commands := range ds.Stacks {
			var publicCmds []Command
			for _, cmd := range commands {
				if !cmd.IsPrivate {
					publicCmds = append(publicCmds, cmd)
				}
			}
			if len(publicCmds) > 0 {
				ds.Stacks[stackName] = publicCmds
			} else {
				delete(ds.Stacks, stackName)
			}
		}
		return nil
	}

	configDir := filepath.Dir(ds.path)
	keysDir := filepath.Join(configDir, ".keys")
	privKeyPath := filepath.Join(keysDir, "private_key.pem")

	if _, err := os.Stat(privKeyPath); err == nil {
		for stackName, commands := range ds.Stacks {
			for i, cmd := range commands {
				if cmd.IsPrivate && cmd.Encrypted != "" {
					cipherBytes, err := base64.StdEncoding.DecodeString(cmd.Encrypted)
					if err == nil {
						plain, err := crypto.Decrypt(cipherBytes, privKeyPath)
						if err == nil {
							commands[i].Cmd = string(plain)
						} else {
							commands[i].Cmd = "[DECRYPTION FAILED]"
						}
					}
				}
			}
			ds.Stacks[stackName] = commands
		}
	}

	return nil
}

func (ds *DataStore) SaveData() error {
	ds.mu.RLock()

	saveStacks := make(map[string][]Command)

	configDir := filepath.Dir(ds.path)
	keysDir := filepath.Join(configDir, ".keys")
	pubKeyPath := filepath.Join(keysDir, "public_key.pem")

	for k, v := range ds.Stacks {
		saveCmds := make([]Command, len(v))
		for i, c := range v {
			saveCmds[i] = c
			if c.IsPrivate {
				if c.Cmd != "" {
					cipherBytes, err := crypto.Encrypt([]byte(c.Cmd), pubKeyPath)
					if err == nil {
						saveCmds[i].Encrypted = base64.StdEncoding.EncodeToString(cipherBytes)
						saveCmds[i].Cmd = ""
					} else {
						fmt.Fprintf(os.Stderr, "Warning: Failed to encrypt command: %v\n", err)
					}
				}
			}
		}
		saveStacks[k] = saveCmds
	}

	data, err := json.MarshalIndent(saveStacks, "", "  ")
	ds.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	dir := filepath.Dir(ds.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(ds.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
}

func (ds *DataStore) AddCommand(stackName string, cmdStr string, tags []string, isPrivate bool) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if isPrivate {
		configDir := filepath.Dir(ds.path)
		keysDir := filepath.Join(configDir, ".keys")
		if err := os.MkdirAll(keysDir, 0700); err != nil {
			return fmt.Errorf("failed to create keys directory: %w", err)
		}
		if err := crypto.EnsureKeysExists(keysDir); err != nil {
			return fmt.Errorf("failed to ensure encryption keys: %w", err)
		}
	}

	newCmd := Command{
		Cmd:       cmdStr,
		Tags:      tags,
		IsPrivate: isPrivate,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	currentStack := ds.Stacks[stackName]
	ds.Stacks[stackName] = append([]Command{newCmd}, currentStack...)
	return nil
}

func (ds *DataStore) GetStack(stackName string) []Command {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	stack, exists := ds.Stacks[stackName]
	if !exists {
		return nil
	}
	result := make([]Command, len(stack))
	copy(result, stack)
	return result
}

func (ds *DataStore) RemoveCommand(stackName string, index int) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	stack, exists := ds.Stacks[stackName]
	if !exists {
		return fmt.Errorf("stack '%s' does not exist", stackName)
	}

	if index < 0 || index >= len(stack) {
		return fmt.Errorf("index %d out of bounds for stack '%s'", index, stackName)
	}

	ds.Stacks[stackName] = append(stack[:index], stack[index+1:]...)

	return nil
}

func (ds *DataStore) RemoveStack(stackName string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, exists := ds.Stacks[stackName]; !exists {
		return fmt.Errorf("stack '%s' does not exist", stackName)
	}

	delete(ds.Stacks, stackName)
	return nil
}

func (ds *DataStore) Swap(stackName string, i, j int) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	stack, exists := ds.Stacks[stackName]
	if !exists {
		return fmt.Errorf("stack '%s' does not exist", stackName)
	}

	if i < 0 || i >= len(stack) {
		return fmt.Errorf("index %d out of bounds for stack '%s'", i, stackName)
	}
	if j < 0 || j >= len(stack) {
		return fmt.Errorf("index %d out of bounds for stack '%s'", j, stackName)
	}

	stack[i], stack[j] = stack[j], stack[i]
	ds.Stacks[stackName] = stack

	return nil
}
