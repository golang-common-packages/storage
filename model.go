package database

// Database model for database config
type Database struct {
	MongoDB MongoDB `json:"mongodb"`
}

// MongoDB model for MongoDB config
type MongoDB struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Hosts    []string `json:"hosts"`
	DB       string   `json:"db"`
	Options  []string `json:"options"`
}

type MatchLookup struct {
	Match  []Match  `json:"match"`
	Lookup []Lookup `json:"lookup"`
}

type Match struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Lookup struct {
	From         string `json:"From"`
	LocalField   string `json:"localField"`
	ForeignField string `json:"foreignField"`
	As           string `json:"as"`
}
