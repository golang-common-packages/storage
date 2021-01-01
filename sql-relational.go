package storage

// ISQLRelational factory pattern interface
type ISQLRelational interface {
	Execute(query string, dataModel interface{}) (interface{}, error)
}

const (
	// SQLLike database (common relational database)
	SQLLike = iota
)

// newSQLRelational factory pattern
func newSQLRelational(databaseCompany int, config *Config) interface{} {

	switch databaseCompany {
	case SQLLike:
		return newSQLLike(&config.LIKE)
	}

	return nil
}
