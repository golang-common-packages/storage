package sql

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// Client manage all SQL API
type Client struct {
	db *sql.DB
}

var ctx = context.Background()

// NewSQL function return a new SQL client
func NewSQL(
	driverName,
	dataSourceName string,
	maxConnectionLifetime time.Duration,
	maxConnectionIdle,
	maxConnectionOpen int) *Client {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	currentClient := &Client{nil}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		cancel()
		log.Println("Error when try to init SQL server: ", err)
		panic(err)
	}

	db.SetConnMaxLifetime(maxConnectionLifetime)
	db.SetMaxIdleConns(maxConnectionIdle)
	db.SetMaxOpenConns(maxConnectionOpen)

	if err := db.PingContext(ctx); err != nil {
		cancel()
		log.Println("Error when try to connect to SQL server: ", err)
		panic(err)
	}

	currentClient.db = db

	return currentClient
}

// Execute return results based on 'query' and 'dataModel'
func (c *Client) Execute(
	query string,
	dataModel interface{}) (interface{}, error) {

	var results []interface{}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Go through each row to get the result
	for rows.Next() {
		err = rows.Scan(&dataModel)
		if err != nil {
			return nil, err
			break
		}
		results = append(results, dataModel)
	}

	// Check for errors during rows "Close".
	// This may be more important if multiple statements are executed
	// in a single batch and rows were written as well as read.
	if err := rows.Close(); err != nil {
		return nil, err
	}

	// Check for errors during row iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
