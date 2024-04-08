package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func createNonRecurrenceEvent(
	conn *pgx.Conn,
	startDate time.Time,
	endDate time.Time,
	startTime time.Time,
	endTime time.Time,
) (*int, error) {
	sqlQ := "insert into events (title, start_date, end_date, start_time, end_time, is_full_day_event, is_recurring) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	var id *int
	err := conn.QueryRow(
		context.Background(),
		sqlQ,
		"random-title",
		startDate,
		endDate,
		startTime,
		endTime,
		false,
		false,
	).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		panic(err)
	}

	sqlQ = "select id, start_date from events"
	rows, err := conn.Query(
		context.Background(),
		sqlQ,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id *int
		var start_date *time.Time
		err = rows.Scan(&id, &start_date)
		if err != nil {
			return nil, err
		}

		fmt.Println(*id, *start_date, (*start_date).UTC())
	}

	return id, nil
}
