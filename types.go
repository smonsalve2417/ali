package main

import "time"

type Alert struct {
	DeviceID string    `json:"device_id" bson:"device_id"`
	Estado   string    `json:"estado"`
	PercLOS  float64   `json:"perclos"`
	Blinks   int       `json:"blinks"`
	Yawns    int       `json:"yawns"`
	Ts       time.Time `json:"time"`
	New      bool      `json:"new"`
}

type HistoryPayload struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type AlertStats struct {
	AvgPerclos    float64        `json:"avg_perclos"`
	TotalBlinks   int            `json:"total_blinks"`
	TotalYawns    int            `json:"total_yawns"`
	EstadoCount   map[string]int `json:"estado_count"`
	EstadoPercent map[string]int `json:"estado_percent"`
	TotalRecords  int            `json:"total_records"`
}

type StatsPayload struct {
	Alerts     []Alert    `json:"alerts"`
	AlertStats AlertStats `json:"alert_stats"`
}
