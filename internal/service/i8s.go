package service

type Inserter interface {
	Push(q, params string)
}

type Updater interface {
	Push(q string)
}

type Parser interface {
	Parse(body string) (params, content string, insert, update bool, err error)
}
