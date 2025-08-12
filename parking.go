package main

import (
	"fmt"
	"time"
)

const (
	ParkingRowsCount    = 10
	ParkingColumnsCount = 10
	MotorCycleRowsCount = 2
	CarRowsCount        = 6
	TruckRowsCount      = 2
)

type parking struct {
	places  [ParkingRowsCount][ParkingColumnsCount]*parkingPlace
	tickets map[string]*ticket
}

type parkingStat struct {
	ParkingSpacesCount               int `json:"Parking Spaces Count"`
	FreeParkingSpacesCount           int `json:"Free Parking Spaces Count"`
	MotorCycleParkingSpacesCount     int `json:"MotorCycle Parking Spaces Count"`
	FreeMotorCycleParkingSpacesCount int `json:"Free MotorCycle Parking Spaces Count"`
	CarParkingSpacesCount            int `json:"Car Parking Spaces Count"`
	FreeCarParkingSpacesCount        int `json:"Free Car Parking Spaces Count"`
	TruckParkingSpacesCount          int `json:"Truck Parking Spaces Count"`
	FreeTruckParkingSpacesCount      int `json:"Free Truck Parking Spaces Count"`
}

type ticket struct {
	enteryTime   int
	parkingPlace *parkingPlace
}

type parkingBill struct {
	VehiclePlateNumber string `json:"VehiclePlateNumber"`
	VehicleType        string `json:"VehicleType"`
	Duration           int    `json:"Duration"`
	ParkingHourlyFee   int    `json:"ParkingHourlyFee"`
	ParkingTotalCost   int    `json:"parkingTotalCost"`
}

func populate(parking *parking) {
	parking.submitCustomer("1938251", 1754068008213, parking.places[1][1])
	parking.submitCustomer("4284967", 1754068061928, parking.places[0][5])
	parking.submitCustomer("5230194", 1754068051734, parking.places[4][3])
	parking.submitCustomer("8093116", 1754068000007, parking.places[3][2])
	parking.submitCustomer("7392840", 1754068033801, parking.places[3][6])
	parking.submitCustomer("1175362", 1754068076302, parking.places[5][9])
	parking.submitCustomer("6548937", 1754068043519, parking.places[6][5])
	parking.submitCustomer("3849021", 1754068066452, parking.places[9][1])
}

func (parking *parking) submitCustomer(vehiclePlateNumber string, enteryTime int, parkingPlace *parkingPlace) {
	newTicket := ticket{
		enteryTime:   enteryTime,
		parkingPlace: parkingPlace,
	}
	parking.tickets[vehiclePlateNumber] = &newTicket
	parkingPlace.toggle_empty()
}

func (parking *parking) processCustomer(car_type, vehiclePlateNumber string) {
	var row_min, row_max int
	switch car_type {
	case "motorcycle":
		row_min, row_max = 0, MotorCycleRowsCount
	case "car":
		row_min, row_max = MotorCycleRowsCount, MotorCycleRowsCount + CarRowsCount
	case "truck":
		row_min, row_max = ParkingRowsCount - TruckRowsCount, ParkingRowsCount
	}

	for row := row_min; row < row_max; row++ {
		for column := range ParkingColumnsCount {
			if parking.places[row][column].empty {
				parking.submitCustomer(
					vehiclePlateNumber,
					int(time.Now().UnixMilli()),
					parking.places[row][column],
				)
				return
			}
		}
	}
}

func (parking *parking) exitParking(parkingBill *parkingBill) {
	parkingFee := map[string]int{
		"motorcycle": 1,
		"car":        5,
		"truck":      10,
	}

	ticket := parking.tickets[parkingBill.VehiclePlateNumber]
	elapsedTime := time.Now().UnixMilli() - int64(ticket.enteryTime)

	parkingBill.VehicleType = ticket.parkingPlace.vehicleType
	parkingBill.Duration = int(elapsedTime / 3_600_000) //milliseconds -> hours
	parkingBill.ParkingHourlyFee = parkingFee[parkingBill.VehicleType]
	parkingBill.ParkingTotalCost = parkingBill.Duration * parkingBill.ParkingHourlyFee

	defer ticket.parkingPlace.toggle_empty()

}

func (parking *parking) getStat() parkingStat {
	parking_stat := parkingStat{
		ParkingRowsCount * ParkingColumnsCount,
		ParkingRowsCount*ParkingColumnsCount - len(parking.tickets),
		MotorCycleRowsCount * ParkingColumnsCount,
		MotorCycleRowsCount * ParkingColumnsCount,
		CarRowsCount * ParkingColumnsCount,
		CarRowsCount * ParkingColumnsCount,
		TruckRowsCount * ParkingColumnsCount,
		TruckRowsCount * ParkingColumnsCount,
	}

	for _, ticket := range parking.tickets {
		switch ticket.parkingPlace.vehicleType {
		case "motorcycle":
			parking_stat.FreeMotorCycleParkingSpacesCount--
		case "car":
			parking_stat.FreeCarParkingSpacesCount--
		case "truck":
			parking_stat.FreeTruckParkingSpacesCount--
		}
	}

	return parking_stat
}

func (parking *parking) deleteTicket(vehiclePlateNumber string) {
	parking.tickets[vehiclePlateNumber].parkingPlace.toggle_empty()
	delete(parking.tickets, vehiclePlateNumber)
}

func (parking *parking) updateTicket(vehiclePlateNumber string, enteryTime int, parking_row, parking_column int) {
	parking.deleteTicket(vehiclePlateNumber)
	parking.submitCustomer(vehiclePlateNumber, enteryTime, parking.places[parking_row][parking_column])
}

func (parking *parking) String() string {
	var parkingString [ParkingRowsCount][ParkingColumnsCount]string
	for vehiclePlateNumber, ticket := range parking.tickets {
		parkingString[ticket.parkingPlace.row][ticket.parkingPlace.column] = vehiclePlateNumber
	}

	finalString := ""
	for row := range ParkingRowsCount {
		var carType string
		if row < MotorCycleRowsCount {
			carType = "m:"
		} else if row < MotorCycleRowsCount + CarRowsCount {
			carType = "c:"
		} else {
			carType = "t:"
		}
		finalString += fmt.Sprintf("%s\t| ", carType)

		for column := range ParkingColumnsCount {
			if parking.places[row][column].empty {
				finalString += "------- | "
			} else {
				finalString += fmt.Sprintf("%s | ", parkingString[row][column])
			}

		}
		finalString += "\n"
	}
	return finalString + "\n"
}

func newParking() *parking {
	parking := &parking{
		tickets: make(map[string]*ticket),
		places:  [ParkingRowsCount][ParkingColumnsCount]*parkingPlace{},
	}

	for row := range 10 {
		for column := range 10 {
			var carType string
			if row < 2 {
				carType = "motorcycle"
			} else if row < 8 {
				carType = "car"
			} else {
				carType = "truck"
			}

			parking.places[row][column] = &parkingPlace{carType, true, row, column}
		}
	}

	populate(parking)
	
	return parking
}
