package store

import (
	"context"
	"database/sql"
	"fmt"
)

type TicketsStore struct {
	db *sql.DB
}

type Ticket struct {
	Id            int
	VehicleId     int
	ParkingSpotId int
	Content       string
	EnteryTime    int
}

type ParkingBill struct {
	VehiclePlateNumber string `json:"VehiclePlateNumber"`
	VehicleType        string `json:"VehicleType"`
	Duration           int    `json:"Duration"`
	ParkingHourlyFee   int    `json:"ParkingHourlyFee"`
	ParkingTotalCost   int    `json:"parkingTotalCost"`
}

func (ticketStore *TicketsStore) InsertTicket(ticket *Ticket, ctx context.Context) error {
	query := `INSERT INTO tickets (vehicle_id, parking_spot_id, content, entery_time) VALUES ($1, $2, $3, $4) RETURNING id`
	err := ticketStore.db.QueryRowContext(
		ctx,
		query,
		ticket.VehicleId,
		ticket.ParkingSpotId,
		ticket.Content,
		ticket.EnteryTime,
	).Scan(
		&ticket.Id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (ticketStore *TicketsStore) GetByPlateNumber(parkingBill *ParkingBill, ctx context.Context) Ticket {
	var ticket Ticket
	query := `SELECT t.id, t.vehicle_id, t.parking_spot_id, t.content, t.entery_time, v.type FROM tickets t JOIN vehicle v ON t.vehicle_id=v.id WHERE v.plate_number = $1`
	err := ticketStore.db.QueryRowContext(
		ctx,
		query,
		parkingBill.VehiclePlateNumber,
	).Scan(
		&ticket.Id,
		&ticket.VehicleId,
		&ticket.ParkingSpotId,
		&ticket.Content,
		&ticket.EnteryTime,
		&parkingBill.VehicleType,
	)

	if err != nil {
		panic(err)
	} else {
		return ticket
	}
}

func (ticketStore *TicketsStore) DeleteTicketById(ticketId int, ctx context.Context) error {
	query := `DELETE FROM tickets WHERE id = $1`
	fmt.Println(ticketId)
	_, err := ticketStore.db.ExecContext(
		ctx,
		query,
		ticketId,
	)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (ticketStore *TicketsStore) GetTticketsCountData(ctx context.Context) map[string]int {
	var motorcyclesCount int
	var carsCount int
	var trucksCount int
	query := `
    SELECT
        COUNT(*) FILTER (WHERE v.type = 'motorcycle') AS motorcycles_count,
        COUNT(*) FILTER (WHERE v.type = 'car') AS cars_count,
        COUNT(*) FILTER (WHERE v.type = 'truck') AS trucks_count
    FROM
        tickets t JOIN vehicle v ON t.vehicle_id = v.id;`
	err := ticketStore.db.QueryRowContext(
		ctx,
		query,
	).Scan(
		&motorcyclesCount,
		&carsCount,
		&trucksCount,
	)

	sum := motorcyclesCount + carsCount + trucksCount

	if err != nil {
		panic(err)
	} else {
		return map[string]int{
			"motorcyclesCount": motorcyclesCount,
			"carsCount":        carsCount,
			"trucksCount":      trucksCount,
			"sum":              sum,
		}
	}
}
