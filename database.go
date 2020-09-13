package database

const (
	// SQL database
	SQL = iota
	// NOSQL database
	NOSQL
	// CACHING databse
	CACHING
)

// New database with abstract factory pattern
func New(databaseType int) func(databaseCompany int, config *Config) interface{} {
	switch databaseType {
	case SQL:
		return NewSQL
	case NOSQL:
		return NewNoSQL
	case CACHING:
		return NewCaching
	default:
		return nil
	}
}
