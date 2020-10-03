package storage

// ISQLRelation factory pattern interface
type ISQLRelation interface {
	Execute(query string, dataModel interface{}) (interface{}, error)
}

const (
	// SQLLike database (common relational database)
	SQLLike = iota
)

// NewSQLRelational factory pattern
func NewSQLRelational(databaseCompany int, config *Config) interface{} {

	switch databaseCompany {
	case SQLLike:
		return NewSQLLike(&config.LIKE)
	}

	return nil
}
