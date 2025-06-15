package main

import (
	"context"
	"fmt"
	"log"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Writer()))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	src := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{".git"},
	})

	base := client.Container().
		From("ubuntu:22.04").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "golang-go", "openjdk-17-jdk", "nodejs", "npm", "curl"})

	goPath := "/root/go/bin"
	base = base.
		WithExec([]string{"go", "install", "github.com/magefile/mage@latest"}).
		WithEnvVariable("PATH", fmt.Sprintf("%s:/usr/local/go/bin:/usr/bin:/bin", goPath)).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")

	dep := base.WithExec([]string{"go", "mod", "download"})

	antlrJar := "/src/syntax/hulo/parser/antlr.jar"
	antlr := dep.
		WithExec([]string{"mkdir", "-p", "/src/syntax/hulo/parser"}).
		WithExec([]string{"curl", "-L", "https://www.antlr.org/download/antlr-4.13.2-complete.jar", "-o", antlrJar})

	genParser := antlr.WithExec([]string{
		"java", "-jar", antlrJar,
		"-Dlanguage=Go", "-visitor", "-no-listener",
		"-o", "/src/syntax/hulo/parser/generated",
		"-package", "generated",
		"/src/syntax/hulo/parser/grammar/*.g4",
	})

	build := genParser.WithExec([]string{"goreleaser", "build", "--snapshot", "--clean"})

	exportDir := build.Directory("/src/dist")

	_, err = exportDir.Export(ctx, "./dist")
	if err != nil {
		panic(err)
	}

	fmt.Println("Build completed and artifacts exported to ./dist")
}
