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
	URL = "https://storage.googleapis.com/golang/"
	// expects go version, source file and line number
	SOURCE_URL = "https://github.com/golang/go/blob/go%s/%s#L%d"
)

type Interface struct {
	Name    string
	Package string
}

type Location struct {
	SourceFile string
	LineNumber string
	Link       string
}

type InterfaceList map[Interface]map[string]Location

func NewInterfaceList() InterfaceList {
	return make(map[Interface]map[string]Location)
}

func (il InterfaceList) AddInterface(name, pkg, version, sourceFile, lineNumber string) {
	interf := Interface{
		Name:    name,
		Package: pkg,
	}
	link := fmt.Sprintf(SOURCE_URL, version, sourceFile, lineNumber)
	location := Location{
		SourceFile: sourceFile,
		LineNumber: lineNumber,
		Link:       link,
	}
	il[interf][version] = location
}

type ByName []Interface

func (b ByName) Len() int           { return len(b) }
func (b ByName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByName) Less(i, j int) bool { return b[i].Name < b[j].Name }

func versionMajorMinor(v string) (int, int) {
	array := strings.Split(strings.Split(v, "rc")[0], ".")
	major, err := strconv.Atoi(array[0])
	if err != nil {
		major = 0
	}
	minor, err := strconv.Atoi(array[1])
	if err != nil {
		minor = 0
	}
	return major, minor
}

func parseSourceFile(filename string, source io.Reader, sourceDir string, version string, interfaces InterfaceList) {
	regexpInterface := regexp.MustCompile(`\s*type\s+([A-Z]\w*)\s+interface\s+{`)
	reader := bufio.NewReader(source)
	pack := filename[len(sourceDir)+1 : strings.LastIndex(filename, "/")]
	if strings.HasSuffix(pack, "testdata") {
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
		lineNumber += 1
	}
}

func addInterfaces(version string, interfaces InterfaceList) {
	println(fmt.Sprintf("Generating interface list for version %s...", version))
	// source directory changed from 1.4
	major, minor := versionMajorMinor(version)
	sourceDir := "go/src"
	if major <= 1 && minor < 4 {
		sourceDir = "go/src/pkg"
	}
	// download compressed archive
	response, err := http.Get(URL + "go" + version + ".src.tar.gz")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	// gunzip the archive stream
	gzipReader, err := gzip.NewReader(response.Body)
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
		if strings.HasPrefix(header.Name, sourceDir) &&
			strings.HasSuffix(header.Name, ".go") &&
			!strings.HasSuffix(header.Name, "doc.go") &&
			!strings.HasSuffix(header.Name, "_test.go") {
			parseSourceFile(header.Name, tarReader, sourceDir, version, interfaces)
		}
	}
}

func printInterfaces(interfaceList InterfaceList, versions []string) {
	interfaces := make([]Interface, 0)
	for _, i := range interfaces {
		interfaces = append(interfaces, i)
	}
	sort.Sort(ByName(interfaces))
	lenName := 0
	lenVersions := make(map[string]int)
	for _, i := range interfaces {
		if len(i.Name) > lenName {
			lenName = len(i.Name)
		}
		for _, version := range versions {
			loc := interfaceList[i][version]
			lenVersion := len(loc.SourceFile) + len(loc.LineNumber) + len(loc.Link) + 8
			if lenVersions[version] < lenVersion {
				lenVersions[version] = lenVersion
			}
		}
	}
	formatLine := "%-" + strconv.Itoa(lenName) + "s"
	for _, v := range versions {
		formatLine += " %-" + strconv.Itoa(lenVersions[v])
	}
	args := []interface{}{"Interface"}
	for _, v := range versions {
		args = append(args, v)
	}
	fmt.Printf(formatLine, args...)
	separator := strings.Repeat("-", lenName) + "  "
	for _, v := range versions {
		separator += strings.Repeat("-", lenVersions[v]) + "  "
	}
	fmt.Println(separator)
	for _, i := range interfaces {
		versionLink := make(map[string]string)
		for _, v := range versions {
			versionLink[v] = "[" + interfaceList[i][v].SourceFile + " l." +
				interfaceList[i][v].LineNumber + "](" +
				interfaceList[i][v].Link + ")"
		}
		args := []interface{}{i.Name}
		for _, vl := range versionLink {
			args = append(args, vl)
		}
		fmt.Printf(formatLine, args...)
	}
}

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
