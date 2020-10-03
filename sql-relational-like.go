package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/golang-common-packages/hash"
)

// SQLLikeClient manage all SQL-Like actions
type SQLLikeClient struct {
	Client *sql.DB
	Config *LIKE
}

var (
	// sqlLikeClientSessionMapping singleton pattern
	sqlLikeClientSessionMapping = make(map[string]*SQLLikeClient)
)

// NewSQLLike init new instance
// The sql package must be used in conjunction with a database driver. See https://golang.org/s/sqldrivers for a list of driverNames.
func NewSQLLike(config *LIKE) *SQLLikeClient {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentSQLLikeSession := sqlLikeClientSessionMapping[configAsString]
	if currentSQLLikeSession == nil {
		currentSQLLikeSession = &SQLLikeClient{nil, nil}

		client, err := sql.Open(config.DriverName, config.DataSourceName)
		if err != nil {
			log.Println("Error when try to init SQL server: ", err)
			panic(err)
		}

		client.SetConnMaxLifetime(config.MaxConnectionLifetime)
		client.SetMaxIdleConns(config.MaxConnectionIdle)
		client.SetMaxOpenConns(config.MaxConnectionOpen)

		if err := client.PingContext(context.TODO()); err != nil {
			log.Println("Error when try to connect to SQL server: ", err)
			panic(err)
		}

		currentSQLLikeSession.Client = client
		currentSQLLikeSession.Config = config
		sqlLikeClientSessionMapping[configAsString] = currentSQLLikeSession
		log.Println("Connected to SQL-Like Server")
	}

	return currentSQLLikeSession
}

// Execute return results based on 'query' and 'dataModel'
func (c *SQLLikeClient) Execute(
	query string,
	dataModel interface{}) (interface{}, error) {

	var results []interface{}
	_, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	rows, err := c.Client.QueryContext(context.TODO(), query)
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
