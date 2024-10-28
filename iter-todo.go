package mcfile

import(
	"fmt"
	"iter"
	// ON "github.com/fbaube/orderednodes"
)

// https://pkg.go.dev/iter@master
// func (s *mcfile) All() iter.Seq[*mcfile] 
// If an iterator requires additional configuration, the 
// constructor function can take additional config args:
// Scan returns an iterator over key-value pairs with min ≤ key ≤ max.
// func (m *Map[K, V]) Scan(min, max K) iter.Seq2[K, V]]
// When there are multiple possible iteration orders,
// the method name may indicate that order:
// Preorder returns an iterator over all nodes of the syntax tree
// beneath (and including) the specified root, in depth-first preorder,
// visiting a parent node before its children.
// func Preorder(root Node) iter.Seq[Node]

func (p *CNord) PreorderDF(f func(*CNord) bool) { }


// https://github.com/golang/go/issues/66339

// Preorder returns a go1.23 iterator over the nodes of the syntax
// tree beneath the specified root, in depth-first order.
// Each node is produced before its children.
//
// For greater control over the traversal of each subtree, use [Inspect].
// func Preorder(root Node) iter.Seq[Node]

// https://antonz.org/go-1-23/

// You can define a function that returns an iterator:

// Reversed returns an iterator that loops over a slice in reverse order.
func Reversed[V any](s []V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i]) {
				return
			}
		}
	}
}

// And a function that consumes an iterator:

// PrintAll prints all elements in a sequence.
func PrintAll[V any](s iter.Seq[V]) {
	for v := range s {
		fmt.Print(v, " ")
	}
	fmt.Println()
}

// And compose them in a convenient way:

func main() {
	s := []int{1, 2, 3, 4, 5}
	PrintAll(Reversed(s))
}

