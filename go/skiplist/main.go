package main

import "fmt"

func main() {
	s := NewSkipList()
	s.Insert(50, "5")
	s.Insert(40, "4")
	s.Insert(70, "7")
	s.Insert(100, "10")
	s.Insert(12, "12")
	s.Insert(22, "19")
	s.Insert(30, "8")

	s.DisplayAll()

	ret, err := s.Search(70)
	if err == nil {
		fmt.Println("key 50: val ->", ret)
	} else {
		fmt.Println("notfound", err)
	}
}
