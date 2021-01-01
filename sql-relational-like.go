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

// newSQLLike init new instance
// The sql package must be used in conjunction with a database driver. See https://golang.org/s/sqldrivers for a list of driverNames.
func newSQLLike(config *LIKE) *SQLLikeClient {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal SQL-Like configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentSQLLikeSession := sqlLikeClientSessionMapping[configAsString]
	if currentSQLLikeSession == nil {
		currentSQLLikeSession = &SQLLikeClient{nil, nil}

		client, err := sql.Open(config.DriverName, config.DataSourceName)
		if err != nil {
			log.Fatalln("Unable to connect to SQL-Like: ", err)
		}

		client.SetConnMaxLifetime(config.MaxConnectionLifetime)
		client.SetMaxIdleConns(config.MaxConnectionIdle)
		client.SetMaxOpenConns(config.MaxConnectionOpen)

		if err := client.PingContext(context.TODO()); err != nil {
			log.Fatalln("Unable to ping to SQL-Like: ", err)
		}

		currentSQLLikeSession.Client = client
		currentSQLLikeSession.Config = config
		sqlLikeClientSessionMapping[configAsString] = currentSQLLikeSession
		log.Println("Connected to SQL-Like")
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
		log.Println("Unable to execute query: ", err)
		return nil, err
	}

	// Go through each row to get the result
	for rows.Next() {
		err = rows.Scan(&dataModel)
		if err != nil {
			log.Println("Unable to scan rows data: ", err)
			return nil, err
		}
		results = append(results, dataModel)
	}
	/*
		Check for errors during rows "Close"
		his may be more important if multiple statements are executed
		in a single batch and rows were written as well as read.
	*/
	if err := rows.Close(); err != nil {
		return nil, err
	}

	// Check for errors during row iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
