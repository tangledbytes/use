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



### Benchmarks
**Set Operation - No Sync**
```
goos: darwin
goarch: arm64
pkg: github.com/utkarsh-pro/use/pkg/storage/stupid
BenchmarkStupidSetNoSync/128B-10         	  510524	      2940 ns/op	  43.54 MB/s	     528 B/op	      14 allocs/op
BenchmarkStupidSetNoSync/256B-10         	  337562	      3178 ns/op	  80.56 MB/s	     968 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/1K-10           	  331254	      3702 ns/op	 276.60 MB/s	    3560 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/2K-10           	  277218	      4347 ns/op	 471.09 MB/s	    7272 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/4K-10           	  247633	      4940 ns/op	 829.10 MB/s	   14440 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/8K-10           	  189014	      6698 ns/op	1223.11 MB/s	   28008 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/16K-10          	  118286	      9638 ns/op	1699.89 MB/s	   59496 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/32K-10          	   81962	     15055 ns/op	2176.56 MB/s	  114792 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/64K-10          	   48633	     25361 ns/op	2584.15 MB/s	  213097 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/128K-10         	   26365	     54351 ns/op	2411.60 MB/s	  409707 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/256K-10         	    8853	    128416 ns/op	2041.36 MB/s	  802924 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/512K-10         	    6817	    194201 ns/op	2699.72 MB/s	 1589360 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/1M-10           	    3225	    425407 ns/op	2464.88 MB/s	 3162225 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/2M-10           	    1670	    661528 ns/op	3170.16 MB/s	 6307950 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/4M-10           	     760	   1405073 ns/op	2985.11 MB/s	12599407 B/op	      15 allocs/op
BenchmarkStupidSetNoSync/8M-10           	     434	   2839935 ns/op	2953.80 MB/s	25182318 B/op	      15 allocs/op
PASS
ok  	github.com/utkarsh-pro/use/pkg/storage/stupid	27.126s
```

**Set Operation - Sync**
```
goos: darwin
goarch: arm64
pkg: github.com/utkarsh-pro/use/pkg/storage/stupid
BenchmarkStupidSetSync/128B-10         	      55	  18468191 ns/op	   0.01 MB/s	     528 B/op	      14 allocs/op
BenchmarkStupidSetSync/256B-10         	      63	  18528600 ns/op	   0.01 MB/s	     968 B/op	      15 allocs/op
BenchmarkStupidSetSync/1K-10           	      62	  18497663 ns/op	   0.06 MB/s	    3560 B/op	      15 allocs/op
BenchmarkStupidSetSync/2K-10           	      63	  19417550 ns/op	   0.11 MB/s	    7272 B/op	      15 allocs/op
BenchmarkStupidSetSync/4K-10           	      63	  19855822 ns/op	   0.21 MB/s	   14440 B/op	      15 allocs/op
BenchmarkStupidSetSync/8K-10           	      63	  18257138 ns/op	   0.45 MB/s	   28008 B/op	      15 allocs/op
BenchmarkStupidSetSync/16K-10          	      63	  18745929 ns/op	   0.87 MB/s	   59496 B/op	      15 allocs/op
BenchmarkStupidSetSync/32K-10          	      64	  18535575 ns/op	   1.77 MB/s	  114792 B/op	      15 allocs/op
BenchmarkStupidSetSync/64K-10          	      64	  18512802 ns/op	   3.54 MB/s	  213096 B/op	      15 allocs/op
BenchmarkStupidSetSync/128K-10         	      62	  19197916 ns/op	   6.83 MB/s	  409706 B/op	      15 allocs/op
BenchmarkStupidSetSync/256K-10         	      63	  18500382 ns/op	  14.17 MB/s	  802956 B/op	      15 allocs/op
BenchmarkStupidSetSync/512K-10         	      58	  18841219 ns/op	  27.83 MB/s	 1589357 B/op	      15 allocs/op
BenchmarkStupidSetSync/1M-10           	      58	  19661541 ns/op	  53.33 MB/s	 3162267 B/op	      15 allocs/op
BenchmarkStupidSetSync/2M-10           	      58	  19847116 ns/op	 105.67 MB/s	 6307953 B/op	      15 allocs/op
BenchmarkStupidSetSync/4M-10           	      56	  22466036 ns/op	 186.70 MB/s	12599409 B/op	      15 allocs/op
BenchmarkStupidSetSync/8M-10           	      48	  24237001 ns/op	 346.11 MB/s	25182323 B/op	      15 allocs/op
PASS
ok  	github.com/utkarsh-pro/use/pkg/storage/stupid	19.562s
```

