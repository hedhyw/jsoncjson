package main

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const inputJSONC = `{
	// Comment.
	"Hello": "world" /* Comment. */
}`

const outputJSON = "{\n\t\t\"Hello\": \"world\" \n}"

func TestRunStdin(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}

	err := run(nil, strings.NewReader(inputJSONC), stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}

	if got := stdout.String(); got != outputJSON {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inPath := filepath.Join(dir, "in.jsonc")

	err := os.WriteFile(inPath, []byte(inputJSONC), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	stdout := &bytes.Buffer{}

	err = run([]string{inPath, "-"}, strings.NewReader(inputJSONC), stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}

	if got := stdout.String(); got != outputJSON+outputJSON {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunOutputFile(t *testing.T) {
	t.Parallel()

	outPath := filepath.Join(t.TempDir(), "out.json")

	err := run([]string{"-o", outPath}, strings.NewReader(inputJSONC), &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}

	if got := string(data); got != outputJSON {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunVersion(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}

	err := run([]string{"-version"}, strings.NewReader(""), stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}

	if stdout.Len() == 0 {
		t.Fatal("version is not printed")
	}
}

func TestRunHelp(t *testing.T) {
	t.Parallel()

	stderr := &bytes.Buffer{}

	err := run([]string{"-h"}, strings.NewReader(""), &bytes.Buffer{}, stderr)
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("unexpected usage: %q", stderr.String())
	}
}

func TestRunMissingInput(t *testing.T) {
	t.Parallel()

	err := run([]string{"not_found.jsonc"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunBadOutput(t *testing.T) {
	t.Parallel()

	outPath := filepath.Join(t.TempDir(), "not_found", "out.json")

	err := run([]string{"-o", outPath}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	if getVersion() == "" {
		t.Fatal("version is empty")
	}
}
