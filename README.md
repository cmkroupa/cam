# cam (Command Action Manager)

A persistent, command manager for your console with **RSA encrypted private storage** and **integrated Gemini AI responses**.

## Install

```bash
go build -o cam
mv cam /usr/local/bin/
```

## Usage

| Command | Description | Example |
| :--- | :--- | :--- |
| **`push`** | Save command (`-p` for private) | `cam push git "git commit"` |
| **`ls`** | List stacks (`-p` for private) | `cam ls git` |
| **`cp`** | Copy to clipboard | `cam cp git 1` |
| **`mv`** | Copy & remove (Cut) | `cam mv git 1` |
| **`swap`** | Swap two commands | `cam swap git 0 2` |
| **`clear`** | Delete a cmd / stack  | `cam clear git` |
| **`f`** | Fuzzy search (public cmds only) | `cam f commit` |
| **`ask`** | Ask Gemini AI a question | `cam ask "how to undo git commit"` |
| **`config`** | Configure settings | `cam config api-key <KEY>` |

**All Data Stored in :** `~/.config/cam/data.json`

### Private Commands

Use `cam push -p` to encrypt a command. It will only be visible with `cam ls -p`.

**Note:** Currently private commands are not a secure method of storing sensitive information, as the RSA keys are stored in a file nearby. It is more of a preventative measure against snooping. Future upgrade will include a more secure method of storing private commands.

- **Encryption:** Uses RSA keys stored in `~/.config/cam/.keys/`.
- **Search:** Private commands are always excluded from fuzzy search (`cam f`).

### AI Assistant (Gemini)

`cam` integrates with Google's Gemini models to answer questions extremely concisely from the command line.

**Setup:**
To use the `ask` command, you need a Google Gemini API key.

1. Get a key from Google AI Studio
**Note:** Gemini offers a generous free tier that is more than sufficient for daily personal usage.
2. Configure it in `cam config`:

   ```bash
   cam config api-key <YOUR_API_KEY>
   ```

   - This will RSA encrypt your API key and store it in `~/.config/cam/.config.json`.

**Usage:**

- **Default:** concise, readable explanations.

  ```bash
  cam ask "how do I extract a tar.gz file"
  ```

- **One-liner (`-o`):** raw command string (great for copying).

  ```bash
  cam ask -o "how do I extract a tar.gz file"
  ```

## Roadmap

- [ ] **Higher Security Storage**
- [ ] **Session Storage**: Ability to save a session of commands.
- [ ] **Run**: Ability to run a command from a stack.
- [ ] **Chained Commands**: Ability to chain multiple commands from a stack.
- [ ] **Local LLM Support**
