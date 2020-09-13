package database

// ISQL factory pattern interface
type ISQL interface {
	Execute(
		query string,
		dataModel interface{}) (interface{}, error)
}

const (
	// SQLLike database (common relational database)
	SQLLike = iota
)

// NewSQL factory pattern
func NewSQL(
	databaseCompany int,
	config *Config) interface{} {

	switch databaseCompany {
	case SQLLike:
		return NewSQLLike(&config.LIKE)
	}

	return nil
}
