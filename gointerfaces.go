package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	oldSrcURL = "https://storage.googleapis.com/google-code-archive-downloads/v2/code.google.com/go/"
	newSrcURL = "https://storage.googleapis.com/golang/"
	oldSrcDir = "src/pkg"
	newSrcDir = "src"
	// expects go version, source file and line number
	sourceURL       = "https://github.com/golang/go/blob/go%s/%s#L%s"
	interfaceRegexp = `^\s*type\s+([A-Z]\w*)\s+interface\s*{`
)

// Interface is an interface
type Interface struct {
	Name    string
	Package string
}

// Location is the location in sources
type Location struct {
	SourceFile string
	LineNumber string
	Link       string
}

// InterfaceList is a map of interfaces to their location
type InterfaceList map[Interface]map[string]Location

// NewInterfaceList builds a list of interfaces
func NewInterfaceList() InterfaceList {
	return make(map[Interface]map[string]Location)
}

// AddInterface adds an interface to a list
func (il InterfaceList) AddInterface(name, pkg, version, sourceFile, lineNumber string) {
	interf := Interface{
		Name:    name,
		Package: pkg,
	}
	link := fmt.Sprintf(sourceURL, version, sourceFile, lineNumber)
	location := Location{
		SourceFile: sourceFile,
		LineNumber: lineNumber,
		Link:       link,
	}
	if il[interf] == nil {
		il[interf] = make(map[string]Location)
	}
	il[interf][version] = location
}

// ByName is a list of interfaces
type ByName []Interface

// Len returns the length of the list
func (b ByName) Len() int { return len(b) }

// Swap swaps two interfaces
func (b ByName) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// Less tells if i is less than j
func (b ByName) Less(i, j int) bool { return b[i].Name < b[j].Name }

// srcDirUrl returns the URL of source directory
func srcDirURL(v string) (string, string) {
	array := strings.Split(strings.Split(strings.Split(v, "beta")[0], "rc")[0], ".")
	major, err := strconv.Atoi(array[0])
	if err != nil {
		major = 0
	}
	minor, err := strconv.Atoi(array[1])
	if err != nil {
		minor = 0
	}
	srcDir := ""
	srcURL := ""
	if major <= 1 && minor < 4 {
		srcDir = oldSrcDir
	} else {
		srcDir = newSrcDir
	}
	if major <= 1 && minor < 2 {
		srcURL = oldSrcURL
	} else {
		srcURL = newSrcURL
	}
	return srcDir, srcURL
}

// parseSourceFile parses a source file and populates the interface list
func parseSourceFile(filename string, source io.Reader, sourceDir string, version string, interfaces InterfaceList) {
	regexpInterface := regexp.MustCompile(interfaceRegexp)
	reader := bufio.NewReader(source)
	pack := filename[len(sourceDir)+4 : strings.LastIndex(filename, "/")]
	if strings.HasSuffix(pack, "testdata") || strings.HasPrefix(pack, "cmd") {
		return
	}
	lineNumber := 1
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			panic("Error parsing source file")
		}
		matches := regexpInterface.FindSubmatch(line)
		if len(matches) > 0 {
			name := string(matches[1])
			pkg := string(pack)
			sourceFile := filename[3:]
			lineNumber := strconv.Itoa(lineNumber)
			interfaces.AddInterface(name, pkg, version, sourceFile, lineNumber)
		}
		if err == io.EOF {
			break
		}
		lineNumber++
	}
}

// addInterfaces generates interface list for given version
func addInterfaces(version string, interfaces InterfaceList) {
	println(fmt.Sprintf("Generating interface list for version %s...", version))
	srcDir, srcURL := srcDirURL(version)
	// download compressed archive
	response, err := http.Get(srcURL + "go" + version + ".src.tar.gz")
	if err != nil {
		panic(err)
	}
	reader := response.Body
	defer response.Body.Close()
	// gunzip the archive stream
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		panic(err)
	}
	// parse tar source files in source dir
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			break
		}
		if strings.HasPrefix(header.Name, "go/"+srcDir) &&
			strings.HasSuffix(header.Name, ".go") &&
			!strings.HasSuffix(header.Name, "doc.go") &&
			!strings.HasSuffix(header.Name, "_test.go") {
			parseSourceFile(header.Name, tarReader, srcDir, version, interfaces)
		}
	}
}

// printInterfaces prints interfaces for given versions
func printInterfaces(interfaceList InterfaceList, versions []string) {
	interfaces := make([]Interface, 0)
	for i := range interfaceList {
		interfaces = append(interfaces, i)
	}
	sort.Sort(ByName(interfaces))
	lenName := 0
	lenPackage := 0
	lenVersions := make(map[string]int)
	for _, i := range interfaces {
		if len(i.Name) > lenName {
			lenName = len(i.Name)
		}
		if len(i.Package) > lenPackage {
			lenPackage = len(i.Package)
		}
		for _, version := range versions {
			loc := interfaceList[i][version]
			lenVersion := len(loc.Link) + 10
			if lenVersions[version] < lenVersion {
				lenVersions[version] = lenVersion
			}
		}
	}
	formatLine := "%-" + strconv.Itoa(lenName) + "s" + " | %-" + strconv.Itoa(lenPackage) + "s"
	for _, v := range versions {
		formatLine += " | %-" + strconv.Itoa(lenVersions[v]) + "s"
	}
	args := []interface{}{"Interface", "Package"}
	for _, v := range versions {
		args = append(args, v)
	}
	fmt.Println(fmt.Sprintf(formatLine, args...))
	separator := ":" + strings.Repeat("-", lenName-1) + " | :" + strings.Repeat("-", lenPackage-1)
	for _, v := range versions {
		separator += " | " + strings.Repeat("-", lenVersions[v])
	}
	fmt.Println(separator)
	for _, i := range interfaces {
		versionLink := make(map[string]string)
		for _, v := range versions {
			if len(interfaceList[i][v].SourceFile) > 0 {
				versionLink[v] = "[source](" + interfaceList[i][v].Link + ")"
			} else {
				versionLink[v] = "-"
			}
		}
		args := []interface{}{i.Name, i.Package}
		for _, v := range versions {
			args = append(args, versionLink[v])
		}
		fmt.Println(fmt.Sprintf(formatLine, args...))
	}
}

// main is the program entry point
func main() {
	// read versions on command line
	if len(os.Args) < 2 {
		panic("Must pass go version(s) on command line")
	}
	versions := os.Args[1:]
	// iterate on versions
	interfaces := NewInterfaceList()
	for _, version := range versions {
		addInterfaces(version, interfaces)
	}
	// print the result
	println("Printing table...")
	printInterfaces(interfaces, versions)
}
