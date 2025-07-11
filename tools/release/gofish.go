// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package release

type Resource struct {
	Path        string
	InstallPath string
	Executable  bool
}

type Package struct {
	OS        string
	Arch      string
	URL       string
	SHA256    string
	Resources []Resource
}

type Food struct {
	Name        string
	Description string
	Homepage    string
	Version     string
	License     string
	Packages    []Package
}
