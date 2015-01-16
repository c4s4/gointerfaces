GO Interfaces
=============

This program lists all public GO interfaces. To run it, just type:

```sh
go run gointerfaces.go <versions>
```

Where *&lt;versions>* is a list of GO versions, for instance *1.0.3 1.1.2 1.2.2 1.3.3 1.4*.

This compiles and runs the program that will:

- Download the GO source tarballs.
- Parse all GO source files.
- Extract all interface declarations for GO versions.
- Print them on the console in markdown table format.

To get result in HTML, you can pipe the output to *pandoc*:

```sh
go run gointerfaces.go 1.0.3 1.1.2 1.2.2 1.3.3 1.4.1 | pandoc -f markdown -t html
```

You may see the result on this page: <http://sweetohm.net/html/gointerfaces.en.html>.

*Enjoy!*

History
-------

- **1.2.1** (*2015-01-16*): Updated for GO *1.4.1*.
- **1.2.0** (*2015-01-12*): Added generation for a list of versions.
- **1.1.0** (*2014-12-24*): Fixed for GO 1.4 and links point to sources on Github.
- **1.0.0** (*2014-10-29*): First release.
