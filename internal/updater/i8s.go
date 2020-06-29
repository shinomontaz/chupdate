package updater

type Inserter interface {
	Push(q, params string)
}
