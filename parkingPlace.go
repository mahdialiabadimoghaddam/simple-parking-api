package main

type parkingPlace struct {
	vehicleType string
	empty       bool
	row, column int
}

func (parkingPlace *parkingPlace) toggle_empty() {
	parkingPlace.empty = !parkingPlace.empty
}
