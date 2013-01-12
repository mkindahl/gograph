Package `gograph`
=================

Library for working with graphs using the Go language.

The library is distributed under the BSD License. See license text
below.

You can use the `go get` command to download and install the packages:

>   go get github.com/mkindahl/gograph/directed
>   go get github.com/mkindahl/gograph/djs

Description
===========

The library is intended to support working with graphs, as in [Graph
Theory](https://en.wikipedia.org/wiki/Graph_theory). Each graph
consists of a set of vertices and a set of edges, where each edge
connect two vertices. If the order of the vertices for an edge is
important, the graph is *directed*, otherwise it is *undirected*.

Currently, there is only support for directed graphs and
disjoint-sets, but the plan is to add support for undirected graphs as
well.


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


BSD License Text
================

Copyright (c) 2013, Mats Kindahl.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

>    Redistributions of source code must retain the above copyright
>    notice, this list of conditions and the following disclaimer.
>    Redistributions in binary form must reproduce the above copyright
>    notice, this list of conditions and the following disclaimer in
>    the documentation and/or other materials provided with the
>    distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
