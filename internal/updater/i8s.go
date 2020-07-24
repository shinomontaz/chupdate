package updater

type Inserter interface {
	Push(q, params string)
}

type Parser interface {
	Updateparse(body string) (table, where string, cols []string, values []string, condition_cols []string)
}
