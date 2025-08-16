package store

import (
	"context"
	"database/sql"
	"fmt"
)

type ParkingStore struct {
	db *sql.DB
}

func Create(ctx context.Context) error {
	return nil
}

func (parkingStore *ParkingStore) AssignSpot(vehicle *Vehicle, ctx context.Context) (int, error) {
	var rowID int
	query := `UPDATE parking_spot SET empty = FALSE WHERE id = (
			SELECT id
			FROM parking_spot
			WHERE empty = TRUE AND vehicle_type = $1
			ORDER BY id ASC
			LIMIT 1
		)
		RETURNING id;`

	err := parkingStore.db.QueryRowContext(
		ctx,
		query,
		vehicle.VehicleType,
	).Scan(
		&rowID,
	)

	return rowID, err
}

func (parkingStore *ParkingStore) UpdateParkingSpotAsEmpty(ParkingSpotId int, ctx context.Context) error {
	query := `UPDATE parking_spot SET empty = TRUE WHERE id = $1`

	_, err := parkingStore.db.ExecContext(
		ctx,
		query,
		ParkingSpotId,
	)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}
