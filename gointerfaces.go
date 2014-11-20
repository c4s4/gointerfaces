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
    SOURCE_URL = "https://code.google.com/p/go/source/browse/"
)

type Interface struct {
    Name       string
    Package    string
    SourceFile string
    LineNumber string
    Link       string
}

type ByName []Interface

func (b ByName) Len() int           {return len(b)}
func (b ByName) Swap(i, j int)      {b[i], b[j] = b[j], b[i]}
func (b ByName) Less(i, j int) bool {return b[i].Name < b[j].Name}

func majMin(v string) (int, int) {
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

func parseSourceFile(filename string, source io.Reader, sourceDir string) []Interface {
    regexpInterface := regexp.MustCompile(`\s*type\s+([A-Z]\w*)\s+interface\s+{`)
    interfaces := make([]Interface, 0)
    reader := bufio.NewReader(source)
    pack := filename[len(sourceDir)+1:strings.LastIndex(filename, "/")]
    if strings.HasSuffix(pack, "testdata") {
        return nil
    }
    lineNumber := 1
    for {
        line, err := reader.ReadBytes('\n')
        if err != nil && err != io.EOF {
            panic("Error parsing source file")
        }
        matches := regexpInterface.FindSubmatch(line)
        if len(matches) > 0 {
            interf := Interface {
                Name:       string(matches[1]),
                Package:    string(pack),
                SourceFile: filename[3:],
                LineNumber: strconv.Itoa(lineNumber),
                Link:       SOURCE_URL+filename[3:]+"?name=release#"+strconv.Itoa(lineNumber),
            }
            interfaces = append(interfaces, interf)
        }
        if err == io.EOF {
            break
        }
        lineNumber += 1
    }
    return interfaces
}

func printInterfaces(interfaces []Interface) {
    lenName := 0
    lenPackage := 0
    lenSourceFile := 0
    lenLineNumber := 0
    for _, i := range interfaces {
        if len(i.Name)+len(i.Link)+4 > lenName {
            lenName = len(i.Name)+len(i.Link)+4
        }
        if len(i.Package) > lenPackage {
            lenPackage = len(i.Package)
        }
        if len(i.SourceFile) > lenSourceFile {
            lenSourceFile = len(i.SourceFile)
        }
        if len(i.LineNumber) > lenLineNumber {
            lenLineNumber = len(i.LineNumber)
        }
     }
    formatLine := "%-" + strconv.Itoa(lenName) + "s  %-" + strconv.Itoa(lenPackage) +
        "s  %-" + strconv.Itoa(lenSourceFile) + "s  %-" + strconv.Itoa(lenLineNumber) +
        "s\n"
    fmt.Printf(formatLine, "Interface", "Package", "Source File", "Line")
    separator := strings.Repeat("-", lenName) + "  " + strings.Repeat("-", lenPackage) +
        "  " + strings.Repeat("-", lenSourceFile) + "  " + strings.Repeat("-", lenLineNumber)
    fmt.Println(separator)
    for _, i := range interfaces {
        link := "["+i.Name+"]("+i.Link+")"
        fmt.Printf(formatLine, link, i.Package, i.SourceFile, i.LineNumber)
    }
}

func main() {
    // read version on command line
    if len(os.Args) != 2 {
        panic("Must pass go version on command line")
    }
    version := os.Args[1]
    // source directory changed from 1.4
    major, minor := majMin(version)
    sourceDir := "go/src"
    if major <= 1 && minor < 4 {
        sourceDir = "go/src/pkg"
    }
    // download compressed archive
    println("Downloading archive...")
    response, err := http.Get(URL+"go"+version+".src.tar.gz")
    if err!=nil {
        panic(err)
    }
    defer response.Body.Close()
    // gunzip the archive stream
    gzipReader, err := gzip.NewReader(response.Body)
    if err != nil {
        panic(err)
    }
    // parse tar source files in source dir
    println("Parsing archive...")
    tarReader := tar.NewReader(gzipReader)
    interfaces := make([]Interface, 0)
    for {
        header, err := tarReader.Next()
        if err != nil {
            break
        }
        if strings.HasPrefix(header.Name, sourceDir) &&
           strings.HasSuffix(header.Name, ".go") &&
           !strings.HasSuffix(header.Name, "doc.go") &&
           !strings.HasSuffix(header.Name, "_test.go") {
            newInterfaces := parseSourceFile(header.Name, tarReader, sourceDir)
            interfaces = append(interfaces, newInterfaces...)
        }
    }
    // print the result
    sort.Sort(ByName(interfaces))
    printInterfaces(interfaces)
}
