// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/hulo-lang/hulo/cmd/hulo-repl/exec"
	"github.com/hulo-lang/hulo/cmd/hulo-repl/suggest"
	"github.com/hulo-lang/hulo/cmd/hulo-repl/theme"
	"github.com/hulo-lang/hulo/cmd/meta"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/opencommand/tinge"
	"github.com/spf13/cobra"
)

type replParameters struct {
	DryRun         bool
	Verbose        bool
	Config         string
	MaxSuggestions int
	HistoryFile    string
	Theme          string
	ShowVersion    bool
}

var replParams replParameters

var rootCmd = &cobra.Command{
	Use:   "hulo-repl",
	Short: "Hulo REPL - Interactive Hulo language shell",
	Long: `Hulo REPL is an interactive shell for the Hulo programming language.
It provides a modern, feature-rich environment for testing and experimenting with Hulo code.

Features:
  • Syntax highlighting and autocompletion
  • Command history
  • Built-in help system
  • Cross-platform support
  • Configurable themes and settings`,
	Run: runREPL,
}

func init() {
	rootCmd.Flags().BoolVarP(&replParams.DryRun, "dry-run", "d", false, "Run in dry-run mode (don't execute commands)")
	rootCmd.Flags().BoolVarP(&replParams.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&replParams.Config, "config", "c", "", "Path to configuration file (optional)")
	rootCmd.Flags().IntVarP(&replParams.MaxSuggestions, "max-suggestions", "m", 10, "Maximum number of suggestions to show")
	rootCmd.Flags().StringVarP(&replParams.HistoryFile, "history", "H", "", "Path to history file")
	rootCmd.Flags().StringVarP(&replParams.Theme, "theme", "t", "default", "Color theme (default, dark, light, colorful)")
	rootCmd.Flags().BoolVarP(&replParams.ShowVersion, "version", "V", false, "Show version information")

	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)
}

func runREPL(cmd *cobra.Command, args []string) {
	// 显示版本信息
	if replParams.ShowVersion {
		fmt.Println(tinge.Styled().
			With(tinge.Bold, tinge.Green).
			Text("Hulo REPL v1.0.0").
			Newline().
			String())
		return
	}

	// 显示配置信息
	if replParams.Verbose {
		fmt.Println(tinge.Styled().
			With(tinge.Bold, tinge.Blue).
			Text("Configuration:").
			Newline().
			String())
		fmt.Printf("  Dry Run: %v\n", replParams.DryRun)
		fmt.Printf("  Verbose: %v\n", replParams.Verbose)
		fmt.Printf("  Config File: %s\n", replParams.Config)
		fmt.Printf("  Max Suggestions: %d\n", replParams.MaxSuggestions)
		fmt.Printf("  History File: %s\n", replParams.HistoryFile)
		fmt.Printf("  Theme: %s\n", replParams.Theme)
		fmt.Println()
	}

	// 设置signal处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Print("\r\033[K")
		printGoodbye()
		os.Exit(0)
	}()

	// 显示启动信息
	fmt.Println(tinge.Styled().
		Newline().
		Space(2).
		With(tinge.Bold, tinge.Italic, tinge.Green).
		Text("Hulo-REPL").
		Space().
		Green(meta.Version).
		Newline().
		String())

	if replParams.DryRun {
		fmt.Println(tinge.Styled().
			Space(2).
			With(tinge.Bold, tinge.Yellow).
			Text("DRY-RUN MODE").
			Space().
			With(tinge.Italic).
			Text("(Commands will not be executed)").
			Newline().
			String())
	}

	fmt.Println(tinge.Styled().
		Space(2).
		Green("➜").
		Space(2).
		Grey("Type").
		Space().
		Bold("help").
		Space().
		Grey("for commands,").
		Space().
		Bold("exit").
		Space().
		Grey("to quit").
		Newline().
		String())

	cfgPath := config.HuloReplFileName

	if cmd.Flags().Changed("config") {
		cfgPath = replParams.Config
	}

	cfg, err := util.LoadConfigure[config.HuloRepl](cfgPath)
	if err != nil {
		panic(err)
	}

	options := []prompt.Option{
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("Hulo REPL"),
		prompt.OptionMaxSuggestion(cfg.MaxSuggestions),
	}

	for name, file := range cfg.ThemeFiles {
		if cfg.Theme == name {
			themeFile, err := util.LoadConfigure[config.PromptTheme](file)
			if err != nil {
				panic(err)
			}
			options = append(options, theme.Parse(themeFile)...)
		}
	}

	p := prompt.New(
		executor,
		completer,
		options...,
	)

	p.Run()
}

