package tarantella

import "math"

// from https://github.com/tarantool/tarantool/blob/5d658e7e1aceba1daef8491d321941f08bbd7cfd/src/box/iproto_constants.h#L68
// enum iproto_key
const (
	IPROTO_REQUEST_TYPE uint64 = 0x00
	IPROTO_SYNC         uint64 = 0x01

	/* Replication keys (header) */
	IPROTO_REPLICA_ID     uint64 = 0x02
	IPROTO_LSN            uint64 = 0x03
	IPROTO_TIMESTAMP      uint64 = 0x04
	IPROTO_SCHEMA_VERSION uint64 = 0x05
	IPROTO_SERVER_VERSION uint64 = 0x06
	IPROTO_GROUP_ID       uint64 = 0x07
	IPROTO_TSN            uint64 = 0x08
	IPROTO_FLAGS          uint64 = 0x09
	IPROTO_STREAM_ID      uint64 = 0x0a
	/* Leave a gap for other keys in the header. */
	IPROTO_SPACE_ID   uint64 = 0x10
	IPROTO_INDEX_ID   uint64 = 0x11
	IPROTO_LIMIT      uint64 = 0x12
	IPROTO_OFFSET     uint64 = 0x13
	IPROTO_ITERATOR   uint64 = 0x14
	IPROTO_INDEX_BASE uint64 = 0x15
	/* Leave a gap between integer values and other keys */
	/**
	 * Flag indicating the need to send position of
	 * last selected tuple in response.
	 */
	IPROTO_FETCH_POSITION uint64 = 0x1f
	IPROTO_KEY            uint64 = 0x20
	IPROTO_TUPLE          uint64 = 0x21
	IPROTO_FUNCTION_NAME  uint64 = 0x22
	IPROTO_USER_NAME      uint64 = 0x23

	/*
	 * Replication keys (body).
	 * Unfortunately, there is no gap between request and
	 * replication keys (between USER_NAME and INSTANCE_UUID).
	 * So imagine, that OPS, EXPR and FIELD_NAME keys follows
	 * the USER_NAME key.
	 */
	IPROTO_INSTANCE_UUID   uint64 = 0x24
	IPROTO_REPLICASET_UUID uint64 = 0x25
	IPROTO_VCLOCK          uint64 = 0x26

	/* Also request keys. See the comment above. */
	IPROTO_EXPR       uint64 = 0x27 /* EVAL */
	IPROTO_OPS        uint64 = 0x28 /* UPSERT but not UPDATE ops, because of legacy */
	IPROTO_BALLOT     uint64 = 0x29
	IPROTO_TUPLE_META uint64 = 0x2a
	IPROTO_OPTIONS    uint64 = 0x2b
	/** Old tuple (i.e. before DML request is applied). */
	IPROTO_OLD_TUPLE uint64 = 0x2c
	/** New tuple (i.e. result of DML request). */
	IPROTO_NEW_TUPLE uint64 = 0x2d
	/** Position of last selected tuple to start iteration after it. */
	IPROTO_AFTER_POSITION uint64 = 0x2e
	/** Last selected tuple to start iteration after it. */
	IPROTO_AFTER_TUPLE uint64 = 0x2f

	/** Response keys. */
	IPROTO_DATA     uint64 = 0x30
	IPROTO_ERROR_24 uint64 = 0x31
	/**
	 * IPROTO_METADATA: [
	 *      { IPROTO_FIELD_NAME: name },
	 *      { ... },
	 *      ...
	 * ]
	 */
	IPROTO_METADATA      uint64 = 0x32
	IPROTO_BIND_METADATA uint64 = 0x33
	IPROTO_BIND_COUNT    uint64 = 0x34
	/** Position of last selected tuple in response. */
	IPROTO_POSITION uint64 = 0x35

	/* Leave a gap between response keys and SQL keys. */
	IPROTO_SQL_TEXT uint64 = 0x40
	IPROTO_SQL_BIND uint64 = 0x41
	/**
	 * IPROTO_SQL_INFO: {
	 *     SQL_INFO_ROW_COUNT: number
	 * }
	 */
	IPROTO_SQL_INFO uint64 = 0x42
	IPROTO_STMT_ID  uint64 = 0x43
	/* Leave a gap between SQL keys and additional request keys */
	IPROTO_REPLICA_ANON uint64 = 0x50
	IPROTO_ID_FILTER    uint64 = 0x51
	IPROTO_ERROR        uint64 = 0x52
	/**
	 * Term. Has the same meaning as IPROTO_RAFT_TERM, but is an iproto
	 * key, rather than a raft key. Used for PROMOTE request, which needs
	 * both iproto (e.g. REPLICA_ID) and raft (RAFT_TERM) keys.
	 */
	IPROTO_TERM uint64 = 0x53
	/** Protocol version. */
	IPROTO_VERSION uint64 = 0x54
	/** Protocol features. */
	IPROTO_FEATURES uint64 = 0x55
	/** Operation timeout. Specific to request type. */
	IPROTO_TIMEOUT uint64 = 0x56
	/** Key name and data sent to a remote watcher. */
	IPROTO_EVENT_KEY  uint64 = 0x57
	IPROTO_EVENT_DATA uint64 = 0x58
	/** Isolation level, is used only by IPROTO_BEGIN request. */
	IPROTO_TXN_ISOLATION uint64 = 0x59
	/** A vclock synchronisation request identifier. */
	IPROTO_VCLOCK_SYNC uint64 = 0x5a
	/**
	 * Name of the authentication method that is currently used on
	 * the server (value of box.cfg.auth_type). It's sent in reply
	 * to IPROTO_ID request. A client can use it as the default
	 * authentication method.
	 */
	IPROTO_AUTH_TYPE uint64 = 0x5b
	/*
	 * Be careful to not extend iproto_key values over 0x7f.
	 * iproto_keys are encoded in msgpack as positive fixnum, which ends at
	 * 0x7f, and we rely on this in some places by allocating a uint8_t to
	 * hold a msgpack-encoded key value.
	 */
	IPROTO_KEY_MAX uint64 = IPROTO_AUTH_TYPE + 1
)

