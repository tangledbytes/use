# USE

USE is an stupid storage engine. The sole purpose of the existence of this storage engine is so that I can play with multiple things in a single project. Plan for iteration 2 (iteration 1 of the project can be found under tag `iter1`).

What does USE stand for? It's a recursive name U -> USE, S -> Storage, E -> Engine
## Goals
- [ ] A storage engine which allows support for storing `<key, value>` pairs. The storage engine should use LSM tree backed by SSTable (on disk) and memtable (in memory).
  - [ ] memtable won't be multidimensional in nature and implemented using simple skiplists.
  - [ ] Should support primary and secondary indexing.
	- [ ] primary key indexing would be clustered in nature. No Heap File.
	- [ ] seconday key indexing will not be clustered in nature. They should rather point to the primary key in the primary key index.
- [ ] A simple transport layer which can perform operations like:
  - [ ] `GET <key>`.
  - [ ] `SET <key> <value>`.
  - [ ] `DELETE [key]`.
  - [ ] `LIST`.

## Out of Scope
These are several items that are out of scope of iteration 2 and might be (hopefully) will be added in later iteration.
1. Transactions.
2. Storage space optimisations.
3. SQL or any other sophisticated querying mechanism.
4. Distributed anything.
5. Go GC would be buzz killer but will not be optimised in this iteration (and most likely in any coming iterations).

# Updates