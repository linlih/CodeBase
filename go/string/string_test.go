package string

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// 关于字符串的性能还是要根据实际情况进行选择，一般是不能使用Sprintf来进行拼接，这个方法使用到了反射性能较差
//

/*
goos: linux
goarch: amd64
cpu: AMD Ryzen 7 4800H with Radeon Graphics
BenchmarkStringConcat-8           144367            137787 ns/op          364876 B/op          1 allocs/op
BenchmarkStringSprintf-8          124742            229409 ns/op          628520 B/op          5 allocs/op
BenchmarkStringBytesBuffer-8    93628195             50.06 ns/op              12 B/op          0 allocs/op
BenchmarkStringBuilder-8        100000000            387.4 ns/op              29 B/op          0 allocs/op
BenchmarkStringJoin-8              78795             67996 ns/op          200825 B/op          1 allocs/op
*/

var chartoappend = "hello"

func BenchmarkStringConcat(b *testing.B) {
	var str string

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str += chartoappend
	}
}

func BenchmarkStringSprintf(b *testing.B) {
	var str string

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		str = fmt.Sprintf("%v%v", str, chartoappend)
	}
}

func BenchmarkStringBytesBuffer(b *testing.B) {
	var buffer bytes.Buffer

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buffer.WriteString(chartoappend)
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	var strBuilder strings.Builder

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		strBuilder.WriteString(chartoappend)
	}
}

func BenchmarkStringJoin(b *testing.B) {
	var str string

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		str = strings.Join([]string{str, chartoappend}, "")
	}
}
