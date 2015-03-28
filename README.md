GO Interfaces
=============

This program lists all public GO interfaces. To run it, just type:

    go run gointerfaces.go <versions>

Where *&lt;versions>* is a list of GO versions, for instance *1.0.3 1.1.2 1.2.2 1.3.3 1.4*.

This compiles and runs the program that will:

- Download the GO source tarballs.
- Parse all GO source files.
- Extract all interface declarations for GO versions.
- Print them on the console in markdown table format.

To get result in HTML, you can pipe the output to *pandoc*:

    go run gointerfaces.go 1.4.1 | pandoc -f markdown -t html

You may see the result on this page: <http://sweetohm.net/html/gointerfaces.en.html>.

*Enjoy!*
