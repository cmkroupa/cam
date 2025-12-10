# cam (Command Action Manager)

A persistent, command manager for your console with **RSA encrypted private storage** and **integrated local Ollama AI responses**.

## Install

```bash
go build -o cam
mv cam /usr/local/bin/
```

## Usage

| Command | Description | Example |
| :--- | :--- | :--- |
| **`pin`** | Save command (`-p` for private) | `cam pin git "git commit"` |
| **`ls`** | List stacks (`-p` for private) | `cam ls git` |
| **`cp`** | Copy to clipboard | `cam cp git 1` |
| **`mv`** | Copy & remove (Cut) | `cam mv git 1` |
| **`swap`** | Swap two commands | `cam swap git 0 2` |
| **`rm`** | Delete a cmd / stack  | `cam rm git` |
| **`f`** | Fuzzy search (public cmds only) | `cam f commit` |
| **`ask`** | Ask your local AI a question | `cam ask "how to undo git commit"` |
| **`run`** | Run command from stack | `cam run git 0` |
| **`cmdr`** | Generate shell command | `cam cmdr "list files sorted by size"` |
| **`config`** | Configure settings | `cam config model llama3` |

**All Data Stored in :** `~/.config/cam/data.json`

### Private Commands

Use `cam pin -p` to encrypt a command. It will only be visible with `cam ls -p`.

**Note:** Currently private commands are not a secure method of storing sensitive information, as the RSA keys are stored in a file nearby. It is more of a preventative measure against snooping. Future upgrade will include a more secure method of storing private commands.

### AI Assistant (Ollama)

`cam` uses [Ollama](https://ollama.com) to generate commands and answer questions directly from your terminal.

**Setup:**

1. **Install Ollama:** Download from [ollama.com](https://ollama.com) and ensure it is running (`ollama serve`).
2. **Configure Model:** Tell `cam` which model to use.

    ```bash
    cam config model qwen2.5
    # or
    cam config model llama3
    ```

    *Find more models at [ollama.com/library](https://ollama.com/library).*

**Usage:**

**1. `cam ask` - Explanations & Help**
Best for general queries or when you need a concise explanation.

- **Default:** concise, readable markdown explanations.

  ```bash
  cam ask "how do I extract a tar.gz file"
  ```

- **One-liner (`-o`):** text-only short answer.
- **With Context (`-C`):** includes a tree view of the current directory (depth 2) to help the AI understand your file structure.

  ```bash
  cam ask -C "what is in this project?"
  ```

**2. `cam cmdr` - Pure Command Generation**
Best when you just want the command string to run immediately.

- **Auto Context:** Automatically reads the current directory (tree view) to provide accurate file/folder names.
- **Copy to Clipboard (`-c`):** Generates the command and copies it for you.

  ```bash
  # Asks how to move files, AI sees your actual folder names
  cam cmdr "move all images to the assets folder"
  
  # Copies result to clipboard
  cam cmdr -c "find all .go files"
  ```

## Roadmap

- [ ] **Session Storage**: Ability to save a session of commands.
- [ ] **Run**: Ability to run a command from a stack.
- [ ] **Chained Commands**: Ability to chain multiple commands from a stack.
