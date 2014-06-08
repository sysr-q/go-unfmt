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
	"go/ast"
	"go/token"
	"go/parser"
	"io/ioutil"
	"os"
)

func unfmtNode(node ast.Node) error {
	return nil
}

func unfmt(src []byte, srcName string) (<-chan string, <-chan error, <-chan struct{}) {
	var (
		srcOut = make(chan string)
		srcErr = make(chan error)
		done   = make(chan struct{})
	)

	go func() {
		defer close(done)

		fset := token.NewFileSet()
		fset.AddFile("<stdin>", fset.Base(), len(src))

		astree, err := parser.ParseFile(fset, srcName, src, parser.ParseComments | parser.DeclarationErrors | parser.AllErrors)
		if err != nil {
			srcErr <- err
			return
		}

		// TODO(cyphar): Replace with proper unfmt-ing code
		srcOut <- fmt.Sprintln(astree)
	}()

	return srcOut, srcErr, done
}

var (
	oVars       = flag.Bool("fix-vars", false, "'Correct' variable names according to the unfmt spec. -- NYI")
	oIndent     = flag.Bool("fix-indent", true, "'Correct' indentation according to the unfmt spec. -- NYI")
	oWhitespace = flag.Bool("fix-whitespace", true, "'Correct' whitespace in control structures according to the unfmt spec. -- NYI")
	oParameters = flag.Bool("fix-parameters", true, "'Correct' function parameters according to the unfmt spec. -- NYI")
	oStruct     = flag.Bool("fix-structs", true, "'Correct' struct instantiations according to the unfmt spec. -- NYI")
	oImports    = flag.Bool("fix-imports", true, "'Correct' imports according to the unfmt spec. -- NYI")
)

func main() {
	flag.Parse()

	for _, fname := range flag.Args() {
		goFile, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go unfmt: failed to ruin the source file '%s': %s\n", fname, err.Error())
			continue
		}
		defer goFile.Close()

		goSrc, err := ioutil.ReadAll(goFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go unfmt: failed to read the source file '%s': %s\n", fname, err.Error())
			continue
		}

		srcChan, errChan, done := unfmt(goSrc, fname)

		for {
			select {
			case srcToken := <-srcChan:
				// TODO(cyphar): Replace this with proper printing code
				fmt.Print(srcToken)
			case err := <-errChan:
				fmt.Fprintf(os.Stderr, "go unfmt: error encountered while ruining source file '%s': %s\n", fname, err.Error())
				goto fail
			case <-done:
				goto fail
			}
		}

		// TODO(cyphar): remove this f*kkin thing.
		fail:
	}
}
