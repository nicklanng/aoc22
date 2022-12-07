package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"math"
	"strings"
)

type FSMode int

const (
	FSModeListing = iota
)

type FileSystem struct {
	TotalSpace       int
	CurrentSpace     int
	Mode             FSMode
	CurrentDirectory *File
	Root             *File
}

func NewFileSystem(space int) FileSystem {
	root := &File{Name: "/"}
	return FileSystem{
		TotalSpace:       space,
		CurrentSpace:     space,
		Root:             &File{Name: "/"},
		CurrentDirectory: root,
	}
}

func (fs *FileSystem) ChangeDirectory(target string) {
	switch target {
	case "/":
		fs.CurrentDirectory = fs.Root
	case "..":
		fs.CurrentDirectory = fs.CurrentDirectory.parent
	default:
		for _, c := range fs.CurrentDirectory.children {
			if c.Name == target {
				fs.CurrentDirectory = c
				return
			}
		}
		panic("unknown folder")
	}
}

func (fs *FileSystem) CreateDirectory(name string) {
	for _, c := range fs.CurrentDirectory.children {
		if c.Name == name {
			panic("duplicate name")
		}
	}
	fs.CurrentDirectory.children = append(fs.CurrentDirectory.children, &File{
		Name:      name,
		parent:    fs.CurrentDirectory,
		Directory: true,
	})
}

func (fs *FileSystem) CreateFile(name string, size int) {
	for _, c := range fs.CurrentDirectory.children {
		if c.Name == name {
			panic("duplicate name")
		}
	}
	if size > fs.CurrentSpace {
		panic("insufficient space")
	}
	fs.CurrentSpace -= size
	fs.CurrentDirectory.children = append(fs.CurrentDirectory.children, &File{
		Name:   name,
		size:   size,
		parent: fs.CurrentDirectory,
	})
}

func (fs *FileSystem) Walk(f func(*File)) {
	fs.walk(fs.Root, f)
}

func (fs *FileSystem) walk(file *File, f func(*File)) {
	f(file)
	for _, c := range file.children {
		fs.walk(c, f)
	}
}

type File struct {
	Name      string
	parent    *File
	size      int
	Directory bool
	children  []*File
}

func (f *File) Size() int {
	if f.Directory {
		var total int
		for _, c := range f.children {
			total += c.Size()
		}
		return total
	}

	return f.size
}

//go:embed input
var input []byte

func main() {
	scanner := bufio.NewScanner(bytes.NewReader(input))

	fs := NewFileSystem(70000000)

	for scanner.Scan() {
		s := scanner.Text()
		if s[0] == '$' {
			fields := strings.Fields(s[1:])
			switch fields[0] {
			case "cd":
				fs.ChangeDirectory(fields[1])
			case "ls":
				fs.Mode = FSModeListing
			}
		} else {
			switch fs.Mode {
			case FSModeListing:
				fields := strings.Fields(s)
				if fields[0] == "dir" {
					fs.CreateDirectory(fields[1])
				} else {
					size := lib.MustParseInt(fields[0])
					fs.CreateFile(fields[1], size)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var partOneTotal int
	fs.Walk(func(file *File) {
		if !file.Directory {
			return
		}
		size := file.Size()
		if size <= 100000 {
			partOneTotal += size
		}
	})
	fmt.Printf("Part 1 Total: %d\n", partOneTotal)

	spaceToFree := 30000000 - fs.CurrentSpace
	partTwoTargetDirectorySize := math.MaxInt
	fs.Walk(func(file *File) {
		if !file.Directory {
			return
		}
		size := file.Size()
		if size < spaceToFree {
			return
		}
		if size < partTwoTargetDirectorySize {
			partTwoTargetDirectorySize = size
		}
	})
	fmt.Printf("Part 2 Size: %d\n", partTwoTargetDirectorySize)
}
