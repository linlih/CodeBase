package main

import (
	"sync"
	"testing"
)

// 参考自：https://medium.com/swlh/go-the-idea-behind-sync-pool-32da5089df72

type Person struct {
	Age int
}

var personPool = sync.Pool{New: func() interface{} { return new(Person) }}

// 测试结果样例：
// BenchmarkWithoutPool-8   	    7575	    133696 ns/op	   80000 B/op	   10000 allocs/op
func BenchmarkWithoutPool(b *testing.B) {
	var p *Person
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			p = new(Person)
			p.Age = 23
		}
	}
}

// 测试结果样例：
// BenchmarkWithPool-8   	   10000	    112084 ns/op	       0 B/op	       0 allocs/op
func BenchmarkWithPool(b *testing.B) {
	var p *Person
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			p = personPool.Get().(*Person)
			p.Age = 23
			personPool.Put(p)
		}
	}
}

// 可以看到在一定程度上，Pool会有一定的性能损失，这个就是 trade-off 了
// BenchmarkPool-8   	642773518	         1.763 ns/op
// BenchmarkAllocation-8   	1000000000	         0.1124 ns/op
func BenchmarkPool(b *testing.B) {
	var p sync.Pool
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Put(1)
			p.Get()
		}
	})
}

func BenchmarkAllocation(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := 0
			i = i
		}
	})
}
