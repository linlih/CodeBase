package main

import "fmt"

type windowsAdapter struct {
	windowMachine *windows
}

func (w *windowsAdapter) insertIntoLightningPort() {
	fmt.Println("Adapter convert Lightning port to USB")
	w.windowMachine.insertIntoUSBPort()
}
