package store

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	ParkingRowsCount    = 10
	ParkingColumnsCount = 10
	MotorCycleRowsCount = 2
	CarRowsCount        = 6
	TruckRowsCount      = 2
)

type Storage struct {
	ParkingStore interface {
		AssignSpot(*Vehicle, context.Context) (int, error)
		UpdateParkingSpotAsEmpty(int, context.Context) error
	}
	TicketsStore interface {
		InsertTicket(*Ticket, context.Context) error
		GetByPlateNumber(parkingBill *ParkingBill, ctx context.Context) Ticket
		DeleteTicketById(int, context.Context) error
		GetTticketsCountData(context.Context) map[string]int
	}
	VehicleStore interface {
		InsertVehicle(*Vehicle, context.Context) error
		DeleteVehicleById(int, context.Context) error
	}
}

func NewStorage(db *sql.DB) Storage {
	initializeDB(db)

	return Storage{
		ParkingStore: &ParkingStore{db},
		TicketsStore: &TicketsStore{db},
		VehicleStore: &VehicleStore{db},
	}
}

func initializeDB(db *sql.DB) {
	var parking_spot_populating string
	for row := range ParkingRowsCount {
		var vehicleType string
		if row < MotorCycleRowsCount {
			vehicleType = "motorcycle"
		} else if row < MotorCycleRowsCount+CarRowsCount {
			vehicleType = "car"
		} else {
			vehicleType = "truck"
		}

		for column := range ParkingColumnsCount {
			parking_spot_populating += fmt.Sprintf("(%d, %d, '%s', TRUE),", row, column, vehicleType)
		}
	}
	parking_spot_populating = parking_spot_populating[:len(parking_spot_populating)-1]+";"
	db.QueryRowContext(
		context.Background(),
		fmt.Sprintf("INSERT INTO parking_spot (row_number, column_number, vehicle_type, empty) VALUES %s", parking_spot_populating),
	)
}