var executors = []exec.Executor{
	&exec.ExitExecutor{},
	&exec.HelpExecutor{},
	&exec.ClearExecutor{},
	&exec.ConfigExecutor{},
	&exec.VersionExecutor{},
}

func executor(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}

	for _, executor := range executors {
		if executor.CanHandle(in) {
			executor.Execute(in)
			return
		}
	}

	// 根据dry-run模式处理命令
	if replParams.DryRun {
		fmt.Printf(tinge.Styled().
			With(tinge.Bold, tinge.Yellow).
			Text("🔍 [DRY-RUN] Would execute: %s").
			Newline().
			String(), in)
	} else {
		fmt.Printf("Executing: %s\n", in)
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	word := d.GetWordBeforeCursor()

	// 如果当前行为空，返回所有关键字
	if word == "" {
		return suggest.KeyWords
	}

	// 过滤匹配的关键字
	suggestions := prompt.FilterHasPrefix(suggest.KeyWords, word, true)

	// 如果输入以特定字符开头，添加相应的建议
	if strings.HasPrefix(word, "import") {
		suggestions = append(suggestions, prompt.Suggest{Text: "import * from \"std\"", Description: "Import all from standard library"})
		suggestions = append(suggestions, prompt.Suggest{Text: "import { echo } from \"std\"", Description: "Import specific function"})
	}

	if strings.HasPrefix(word, "let") || strings.HasPrefix(word, "var") || strings.HasPrefix(word, "const") {
		suggestions = append(suggestions, prompt.Suggest{Text: "let x: num = 42", Description: "Declare a number variable"})
		suggestions = append(suggestions, prompt.Suggest{Text: "let name: str = \"Hulo\"", Description: "Declare a string variable"})
		suggestions = append(suggestions, prompt.Suggest{Text: "let flag: bool = true", Description: "Declare a boolean variable"})
	}

	if strings.HasPrefix(word, "fn") {
		suggestions = append(suggestions, prompt.Suggest{Text: "fn hello() { echo \"Hello, World!\" }", Description: "Define a simple function"})
		suggestions = append(suggestions, prompt.Suggest{Text: "fn add(a: num, b: num) -> num { return a + b }", Description: "Define a function with parameters and return type"})
	}

	if strings.HasPrefix(word, "if") {
		suggestions = append(suggestions, prompt.Suggest{Text: "if true { echo \"condition is true\" }", Description: "Simple if statement"})
		suggestions = append(suggestions, prompt.Suggest{Text: "if x > 10 { echo \"x is greater than 10\" } else { echo \"x is 10 or less\" }", Description: "If-else statement"})
	}

	if strings.HasPrefix(word, "for") {
		suggestions = append(suggestions, prompt.Suggest{Text: "for i in range(10) { echo i }", Description: "For loop with range"})
		suggestions = append(suggestions, prompt.Suggest{Text: "for item in list { echo item }", Description: "For loop with list"})
	}

	// 添加命令建议
	commandSuggestions := prompt.FilterHasPrefix(suggest.Commands, word, true)
	suggestions = append(suggestions, commandSuggestions...)

	return suggestions
}

func printGoodbye() {
	fmt.Print("\r\033[K")

	fmt.Println(tinge.Styled().
		With(tinge.Bold, tinge.Green).
		Text("✨ Goodbye! See you next time!").
		Newline().
		String())
}

func printConfig() {
	fmt.Println(tinge.Styled().
		With(tinge.Bold, tinge.Blue).
		Text("⚙️  Current Configuration").
		Newline().
		String())

	fmt.Printf("  Dry Run: %v\n", replParams.DryRun)
	fmt.Printf("  Verbose: %v\n", replParams.Verbose)
	fmt.Printf("  Config File: %s\n", replParams.Config)
	fmt.Printf("  Max Suggestions: %d\n", replParams.MaxSuggestions)
	fmt.Printf("  History File: %s\n", replParams.HistoryFile)
	fmt.Printf("  Theme: %s\n", replParams.Theme)
	fmt.Println()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
