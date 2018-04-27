package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Editor(t *testing.T) {
	ed := newEditor()
	// const expectedOutput = `{a:1}`
	evs := []event{
		{keyPress, '{'},
		// {keyPress, 'i'},
		{keyPress, 'a'},
		{keyPress, ':'},
		{keyPress, '1'},
	}

	// Execute Events in Editor
	for _, ev := range evs {
		t.Logf("editor received event: %v\n", ev)
		t.Logf("editor before event: %v\n", ed.root)
		if err := ed.handle(ev); err != nil {
			t.Error(err)
		}
		t.Logf("editor after event: %v\n", ed.root)
		t.Logf("editor.node: %v\n", ed.node)
		t.Logf("\n")
	}

	// Validate Output Tree Structure
	var ok bool
	assert.NotNil(t, ed.root)
	json, ok := ed.root.val.(*jsonNode)
	assert.True(t, ok)
	assert.NotNil(t, json)

	// Verify Object
	assert.NotNil(t, json.val)
	obj, ok := json.val.(*objectNode)
	assert.True(t, ok, "unexpected type for ed.root.val: %s", obj)

	// Verify Correct Key / Value Pair
	assert.Equal(t, len(obj.pairs), 1)
	pair := obj.pairs[0]
	assert.NotNil(t, pair)

	// Verify Identifier Key
	assert.NotNil(t, pair.key)
	assert.NotNil(t, pair.key.val)
	id, ok := pair.key.val.(*identNode)
	assert.True(t, ok)
	assert.NotNil(t, id.val)
	assert.Equal(t, string(id.val.val), "a")

	// Verify Number Value
	assert.NotNil(t, pair.val)
	assert.NotNil(t, pair.val.val)
	val, ok := pair.val.val.(*valueNode)
	assert.True(t, ok)

	// Verify Number
	assert.NotNil(t, val.val)
	val2, ok := val.val.(*numberNode)
	assert.True(t, ok)
	assert.NotNil(t, val.val)
	assert.Equal(t, string(val2.val.val), "1")
}
