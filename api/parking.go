package main

import (
	"math"
	"context"
	"fmt"
	"parking_app/store"
	"time"
)

type ParkingStat struct {
	ParkingSpacesCount               int `json:"Parking Spaces Count"`
	FreeParkingSpacesCount           int `json:"Free Parking Spaces Count"`
	MotorCycleParkingSpacesCount     int `json:"MotorCycle Parking Spaces Count"`
	FreeMotorCycleParkingSpacesCount int `json:"Free MotorCycle Parking Spaces Count"`
	CarParkingSpacesCount            int `json:"Car Parking Spaces Count"`
	FreeCarParkingSpacesCount        int `json:"Free Car Parking Spaces Count"`
	TruckParkingSpacesCount          int `json:"Truck Parking Spaces Count"`
	FreeTruckParkingSpacesCount      int `json:"Free Truck Parking Spaces Count"`
}

func (app *application) processCustomer(vehicle *store.Vehicle, ctx context.Context) store.Ticket{
	app.store.VehicleStore.InsertVehicle(vehicle, ctx)
	parkingSpotId, _ := app.store.ParkingStore.AssignSpot(vehicle, ctx)
	ticket := store.Ticket{
		VehicleId:     vehicle.Id,
		ParkingSpotId: parkingSpotId,
		Content:       "",
		EnteryTime:    int(time.Now().UnixMilli()),
	}

	err := app.store.TicketsStore.InsertTicket(&ticket, ctx)
	if err != nil {
		fmt.Println(err)
	}

	return ticket
}

func (app *application) exitParking(parkingBill *store.ParkingBill, ctx context.Context) string {
	parkingFee := map[string]int{
		"motorcycle": 1,
		"car":        5,
		"truck":      10,
	}
	
	ticket := app.store.TicketsStore.GetByPlateNumber(parkingBill, ctx)
	app.store.TicketsStore.DeleteTicketById(ticket.Id, ctx)
	app.store.VehicleStore.DeleteVehicleById(ticket.VehicleId, ctx)
	app.store.ParkingStore.UpdateParkingSpotAsEmpty(ticket.ParkingSpotId, ctx)

	elapsedTime := time.Now().UnixMilli() - int64(ticket.EnteryTime)
	parkingBill.Duration = int(math.Max(1, float64(int(elapsedTime / 3_600_000))))//milliseconds -> hours -- minimum is 1 hour
	parkingBill.ParkingHourlyFee = parkingFee[parkingBill.VehicleType]
	parkingBill.ParkingTotalCost = parkingBill.Duration * parkingBill.ParkingHourlyFee

	return "success"
}

func (app *application) getParkingStat(ctx context.Context) ParkingStat {
	ticketsCountData := app.store.TicketsStore.GetTticketsCountData(ctx)

	parking_stat := ParkingStat{
		store.ParkingRowsCount * store.ParkingColumnsCount,
		store.ParkingRowsCount*store.ParkingColumnsCount - ticketsCountData["sum"],
		store.MotorCycleRowsCount * store.ParkingColumnsCount,
		store.MotorCycleRowsCount * store.ParkingColumnsCount,
		store.CarRowsCount * store.ParkingColumnsCount,
		store.CarRowsCount * store.ParkingColumnsCount,
		store.TruckRowsCount * store.ParkingColumnsCount,
		store.TruckRowsCount * store.ParkingColumnsCount,
	}

	parking_stat.FreeMotorCycleParkingSpacesCount -= ticketsCountData["motorcyclesCount"]
	parking_stat.FreeCarParkingSpacesCount -= ticketsCountData["carsCount"]
	parking_stat.FreeTruckParkingSpacesCount -= ticketsCountData["trucksCount"]

	return parking_stat
}

func (app *application) deleteTicket(vehiclePlateNumber string, ctx context.Context) string {
	return app.exitParking(&store.ParkingBill{VehiclePlateNumber: vehiclePlateNumber}, ctx)
}
