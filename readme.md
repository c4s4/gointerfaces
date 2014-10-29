GO Interfaces
=============

This program lists all public GO interfaces. To run it, just type:

```sh
go run gointerfaces.go <version>
```

Where *&lt;version>* is the GO version to parse, for instance *1.3.3*.

This compiles and runs the program that will:

- Download the GO source tarball.
- Parse all GO source files.
- Extract all interface declarations.
- Print them on the console in markdown table format.

To get result in HTML, you can pipe tge output to *pandoc* to convert to HTML :

```
go run gointerfaces.go 1.3.3 | pandoc -f markdown -t html
```

You may see the result for *1.3.3* release on this page: <http://sweetohm.net/html/gointerfaces.en.html>.

*Enjoy!*