package sql

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/golang-common-packages/database/model"
)

// Client manage all SQL-Like API
type Client struct {
	db     *sql.DB
	Config *model.LIKE
}

var (
	ctx    context.Context
	cancel context.CancelFunc
)

// NewSQLLike function return a new SQL-Like client
// The sql package must be used in conjunction with a database driver. See https://golang.org/s/sqldrivers for a list of driverNames.
func NewSQLLike(config *model.LIKE) *Client {

	ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
	currentClient := &Client{nil, nil}

	db, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		cancel()
		log.Println("Error when try to init SQL server: ", err)
		panic(err)
	}

	db.SetConnMaxLifetime(config.MaxConnectionLifetime)
	db.SetMaxIdleConns(config.MaxConnectionIdle)
	db.SetMaxOpenConns(config.MaxConnectionOpen)

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
