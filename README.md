### GO tree

### AVL tree implementation
 * very fast look ups
 * very fast iteration
 * user friendly library
 * handles custom value comparisons
 
#### Usage
```
import (
    "fmt"
	"github.com/ross-oreto/go-tree"
)

tree := tree.New()
tree.Insert('Oreto').Insert('Michael').Insert('Ross')
fmt.Println(tree.Values())
```
