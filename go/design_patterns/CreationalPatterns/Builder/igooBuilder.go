package main

type iglooBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func (b *iglooBuilder) setWindowType() {
	b.windowType = "Snow Window"
}

func (b *iglooBuilder) setDoorType() {
	b.doorType = "Snow door"
}

func (b *iglooBuilder) setNumFloor() {
	b.floor = 1
}

func (b *iglooBuilder) getHouse() house {
	return house{
		windowType: b.windowType,
		doorType:   b.doorType,
		floor:      b.floor,
	}
}

func newIglooBuilder() *iglooBuilder {
	return &iglooBuilder{}
}
