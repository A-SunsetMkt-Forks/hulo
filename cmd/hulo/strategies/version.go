// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package strategies

import (
	"fmt"
	"runtime"

	"github.com/hulo-lang/hulo/cmd/meta"
)

var _ Strategy = (*VersionStrategy)(nil)

type VersionStrategy struct{}

func (v *VersionStrategy) CanHandle(params *Parameters) bool {
	return params.Version
}

func (v *VersionStrategy) Execute(params *Parameters, args []string) error {
	fmt.Printf("hulo %s\n", meta.Version)
	fmt.Printf("Hulo Runtime Environment (build %s)\n", meta.Date)
	fmt.Printf("Hulo %s %s Compiler (build %s, %s)\n", runtime.GOOS, runtime.GOARCH, meta.Commit, meta.GoVersion)
	return nil
}

func (v *VersionStrategy) Priority() int {
	return PriorityVersion
}
