CREATE TABLE history.comment
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
`fia_comment_member_id` Nullable(Int32),
`fsa_comment` Nullable(String),
`faa_comment_created` Nullable(DateTime('Europe/Moscow')),
`fia_comment_id` Nullable(Int32),
`fib_comment_member_id` Nullable(Int32),
`fsb_comment` Nullable(String),
`fab_comment_created` Nullable(DateTime('Europe/Moscow')),
`fib_comment_id` Nullable(Int32),
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(facontext_create)
ORDER BY (fiproduct_id, facontext_create, fia_call_id, fib_call_id)
SETTINGS index_granularity = 256