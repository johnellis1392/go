package main

import (
	"fmt"
	"unicode/utf8"
)

type scanner struct {
	queue   queue
	weights map[rune]float32
	n       int
}

type queue []*node

func (q *queue) sort() {
	for i := 0; i < len(q); i++ {
	}
}

func (q *queue) enqueue(n *node) {
	q = append(q, n)
	q.sort()
}

func (q *queue) dequeue() *node {
	if len(q) == 0 {
		return nil
	}

	if len(q) == 1 {
		n := q[0]
		q[0] = nil
		return n
	}

	n := q[0]
	ni := q[len(q)-1]
	q[len(q)-1] = nil
	q[0] = ni
	q.sort()
	return n
}

type node struct {
	w   float32
	val rune
}

type huffmanTree struct {
	root node
}

func runify(input string) []rune {
	var rs []rune
	for pos := 0; pos <= len(input); {
		r, dw := utf8.DecodeRuneInString(input[pos:])
		pos += dw
		rs = append(rs, r)
	}
	return rs
}

func scan(input []rune) scanner {
	var s scanner
	s.n = len(input)
	for _, r := range input {
		if val, ok := s.weights[r]; ok {
			s.weights[r] = val + 1
		} else {
			s.weights[r] = 1
		}
	}

	// Convert Counts into Probabilities
	for k, v := range s.weights {
		s.weights[k] = v / float32(s.n)
	}

	return s
}

// encode calculates the Huffman Tree for the given token sequence.
func encode(input string) (*huffmanTree, error) {
	// rs := runify(input)
	return nil, nil
}

func main() {
	fmt.Println("Running...")
}
