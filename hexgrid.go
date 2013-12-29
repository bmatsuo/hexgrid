/* 
*  File: hexgrid.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 15:16:00 PDT 2011
*  Description: This file only contains godoc info.
*/

/*
Package hexgrid implements an n by m hex-tile grid for use in games.

The Grid object has indexable tiles, vertices, and edges. The
tiles, vertices, and edges can be used to hold arbitrary objects.
And is navigatable as a graph no matter whether a game uses tile
connections, vertex connections, or both.

The basic idea behind connections, as far as the API is concerned, is
that of two core concepts in graph theory *incidence* and *adjacency*.
Adjacency is a binary relation on objects of the same type. While,
incidence is a symmetric binary relation on tiles and objects of other
strictly different types.

For example, tiles are adjacent other tiles with which they share one
edge. An object shared by two adjacent objects is incident with both
those objects seperately. So, continuing the example, the shared edge
between two adjacent tiles is incident with each tile. Similarly,
the endpoints of that edge are shared between the adjacent tiles.
So the end points (vertices) of the edge are also incident with both
tiles.

Because tiles are incident to both edges and vertices, instead of
referring to them as the 'incident edges' and 'incident vertices' of
the tile, they are simply the 'edges' and the 'vertices' of the tile.

Furthermore, instead of forcing the ideas of incidence and adjacency
onto edges and vertices, we simply say edges have 'ends' and vertices
have 'edges'. The typical graph theoretic notion of adjacency and
incidence does not align with the hexgrid-specific notion defined here.
*/
package hexgrid
