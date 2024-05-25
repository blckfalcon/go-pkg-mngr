package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListPackagesNonExistingDir(t *testing.T) {
	nonExistingDir := "/non/existing/dir"

	_, err := listPackages(nonExistingDir)
	if err == nil {
		t.Fatalf("listPackages: got %s, want nil", err)
	}
}

func TestListPackagesNoPackages(t *testing.T) {
	dst := t.TempDir()
	err := os.Mkdir(dst+"/bin", 0755)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	pkgs, err := listPackages(dst)
	if err != nil {
		t.Fatalf("listPackages returned an error: %v", err)
	}

	if len(pkgs) != 0 {
		t.Fatalf("listPackages: got %d, want 0", len(pkgs))
	}
}

func TestListPackages(t *testing.T) {
	goPath := "./testdata/"
	pkgs, err := listPackages(goPath)
	if err != nil {
		t.Fatalf("listPackages returned an error: %v", err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("listPackages: got %d, want 1", len(pkgs))
	}

	if pkgs[0] != "cowsay" {
		t.Fatalf("listPackages: got %s, want cowsay", pkgs[0])
	}
}

func TestPrintUpdateInfo(t *testing.T) {
	goPath := "./testdata/"
	pkgs, err := listPackages(goPath)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	var out bytes.Buffer
	err = printUpdateInfo(&out, goPath, pkgs[0])
	if err != nil {
		t.Fatalf("unexpected error")
	}

	got := strings.Fields(out.String())
	want := []string{"cowsay", "github.com/Code-Hex/Neo-cowsay/cmd/v2/cowsay", "v2.0.4", "v2.0.4"}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("printUpdateInfo (-got +want):\n%s", diff)
	}
}