// from https://github.com/tarantool/tarantool/blob/5d658e7e1aceba1daef8491d321941f08bbd7cfd/src/box/iproto_constants.h#L233
// command codes
// enum iproto_type
const (
	/** Acknowledgement that request or command is successful */
	IPROTO_OK uint64 = 0

	/** SELECT request */
	IPROTO_SELECT uint64 = 1
	/** INSERT request */
	IPROTO_INSERT uint64 = 2
	/** REPLACE request */
	IPROTO_REPLACE uint64 = 3
	/** UPDATE request */
	IPROTO_UPDATE uint64 = 4
	/** DELETE request */
	IPROTO_DELETE uint64 = 5
	/** CALL request - wraps result into [tuple, tuple, ...] format */
	IPROTO_CALL_16 uint64 = 6
	/** AUTH request */
	IPROTO_AUTH uint64 = 7
	/** EVAL request */
	IPROTO_EVAL uint64 = 8
	/** UPSERT request */
	IPROTO_UPSERT uint64 = 9
	/** CALL request - returns arbitrary MessagePack */
	IPROTO_CALL uint64 = 10
	/** Execute an SQL statement. */
	IPROTO_EXECUTE uint64 = 11
	/** No operation. Treated as DML, used to bump LSN. */
	IPROTO_NOP uint64 = 12
	/** Prepare SQL statement. */
	IPROTO_PREPARE uint64 = 13
	/* Begin transaction */
	IPROTO_BEGIN uint64 = 14
	/* Commit transaction */
	IPROTO_COMMIT uint64 = 15
	/* Rollback transaction */
	IPROTO_ROLLBACK uint64 = 16
	/** The maximum typecode used for box.stat() */
	IPROTO_TYPE_STAT_MAX uint64 = IPROTO_ROLLBACK + 1 // was hardcoded by me

	IPROTO_RAFT uint64 = 30
	/** PROMOTE request. */
	IPROTO_RAFT_PROMOTE uint64 = 31
	/** DEMOTE request. */
	IPROTO_RAFT_DEMOTE uint64 = 32

	/** A confirmation message for synchronous transactions. */
	IPROTO_RAFT_CONFIRM uint64 = 40
	/** A rollback message for synchronous transactions. */
	IPROTO_RAFT_ROLLBACK uint64 = 41

	/** PING request */
	IPROTO_PING uint64 = 64
	/** Replication JOIN command */
	IPROTO_JOIN uint64 = 65
	/** Replication SUBSCRIBE command */
	IPROTO_SUBSCRIBE uint64 = 66
	/** DEPRECATED: use IPROTO_VOTE instead */
	IPROTO_VOTE_DEPRECATED uint64 = 67
	/** Vote request command for master election */
	IPROTO_VOTE uint64 = 68
	/** Anonymous replication FETCH SNAPSHOT. */
	IPROTO_FETCH_SNAPSHOT uint64 = 69
	/** REGISTER request to leave anonymous replication. */
	IPROTO_REGISTER      uint64 = 70
	IPROTO_JOIN_META     uint64 = 71
	IPROTO_JOIN_SNAPSHOT uint64 = 72
	/** Protocol features request. */
	IPROTO_ID uint64 = 73
	/**
	 * The following three request types are used by the remote watcher
	 * protocol (box.watch over network), which operates as follows:
	 *
	 *  1. The client sends an IPROTO_WATCH packet to subscribe to changes
	 *     of a specified key defined on the server.
	 *  2. The server sends an IPROTO_EVENT packet to the subscribed client
	 *     with the key name and its current value unconditionally after
	 *     registration and then every time the key value is updated
	 *     provided the last notification was acknowledged (see below).
	 *  3. Upon receiving a notification, the client sends an IPROTO_WATCH
	 *     packet to acknowledge the notification.
	 *  4. When the client doesn't want to receive any more notifications,
	 *     it unsubscribes by sending an IPROTO_UNWATCH packet.
	 *
	 * All the three request types are fully asynchronous - a receiving end
	 * doesn't send a packet in reply to any of them (therefore neither of
	 * them has a sync number).
	 */
	IPROTO_WATCH   uint64 = 74
	IPROTO_UNWATCH uint64 = 75
	IPROTO_EVENT   uint64 = 76

	/** Vinyl run info stored in .index file */
	VY_INDEX_RUN_INFO uint64 = 100
	/** Vinyl page info stored in .index file */
	VY_INDEX_PAGE_INFO uint64 = 101
	/** Vinyl row index stored in .run file */
	VY_RUN_ROW_INDEX uint64 = 102

	/** Non-final response type. */
	IPROTO_CHUNK uint64 = 128

	/**
	 * Error codes uint64 =(IPROTO_TYPE_ERROR | ER_XXX from errcode.h)
	 */
	IPROTO_TYPE_ERROR uint64 = 1 << 15

	/**
	 * Used for overriding the unknown request handler.
	 */
	IPROTO_UNKNOWN uint64 = math.MaxUint64
)
