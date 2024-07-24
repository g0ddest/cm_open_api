package postgres

import (
	"cm_open_api/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
	"time"
)

func GetOutages(connStr string) ([]models.Outage, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to PostgreSQL: %v", err)
	}
	defer conn.Close(context.Background())

	query := `SELECT
        message_id, incident_id, service, organization, short_description, event, event_start, event_stop,
        NULLIF(region_kladr, ''), NULLIF(region_type, ''), NULLIF(region_name, ''),
        city_kladr, COALESCE(NULLIF(city_name, ''), city) as city_name, city_type,
        NULLIF(street_kladr, ''), COALESCE(NULLIF(street_name, ''), street) as street_name,
        COALESCE(NULLIF(street_type, ''), street_type_raw) as street_type,
        NULLIF(house_numbers, ''), NULLIF(house_ranges, '')
        FROM communal_outages
        WHERE event = 'shutdown' AND event_start >= date_trunc('day', current_date)
        ORDER BY event_start`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var outages []models.Outage
	for rows.Next() {
		var outage models.Outage
		var regionKladr, regionType, regionName, streetKladr, houseNumbers, houseRanges *string
		var eventStart, eventStop *time.Time

		err := rows.Scan(
			&outage.MessageID, &outage.IncidentID, &outage.Service, &outage.Organization,
			&outage.ShortDescription, &outage.Event, &eventStart, &eventStop,
			&regionKladr, &regionType, &regionName,
			&outage.CityKladr, &outage.CityName, &outage.CityType,
			&streetKladr, &outage.StreetName, &outage.StreetType,
			&houseNumbers, &houseRanges,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if eventStart != nil {
			outage.EventStart = *eventStart
		}
		if eventStop != nil {
			outage.EventStop = eventStop
		}
		if regionKladr != nil {
			outage.RegionKladr = regionKladr
		}
		if regionType != nil {
			outage.RegionType = regionType
		}
		if regionName != nil {
			outage.RegionName = regionName
		}
		if streetKladr != nil {
			outage.StreetKladr = streetKladr
		}
		if houseNumbers != nil {
			outage.HouseNumbers = strings.Split(*houseNumbers, ",")
		}
		if houseRanges != nil {
			outage.HouseRanges = strings.Split(*houseRanges, ",")
		}

		outages = append(outages, outage)
	}

	return outages, nil
}
