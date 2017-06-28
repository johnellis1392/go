package btree

// Go Build Template:
// https://github.com/thockin/go-build-template

// Slides for Golang best practices
// https://talks.golang.org/2013/bestpractices.slide#12

// Simple analog to pattern-matching?
// func example(value interface{}) {
// 	switch x := value.(type) {
// 	case string:
// 		fmt.Println(x)
// 	case int:
// 		fmt.Println(x)
// 	default:
// 		// ...
// 	}
// }

// Types wrapping functions:
// type TestFunc func(float64) float64
// var (
//     ident = TestFunc(func(x float64) float64 { return x })
//     sin = TestFun(math.sin)
// )

// New Construct a new BTree root node
func New(n int32) BTree {
	return BLeaf{
		Order:    n,
		Elements: make(map[string]interface{}, n),
	}
}

// BTree BTree utility function interface
type BTree interface {
	Insert(key string, value interface{}) BTree
	Delete(key string) BTree
	Find(key string) interface{}
}

// BNode B-Tree node
type BNode struct {
	Order    int32
	Elements map[string]BTree
}

// Insert ...
func (b BNode) Insert(key string, value interface{}) BTree {
	return b
}

// Delete ...
func (b BNode) Delete(key string) BTree {
	return b
}

// Find ...
func (b BNode) Find(key string) interface{} {
	return nil
}

// BLeaf B-Tree Leaf Node
type BLeaf struct {
	Order    int32
	Elements map[string]interface{}
}

// Insert ...
func (b BLeaf) Insert(key string, value interface{}) BTree {
	b.Elements[key] = value
	return b
}

// Delete ...
func (b BLeaf) Delete(key string) BTree {
	return b
}

// Find ...
func (b BLeaf) Find(key string) interface{} {
	return nil
}
