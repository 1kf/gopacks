package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

var (
	dir     string
	verbose bool

	packages map[string]bool
)

func init() {

	packages = make(map[string]bool)

	flag.StringVar(&dir, "dir", "", "search folder")
	flag.BoolVar(&verbose, "verbose", false, "verbose information")

	flag.Parse()

	if dir == "" {
		dir, _ = os.Getwd()
	}

	dir, _ = filepath.Abs(dir)
}

func extractPackages(path string) []string {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		log.Fatal(err)
	}

	var packages []string
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, s := range file.Imports {
				packages = append(packages, s.Path.Value)
			}
		}
	}

	return packages
}

func extractFile(path string) ([]string, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var pkgs []string
	for _, s := range f.Imports {
		pkgs = append(pkgs, s.Path.Value)
	}
	return pkgs, nil
}

func walk(dir string) {
	recurse := true
	walkFn := func(path string, info os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.IsDir() && path != dir && !recurse {
			return filepath.SkipDir
		}

		if stat.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext == ".go" {
			pkgs, err := extractFile(path)
			if err != nil {
				return err
			}

			for _, pkg := range pkgs {
				v := pkg[1 : len(pkg)-1]
				fmt.Println(v)
			}
		}

		return nil
	}

	err := filepath.Walk(dir, walkFn)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	walk(dir)
}
