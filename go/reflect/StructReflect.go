package main

import (
	"fmt"
	"reflect"
)

// 使用反射输出结构体的名字以及各个字段的名字类型以及Tag信息

type ReflectStruct struct {
	Name   string
	Age    int
	saving float64
	Cars   []string `flagTest:"bmw"`
}

func main() {
	structType := reflect.ValueOf(ReflectStruct{}).Type()

	// 获取结构体的名字
	fmt.Println(structType.Name())

	// 获取结构体的每个字段
	for i := 0; i < structType.NumField(); i++ {
		subField := structType.Field(i)
		fmt.Println(subField.Type, " ", subField.Name, " ", subField.IsExported())
		if subField.Tag != "" {
			fmt.Println(subField.Tag.Get("flagTest"))
		}
	}

	testSturct := ReflectStruct{
		Name:   "test",
		Age:    12,
		saving: 1000.0,
		Cars:   []string{"bmw", "other"},
	}
	// 将反射转变成原有的结构体，反射空间到原有的变量空间的转换
	structValue := reflect.ValueOf(testSturct)
	test := structValue.Interface().(ReflectStruct)
	fmt.Println(test.Age)
}
