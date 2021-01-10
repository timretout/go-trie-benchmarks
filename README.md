# go-trie-benchmarks

## Aim

Benchmarking the performance of different Golang trie implementations for the
specific purpose of UK postcode autocompletion.

Currently it profiles these implementations:

- badgerodon "github.com/badgerodon/collections/trie"
- claudiu "github.com/claudiu/trie"
- derekparker "github.com/derekparker/trie"
- dghubble "github.com/dghubble/trie"
- timretout "github.com/timretout/trie" (mine!)
- viant "github.com/viant/ptrie"

## Method

These benchmarks look specifically at the problem of indexing all UK postcodes.
To run these you need a copy of ONSPD (currently the May 2020 release) extracted
in an ONSPD directory next to this one.

This consists of approximately 2.6 million UK postcodes (short strings of 5 to 7
bytes using alphanumeric characters, sharing many common prefixes).

Significantly faster results were seen with sequential accesses (i.e. testing
keys in insertion order), so random access was also profiled.

### Results

Results from my machine.

Creation of the trie:

 BenchmarkImportONSPD/badgerodon-8  	       1	10333619248 ns/op	6595485424 B/op	 8122542 allocs/op
 BenchmarkImportONSPD/claudiu-8     	       1	1133079749 ns/op	431894488 B/op	 6643031 allocs/op
 BenchmarkImportONSPD/derekparker-8 	       1	1958739912 ns/op	1205441648 B/op	14128892 allocs/op
 BenchmarkImportONSPD/dghubble-8    	       2	 772355485 ns/op	220278120 B/op	 3553773 allocs/op
 BenchmarkImportONSPD/timretout-8   	       4	 340029844 ns/op	148074024 B/op	 3004380 allocs/op
 BenchmarkImportONSPD/viant-8       	       1	1673385972 ns/op	406276248 B/op	 9251588 allocs/op

Sequential exists checks:

 BenchmarkONSPDSequentialExists/badgerodon-8         	 4435724	       248 ns/op	      24 B/op	       2 allocs/op
 BenchmarkONSPDSequentialExists/claudiu-8            	 8300227	       146 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDSequentialExists/derekparker-8        	 4512182	       284 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDSequentialExists/dghubble-8           	 6729548	       177 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDSequentialExists/timretout-8          	16112839	        74.6 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDSequentialExists/viant-8              	 3550863	       336 ns/op	       8 B/op	       1 allocs/op

Random exists checks:

 BenchmarkONSPDRandomExists/badgerodon-8             	 1480521	       769 ns/op	      24 B/op	       2 allocs/op
 BenchmarkONSPDRandomExists/claudiu-8                	  915116	      1272 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDRandomExists/derekparker-8            	  619204	      1768 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDRandomExists/dghubble-8               	  908576	      1215 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDRandomExists/timretout-8              	 1358100	       827 ns/op	       0 B/op	       0 allocs/op
 BenchmarkONSPDRandomExists/viant-8                  	  513060	      3127 ns/op	       8 B/op	       1 allocs/op

## Analysis

My implementation ("timretout") uses least memory (148MB, compared to 220MB for
the nearest competitor), has fastest sequential access (~75ms) and
second-fastest random access (827ms).  It cheats by using knowledge of the
alphabet - this is probably the main difference between this and the other
implementations.  It should be possible to generalize the implementation, and
then I could come back and rebenchmark...

Badgerodon (a.k.a. go-collections) has surprisingly fast random access, through
a trick where it converts the input string to bytes and then using a simple loop
to look up the nodes.  On the other hand it uses 6.5GB memory, which is
excessive.  Sequential access seems less fast.

Claudiu/dghubble's implementations seem to have comparable performance for this
use case, with claudiu's using more memory.

derekparker uses 1.2GB RAM and has relatively low performance. viant seems
slowest.

## Evaluation

Beware, this is a very specific use case, tested on one machine!

An analysis of tries for general use would look at more varied data sets.

## Conclusions

Knowing your data lets you write better data structures.
