package types

import "math/big"

const (
	SuccessTxStatus = "Success"
	PendingTxStatus = "Pending"

	MAX_CHUNK_SIZE = 256 * 1024
	MIN_CHUNK_SIZE = 32 * 1024
	NOTE_SIZE      = 32
	HASH_SIZE      = 32

	// number of bits in a big.Word
	WordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	WordBytes      = WordBits / 8
	BranchNodeType = "branch"
	LeafNodeType   = "leaf"

	// Maximum amount of chunks we will upload in the body.
	MAX_CHUNKS_IN_BODY = 1

	// We assume these errors are intermitment and we can try again after a delay:
	// - not_joined
	// - timeout
	// - data_root_not_found (we may have hit a node that just hasn't seen it yet)
	// - exceeds_disk_pool_size_limit
	// We also try again after any kind of unexpected network errors

	// Amount we will delay on receiving an error response but do want to continue.
	ERROR_DELAY = 1000 * 40
)

// Errors from /chunk we should never try and continue on.
var FATAL_CHUNK_UPLOAD_ERRORS = map[string]struct{}{
	"{\"error\":\"disk_full\"}":                        struct{}{},
	"{\"error\":\"invalid_json\"}":                     struct{}{},
	"{\"error\":\"chunk_too_big\"}":                    struct{}{},
	"{\"error\":\"data_path_too_big\"}":                struct{}{},
	"{\"error\":\"offset_too_big\"}":                   struct{}{},
	"{\"error\":\"data_size_too_big\"}":                struct{}{},
	"{\"error\":\"chunk_proof_ratio_not_attractive\"}": struct{}{},
	"{\"error\":\"invalid_proof\"}":                    struct{}{},
}

// about bundle
const (
	BUNDLER_HOST           = "https://node1.bundlr.network"
	MIN_BUNDLE_BINARY_SIZE = 1044
)
