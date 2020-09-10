package updater

import "github.com/shinomontaz/chupdate/internal/types"

type Inserter interface {
	Push(pq *types.ParsedQuery)
}

type Parser interface {
}
