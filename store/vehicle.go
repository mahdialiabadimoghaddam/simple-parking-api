package store

import (
	"context"
	"database/sql"
	"fmt"
)

type VehicleStore struct {
	db *sql.DB
}

type Vehicle struct {
	Id                 int    `json:"id,omitempty"`
	VehicleType        string `json:"vehicleType"`
	VehiclePlateNumber string `json:"vehiclePlateNumber"`
}

func (VehicleStore *VehicleStore) InsertVehicle(vehicle *Vehicle, ctx context.Context) error {
	query := `INSERT INTO vehicle (type, plate_number) VALUES ($1, $2) RETURNING id`
	err := VehicleStore.db.QueryRowContext(
		ctx,
		query,
		vehicle.VehicleType,
		vehicle.VehiclePlateNumber,
	).Scan(
		&vehicle.Id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (VehicleStore *VehicleStore) DeleteVehicleById(vehicleId int, ctx context.Context) error {
	query := `DELETE FROM vehicle WHERE id = $1`
	_, err := VehicleStore.db.ExecContext(
		ctx,
		query,
		vehicleId,
	)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}
