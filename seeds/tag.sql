CREATE TABLE history.tag
(
`fia_call_id` Nullable(Int64),
`fia_member_id` Nullable(Int32),
`fia_strunk_id` Nullable(Int32),
`fib_call_id` Nullable(Int64),
`fib_member_id` Nullable(Int32),
`fib_strunk_id` Nullable(Int32),
`ficonference_id` Nullable(Int32),
`facontext_create` DateTime('Europe/Moscow'),   /* ? */
`facreate` Nullable(DateTime('Europe/Moscow')), /* ? */
`fiproduct_id` Int32,
`fia_tag_ids` Array(Int32),
`fia_tag_modified` Nullable(DateTime('Europe/Moscow')),
`fib_tag_ids` Array(Int32),
`fib_tag_modified` Nullable(DateTime('Europe/Moscow')),
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(facontext_create)
ORDER BY (fiproduct_id, facontext_create, fia_call_id, fib_call_id)
SETTINGS index_granularity = 256