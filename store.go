package main

//Son todas las funciones que interactúan con la base de datos y Docker

import (
	"database/sql"
	"strings"
	"time"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db: db}
}

func (s *store) GetAlerts() ([]Alert, error) {
	query := `SELECT device_id, estado, perclos, blinks, yawns, ts FROM alerts ORDER BY ts DESC LIMIT 10`

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

		a.Ts = AjustaHora(a.Ts)

		alerts = append(alerts, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (s *store) GetLastAlert() ([]Alert, error) {
	query := `SELECT device_id, estado, perclos, blinks, yawns, ts FROM alerts ORDER BY ts DESC LIMIT 1`

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

		a.Ts = AjustaHora(a.Ts)

		if a.Estado == "" {
			a.Estado = "NORMAL"
		}

		alerts = append(alerts, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (s *store) GetAlertsByRange(startDate, endDate time.Time) ([]Alert, error) {
	query := `
		SELECT device_id, estado, perclos, blinks, yawns, ts 
		FROM alerts 
		WHERE ts BETWEEN ? AND ?
		ORDER BY ts ASC
	`

	startDate = AjustaHoraAdd5(startDate)
	endDate = AjustaHoraAdd5(endDate)

	rows, err := s.db.Query(query, startDate, endDate)
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

		// Ajustar zona horaria (-5h)
		a.Ts = AjustaHora(a.Ts)

		if a.Estado == "" {
			a.Estado = "NORMAL"
		}

		alerts = append(alerts, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

func AjustaHora(t time.Time) time.Time {
	return t.Add(-5 * time.Hour)
}
func AjustaHoraAdd5(t time.Time) time.Time {
	return t.Add(5 * time.Hour)
}

func (s *store) CalcularEstadisticas(alerts []Alert) AlertStats {
	stats := AlertStats{
		EstadoCount: map[string]int{
			"NORMAL":       0,
			"FATIGA":       0,
			"SOMNOLIENCIA": 0,
			"MICROSUEÑO":   0,
		},
	}

	if len(alerts) == 0 {
		return stats
	}

	var sumPerclos float64
	for _, a := range alerts {
		sumPerclos += a.PercLOS
		stats.TotalBlinks += a.Blinks
		stats.TotalYawns += a.Yawns

		// contar estados (ignora mayúsculas/minúsculas)
		switch strings.ToUpper(a.Estado) {
		case "NORMAL":
			stats.EstadoCount["NORMAL"]++
		case "FATIGA":
			stats.EstadoCount["FATIGA"]++
		case "SOMNOLENCIA":
			stats.EstadoCount["SOMNOLENCIA"]++
		case "MICROSUEÑO":
			stats.EstadoCount["MICROSUEÑO"]++
		}
	}

	stats.AvgPerclos = sumPerclos / float64(len(alerts))
	stats.TotalRecords = len(alerts)
	stats.EstadoPercent = make(map[string]int)
	for estado, count := range stats.EstadoCount {
		percent := (float64(count) / float64(len(alerts))) * 100
		stats.EstadoPercent[estado] = int(percent)
	}

	return stats
}
