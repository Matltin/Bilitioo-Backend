// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: city.sql

package db

import (
	"context"
	"time"
)

const getCities = `-- name: GetCities :many
SELECT 
    "province", 
    "county"
FROM "city"
`

type GetCitiesRow struct {
	Province string `json:"province"`
	County   string `json:"county"`
}

func (q *Queries) GetCities(ctx context.Context) ([]GetCitiesRow, error) {
	rows, err := q.db.QueryContext(ctx, getCities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCitiesRow{}
	for rows.Next() {
		var i GetCitiesRow
		if err := rows.Scan(&i.Province, &i.County); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchTicketsByCities = `-- name: SearchTicketsByCities :many
SELECT 
    t.id AS ticket_id,
    t.departure_time,
    t.arrival_time,
    t.amount,
    t.status,
    r.id AS route_id,
    r.origin_city_id,
    r.destination_city_id,
    c1.province AS origin_province,
    c1.county AS origin_county,
    c2.province AS destination_province,
    c2.county AS destination_county,
    v.vehicle_type,
    v.capacity,
    comp.name AS company_name
FROM ticket t
JOIN route r ON t.route_id = r.id
JOIN city c1 ON r.origin_city_id = c1.id
JOIN city c2 ON r.destination_city_id = c2.id
JOIN vehicle v ON t.vehicle_id = v.id
JOIN company comp ON v.company_id = comp.id
WHERE r.origin_city_id = $1 AND r.destination_city_id = $2 AND v.vehicle_type = $3
ORDER BY t.departure_time ASC
`

type SearchTicketsByCitiesParams struct {
	OriginCityID      int64       `json:"origin_city_id"`
	DestinationCityID int64       `json:"destination_city_id"`
	VehicleType       VehicleType `json:"vehicle_type"`
}

type SearchTicketsByCitiesRow struct {
	TicketID            int64                        `json:"ticket_id"`
	DepartureTime       time.Time                    `json:"departure_time"`
	ArrivalTime         time.Time                    `json:"arrival_time"`
	Amount              int64                        `json:"amount"`
	Status              CheckReservationTicketStatus `json:"status"`
	RouteID             int64                        `json:"route_id"`
	OriginCityID        int64                        `json:"origin_city_id"`
	DestinationCityID   int64                        `json:"destination_city_id"`
	OriginProvince      string                       `json:"origin_province"`
	OriginCounty        string                       `json:"origin_county"`
	DestinationProvince string                       `json:"destination_province"`
	DestinationCounty   string                       `json:"destination_county"`
	VehicleType         VehicleType                  `json:"vehicle_type"`
	Capacity            int32                        `json:"capacity"`
	CompanyName         string                       `json:"company_name"`
}

func (q *Queries) SearchTicketsByCities(ctx context.Context, arg SearchTicketsByCitiesParams) ([]SearchTicketsByCitiesRow, error) {
	rows, err := q.db.QueryContext(ctx, searchTicketsByCities, arg.OriginCityID, arg.DestinationCityID, arg.VehicleType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchTicketsByCitiesRow{}
	for rows.Next() {
		var i SearchTicketsByCitiesRow
		if err := rows.Scan(
			&i.TicketID,
			&i.DepartureTime,
			&i.ArrivalTime,
			&i.Amount,
			&i.Status,
			&i.RouteID,
			&i.OriginCityID,
			&i.DestinationCityID,
			&i.OriginProvince,
			&i.OriginCounty,
			&i.DestinationProvince,
			&i.DestinationCounty,
			&i.VehicleType,
			&i.Capacity,
			&i.CompanyName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