**Set Operation - Async Syncing**
```
goos: darwin
goarch: arm64
pkg: github.com/utkarsh-pro/use/pkg/storage/stupid
BenchmarkStupidSetAsyncSync/128B-10         	     100	  12918243 ns/op	   0.01 MB/s	     673 B/op	      15 allocs/op
BenchmarkStupidSetAsyncSync/256B-10         	     104	  12936466 ns/op	   0.02 MB/s	    1028 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/1K-10           	      91	  14281883 ns/op	   0.07 MB/s	    3626 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/2K-10           	      79	  15608401 ns/op	   0.13 MB/s	    7324 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/4K-10           	      98	  14038768 ns/op	   0.29 MB/s	   14489 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/8K-10           	      90	  14180032 ns/op	   0.58 MB/s	   28051 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/16K-10          	      81	  14509683 ns/op	   1.13 MB/s	   59523 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/32K-10          	      68	  16905240 ns/op	   1.94 MB/s	  114808 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/64K-10          	      64	  16283159 ns/op	   4.02 MB/s	  213112 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/128K-10         	      67	  16624864 ns/op	   7.88 MB/s	  409727 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/256K-10         	      64	  16924340 ns/op	  15.49 MB/s	  802960 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/512K-10         	      73	  17693342 ns/op	  29.63 MB/s	 1589403 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/1M-10           	      64	  19745871 ns/op	  53.10 MB/s	 3162248 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/2M-10           	      70	  20500542 ns/op	 102.30 MB/s	 6307972 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/4M-10           	     118	  10283815 ns/op	 407.85 MB/s	12599434 B/op	      16 allocs/op
BenchmarkStupidSetAsyncSync/8M-10           	      68	  16451892 ns/op	 509.89 MB/s	25182336 B/op	      16 allocs/op
PASS
ok  	github.com/utkarsh-pro/use/pkg/storage/stupid	27.076s
```

**Get Operation**
```
goos: darwin
goarch: arm64
pkg: github.com/utkarsh-pro/use/pkg/storage/stupid
BenchmarkStupidGet/128B-10         	  248076	      4941 ns/op	  25.91 MB/s	     360 B/op	      20 allocs/op
BenchmarkStupidGet/256B-10         	  229862	      4884 ns/op	  52.42 MB/s	     488 B/op	      20 allocs/op
BenchmarkStupidGet/1K-10           	  222211	      5003 ns/op	 204.69 MB/s	    1256 B/op	      20 allocs/op
BenchmarkStupidGet/2K-10           	  218511	      5123 ns/op	 399.79 MB/s	    2280 B/op	      20 allocs/op
BenchmarkStupidGet/4K-10           	  211711	      5297 ns/op	 773.32 MB/s	    4328 B/op	      20 allocs/op
BenchmarkStupidGet/8K-10           	  192913	      5811 ns/op	1409.83 MB/s	    8424 B/op	      20 allocs/op
BenchmarkStupidGet/16K-10          	  169510	      6750 ns/op	2427.38 MB/s	   16616 B/op	      20 allocs/op
BenchmarkStupidGet/32K-10          	  135506	      8380 ns/op	3910.19 MB/s	   33000 B/op	      20 allocs/op
BenchmarkStupidGet/64K-10          	   99961	     11797 ns/op	5555.51 MB/s	   65768 B/op	      20 allocs/op
BenchmarkStupidGet/128K-10         	   66978	     17901 ns/op	7322.18 MB/s	  131305 B/op	      20 allocs/op
BenchmarkStupidGet/256K-10         	   39435	     29523 ns/op	8879.21 MB/s	  262377 B/op	      20 allocs/op
BenchmarkStupidGet/512K-10         	   21183	     56949 ns/op	9206.26 MB/s	  524523 B/op	      20 allocs/op
BenchmarkStupidGet/1M-10           	   10000	    104689 ns/op	10016.09 MB/s	 1048813 B/op	      20 allocs/op
BenchmarkStupidGet/2M-10           	    6278	    187423 ns/op	11189.42 MB/s	 2097390 B/op	      20 allocs/op
BenchmarkStupidGet/4M-10           	    2946	    407927 ns/op	10281.99 MB/s	 4194539 B/op	      20 allocs/op
BenchmarkStupidGet/8M-10           	    1182	   1018820 ns/op	8233.65 MB/s	 8388843 B/op	      20 allocs/op
PASS
ok  	github.com/utkarsh-pro/use/pkg/storage/stupid	23.241s
```

## Updates
### 23rd October 2022
- Improved "stupid" storage type reading. Now supports parallel reads without locking.
- TLV reader and writer decoupled for ergonomics reason.
- Add HTTP transport layer.
- Add physical snapshot API without locking writes.
- Add IDs to the DB writes.

### 25th October 2022
- Add custom logger.
- Add support for data corruption detection and auto fixing during the startup.
- Add `ForEach` method to minimize code duplication and better control over the iteration.

### 4th November 2022
- Add bitset data structure.

### 12th November 2022
- Add a standard bloom filter implementation.