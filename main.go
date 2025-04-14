package main

import (
	"flag"
	"fmt"
)

const (
	numElements = 1000000
	numWorkers  = 10
)

func main() {
	var seed, size uint64
	flag.Uint64Var(&seed, "seed", 1234, "Seed value for random generator")
	flag.Uint64Var(&size, "dataset-size", 1_000_000, "Number of elements in the dataset")
	flag.Parse()

	b := New[int, int](size, seed, func() Map[int, int] { return NewSimpleMap[int, int]() })

	fmt.Println("Running Map Benchmarks")

	b.Run()
}

type SimpleMap[K comparable, V any] struct {
	data map[K]V
}

func NewSimpleMap[K comparable, V any]() *SimpleMap[K, V] {
	return &SimpleMap[K, V]{data: make(map[K]V)}
}

func (m *SimpleMap[K, V]) Get(key K) (V, bool) {
	value, ok := m.data[key]
	return value, ok
}

func (m *SimpleMap[K, V]) Set(key K, value V) {
	m.data[key] = value
}

func (m *SimpleMap[K, V]) Delete(key K) {
	delete(m.data, key)
}
