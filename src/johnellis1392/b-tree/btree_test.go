package btree

import (
	"strconv"
	"testing"
)

func Test_BTree_New_ShouldCreateNewBTree(t *testing.T) {
	b := New(10)
	switch x := b.(type) {
	case BLeaf:
		if x.Order != 10 {
			t.Errorf("Expected order=10, got: %q", x.Order)
		}
		break
	default:
		t.Errorf("Expected new BLeaf{}, got: %q", b)
	}
}

func Test_BLeaf_Insert_ShouldInsertANewValue(t *testing.T) {
	b := New(10).(BLeaf)
	const (
		key   = "key1"
		value = "value1"
	)

	b.Insert(key, value)
	if b.Elements[key] != value {
		t.Errorf("Insert value failed, values: %q", b.Elements)
	}
}

func Test_Bleaf_Insert_ShouldReturnNewBNodeWhenOrderIsExceeded(t *testing.T) {
	var b BTree = New(2)
	for i := 1; i <= 3; i++ {
		b = b.Insert(strconv.Itoa(i), i)
	}

	switch b.(type) {
	case BNode:
		// Pass
		break
	default:
		t.Errorf("Expected element of type BNode, got: %q", b)
	}
}
