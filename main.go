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
	_ "go/ast" // TODO(cyphar): use this import
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)

func unfmt(src []byte, srcName string) (string, error) {
	fset := token.NewFileSet()
	fset.AddFile("<stdin>", fset.Base(), len(src))

	astree, err := parser.ParseFile(fset, srcName, src, parser.ParseComments|parser.DeclarationErrors|parser.AllErrors)
	if err != nil {
		return "", err
	}

	// TODO(cyphar): Replace with proper unfmt-ing code
	return fmt.Sprintln(astree), nil
}

var (
	oVars       = flag.Bool("vars", false, "'Correct' variable names according to the unfmt spec. -- NYI")
	oIndent     = flag.Bool("indent", true, "'Correct' indentation according to the unfmt spec. -- NYI")
	oWhitespace = flag.Bool("whitespace", true, "'Correct' whitespace in control structures according to the unfmt spec. -- NYI")
	oParameters = flag.Bool("parameters", true, "'Correct' function parameters according to the unfmt spec. -- NYI")
	oStruct     = flag.Bool("structs", true, "'Correct' struct instantiations according to the unfmt spec. -- NYI")
	oImports    = flag.Bool("imports", true, "'Correct' imports according to the unfmt spec. -- NYI")
)

func main() {
	flag.Parse()

	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go unfmt: failed to ruin the source file '%s': %s\n", fname, err.Error())
			continue
		}
		defer f.Close()

		src, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go unfmt: failed to read the source file '%s': %s\n", fname, err.Error())
			continue
		}

		fmtSrc, err := unfmt(src, fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go unfmt: error encountered while ruining source file '%s'.\n", fname)
			panic(err) // TODO(sysr_q): use log.Fatalf or something.
		}

		fmt.Println(fmtSrc)
	}
}
