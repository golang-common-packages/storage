package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// Client manage all SQL-Like actions
type Client struct {
	db     *sql.DB
	Config *LIKE
}

// NewSQLLike init new instance
// The sql package must be used in conjunction with a database driver. See https://golang.org/s/sqldrivers for a list of driverNames.
func NewSQLLike(config *LIKE) *Client {
	currentClient := &Client{nil, nil}

	db, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		log.Println("Error when try to init SQL server: ", err)
		panic(err)
	}

	db.SetConnMaxLifetime(config.MaxConnectionLifetime)
	db.SetMaxIdleConns(config.MaxConnectionIdle)
	db.SetMaxOpenConns(config.MaxConnectionOpen)

	if err := db.PingContext(context.TODO()); err != nil {
		log.Println("Error when try to connect to SQL server: ", err)
		panic(err)
	}

	currentClient.db = db
	log.Println("Connected to SQL-Like Server")

	return currentClient
}

// Execute return results based on 'query' and 'dataModel'
func (c *Client) Execute(
	query string,
	dataModel interface{}) (interface{}, error) {

	var results []interface{}
	_, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	rows, err := c.db.QueryContext(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	// Go through each row to get the result
	for rows.Next() {
		err = rows.Scan(&dataModel)
		if err != nil {
			return nil, err
		}
		results = append(results, dataModel)
	}

	// Check for errors during rows "Close"
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
