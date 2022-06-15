package main

import "fmt"

func main() {
	s := NewSkipList()
	s.Insert(1, "1")
	s.Insert(2, "2")
	s.Insert(3, "3")
	s.Insert(4, "4")
	s.Insert(5, "5")
	s.Insert(6, "6")
	s.Insert(7, "7")
	s.Insert(8, "8")
	s.Insert(9, "9")
	s.Insert(10, "10")

	s.DisplayAll()

	ret, err := s.Search(8)
	if err == nil {
		fmt.Println("key 6: val ->", ret)
	} else {
		fmt.Println("notfound", err)
	}
	s.Delete(6)
	s.DisplayAll()
}
