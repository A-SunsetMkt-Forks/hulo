// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package meta

// These variables are set by the -ldflags command.
var (
	// Version is the version of the compiler.
	// GitVersion defines the version string of the compiler.
	// It is set by the -ldflags command.
	Version = "dev"
	// Date is the date of the compiler
	Date = "unknown"
	// Commit is the commit of the compiler
	Commit = "none"
	// GoVersion is the version of the Go compiler
	GoVersion = "unknown"
)

const HULOLOGO = `
██╗  ██╗██╗   ██╗██╗      ██████╗
██║  ██║██║   ██║██║     ██╔═══██╗
███████║██║   ██║██║     ██║   ██║
██╔══██║██║   ██║██║     ██║   ██║
██║  ██║╚██████╔╝███████╗╚██████╔╝
╚═╝  ╚═╝ ╚═════╝ ╚══════╝ ╚═════╝
`

const HLPMLOGO = `
██╗  ██╗██╗     ██████╗ ███╗   ███╗
██║  ██║██║    ██╔═══██╗████╗ ████║
███████║██║    ██║   ██║██╔████╔██║
██╔══██║██║    ██║   ██║██║╚██╔╝██║
██║  ██║███████╗╚██████╔╝██║ ╚═╝ ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝     ╚═╝
`
