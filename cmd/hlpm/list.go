package main

import (
	"fmt"
	"path/filepath"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		if !util.Exists(config.HuloPkgFileName) {
			log.Fatal("package file not found, please run `hlpm init` first")
		}

		modulesDir, err := getHuloModulesDir()
		if err != nil {
			log.WithError(err).Fatal("failed to get hulo modules directory")
		}

		pkg, err := util.LoadConfigure[config.HuloPkg](config.HuloPkgFileName)
		if err != nil {
			log.WithError(err).Fatal("failed to load package")
		}

		if len(pkg.Dependencies) == 0 {
			fmt.Println("No dependencies found.")
			return
		}

		fmt.Printf("Dependency tree for %s@%s:\n\n", pkg.Name, pkg.Version)

		// Build dependency tree
		visited := make(map[string]bool)
		tree, err := buildDependencyTree(modulesDir, pkg.Dependencies, 0, visited)
		if err != nil {
			log.WithError(err).Fatal("failed to build dependency tree")
		}

		// Print the tree
		printDependencyTree(tree, "", true)
	},
}

// DependencyNode represents a node in the dependency tree
type DependencyNode struct {
	Name     string                     // owner/repo
	Version  string                     // version
	Children map[string]*DependencyNode // sub-dependencies
	Level    int                        // depth level
}

// DependencyTree represents the complete dependency tree
type DependencyTree struct {
	Root *DependencyNode
}

// loadPackageDependencies loads dependencies from a package's hulo.pkg.yaml file
func loadPackageDependencies(modulesDir, owner, repo, version string) (map[string]string, error) {
	packageDir := filepath.Join(modulesDir, owner, repo, version)
	packageFile := filepath.Join(packageDir, config.HuloPkgFileName)

	if !util.Exists(packageFile) {
		return make(map[string]string), nil // No dependencies
	}

	pkg, err := util.LoadConfigure[config.HuloPkg](packageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load package %s/%s@%s: %w", owner, repo, version, err)
	}

	return pkg.Dependencies, nil
}

// buildDependencyTree recursively builds the dependency tree
func buildDependencyTree(modulesDir string, dependencies map[string]string, level int, visited map[string]bool) (*DependencyNode, error) {
	if level > 10 { // Prevent infinite recursion
		return nil, fmt.Errorf("dependency depth too deep, possible circular dependency")
	}

	root := &DependencyNode{
		Children: make(map[string]*DependencyNode),
		Level:    level,
	}

	for pkgName, version := range dependencies {
		owner, repo, _, err := resolvePackage(pkgName)
		if err != nil {
			log.WithError(err).Warnf("failed to resolve package %s", pkgName)
			continue
		}

		// Create node for this dependency
		node := &DependencyNode{
			Name:     pkgName,
			Version:  version,
			Children: make(map[string]*DependencyNode),
			Level:    level + 1,
		}

		// 检查是否访问过
		key := fmt.Sprintf("%s@%s", pkgName, version)
		if !visited[key] {
			visited[key] = true
			// 只在没访问过时递归子依赖
			subDeps, err := loadPackageDependencies(modulesDir, owner, repo, version)
			if err != nil {
				log.WithError(err).Warnf("failed to load dependencies for %s@%s", pkgName, version)
			} else if len(subDeps) > 0 {
				subTree, err := buildDependencyTree(modulesDir, subDeps, level+1, visited)
				if err != nil {
					log.WithError(err).Warnf("failed to build sub-tree for %s@%s", pkgName, version)
				} else {
					node.Children = subTree.Children
				}
			}
		}
		// 无论是否访问过，都要加到父节点
		root.Children[pkgName] = node
	}

	return root, nil
}

// printDependencyTree prints the dependency tree in a tree-like format
func printDependencyTree(node *DependencyNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Print current node
	if node.Name != "" {
		fmt.Printf("%s%s@%s\n", prefix, node.Name, node.Version)
	}

	// Print children
	children := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		children = append(children, childName)
	}

	for i, childName := range children {
		child := node.Children[childName]
		isLastChild := i == len(children)-1

		// Determine the prefix for the child
		childPrefix := prefix
		if node.Name != "" { // Not root node
			if isLast {
				childPrefix += "    "
			} else {
				childPrefix += "│    "
			}

		}

		// Add connector
		if node.Name != "" { // Not root node
			if isLastChild {
				childPrefix += "└── "
			} else {
				childPrefix += "├── "
			}
		}

		printDependencyTree(child, childPrefix, isLastChild)
	}
}
