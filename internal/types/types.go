package types

type ParsedQuery struct {
	Insert        bool
	Update        bool
	Table         string
	Query         string
	InsertContent string
	UpdateMap     map[string]string
	Conditions    []string
}
