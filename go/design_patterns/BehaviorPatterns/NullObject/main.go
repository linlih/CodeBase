package main

import "fmt"

type IEmployee interface {
	GetName() string
	IsNULL() bool
}

type Employee struct {
	Name string
	Null bool
}

func NewEmployee(name string) IEmployee {
	return &Employee{Name: name, Null: true}
}

func (e *Employee) GetName() string {
	return e.Name
}

func (e *Employee) IsNULL() bool {
	return e.Null
}

type NullObject struct {
	Name string
	Null bool
}

func NewNullObject() IEmployee {
	return &Employee{Name: "null object", Null: false}
}

func (n *NullObject) GetName() string {
	return n.Name
}

func (n *NullObject) IsNULL() bool {
	return n.Null
}

func GetObject(name string) IEmployee {
	namelist := []string{"jay", "john", "ian"}
	for i := 0; i < 3; i++ {
		if name == namelist[i] {
			return NewEmployee(name)
		}
	}
	return NewNullObject()
}

func main() {
	emp1 := GetObject("john")
	fmt.Println("Name:", emp1.GetName(), ", is exist:", emp1.IsNULL())
	emp2 := GetObject("josh")
	fmt.Println("Name:", emp2.GetName(), ", is exist:", emp2.IsNULL())
}
