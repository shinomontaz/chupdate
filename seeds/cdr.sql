CREATE TABLE history.cdr
(
    `ficdr_id` Int64, 
    `fiorder` Nullable(Int16), 
    `fiparams` Nullable(Int16), 
    `fsani` Nullable(String), 
    `fsdnis` Nullable(String), 
    `fiswitch_id` Nullable(Int16), 
    `fia_call_id` Nullable(Int64), 
    `fia_member_id` Nullable(Int32), 
    `fia_strunk_id` Nullable(Int32), 
    `fib_call_id` Nullable(Int64), 
    `fib_member_id` Nullable(Int32), 
    `fib_strunk_id` Nullable(Int32), 
    `fiend_reason` Nullable(Int16), 
    `ficonference_id` Nullable(Int32), 
    `fiacd_group_id` Nullable(Int32), 
    `fiacd_session_id` Nullable(Int32), 
    `faacd_create` Nullable(DateTime('Europe/Moscow')), 
    `ficampaign_id` Nullable(Int32), 
    `ficampaign_task_id` Nullable(Int32), 
    `facreate` Nullable(DateTime('Europe/Moscow')), 
    `fatalk` Nullable(DateTime('Europe/Moscow')), 
    `faend` Nullable(DateTime('Europe/Moscow')), 
    `fiwait_duration` Nullable(Int16), 
    `fitalk_duration` Nullable(Int16), 
    `fiqueue_duration` Nullable(Int16), 
    `fihold_duration` Nullable(Int16), 
    `fsivr_way` Nullable(String), 
    `fiproduct_id` Int32, 
    `ficontext_id` Int64, 
    `fiinit_type` Int8, 
    `ficontext_params` Int16, 
    `fscontext_ani` Nullable(String), 
    `fscontext_dnis` Nullable(String), 
    `fsline` Nullable(String), 
    `facontext_create` DateTime('Europe/Moscow'), 
    `faservice` Nullable(DateTime('Europe/Moscow')), 
    `facontext_talk` Nullable(DateTime('Europe/Moscow')), 
    `facontext_end` DateTime('Europe/Moscow'), 
    `ficontext_end_reason` Int16, 
    `ficontext_member_ids` Array(Int32), 
    `ficontext_acd_group_ids` Array(Int32), 
    `fibilling_user_session_id` Nullable(Int64), 
    `ficallback_button_id` Nullable(Int32), 
    `fivcdr_id` Nullable(Int64), 
    `fistrunk_id_line` Nullable(Int32), 
    `finstrunk_id` Array(Int32), 
    `fia_comment_member_id` Nullable(Int32), 
    `fsa_comment` Nullable(String), 
    `faa_comment_created` Nullable(DateTime('Europe/Moscow')), 
    `fia_comment_id` Nullable(Int32), 
    `fib_comment_member_id` Nullable(Int32), 
    `fsb_comment` Nullable(String), 
    `fab_comment_created` Nullable(DateTime('Europe/Moscow')), 
    `fib_comment_id` Nullable(Int32), 
    `fia_tag_ids` Array(Int32), 
    `fia_tag_modified` Nullable(DateTime('Europe/Moscow')), 
    `fib_tag_ids` Array(Int32), 
    `fib_tag_modified` Nullable(DateTime('Europe/Moscow')), 
    `fiqcs_ctrl_rule_id` Nullable(Int32), 
    `fiqcs_answer_id` Nullable(Int32), 
    `fiqcs_answer_value` Nullable(Int16), 
    `faqci_modified` Nullable(DateTime('Europe/Moscow')), 
    `fiqci_qual_ctrl_form_id` Nullable(Int32), 
    `fiqci_controller_mbr_id` Nullable(Int32), 
    `fiqci_listen_time` Nullable(Int32), 
    `fxqci_mark_avg` Nullable(Decimal(4, 2)), 
    `fjqci_mark` Nullable(String), 
    `fsqci_comment` Nullable(String), 
    `ficancel_recall_member` Nullable(Int32), 
    `facancel_recall_datetime` Nullable(DateTime('Europe/Moscow'))
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(facontext_create)
ORDER BY (fiproduct_id, facontext_create, ficontext_id, ficdr_id)
SETTINGS index_granularity = 256