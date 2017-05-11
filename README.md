### GO tree

### AVL tree implementation
 * very fast look ups
 * very fast iteration
 * user friendly library
 * handles custom value comparisons
 
#### Basic Usage
```
import (
    "fmt"
	"github.com/ross-oreto/go-tree"
)

tree := tree.New()
tree.Insert('Oreto').Insert('Michael').Insert('Ross')
fmt.Println(tree.Values())
```

#### Retrieve
```
tree.Get('Ross')
```

#### Delete
```
tree.Delete('Ross')
```

#### Clear
```
tree.Init()
```

#### Custom Values
Custom structs can be inserted by implementing one of two interfaces
1. tree.Comparer Comp(val interface{}) int8
2. fmt.Stringer String() string
```
type TestKey1 struct {
	Name string
}
func (testkey TestKey1) Comp(val interface{}) int8 {
	var c int8 = 0
	tk := val.(TestKey1)
	if testkey.Name > tk.Name {
		c = 1
	} else if testkey.Name < tk.Name {
		c = -1
	}
	return c
}
type TestKey2 struct {
	Name string
}
func (testkey TestKey2) String() string {
	return testkey.Name
}
```

#### Performance
Benchmarks ran against "github.com/google/btree" using a random tree with 1 million inserts
```
go test -v -benchmem -count 2 -bench .
BenchmarkInsertBtree-4                 5         310280200 ns/op        17600006 B/op    1200000 allocs/op
BenchmarkInsertBtree-4                 3         338451700 ns/op        24000010 B/op    1333333 allocs/op
BenchmarkInsertGtree-4                 5         294047100 ns/op        18477016 B/op    1019763 allocs/op
BenchmarkInsertGtree-4                 5         294281040 ns/op        18477016 B/op    1019763 allocs/op
BenchmarkInsertRandomBtree-4           1        1185917200 ns/op        56000032 B/op    2000000 allocs/op
BenchmarkInsertRandomBtree-4           1        1180782500 ns/op        56000000 B/op    2000000 allocs/op
BenchmarkInsertRandomGtree-4           2         953160650 ns/op        26159948 B/op    1034830 allocs/op
BenchmarkInsertRandomGtree-4           1        1093043600 ns/op        44105176 B/op    1068918 allocs/op
BenchmarkGetBtree-4                    2         921387100 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGetBtree-4                    2         871837150 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGetGtree-4                    2         952132550 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGetGtree-4                    1        1027492200 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkIterationBtree-4           5000            383567 ns/op               0 B/op          0 allocs/op
BenchmarkIterationBtree-4           5000            388986 ns/op               0 B/op          0 allocs/op
BenchmarkIterationGtree-4            200           8280458 ns/op               0 B/op          0 allocs/op
BenchmarkIterationGtree-4            200           8185408 ns/op               0 B/op          0 allocs/op
BenchmarkLenBtree-4             2000000000               0.42 ns/op            0 B/op          0 allocs/op
BenchmarkLenBtree-4             2000000000               0.43 ns/op            0 B/op          0 allocs/op
BenchmarkLenGtree-4             2000000000               0.43 ns/op            0 B/op          0 allocs/op
BenchmarkLenGtree-4             2000000000               0.42 ns/op            0 B/op          0 allocs/op
BenchmarkDeleteBtree-4                 1        1084307700 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkDeleteBtree-4                50          31445050 ns/op         8000034 B/op    1000000 allocs/op
BenchmarkDeleteGtree-4                50          28537434 ns/op         8000007 B/op    1000000 allocs/op
BenchmarkDeleteGtree-4                50          28533440 ns/op         8000007 B/op    1000000 allocs/op
```

 - btree = go-tree
 - gtree = google/btree
 
Results show:
 * Insertions are roughly the same with gtree slightly faster
 * Gets are faster using btree
 * Iteration much faster with btree
 * Length reporting at the same speed
 * Deletions are much faster using gtree
 * A larger memory footprint using btree during insertions

