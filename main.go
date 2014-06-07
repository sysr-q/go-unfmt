/* go-unfmt: write bad code efficiently for good
 * Copyright (C) 2014 sysr-q
 * Copyright (C) 2014 Cyphar
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

func unfmt(r io.Reader) error {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file := fset.AddFile("<stdin>", fset.Base(), len(src))

	s := scanner.Scanner{}
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, str := s.Scan()
		if len(str) == 0 {
			str = tok.String()
		}

		if tok == token.EOF {
			break
		}

		fmt.Println(pos, tok, str)
	}

	return nil
}

var fixVars = flag.Bool("fixVars", true, "'Correct' variable names according to unfmt spec.")

func main() {
	flag.Parse()

	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go unfmt: failed to ruin the source file '%s'\n", fname)
			continue
		}
		defer f.Close()

		unfmt(f)
	}
}
