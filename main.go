package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Node interface {
	fmt.Stringer
}

type Directory struct {
	name  string
	nodes []Node
}

type File struct {
	name string
	size int64
}

func (f File) String() string {
	if f.size == 0 {
		return fmt.Sprintf("%s (empty)", f.name)
	}
	return fmt.Sprintf("%s (%db)", f.name, f.size)
}

func (d Directory) String() string {
	return d.name
}

func walk(path string, nodes []Node, withFiles bool) ([]Node, error) {
	dir, err := ioutil.ReadDir(path)

	for _, fInfo := range dir {
		if !(fInfo.IsDir() || withFiles) {
			continue
		}

		var node Node
		if fInfo.IsDir() {
			nodes, _ := walk(filepath.Join(path, fInfo.Name()), []Node{}, withFiles)
			node = Directory{fInfo.Name(), nodes}
		} else {
			node = File{fInfo.Name(), fInfo.Size()}
		}

		nodes = append(nodes, node)
	}

	return nodes, err
}

func printDir(out io.Writer, nodes []Node, prefixes []string) {
	if len(nodes) == 0 {
		return
	}

	fmt.Fprintf(out, "%s", strings.Join(prefixes, ""))

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node)
		if dir, ok := node.(Directory); ok {
			printDir(out, dir.nodes, append(prefixes, "\t"))
		}
		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node)
	if dir, ok := node.(Directory); ok {
		printDir(out, dir.nodes, append(prefixes, "│\t"))
	}

	printDir(out, nodes[1:], prefixes)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	nodes, err := walk(path, []Node{}, printFiles)
	printDir(out, nodes, []string{})
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
