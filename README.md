# Goodle CLI

A CLI to access Innopolis Moodle from your terminal 
with a nice TUI.

Install SHELL completions:

<details>
<summary>BASH</summary>

This script depends on the `bash-completion` package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

```bash
source <(goodle-cli completion bash)
```

To load completions for every new session, execute once:

#### Linux:

```bash
goodle-cli completion bash > /etc/bash_completion.d/goodle-cli
```

#### macOS:

```bash
goodle-cli completion bash > $(brew --prefix)/etc/bash_completion.d/goodle-cli
```

You will need to start a new shell for this setup to take effect.

</details>

<details>
<summary>Fish</summary>

To load completions in your current shell session:

```shell
goodle-cli completion fish | source
```

To load completions for every new session, execute once:

```shell
goodle-cli completion fish > ~/.config/fish/completions/goodle-cli.fish
```

You will need to start a new shell for this setup to take effect.
</details>

<details>
<summary>ZSH</summary>

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

```zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions in your current shell session:

```zsh
source <(goodle-cli completion zsh)
```

To load completions for every new session, execute once:

#### Linux:

```zsh
goodle-cli completion zsh > "${fpath[1]}/_goodle-cli"
```

#### macOS:

```zsh
goodle-cli completion zsh > $(brew --prefix)/share/zsh/site-functions/_goodle-cli
```

You will need to start a new shell for this setup to take effect.
</details>

<details>
<summary>Powershell</summary>

To load completions in your current shell session:

```powershell
goodle-cli completion powershell | Out-String | Invoke-Expression
```
To load completions for every new session, add the output of the above command
to your powershell profile.
</details>

## Build

You will need [Go compiler](https://go.dev/dl/)
and [just](https://github.com/casey/just)

To compile the app *just* run this command

```bash
just build
```

To *just* run it

```bash
just run
```

To show available recipes

```bash
just --list
```