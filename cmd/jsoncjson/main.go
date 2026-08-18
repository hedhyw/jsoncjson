// Command jsoncjson reads JSONC (JSON with comments) and writes plain JSON.
//
// Usage:
//
//	jsoncjson [flags] [file...]
//
// With no file arguments, or with "-", the input is read from stdin. The
// result is written to stdout, unless -o is given.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/hedhyw/jsoncjson"
)

const usage = `jsoncjson reads JSONC (JSON with comments) and writes plain JSON.

Usage:
  jsoncjson [flags] [file...]

With no file arguments, or with "-", the input is read from stdin.

Flags:
`

// version is set by the release build via -ldflags, it falls back to the
// version recorded by the go tool.
var version = ""

func main() {
	err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}

		fmt.Fprintln(os.Stderr, "jsoncjson:", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) (err error) {
	flagSet := flag.NewFlagSet("jsoncjson", flag.ContinueOnError)
	flagSet.SetOutput(stderr)
	flagSet.Usage = func() {
		_, _ = fmt.Fprint(stderr, usage)
		flagSet.PrintDefaults()
	}

	outPath := flagSet.String("o", "-", `output file, "-" for stdout`)
	showVersion := flagSet.Bool("version", false, "print the version and exit")

	err = flagSet.Parse(args)
	if err != nil {
		return err
	}

	if *showVersion {
		_, _ = fmt.Fprintln(stdout, getVersion())

		return nil
	}

	out, err := openOutput(*outPath, stdout)
	if err != nil {
		return err
	}

	defer func() { err = errors.Join(err, out.Close()) }()

	inPaths := flagSet.Args()
	if len(inPaths) == 0 {
		inPaths = []string{"-"}
	}

	for _, path := range inPaths {
		err = convert(path, stdin, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func convert(path string, stdin io.Reader, out io.Writer) (err error) {
	in, err := openInput(path, stdin)
	if err != nil {
		return err
	}

	defer func() { _ = in.Close() }()

	_, err = io.Copy(out, jsoncjson.NewReader(in))
	if err != nil {
		return fmt.Errorf("converting %s: %w", name(path), err)
	}

	return nil
}

func openInput(path string, stdin io.Reader) (io.ReadCloser, error) {
	if path == "-" {
		return io.NopCloser(stdin), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input: %w", err)
	}

	return file, nil
}

func openOutput(path string, stdout io.Writer) (io.WriteCloser, error) {
	if path == "-" {
		return nopWriteCloser{Writer: stdout}, nil
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("creating output: %w", err)
	}

	return file, nil
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func name(path string) string {
	if path == "-" {
		return "stdin"
	}

	return path
}

func getVersion() string {
	if version != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" {
		return info.Main.Version
	}

	return "dev"
}
