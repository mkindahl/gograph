gograph
=======

Library for working with graphs using the Go language.

The library is intended to support working with graphs, as in [Graph
Theory](https://en.wikipedia.org/wiki/Graph_theory). Each graph
consists of a set of vertices and a set of edges, where each edge
connect two vertices. If the order of the vertices for an edge is
important, the graph is *directed*, otherwise it is *undirected*.

Currently, there is only support for directed graphs, but the plan is
to add support for undirected graphs as well.

Directed Graphs
---------------

Using this library, directed graphs can be constructed by creating a
graph and adding vertices and edges to it. The vertices and edges of
the graphs can be iterated over and processed in different orders.

Currently, there is support for processing vertices in arbitrary
order, in depth-first forest order, or in topological order.

Disjoint-Set
------------

The library also contain an implementation of the [Disjoin-Set
Forest](https://en.wikipedia.org/wiki/Disjoint-set_data_structure)
implementation using union by rank and path compression as described
in the book "Introduction to Algorithms" by Cormen et.al.

Disjoint sets can be used to efficiently compute connected components
of an undirected graph, and is added here since the intention is to
support undirected graphs later.
