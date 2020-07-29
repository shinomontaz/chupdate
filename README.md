# chupdate
clickhouse update middleware

## When to use
1. You use ReplacingMergeTree engine family
2. You need updates in CH ( and reasonable way to do that is to use insert new row and select with *LIMIT BY* )
3. When sending updates to CH in proper way is complicated due to client software limitations ( you don't want to make many changes there)
4. You need to work with one CH master

### Installation

### Inside

Chupdate just simply parse incoming queries and search for updates, when it found it we just make an addiotional select, get actual information and insert full new row.
Also updercover we chupdate keeps connection pool to CH for speed up tasks.

Insired by clickhouse-bulk

### Metrics

#TODO

### Restriction
1. Target tables must have version column with name "time". Now it is hardcoded ( #TODO )
2. Queries must include full qualified table names ( with database name )
