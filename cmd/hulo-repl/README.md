# Hulo REPL

A modern, feature-rich interactive shell for the Hulo programming language.

## Features

-   🚀 **Intelligent Autocompletion**: Get smart suggestions for keywords, functions, and commands as you type.
-   🎨 **Syntax Highlighting**: Your code is highlighted in real-time for improved readability.
-   ⚙️ **Customizable Themes**: Personalize your shell with built-in themes, including `default`, `dark`, `light`, and `colorful`.
-   📝 **Persistent History**: Command history is automatically saved and loaded across sessions.
-   🔍 **Dry-Run Mode**: Safely test commands without actually executing them.
-   🛡️ **Highly Configurable**: Use command-line flags or a `config.yaml` file to tailor the REPL to your workflow.

## Installation

> Go 1.18 or newer is required.

You can build and install the REPL from the source code:

```bash
# Install directly to your $GOPATH/bin
go install ./cmd/repl

# Or, build a local executable
go build -o hulo-repl ./cmd/repl
```

## Usage

### Getting Started

```bash
# Start the REPL with default settings
hulo-repl

# Display all available options
hulo-repl --help

# Show version information
hulo-repl --version
```

### Command-Line Options

| Flag                | Shorthand | Description                               | Default   |
| ------------------- | --------- | ----------------------------------------- | --------- |
| `--config`          | `-c`      | Path to a custom configuration file.      | `""`      |
| `--dry-run`         | `-d`      | Run in dry-run mode (don't execute code). | `false`   |
| `--history`         | `-H`      | Path to the command history file.         | `""`      |
| `--max-suggestions` | `-m`      | Max number of suggestions to display.     | `10`      |
| `--theme`           | `-t`      | Set the color theme for the REPL.         | `default` |
| `--verbose`         | `-v`      | Enable verbose output.                    | `false`   |
| `--version`         | `-V`      | Show version information.                 | `false`   |

Available themes: `default`, `dark`, `light`, `colorful`.

## Configuration File

For a more permanent setup, you can create a `config.yaml` file to define your preferred settings. The REPL will use default settings if no config file is specified.

### Creating a Configuration File

Create a `config.yaml` with your desired settings. Here is an example with all available options:

```yaml
# config.yaml
dry_run: false
verbose: false
max_suggestions: 15
theme: "dark"
history_file: "~/.hulo_repl_history"
```

### Using a Configuration File

Launch the REPL with the `--config` flag to apply your settings:

```bash
hulo-repl --config config.yaml
```

If the specified file doesn't exist, a warning will be displayed, and the REPL will fall back to its default settings.

### Configuration Priority

Command-line flags always take precedence over the configuration file. For example, the following command will use the `light` theme, even if the config file specifies `theme: "dark"`.

```bash
hulo-repl --config config.yaml --theme light
```

## REPL Commands

The following built-in commands are available within the REPL:

| Command       | Description                     |
| ------------- | ------------------------------- |
| `help`        | Shows the help message.         |
| `exit` / `quit` | Exits the REPL session.         |
| `clear`       | Clears the terminal screen.     |
| `config`      | Displays the current configuration. |
| `version`     | Shows the REPL version.         |


## Examples

### Executing Hulo Code

```hulo
>>> let x: num = 42
>>> echo "Hello, World!"
>>> fn add(a: num, b: num) -> num { return a + b }
>>> add(10, 20)
```

### Using Dry-Run Mode

```bash
$ hulo-repl --dry-run
>>> echo "This will not be executed"
🔍 [DRY-RUN] Would execute: echo "This will not be executed"
```

### Using a Different Theme

```bash
# Start with the dark theme
hulo-repl --theme dark
```

## Hotkeys

-   `Tab` - Trigger autocompletion.
-   `Ctrl+C` - Gracefully exit the REPL.
-   `Ctrl+L` - Clear the screen (on most terminals).

## Development

### Building from Source

```bash
# Standard development build
go build -o hulo-repl ./cmd/repl

# Optimized release build
go build -ldflags="-s -w" -o hulo-repl ./cmd/repl
```

### Running Tests

```bash
# Run all tests for the REPL
go test ./cmd/repl/...
```

## Contributing

Contributions are welcome! Please feel free to submit an Issue or Pull Request.

## License

This project is licensed under the MIT License.
