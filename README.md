# USE

USE is an insanely stupid storage engine. The sole purpose of the existence of this storage engine is so that I can play with multiple things in a single project. Plan for iteration 1 is to have:
1. A stupid simple storage engine, which stores key value pairs in a single file. The data is stored in TLV format. Data file is nothing but an append only log. Writes are fast, reads are excrutiatingly slow. No plans to improve any of this in this iteration.
2. A simple HTTP based transporation layer for the storage engine to expose its API. No TLS, no complex query language, nothing. Just a simple HTTP query based API.
3. Single master replication. Nothing fancy here either. Follower nodes will use a join token to join the cluster. Once the token is verified, they can choose either asynchronous or synchronous mode of replication.

That's it for the first iteration. For the future iteration (if that ever happens):
1. I might attempt to improve the core storage engine. Something based on the LSM tree. I also plan to have multiple implementation of the core storage engine, some based on simple SSTable, some on LSM Tree, some on B Tree, some on B+Tree, etc.
2. I would want to have a decent enough query parser (not SQL for this). Maybe a compiled or JIT based.
3. For replication, I don't know man. Multi-leader and leaderless are scary as f**k. I may not be able to muster enough courage to implement it in a learning project.
4. Sharding? Don't know yet.

Whta does USE stand for? It's a recursive name U -> USE, S -> Storage, E -> Engine

## Updates
### 23rd October 2022
- Improved "stupid" storage type reading. Now supports parallel reads without locking.
- TLV reader and writer decoupled for ergonomics reason.
- Add HTTP transport layer.
- Add physical snapshot API without locking writes.
- Add IDs to the DB writes.