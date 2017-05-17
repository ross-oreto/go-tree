### GO tree
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/ross-oreto/go-tree/master/LICENSE)
[![Build Status](https://travis-ci.org/ross-oreto/go-tree.svg?branch=master)](https://travis-ci.org/ross-oreto/go-tree)
[![Go Report Card](https://goreportcard.com/badge/ross-oreto/go-tree)](https://goreportcard.com/report/ross-oreto/go-tree)
[![GoDoc](https://godoc.org/github.com/ross-oreto/go-tree?status.svg)](https://godoc.org/github.com/ross-oreto/go-tree)


### AVL tree implementation
 * very fast insertion, look ups, and deletions
 * very fast iteration
 
#### Basic Usage
```
import (
"fmt"
"github.com/ross-oreto/go-tree"
)

btree := tree.New()
btree.Insert(StringVal("Oreto")).Insert(StringVal("Michael")).Insert(StringVal("Ross"))
fmt.Println(btree.Values())
```

#### Retrieve
```
btree.Get(StringVal("Ross"))
```

#### Delete
```
btree.Delete(StringVal("Ross"))
```

#### Clear
```
btree.Init()
```

#### Val type
Values entered into the tree must implement the Val interface Comp method
 - Comp(val Val) int8

```
type TestKey1 struct {
	Name string
}
// Comp returns 1 if key > val, -1 if key < val and 0 if key equal to val
func (key TestKey1) Comp(val Val) int8 {
	var c int8
	tk := val.(TestKey1)
	if key.Name > tk.Name {
		c = 1
	} else if key.Name < tk.Name {
		c = -1
	}
	return c
}
```

#### Performance
Benchmarks ran against "github.com/google/btree" using a random tree with 1 million inserts
```
BenchmarkBtree_Insert-4                        3         357906766 ns/op        24000016 B/op    1333333 allocs/op
BenchmarkBtree_Insert-4                        3         356572333 ns/op        24000016 B/op    1333333 allocs/op
BenchmarkGtree_Insert-4                        3         364911600 ns/op        25461704 B/op    1032939 allocs/op
BenchmarkGtree_Insert-4                        3         366578900 ns/op        25461698 B/op    1032939 allocs/op
BenchmarkBtree_InsertRandom-4                  1        1354908200 ns/op        56000032 B/op    2000000 allocs/op
BenchmarkBtree_InsertRandom-4                  1        1351906700 ns/op        56000000 B/op    2000000 allocs/op
BenchmarkGtree_InsertRandom-4                  1        1049704500 ns/op        44105176 B/op    1068918 allocs/op
BenchmarkGtree_InsertRandom-4                  1        1070716500 ns/op        44105176 B/op    1068918 allocs/op
BenchmarkBtree_Get-4                           2         984660150 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkBtree_Get-4                           2         969647750 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGtree_Get-4                           1        1065713900 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkGtree_Get-4                           1        1064712600 ns/op         8000000 B/op    1000000 allocs/op
BenchmarkBtree_Iteration-4                  2000            591395 ns/op               0 B/op          0 allocs/op
BenchmarkBtree_Iteration-4                  2000            590395 ns/op               0 B/op          0 allocs/op
BenchmarkGtree_Iteration-4                   100          10356963 ns/op               0 B/op          0 allocs/op
BenchmarkGtree_Iteration-4                   100          10346943 ns/op               0 B/op          0 allocs/op
BenchmarkBtree_Len-4                    2000000000               0.59 ns/op            0 B/op          0 allocs/op
BenchmarkBtree_Len-4                    2000000000               0.59 ns/op            0 B/op          0 allocs/op
BenchmarkGtree_Len-4                    2000000000               0.59 ns/op            0 B/op          0 allocs/op
BenchmarkGtree_Len-4                    2000000000               0.59 ns/op            0 B/op          0 allocs/op
BenchmarkBtree_Delete-4                       10         122482200 ns/op         6400001 B/op     300000 allocs/op
BenchmarkBtree_Delete-4                       10         121581680 ns/op         6400001 B/op     300000 allocs/op
BenchmarkGtree_Delete-4                       10         115377360 ns/op         5217883 B/op     206879 allocs/op
BenchmarkGtree_Delete-4                       10         113776200 ns/op         5217884 B/op     206879 allocs/op
```

 - btree = go-tree
 - gtree = google/btree
 
Results show:
 * Insertions are faster with btree when the list is in order however when random gtree performs faster
 * Gets are faster using btree
 * Iteration much faster with btree
 * Length reporting at the same speed
 * tree deletions are slightly faster using gtree
 * A larger memory footprint using btree during insertions and deletions

