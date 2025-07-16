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
	"github.com/opencommand/tinge"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Configure struct {
	DryRun         bool
	Verbose        bool
	ConfigFile     string
	MaxSuggestions int
	HistoryFile    string
	Theme          string
	ShowVersion    bool
}

func getDefaultConfig() Configure {
	return Configure{
		DryRun:         false,
		Verbose:        false,
		MaxSuggestions: 10,
		HistoryFile:    "",
		Theme:          "default",
		ShowVersion:    false,
	}
}

func loadConfig(configPath string) (Configure, error) {
	config := getDefaultConfig()

	if configPath == "" {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, fmt.Errorf("failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

var config Configure

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
	rootCmd.Flags().BoolVarP(&config.DryRun, "dry-run", "d", false, "Run in dry-run mode (don't execute commands)")
	rootCmd.Flags().BoolVarP(&config.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&config.ConfigFile, "config", "c", "", "Path to configuration file (optional)")
	rootCmd.Flags().IntVarP(&config.MaxSuggestions, "max-suggestions", "m", 10, "Maximum number of suggestions to show")
	rootCmd.Flags().StringVarP(&config.HistoryFile, "history", "H", "", "Path to history file")
	rootCmd.Flags().StringVarP(&config.Theme, "theme", "t", "default", "Color theme (default, dark, light, colorful)")
	rootCmd.Flags().BoolVarP(&config.ShowVersion, "version", "V", false, "Show version information")

	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)
}

func runREPL(cmd *cobra.Command, args []string) {
	loadedConfig, err := loadConfig(config.ConfigFile)
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		fmt.Println("Using default configuration...")
		loadedConfig = getDefaultConfig()
	}

	if cmd.Flags().Changed("dry-run") {
		loadedConfig.DryRun = config.DryRun
	}
	if cmd.Flags().Changed("verbose") {
		loadedConfig.Verbose = config.Verbose
	}
	if cmd.Flags().Changed("max-suggestions") {
		loadedConfig.MaxSuggestions = config.MaxSuggestions
	}
	if cmd.Flags().Changed("history") {
		loadedConfig.HistoryFile = config.HistoryFile
	}
	if cmd.Flags().Changed("theme") {
		loadedConfig.Theme = config.Theme
	}
	loadedConfig.ShowVersion = config.ShowVersion

	// 更新全局配置
	config = loadedConfig

	// 显示版本信息
	if config.ShowVersion {
		fmt.Println(tinge.Styled().
			With(tinge.Bold, tinge.Green).
			Text("Hulo REPL v1.0.0").
			Newline().
			String())
		return
	}

	// 显示配置信息
	if config.Verbose {
		fmt.Println(tinge.Styled().
			With(tinge.Bold, tinge.Blue).
			Text("Configuration:").
			Newline().
			String())
		fmt.Printf("  Dry Run: %v\n", config.DryRun)
		fmt.Printf("  Verbose: %v\n", config.Verbose)
		fmt.Printf("  Config File: %s\n", config.ConfigFile)
		fmt.Printf("  Max Suggestions: %d\n", config.MaxSuggestions)
		fmt.Printf("  History File: %s\n", config.HistoryFile)
		fmt.Printf("  Theme: %s\n", config.Theme)
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
		Green("v1.0.0").
		Newline().
		String())

	// 显示模式信息
	if config.DryRun {
		fmt.Println(tinge.Styled().
			Space(2).
			With(tinge.Bold, tinge.Yellow).
			Text("🔍 DRY-RUN MODE").
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

	// 创建prompt选项
	promptOptions := []prompt.Option{
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("Hulo REPL"),
		prompt.OptionMaxSuggestion(uint16(config.MaxSuggestions)),
	}

	// 根据主题设置颜色
	switch config.Theme {
	case "dark":
		promptOptions = append(promptOptions,
			prompt.OptionPrefixTextColor(prompt.White),
			prompt.OptionPreviewSuggestionTextColor(prompt.Cyan),
			prompt.OptionSelectedSuggestionTextColor(prompt.White),
			prompt.OptionSuggestionBGColor(prompt.DarkGray),
			prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
			prompt.OptionDescriptionBGColor(prompt.DarkGray),
			prompt.OptionDescriptionTextColor(prompt.LightGray),
		)
	case "light":
		promptOptions = append(promptOptions,
			prompt.OptionPrefixTextColor(prompt.Black),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionTextColor(prompt.White),
			prompt.OptionSuggestionBGColor(prompt.LightGray),
			prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
			prompt.OptionDescriptionBGColor(prompt.LightGray),
			prompt.OptionDescriptionTextColor(prompt.Black),
		)
	case "colorful":
		promptOptions = append(promptOptions,
			prompt.OptionPrefixTextColor(prompt.Yellow),
			prompt.OptionPreviewSuggestionTextColor(prompt.Cyan),
			prompt.OptionSelectedSuggestionTextColor(prompt.Yellow),
			prompt.OptionSuggestionBGColor(prompt.DarkGray),
			prompt.OptionSelectedSuggestionBGColor(prompt.Red),
			prompt.OptionDescriptionBGColor(prompt.DarkGray),
			prompt.OptionDescriptionTextColor(prompt.White),
			prompt.OptionSelectedDescriptionTextColor(prompt.White),
			prompt.OptionSelectedDescriptionBGColor(prompt.Red),
			prompt.OptionScrollbarThumbColor(prompt.DarkGray),
			prompt.OptionScrollbarBGColor(prompt.LightGray),
		)
	default: // default theme
		promptOptions = append(promptOptions,
			prompt.OptionPrefixTextColor(prompt.Yellow),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionTextColor(prompt.Yellow),
			prompt.OptionSuggestionBGColor(prompt.DarkGray),
			prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
			prompt.OptionDescriptionBGColor(prompt.DarkGray),
			prompt.OptionDescriptionTextColor(prompt.White),
			prompt.OptionSelectedDescriptionTextColor(prompt.White),
			prompt.OptionSelectedDescriptionBGColor(prompt.Blue),
			prompt.OptionScrollbarThumbColor(prompt.DarkGray),
			prompt.OptionScrollbarBGColor(prompt.LightGray),
		)
	}

	p := prompt.New(
		executor,
		completer,
		promptOptions...,
	)

	p.Run()
}

func executor(in string) {
	// 处理特殊命令
	if strings.TrimSpace(in) == "" {
		return
	}

	// 检查退出命令
	if strings.ToLower(strings.TrimSpace(in)) == "exit" ||
		strings.ToLower(strings.TrimSpace(in)) == "quit" {
		fmt.Print("\r\033[K")
		printGoodbye()
		os.Exit(0)
	}

	// 检查帮助命令
	if strings.ToLower(strings.TrimSpace(in)) == "help" {
		printHelp()
		return
	}

	// 检查清屏命令
	if strings.ToLower(strings.TrimSpace(in)) == "clear" {
		fmt.Print("\033[H\033[2J")
		return
	}

	// 检查配置命令
	if strings.ToLower(strings.TrimSpace(in)) == "config" {
		printConfig()
		return
	}

	// 检查版本命令
	if strings.ToLower(strings.TrimSpace(in)) == "version" {
		fmt.Println(tinge.Styled().
			With(tinge.Bold, tinge.Green).
			Text("Hulo REPL v1.0.0").
			Newline().
			String())
		return
	}

	// 根据dry-run模式处理命令
	if config.DryRun {
		fmt.Printf(tinge.Styled().
			With(tinge.Bold, tinge.Yellow).
			Text("🔍 [DRY-RUN] Would execute: %s").
			Newline().
			String(), in)
	} else {
		fmt.Printf("Executing: %s\n", in)
	}

	// TODO: 集成Hulo编译器
	// result, err := hulo.Compile(in)
	// if err != nil {
	//     fmt.Printf("Error: %v\n", err)
	//     return
	// }
	// fmt.Printf("Result: %s\n", result)
}

func completer(d prompt.Document) []prompt.Suggest {
	word := d.GetWordBeforeCursor()

	// 如果当前行为空，返回所有关键字
	if word == "" {
		return keyWords
	}

	// 过滤匹配的关键字
	suggestions := prompt.FilterHasPrefix(keyWords, word, true)

	// 添加一些常用的命令
	commands := []prompt.Suggest{
		{Text: "help", Description: "Show help information"},
		{Text: "exit", Description: "Exit the REPL"},
		{Text: "quit", Description: "Exit the REPL"},
		{Text: "clear", Description: "Clear the screen"},
		{Text: "config", Description: "Show current configuration"},
		{Text: "version", Description: "Show version information"},
	}

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
	commandSuggestions := prompt.FilterHasPrefix(commands, word, true)
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

	fmt.Printf("  Dry Run: %v\n", config.DryRun)
	fmt.Printf("  Verbose: %v\n", config.Verbose)
	fmt.Printf("  Config File: %s\n", config.ConfigFile)
	fmt.Printf("  Max Suggestions: %d\n", config.MaxSuggestions)
	fmt.Printf("  History File: %s\n", config.HistoryFile)
	fmt.Printf("  Theme: %s\n", config.Theme)
	fmt.Println()
}

func printHelp() {
	fmt.Println(tinge.Styled().
		With(tinge.Bold, tinge.Blue).
		Text("🚀 Hulo REPL Help").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		With(tinge.Bold).
		Text("Commands:").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("help").
		Space(3).
		Text("Show this help message").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("exit").
		Space(3).
		Text("Exit the REPL").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("quit").
		Space(3).
		Text("Exit the REPL").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("clear").
		Space(3).
		Text("Clear the screen").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("config").
		Space(3).
		Text("Show current configuration").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("version").
		Space(3).
		Text("Show version information").
		Newline().
		String())

	fmt.Println()

	fmt.Println(tinge.Styled().
		With(tinge.Bold).
		Text("Examples:").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("let x: num = 42").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("echo \"Hello, World!\"").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("fn add(a: num, b: num) -> num { return a + b }").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("if x > 10 { echo \"x is large\" }").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("import * from \"std\"").
		Newline().
		String())

	fmt.Println()

	fmt.Println(tinge.Styled().
		With(tinge.Italic).
		Text("💡 Use Tab for autocompletion and Ctrl+C to exit.").
		Newline().
		String())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
