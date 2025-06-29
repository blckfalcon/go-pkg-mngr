package main

import (
	"debug/buildinfo"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"golang.org/x/mod/semver"
)

func listPackages(goPath string) ([]string, error) {
	var pkgs []string

	entries, err := os.ReadDir(filepath.Join(goPath, "bin"))
	if err != nil {
		return pkgs, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		pkgs = append(pkgs, entry.Name())
	}

	return pkgs, nil
}

func printUpdateInfo(w io.Writer, goPath, pkg string) error {
	var err error

	b, err := buildinfo.ReadFile(filepath.Join(goPath, "bin", pkg))
	if err != nil {
		return err
	}

	out, err := exec.Command("go", "list", "-m", "-f", "{{.Version}}", b.Main.Path+"@latest").Output()
	if err != nil {
		return err
	}

	latestVersion := strings.TrimSpace(string(out))
	switch {
	case !semver.IsValid(b.Main.Version) || !semver.IsValid(latestVersion):
		color.Set(color.FgYellow)
	case semver.Compare(b.Main.Version, latestVersion) == -1:
		color.Set(color.FgRed)
	case semver.Compare(b.Main.Version, latestVersion) >= 0:
		color.Set(color.FgGreen)
	}

	fmt.Fprintf(w, "%-32s %-64s %-16.15s %-16.15s\n", pkg, b.Path, b.Main.Version, latestVersion)
	color.Unset()

	return err
}

func main() {
	var err error

	goPath := os.Getenv("GOPATH")

	pkgs, err := listPackages(goPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(pkgs)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, pkg := range pkgs {

		wg.Add(1)
		sem <- struct{}{}

		go func(w io.Writer, goPath, pkg string) {
			defer wg.Done()
			defer func() { <-sem }()

			_ = printUpdateInfo(w, goPath, pkg)
		}(os.Stdout, goPath, pkg)

	}
	wg.Wait()
}
