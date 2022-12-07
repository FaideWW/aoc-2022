package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Instruction struct {
	command  string
	arg      string
	response []string
}

type File struct {
	name string
	size int
}

type Directory struct {
	name           string
	parent         *Directory
	size           int
	subdirectories map[string]*Directory
	files          map[string]*File
}

const MAX_SIZE = 100000

const DISK_SIZE = 70000000
const FREE_SPACE_NEEDED = 30000000

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	instructions := parseInput(input)

	fs := buildFileSystem(instructions)
	fs.print()

	candidates := getDirectoriesSmallerThan(fs, MAX_SIZE)
	sum := sumDirectorySizes(candidates)

	fmt.Printf("size of directories smaller than %d: %d\n", MAX_SIZE, sum)

	dirToDelete := findDirectoryToDelete(fs, DISK_SIZE, FREE_SPACE_NEEDED)
	fmt.Printf("directory to delete: %s (size=%d)\n", dirToDelete.name, dirToDelete.getSize())
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) (instructions []Instruction) {
	lines := strings.Split(input, "\n")
	linePointer := 0

	for linePointer < len(lines)-1 {
		currentLine := lines[linePointer]
		if currentLine[0] != '$' {
			panic(errors.New("started command on a non-command line"))
		}

		command := currentLine[2:4]
		arg := strings.TrimSpace(currentLine[4:])

		response := make([]string, 0)

		linePointer++
		nextLine := lines[linePointer]
		for nextLine[0] != '$' {
			response = append(response, nextLine)

			if linePointer >= len(lines)-1 {
				break
			}
			linePointer++
			nextLine = lines[linePointer]
		}

		instruction := Instruction{
			command:  command,
			arg:      arg,
			response: response,
		}
		instructions = append(instructions, instruction)
	}

	return
}

func makeDirectory(name string, parent *Directory) *Directory {
	return &Directory{
		name:           name,
		parent:         parent,
		size:           -1,
		subdirectories: map[string]*Directory{},
		files:          map[string]*File{},
	}

}

func parseFile(input string) *File {
	fields := strings.Split(input, " ")
	name := strings.TrimSpace(fields[1])
	size, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(errors.New("size is not an int"))
	}
	return &File{
		name: name,
		size: size,
	}
}

func buildFileSystem(instructions []Instruction) *Directory {
	root := makeDirectory("/", nil)

	currDir := root

	for _, instruction := range instructions {
		switch instruction.command {
		case "cd":
			{
				switch instruction.arg {
				case "..":
					{

						currDir = currDir.parent
					}
				case currDir.name:
					{

						// do nothing
					}
				default:
					{
						if val, ok := currDir.subdirectories[instruction.arg]; ok {
							currDir = val
						} else {
							newDir := makeDirectory(instruction.arg, currDir)
							currDir.subdirectories[newDir.name] = newDir
							currDir = newDir
						}
					}
				}
			}
		case "ls":
			{
				for _, entityStr := range instruction.response {
					if entityStr[:3] == "dir" {
						dir := makeDirectory(strings.TrimSpace(entityStr[4:]), currDir)
						currDir.subdirectories[dir.name] = dir
					} else {
						file := parseFile(entityStr)
						currDir.files[file.name] = file
					}
				}
			}
		default:
			{
				panic(errors.New("unknown command"))
			}
		}
	}

	return root
}

func printDirectoryTree(root *Directory, prefix string) {
	fmt.Printf("%s- %s (dir)\n", prefix, root.name)
	for _, dir := range root.subdirectories {
		printDirectoryTree(dir, prefix+"  ")
	}
	for _, file := range root.files {
		fmt.Printf("%s  - %s (file, size=%d)\n", prefix, file.name, file.size)
	}
}

func (d *Directory) print() {
	printDirectoryTree(d, "")
}

func (d *Directory) getSize() (sum int) {
	if d.size != -1 {
		return d.size
	}
	for _, dir := range d.subdirectories {
		sum += dir.getSize()
	}
	for _, file := range d.files {
		sum += file.size
	}

	d.size = sum
	return
}

// Crawl the entire file system from the root and return a list of
// directories whose size (including subdirectories) is smaller than the
// provided maximum size
func getDirectoriesSmallerThan(root *Directory, max int) [](*Directory) {
	dirs := make([]*Directory, 0)
	rootSize := root.getSize()
	if rootSize <= max {
		dirs = append(dirs, root)
	}

	for _, child := range root.subdirectories {
		dirs = append(dirs, getDirectoriesSmallerThan(child, max)...)
	}

	return dirs
}

func sumDirectorySizes(dirs []*Directory) (sum int) {
	for _, d := range dirs {
		sum += d.getSize()
	}
	return
}

func findDirectoryToDelete(root *Directory, diskSize int, requiredFreeSpace int) *Directory {
	rootSize := root.getSize()
	spaceAvailable := diskSize - rootSize
	spaceToBeFreed := requiredFreeSpace - spaceAvailable

	allDirs := getDirectoriesSmallerThan(root, math.MaxInt)

	sort.Slice(allDirs, func(i, j int) bool {
		return allDirs[i].getSize() < allDirs[j].getSize()
	})

	dirIndex := findSmallestDirectoryLargerThan(allDirs, 0, len(allDirs), spaceToBeFreed)

	if dirIndex == -1 {
		panic(errors.New("no directory sufficiently small"))
	}

	return allDirs[dirIndex]
}

// Binary search for the next value larger than the target
func findSmallestDirectoryLargerThan(dirs []*Directory, start int, end int, target int) int {
	if start == end {
		if dirs[start].getSize() >= target {
			return start
		} else {
			return -1
		}
	}

	mid := start + (end-start)/2

	if dirs[mid].getSize() < target {
		return findSmallestDirectoryLargerThan(dirs, mid+1, end, target)
	}

	found := findSmallestDirectoryLargerThan(dirs, start, mid, target)
	if found != -1 {
		return found
	}
	return -1
}
