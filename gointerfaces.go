package main

import (
    "archive/tar"
    "bufio"
    "bytes"
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

const URL = "https://storage.googleapis.com/golang/"

type Interface struct {
    Name       string
    Package    string
    SourceFile string
    LineNumber string
}

type ByName []Interface

func (b ByName) Len() int           {return len(b)}
func (b ByName) Swap(i, j int)      {b[i], b[j] = b[j], b[i]}
func (b ByName) Less(i, j int) bool {return b[i].Name < b[j].Name}

func parseSourceFile(filename string, source io.Reader) []Interface {
    regexpInterface := regexp.MustCompile(`\s*type\s+([A-Z]\w*)\s+interface\s+{`)
    interfaces := make([]Interface, 0)
    reader := bufio.NewReader(source)
    pack := filename[11:strings.LastIndex(filename, "/")]
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
                SourceFile: filename,
                LineNumber: strconv.Itoa(lineNumber),
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
    var buffer bytes.Buffer
    buffer.WriteString("<table><tr><th>Interface</th><th>Package</th><th>Source</th><th>Line</th></tr>")
    for _, interf := range interfaces {
        buffer.WriteString("<tr><td>")
        buffer.WriteString(interf.Name)
        buffer.WriteString("</td><td>")
        buffer.WriteString(interf.Package)
        buffer.WriteString("</td><td>")
        buffer.WriteString(interf.SourceFile)
        buffer.WriteString("</td><td>")
        buffer.WriteString(interf.LineNumber)
        buffer.WriteString("</td></tr>")
    }
    buffer.WriteString("</table>")
    fmt.Println(buffer.String())
}

func main() {
    // read version on command line
    if len(os.Args) != 2 {
        fmt.Println("Must pass go version on command line")
        os.Exit(1)
    }
    version := os.Args[1]
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
    // parse tar source files in go/src/pkg
    println("Parsing archive...")
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
    printInterfaces(interfaces)
}
