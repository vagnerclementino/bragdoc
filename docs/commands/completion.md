# Completion

## Available Commands

```bash
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh
```

### bash

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(bragdoc completion bash)

To load completions for every new session, execute once:

#### Linux:

	bragdoc completion bash > /etc/bash_completion.d/bragdoc

#### macOS:

	bragdoc completion bash > $(brew --prefix)/etc/bash_completion.d/bragdoc

You will need to start a new shell for this setup to take effect.

```bash
Usage:
  bragdoc completion bash

Flags:
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### fish

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	bragdoc completion fish | source

To load completions for every new session, execute once:

	bragdoc completion fish > ~/.config/fish/completions/bragdoc.fish

You will need to start a new shell for this setup to take effect.

```bash
Usage:
  bragdoc completion fish [flags]

Flags:
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### powershell

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	bragdoc completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.

```bash
Usage:
  bragdoc completion powershell [flags]

Flags:
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### zsh

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(bragdoc completion zsh)

To load completions for every new session, execute once:

#### Linux:

	bragdoc completion zsh > "${fpath[1]}/_bragdoc"

#### macOS:

	bragdoc completion zsh > $(brew --prefix)/share/zsh/site-functions/_bragdoc

You will need to start a new shell for this setup to take effect.

```bash
Usage:
  bragdoc completion zsh [flags]

Flags:
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```
