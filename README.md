*hexgrid version 0.3_2*

ABOUT hexgrid
=============

Package hexgrid implements an n by m hexagonal grid. A grid can be
used in games with a hex-tile layout such as the games of Hex,
Settlers of Catan, and Sid Meier's Civilization 5.

The Grid object has indexable tiles, vertices, and edges. The
tiles, vertices, and edges can be used to hold arbitrary objects.
And is navigatable as a graph no matter whether a game uses tile
connections, vertex connections, or both.

INSTALLATION
============

Easiest installation is through goinstall

    goinstall github.com/bmatsuo/hexgrid

Or, alternatively, you can clone the repository and install locally.

    git clone git://github.com/bmatsuo/hexgrid.git
    cd hexgrid
    gomake install

DOCUMENTATION
=============

The best way to view the documentation is by running a godoc http
server.

    godoc -http=:6060

Then, in a web browser, visit the url
http://localhost:6060/pkg/github.com/bmatsuo/hexgrid/

AUTHOR
======

**Bryan Matsuo** <bmatsuo@soe.ucsc.edu>

COPYRIGHT & LICENSE
===================

(C) 2011 Bryan Matsuo 

TODO - add licensing information!
