//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"

	"github.com/reviewdog/errorformat/fmts"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, firstlines)
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, pkgcomment)
	fmt.Fprintln(buf, "//")
	fmt.Fprintln(buf, "// Defined formats:")
	fmt.Fprintln(buf, "// ")
	langToFmts := fmts.DefinedFmtsByLang()

	langs := make([]string, 0, len(langToFmts))
	for lang, _ := range langToFmts {
		langs = append(langs, lang)
	}
	sort.Strings(langs)

	for _, lang := range langs {
		nameToFmt := langToFmts[lang]
		names := make([]string, 0, len(nameToFmt))
		for name, _ := range nameToFmt {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Fprintf(buf, "// \t%v\n", lang)
		for _, name := range names {
			f := nameToFmt[name]
			fmt.Fprintf(buf, "// \t\t%s\t%s - %s\n", f.Name, f.Description, f.URL)
		}
	}
	fmt.Fprintln(buf, pkgline)

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile("doc.go", source, 0644)
}

const (
	firstlines = `// Code generated by fmts/gendoc.go; DO NOT EDIT.
// Please run '$ go generate ./...' instead to update this file`
	pkgcomment = `// Package fmts holds defined errorformats.`
	pkgline    = `package fmts`
)
