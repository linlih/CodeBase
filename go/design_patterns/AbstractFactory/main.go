package main

import "fmt"

func main() {
	adidasFactory, _ := getSportsFactory("adidas")
	nikeFactory, _ := getSportsFactory("nike")

	nikeShoe := nikeFactory.makeShoe()
	nikeShirt := nikeFactory.makeShirt()

	adidasShoe := adidasFactory.makeShoe()
	adidasShirt := adidasFactory.makeShirt()

	printShoeDetails(nikeShoe)
	printShirtDetails(nikeShirt)

	printShoeDetails(adidasShoe)
	printShirtDetails(adidasShirt)
}

func printShoeDetails(s iShoe) {
	fmt.Printf("shoe's logo: %s, size: %d\n", s.getLogo(), s.getSize())
}

func printShirtDetails(s iShirt) {
	fmt.Printf("shirt's Logo: %s, size: %d\n", s.getLogo(), s.getSize())
}
