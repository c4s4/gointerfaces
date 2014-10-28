% id:       gointerfaces.en
% date:     2014-10-28
% title:    List of all GO Interfaces
% author:   Michel Casabianca
% email:    michel.casabianca@gmail.com
% keywords: go golang interface
% lang:     en
% toc:      no

While attending dotGo, where the buzzword was clearly *the interface*, I was wondering where I could find a list of all interfaces defined in the GO language. I found nowhere.

Thus I decided to write a little GO program to build this list. You pass the GO version on the command line and this tool:

- Downloads the source tarball.
- Parses source files to extract the interface names and line number where they are defined.
- Write this list on the console in the markdown format.

Here is the result:

INTERFACES

The project is on Github : <https://github.com/c4s4/gointerfaces>.

*Enjoy!*
