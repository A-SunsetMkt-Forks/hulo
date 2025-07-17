// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package theme

import (
	"reflect"

	"github.com/c-bata/go-prompt"
	"github.com/hulo-lang/hulo/internal/config"
)

var textToColor = map[string]prompt.Color{
	"black":      prompt.Black,
	"dark_red":   prompt.DarkRed,
	"dark_green": prompt.DarkGreen,
	"brown":      prompt.Brown,
	"dark_blue":  prompt.DarkBlue,
	"purple":     prompt.Purple,
	"cyan":       prompt.Cyan,
	"light_gray": prompt.LightGray,
	"dark_gray":  prompt.DarkGray,
	"red":        prompt.Red,
	"green":      prompt.Green,
	"yellow":     prompt.Yellow,
	"blue":       prompt.Blue,
	"fuchsia":    prompt.Fuchsia,
	"turquoise":  prompt.Turquoise,
	"white":      prompt.White,
}

var textToOption = map[string]func(prompt.Color) prompt.Option{
	"prefix_text_color":               prompt.OptionPrefixTextColor,
	"prefix_bg_color":                 prompt.OptionPrefixBackgroundColor,
	"input_text_color":                prompt.OptionInputTextColor,
	"input_bg_color":                  prompt.OptionInputBGColor,
	"preview_suggestion_text_color":   prompt.OptionPreviewSuggestionTextColor,
	"preview_suggestion_bg_color":     prompt.OptionPreviewSuggestionBGColor,
	"suggestion_text_color":           prompt.OptionSuggestionTextColor,
	"suggestion_bg_color":             prompt.OptionSuggestionBGColor,
	"selected_suggestion_text_color":  prompt.OptionSelectedSuggestionTextColor,
	"selected_suggestion_bg_color":    prompt.OptionSelectedSuggestionBGColor,
	"description_text_color":          prompt.OptionDescriptionTextColor,
	"description_bg_color":            prompt.OptionDescriptionBGColor,
	"selected_description_text_color": prompt.OptionSelectedDescriptionTextColor,
	"selected_description_bg_color":   prompt.OptionSelectedDescriptionBGColor,
	"scrollbar_thumb_color":           prompt.OptionScrollbarThumbColor,
	"scrollbar_bg_color":              prompt.OptionScrollbarBGColor,
}

func Parse(cfg config.PromptTheme) []prompt.Option {
	options := []prompt.Option{}

	for i := range reflect.TypeOf(cfg).NumField() {
		field := reflect.TypeOf(cfg).Field(i)

		key, ok := field.Tag.Lookup("yaml")
		if !ok {
			continue
		}

		options = append(options, textToOption[key](textToColor[key]))
	}
	return options
}
