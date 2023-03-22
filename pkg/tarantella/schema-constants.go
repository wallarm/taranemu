package tarantella

// from https://github.com/tarantool/tarantool/blob/5d658e7e1aceba1daef8491d321941f08bbd7cfd/src/box/schema_def.h#L66

const (
	/** Start of the reserved range of system spaces. */
	BOX_SYSTEM_ID_MIN uint64 = 256
	/** Space if of _vinyl_deferred_delete. */
	BOX_VINYL_DEFERRED_DELETE_ID uint64 = 257
	/** Space id of _schema. */
	BOX_SCHEMA_ID uint64 = 272
	/** Space id of _collation. */
	BOX_COLLATION_ID uint64 = 276
	/** Space id of _vcollation. */
	BOX_VCOLLATION_ID uint64 = 277
	/** Space id of _space. */
	BOX_SPACE_ID uint64 = 280
	/** Space id of _vspace view. */
	BOX_VSPACE_ID uint64 = 281
	/** Space id of _sequence. */
	BOX_SEQUENCE_ID uint64 = 284
	/** Space id of _sequence_data. */
	BOX_SEQUENCE_DATA_ID uint64 = 285
	/** Space id of _vsequence view. */
	BOX_VSEQUENCE_ID uint64 = 286
	/** Space id of _index. */
	BOX_INDEX_ID uint64 = 288
	/** Space id of _vindex view. */
	BOX_VINDEX_ID uint64 = 289
	/** Space id of _func. */
	BOX_FUNC_ID uint64 = 296
	/** Space id of _vfunc view. */
	BOX_VFUNC_ID uint64 = 297
	/** Space id of _user. */
	BOX_USER_ID uint64 = 304
	/** Space id of _vuser view. */
	BOX_VUSER_ID uint64 = 305
	/** Space id of _priv. */
	BOX_PRIV_ID uint64 = 312
	/** Space id of _vpriv view. */
	BOX_VPRIV_ID uint64 = 313
	/** Space id of _cluster. */
	BOX_CLUSTER_ID uint64 = 320
	/** Space id of _trigger. */
	BOX_TRIGGER_ID uint64 = 328
	/** Space id of _truncate. */
	BOX_TRUNCATE_ID uint64 = 330
	/** Space id of _space_sequence. */
	BOX_SPACE_SEQUENCE_ID uint64 = 340
	/** Space id of _vspace_sequence. */
	BOX_VSPACE_SEQUENCE_ID uint64 = 341
	/** Space id of _fk_constraint. */
	BOX_FK_CONSTRAINT_ID uint64 = 356
	/** Space id of _ck_contraint. */
	BOX_CK_CONSTRAINT_ID uint64 = 364
	/** Space id of _func_index. */
	BOX_FUNC_INDEX_ID uint64 = 372
	/** Space id of _session_settings. */
	BOX_SESSION_SETTINGS_ID uint64 = 380
	/** End of the reserved range of system spaces. */
	BOX_SYSTEM_ID_MAX uint64 = 511
	BOX_ID_NIL        uint64 = 2147483647
)
