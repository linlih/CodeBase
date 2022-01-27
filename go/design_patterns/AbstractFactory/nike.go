package main

type nike struct {
}

func (n *nike) makeShoe() iShoe {
	return &nikeShoe{
		shoe: shoe{logo: "nike", size: 15},
	}
}

func (n *nike) makeShirt() iShirt {
	return &nikeShirt{
		shirt: shirt{logo: "nike", size: 15},
	}
}
