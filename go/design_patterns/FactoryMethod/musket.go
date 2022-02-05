package main

type musket struct {
	gun
}

func newMustket() iGun {
	return &musket{
		gun: gun{name: "Musket gun",
			power: 1},
	}
}
