package main

//Son todas las funciones que interactúan con la base de datos y Docker

import (
	"database/sql"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db: db}
}

func (s *store) GetAlerts() ([]Alert, error) {
	query := `SELECT TOP 10 device_id, estado, perclos, blinks, yawns, ts FROM alerts`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		err := rows.Scan(&a.DeviceID, &a.Estado, &a.PercLOS, &a.Blinks, &a.Yawns, &a.Ts)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}
