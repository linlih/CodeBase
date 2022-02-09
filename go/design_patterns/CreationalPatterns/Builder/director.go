package main

type director struct {
	builder iBuilder
}

func newDirector(b iBuilder) *director {
	return &director{
		builder: b,
	}
}

func (d *director) setBuilder(b iBuilder) {
	d.builder = b
}

func (d *director) buildHouse() house {
	d.builder.setWindowType()
	d.builder.setDoorType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}
