package main

import "github.com/hulo-lang/hulo/internal/version"

func main() {
	vm := version.NewVersionManager(".")
	err := vm.GetQualityReport()
	if err != nil {
		panic(err)
	}
}
