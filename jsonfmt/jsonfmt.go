// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// +build ignore

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dsnet/golib/jsonfmt"
)

func main() {
	minify := flag.Bool("minify", false, "produce minified output")
	standardize := flag.Bool("standardize", false, "produce ECMA-404 compliant output")
	flag.Parse()

	var opts []jsonfmt.Option
	if *minify {
		opts = append(opts, jsonfmt.Minify())
	}
	if *standardize {
		opts = append(opts, jsonfmt.Standardize())
	}

	if len(flag.Args()) > 0 {
		for _, path := range flag.Args() {
			in, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("ioutil.ReadFile error: %v", err)
			}
			out, err := jsonfmt.Format(in, opts...)
			if err != nil {
				log.Fatalf("jsonfmt.Format error: %v", err)
			}
			if err := replaceFile(path, out); err != nil {
				log.Fatalf("replaceFile error: %v", err)
			}
		}
	} else {
		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("ioutil.ReadAll error: %v", err)
		}
		out, err := jsonfmt.Format(in, opts...)
		if err != nil {
			log.Fatalf("jsonfmt.Format error: %v", err)
		}
		if _, err := os.Stdout.Write(out); err != nil {
			log.Fatalf("os.Stdout.Write error: %v", err)
		}
	}
}

func replaceFile(path string, b []byte) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile(filepath.Dir(path), "tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := f.Chmod(fi.Mode()); err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
