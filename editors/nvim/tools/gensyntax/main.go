// Command gensyntax writes the Neovim plugin's Vim syntax files from the
// lexer's keyword list, so highlighting cannot drift from the language.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("out", ".", "directory to write the syntax files into")
	flag.Parse()

	for _, s := range Syntaxes() {
		data, err := Render(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gensyntax: %v\n", err)
			os.Exit(1)
		}
		path := filepath.Join(*dir, s.File)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "gensyntax: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("wrote %s\n", path)
	}
}
