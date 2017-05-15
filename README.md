### GO tree

### AVL tree implementation
 * very fast insertion, look ups, and deletions
 * very fast iteration
 * user friendly library
 * handles custom value comparisons
 
#### Basic Usage
```
import (
    "fmt"
	"github.com/ross-oreto/go-tree"
)

tree := container.NewString()
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
1. tree.CompareTo Comp(val interface{}) int8
2. fmt.Stringer String() string
    Note: Use tree.New() when using these interfaces.
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
BenchmarkInsertBtree-4                 5         285713580 ns/op        17600009 B/op    1200000 allocs/op
BenchmarkInsertBtree-4                 5         281322800 ns/op        17600006 B/op    1200000 allocs/op
BenchmarkInsertGtree-4                 5         296885800 ns/op        18477032 B/op    1019763 allocs/op
BenchmarkInsertGtree-4                 5         300083060 ns/op        18477016 B/op    1019763 allocs/op
BenchmarkInsertRandomBtree-4           1        1166626900 ns/op        56000016 B/op    2000000 allocs/op
BenchmarkInsertRandomBtree-4           1        1161119600 ns/op        56000016 B/op    2000000 allocs/op
BenchmarkInsertRandomGtree-4           2         916495450 ns/op        26159948 B/op    1034830 allocs/op
BenchmarkInsertRandomGtree-4           2         913124250 ns/op        26159948 B/op    1034830 allocs/op
BenchmarkGetBtree-4                    2         863909650 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGetBtree-4                    2         866082800 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGetGtree-4                    2         907872500 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGetGtree-4                    2         929887650 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkIterationBtree-4           5000            366049 ns/op               0 B/op          0 allocs/op
BenchmarkIterationBtree-4           5000            361945 ns/op               0 B/op          0 allocs/op
BenchmarkIterationGtree-4            200           8755904 ns/op               0 B/op          0 allocs/op
BenchmarkIterationGtree-4            200           8674274 ns/op               0 B/op          0 allocs/op
BenchmarkLenBtree-4             2000000000               0.42 ns/op            0 B/op          0 allocs/op
BenchmarkLenBtree-4             2000000000               0.43 ns/op            0 B/op          0 allocs/op
BenchmarkLenGtree-4             2000000000               0.44 ns/op            0 B/op          0 allocs/op
BenchmarkLenGtree-4             2000000000               0.43 ns/op            0 B/op          0 allocs/op
BenchmarkDeleteBtree-4               200           8800903 ns/op         1600048 B/op     200001 allocs/op
BenchmarkDeleteBtree-4               200           8828686 ns/op         1600048 B/op     200001 allocs/op
BenchmarkDeleteGtree-4                20          96146990 ns/op         5232844 B/op     206896 allocs/op
BenchmarkDeleteGtree-4                20          95363905 ns/op         5232846 B/op     206896 allocs/op
```

 - btree = go-tree
 - gtree = google/btree
 
Results show:
 * Insertions are faster with btree when the list is in order however when random gtree performs faster
 * Gets are faster using btree
 * Iteration much faster with btree
 * Length reporting at the same speed
 * tree deletions are faster using btree
 * A larger memory footprint using btree during insertions

