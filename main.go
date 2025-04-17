package main

import (
	"flag"
	"fmt"

	cocroach "github.com/cockroachdb/swiss"
	crn4 "github.com/crn4/swiss"
	dolthub "github.com/dolthub/swiss"
)

func main() {
	var (
		seed, size         uint64
		mapType            string
		keyType, valueType string
	)
	flag.Uint64Var(&seed, "seed", 1234, "Seed value for random generator")
	flag.Uint64Var(&size, "dataset-size", 1_000_000, "Number of elements in the dataset")
	flag.StringVar(&mapType, "map-type", "std", "std/cocroach/crn4/dolthub")
	flag.StringVar(&keyType, "key-type", "int", "int/string/struct{}")
	flag.StringVar(&valueType, "value-type", "int", "int/string/struct{}")
	flag.Parse()

	build := func() Map[int, int] { return NewSimpleMap[int, int]() }
	switch mapType {
	case "cocroach":
		build = func() Map[int, int] { return NewCocroachMap[int, int]() }
	case "crn4":
		build = func() Map[int, int] { return NewCRN4Map[int, int]() }
	case "dolthub":
		build = func() Map[int, int] { return NewDolthubMap[int, int]() }
	}
	b := New[int, int](size, seed, build)

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

type Cocroach[K comparable, V any] struct {
	data *cocroach.Map[K, V]
}

func NewCocroachMap[K comparable, V any]() *Cocroach[K, V] {
	return &Cocroach[K, V]{data: cocroach.New[K, V](0)}
}

func (m *Cocroach[K, V]) Get(key K) (V, bool) {
	value, ok := m.data.Get(key)
	return value, ok
}

func (m *Cocroach[K, V]) Set(key K, value V) {
	m.data.Put(key, value)
}

func (m *Cocroach[K, V]) Delete(key K) {
	m.data.Delete(key)
}

type CRN4[K comparable, V any] struct {
	data *crn4.Map[K, V]
}

func NewCRN4Map[K comparable, V any]() *CRN4[K, V] {
	return &CRN4[K, V]{data: crn4.New[K, V](0)}
}

func (m *CRN4[K, V]) Get(key K) (V, bool) {
	value, ok := m.data.Get(key)
	return value, ok
}

func (m *CRN4[K, V]) Set(key K, value V) {
	m.data.Put(key, value)
}

func (m *CRN4[K, V]) Delete(key K) {
	m.data.Delete(key)
}

type Dolthub[K comparable, V any] struct {
	data *dolthub.Map[K, V]
}

func NewDolthubMap[K comparable, V any]() *Dolthub[K, V] {
	return &Dolthub[K, V]{data: dolthub.NewMap[K, V](0)}
}

func (m *Dolthub[K, V]) Get(key K) (V, bool) {
	value, ok := m.data.Get(key)
	return value, ok
}

func (m *Dolthub[K, V]) Set(key K, value V) {
	m.data.Put(key, value)
}

func (m *Dolthub[K, V]) Delete(key K) {
	m.data.Delete(key)
}
