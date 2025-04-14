package main

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"pgregory.net/rand"
)

func randT[T any](r1 *rand.Rand) T {
	t := reflect.TypeOf((*T)(nil)).Elem()

	switch t.Kind() {
	case reflect.Int:
		v := rand.Int()
		return any(v).(T)
	case reflect.Float64:
		v := rand.Float64()
		return any(v).(T)
	case reflect.String:
		v := randString(r1, 7)
		return any(v).(T)
	case reflect.Struct:
		if t.NumField() == 0 {
			var v T
			return v
		}
		panic("only empty structs are supported")
	default:
		panic("unsupported type")
	}
}

func randString(r *rand.Rand, length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	r.Read(b)
	for i := range length {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

type Map[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Delete(K)
}

type Bench[K comparable, V any] struct {
	m      func() Map[K, V]
	keys   []K
	values []V
}

func New[K comparable, V any](size, seed uint64, m func() Map[K, V]) Bench[K, V] {
	b := Bench[K, V]{m: m, keys: make([]K, size), values: make([]V, size)}
	r := rand.New(seed)
	for i := range size {
		b.keys[i] = randT[K](r)
		b.values[i] = randT[V](r)
	}

	return b
}

func (bench *Bench[K, V]) benchmarkInsert(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		m := bench.m()
		for i, key := range bench.keys {
			m.Set(key, bench.values[i])
		}
	}
}

func (bench *Bench[K, V]) benchmarkLookup(b *testing.B) {
	m := bench.m()
	for i, key := range bench.keys {
		m.Set(key, bench.values[i])
	}
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, _ = m.Get(bench.keys[i%len(bench.keys)])
	}
}

func measureMemoryUsage() {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory Usage: Alloc = %v KB, Sys = %v KB, NumGC = %v\n", m.Alloc/1024, m.Sys/1024, m.NumGC)
}

func (bench *Bench[K, V]) Run() {
	t := testing.Benchmark(bench.benchmarkInsert)
	fmt.Printf("Insert: %v\n", t)

	t = testing.Benchmark(bench.benchmarkLookup)
	fmt.Printf("Lookup: %v\n", t)

	measureMemoryUsage()
}
