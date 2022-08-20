package main

import (
	"fmt"
	"runtime"
)

/*
slice 的扩容测试： slice 类型为 int，float64 也是如下结果
go version: go1.18
slice size: 1 ,capacity: 1
slice size: 2 ,capacity: 2
slice size: 3 ,capacity: 4
slice size: 5 ,capacity: 8
slice size: 9 ,capacity: 16
slice size: 17 ,capacity: 32
slice size: 33 ,capacity: 64
slice size: 65 ,capacity: 128
slice size: 129 ,capacity: 256
slice size: 257 ,capacity: 512
slice size: 513 ,capacity: 848
slice size: 849 ,capacity: 1280
slice size: 1281 ,capacity: 1792
slice size: 1793 ,capacity: 2560

可以看到，在插入的数无法放入已有的空间，空间的增长是翻倍增长的
从 512 的空间大小小之后，扩容就在不再是翻倍增加了。
848 - 512 = 336 （(newcap + 3*threshold)/4=320 最后还会做下内存对齐）
1280 - 848 = 432
1792 - 1280 = 512

1.14 的版本是增加 25%
newcap += newcap / 4
1.18 的版本是这样的：
threshold = 256
newcap += (newcap + 3*threshold) / 4
*/

func main() {
	fmt.Printf("go version: %s\n", runtime.Version())
	var sli []float64
	var oldCap int = cap(sli)
	for i := 0; i < 2000; i++ {
		if oldCap != cap(sli) {
			fmt.Printf("slice size: %d ,capacity: %d\n", len(sli), cap(sli))
			oldCap = cap(sli)
		}
		sli = append(sli, float64(i))
	}
}
