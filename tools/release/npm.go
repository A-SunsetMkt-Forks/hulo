// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package release

type PackageJSON struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Main        string            `json:"main"`
	Bin         map[string]string `json:"bin"`
	Files       []string          `json:"files"`
	Keywords    []string          `json:"keywords"`
	Author      string            `json:"author"`
	License     string            `json:"license"`
	Homepage    string            `json:"homepage"`
	Repository  Repository        `json:"repository"`
	Os          []string          `json:"os"`
	CPU         []string          `json:"cpu"`
	Engines     Engines           `json:"engines"`
}

type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Engines struct {
	Node string `json:"node"`
}
