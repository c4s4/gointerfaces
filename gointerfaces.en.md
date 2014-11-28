% id:       gointerfaces.en
% date:     2014-10-28
% title:    List of all GO Interfaces
% author:   Michel Casabianca
% email:    michel.casabianca@gmail.com
% keywords: go golang interface
% lang:     en
% toc:      no

While attending dotGo, where the buzzword was clearly *the interface*, I was wondering where I could find a list of all interfaces defined in the GO language. I found nowhere.

Thus I decided to write a little GO program that would;

- Downloads the GO source tarball for a given version.
- Parses source files to extract the interface names and line number where they are defined.
- Write this list on the console in the markdown format.

The project is on Github: <https://github.com/c4s4/gointerfaces>.

Here is the result for release *VERSION*:

INTERFACES

You may find a discussion on these interfaces on this page: <http://mwholt.blogspot.fr/2014/08/maximizing-use-of-interfaces-in-go.html>.

*Enjoy!*
