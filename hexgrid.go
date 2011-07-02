/* 
*  File: hexgrid.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 15:16:00 PDT 2011
*  Description: This file only contains godoc info.
*/

//  Package hexgrid implements an n by m hex-tile grid for use in games.
//
//  The Grid object has indexable tiles, vertices, and edges. The
//  tiles, vertices, and edges can be used to hold arbitrary objects.
//  And is navigatable as a graph no matter whether a game uses tile
//  connections, vertex connections, or both.
//  
//  The basic idea behind connections, as far as the API is concerned, is
//  that of two core concepts in graph theory *incidence* and *adjacency*.
//  Adjacency is a binary relation on objects of the same type. While,
//  incidence is a binary relation on objects of strictly different type.
//  
//  For example, tiles are adjacent other tiles with which they share one
//  edge. An object shared by two adjacent objects is incident with both
//  those objects seperately. So, continuing the example, the shared edge
//  between two adjacent tiles is incident with each tile. Similarly,
//  the endpoints of that edge are shared between the adjacent tiles,
//  so the end points (vertices) of the edge are also incident with both
//  tiles.
package hexgrid
