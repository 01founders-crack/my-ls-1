package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"
)

const (
	// ANSI escape codes for colors
	ResetColor  = "\033[0m"
	DirectoryColor = "\033[1;34m"
	SymlinkColor = "\033[1;36m"
)

// FileInfoSlice is a slice of fs.FileInfo
type FileInfoSlice []fs.FileInfo

func (f FileInfoSlice) Len() int { return len(f) }

func (f FileInfoSlice) Less(i, j int) bool { return f[i].Name() < f[j].Name() }

func (f FileInfoSlice) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func getColor(file fs.FileInfo) string {
	if file.IsDir() {
		return DirectoryColor
	}
	if file.Mode()&fs.ModeSymlink != 0 {
		return SymlinkColor
	}
	return ResetColor
}

func lsDir(path string, recursive, longFormat, showHidden, reverse, sortByModTime bool) {
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var files FileInfoSlice
	for _, entry := range dir {
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		files = append(files, info)
	}

	if sortByModTime {
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})
	} else {
		sort.Sort(files)
	}

	if reverse {
		sort.Sort(sort.Reverse(files))
	}

	for _, file := range files {
		color := getColor(file)
		if longFormat {
			sys := file.Sys().(*syscall.Stat_t)
			fmt.Printf("%s%+v %3d %9d %9d %12d %s %s%s\n",
				color,
				file.Mode(),
				file.Sys().(*syscall.Stat_t).Nlink,
				sys.Uid,
				sys.Gid,
				file.Size(),
				file.ModTime().Format(time.Stamp),
				file.Name(),
				ResetColor)
		} else {
			fmt.Printf("%s%s\t%s", color, file.Name(), ResetColor)
		}

		if recursive && file.IsDir() {
			fmt.Println()
			lsDir(fmt.Sprintf("%s/%s", path, file.Name()), recursive, longFormat, showHidden, reverse, sortByModTime)
		}
	}

	fmt.Println()
}

func main() {
	recursive := flag.Bool("R", false, "List directories recursively")
	longFormat := flag.Bool("l", false, "Use long listing format")
	showHidden := flag.Bool("a", false, "Show hidden files")
	reverse := flag.Bool("r", false, "Reverse order while sorting")
	sortByModTime := flag.Bool("t", false, "Sort by modification time")

	flag.Parse()

	// Check for additional arguments
	args := flag.Args()

	if len(args) == 0 {
		// No additional arguments provided, list the current directory
		lsDir(".", *recursive, *longFormat, *showHidden, *reverse, *sortByModTime)
	} else {
		for _, arg := range args {
			// Handle each directory or file name provided as argument
			lsDir(arg, *recursive, *longFormat, *showHidden, *reverse, *sortByModTime)
		}
	}
}
