// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

const HuloReplFileName = "hulo-repl.yaml"

type HuloRepl struct {
	Theme           string                  `yaml:"theme" validate:"oneof=default dark light colorful custom"`
	ThemeFiles      map[string]string       `yaml:"theme_files"`
	HistoryFile     string                  `yaml:"history_file"`
	DevelopmentMode bool                    `yaml:"development_mode"`
	DebugLevel      string                  `yaml:"debug_level" validate:"oneof=debug info warn error"`
	Keybindings     map[string]string       `yaml:"keybindings"`
	MaxSuggestions  uint16                  `yaml:"max_suggestions" validate:"required,min=1,max=100"`
	Parser          DebugSettingsParser     `yaml:"parser" validate:"required"`
	Transpiler      DebugSettingsTranspiler `yaml:"transpiler" validate:"required"`
	Analyzer        DebugSettingsAnalyzer   `yaml:"analyzer" validate:"required"`
	Lexer           DebugSettingsLexer      `yaml:"lexer" validate:"required"`
	Token           DebugSettingsToken      `yaml:"token" validate:"required"`
	Target          string                  `yaml:"target" validate:"required,oneof=bash batch vbs"`
}

type DebugSettingsParser struct {
	Output string `yaml:"output"`
	Enable bool   `yaml:"enable"`
}

type DebugSettingsTranspiler struct {
	Output string `yaml:"output"`
	Enable bool   `yaml:"enable"`
}

type DebugSettingsAnalyzer struct {
	Output string `yaml:"output"`
	Enable bool   `yaml:"enable"`
	Timing bool   `yaml:"timing"`
}

type DebugSettingsLexer struct {
	Output string `yaml:"output"`
	Enable bool   `yaml:"enable"`
}

type DebugSettingsToken struct {
	Output string `yaml:"output"`
	Enable bool   `yaml:"enable"`
}

const HuloReplTemplate = `# Hulo REPL Configuration File
# This file contains default settings for the Hulo REPL

# General settings
dry_run: false
verbose: false
max_suggestions: 15

# Theme settings
theme: "default"  # Options: default, dark, light, colorful, custom
`

type PromptTheme struct {
	PrefixTextColor              string `yaml:"prefix_text_color"`
	PrefixBGColor                string `yaml:"prefix_bg_color"`
	InputTextColor               string `yaml:"input_text_color"`
	InputBGColor                 string `yaml:"input_bg_color"`
	PreviewSuggestionTextColor   string `yaml:"preview_suggestion_text_color"`
	PreviewSuggestionBGColor     string `yaml:"preview_suggestion_bg_color"`
	SuggestionTextColor          string `yaml:"suggestion_text_color"`
	SuggestionBGColor            string `yaml:"suggestion_bg_color"`
	SelectedSuggestionTextColor  string `yaml:"selected_suggestion_text_color"`
	SelectedSuggestionBGColor    string `yaml:"selected_suggestion_bg_color"`
	DescriptionTextColor         string `yaml:"description_text_color"`
	DescriptionBGColor           string `yaml:"description_bg_color"`
	SelectedDescriptionTextColor string `yaml:"selected_description_text_color"`
	SelectedDescriptionBGColor   string `yaml:"selected_description_bg_color"`
	ScrollbarThumbColor          string `yaml:"scrollbar_thumb_color"`
	ScrollbarBGColor             string `yaml:"scrollbar_bg_color"`
}

const PromptThemeTemplate = `prefix_text_color: white
prefix_bg_color:
input_text_color:
input_bg_color:
preview_suggestion_text_color: cyan
preview_suggestion_bg_color: dark_gray
suggestion_text_color: white
suggestion_bg_color: dark_gray
selected_suggestion_text_color: white
selected_suggestion_bg_color: blue
description_text_color: light_gray
description_bg_color: dark_gray
selected_description_text_color: white
selected_description_bg_color: red
scrollbar_thumb_color: gray
scrollbar_bg_color: dark_gray
`
