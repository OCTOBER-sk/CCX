package stream

// Event journal per spec section 89 (lines 3133-3172): bounded journal with sequence, timestamp, event_type, block_id, call_id, payload_hash.
// Enables duplicate/missing detection, debugging, fixture generation, stream reconstruction.
// Never store full payloads by default.

type JournalEntry struct {
	Sequence    int    // sequence number
	Timestamp   int64  // timestamp
	EventType   string // event_type (message_start, content_block_delta, etc.)
	BlockID     string // block_id
	CallID      string // call_id (for tool calls)
	PayloadHash string // hash of payload, not full payload
}

var journal []JournalEntry
var maxEntries = 1000 // bounded per spec 89 (lines 3133-3172) + invariant 8 (spec 515, unbounded buffering forbidden).

func AddEntry(entry JournalEntry) {
	journal = append(journal, entry)
	if len(journal) > maxEntries {
		journal = journal[1:] // drop oldest
	}
}

func GetEntries() []JournalEntry {
	return journal
}
