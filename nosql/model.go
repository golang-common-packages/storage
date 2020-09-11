package nosql

///// MongoDB query model /////

// MatchLookup ...
type MatchLookup struct {
	Match  []Match  `json:"match"`
	Lookup []Lookup `json:"lookup"`
}

// Match ...
type Match struct {
	Field    string              `json:"field"`
	Operator ComparisonOperators `json:"operator"`
	Value    string              `json:"value"`
}

// Lookup ...
type Lookup struct {
	From         string `json:"From"`
	LocalField   string `json:"localField"`
	ForeignField string `json:"foreignField"`
	As           string `json:"as"`
}

// Set ...
type Set struct {
	Operator UpdateOperators `json:"operator"`
	Data     interface{}     `json:"data"`
}

///// MongoDB operator model /////

type ComparisonOperators string

const (
	Equal                ComparisonOperators = "$eq"
	EqualAny             ComparisonOperators = "$in"
	NotEqual             ComparisonOperators = "$ne"
	NotEqualAnyLL        ComparisonOperators = "$nin"
	GreaterThan          ComparisonOperators = "$gt"
	GreaterThanOrEqualTo ComparisonOperators = "$gte"
	LessThan             ComparisonOperators = "$lt"
	LessThanOrEqualTo    ComparisonOperators = "$lte"
)

type UpdateOperators string

const (
	Replaces UpdateOperators = "$set"
)
