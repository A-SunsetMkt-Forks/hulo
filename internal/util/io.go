// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package util

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfigure[T any](path string) (T, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return *new(T), err
	}
	var config T
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return *new(T), err
	}
	return config, nil
}

func SaveConfigure[T any](path string, config T) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func Download(url string, out io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
