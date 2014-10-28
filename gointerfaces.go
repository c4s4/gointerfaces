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
    "strings"
)

const URL = "https://storage.googleapis.com/golang/"

type Interface struct {
    Name       string
    Package    string
    SourceFile string
    LineNumber int
}

type ByName []Interface

func (b ByName) Len() int           {return len(b)}
func (b ByName) Swap(i, j int)      {b[i], b[j] = b[j], b[i]}
func (b ByName) Less(i, j int) bool {return b[i].Name < b[j].Name}

func parseSourceFile(filename string, source io.Reader) []Interface {
    regexpPackage := regexp.MustCompile(`package\s+(.*)`)
    regexpInterface := regexp.MustCompile(`\s*type\s+([A-Z]\w*)\s+interface\s+{`)
    interfaces := make([]Interface, 0)
    reader := bufio.NewReader(source)
    var pack []byte
    lineNumber := 0
    for {
        line, err := reader.ReadBytes('\n')
        if err != nil && err != io.EOF {
            panic("Error parsing source file")
        }
        packageMatch := regexpPackage.FindSubmatch(line)
        if len(packageMatch) > 0 {
            pack = packageMatch[1]
        } else {
            matches := regexpInterface.FindSubmatch(line)
            if len(matches) > 0 {
                interf := Interface {
                    Name:       string(matches[1]),
                    Package:    string(pack),
                    SourceFile: filename,
                    LineNumber: lineNumber,
                }
                interfaces = append(interfaces, interf)
            }
        }
        if err == io.EOF {
            break
        }
        lineNumber += 1
    }
    return interfaces
}

func main() {
    // read version on command line
    if len(os.Args) != 2 {
        panic("Must pass go version on command line")
    }
    version := os.Args[1]
    // download compressed archive
    fmt.Println("Downloading archive...")
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
    // parse tar source files in go/src/pkg
    fmt.Println("Parsing archive...")
    tarReader := tar.NewReader(gzipReader)
    interfaces := make([]Interface, 0)
    for {
        header, err := tarReader.Next()
        if err != nil {
            break
        }
        if strings.HasPrefix(header.Name, "go/src/pkg") &&
           strings.HasSuffix(header.Name, ".go") &&
           !strings.HasSuffix(header.Name, "doc.go") {
            newInterfaces := parseSourceFile(header.Name, tarReader)
            interfaces = append(interfaces, newInterfaces...)
        }
    }
    // print the result
    sort.Sort(ByName(interfaces))
    for _, interf := range interfaces {
        println(interf.Name)
    }
}
