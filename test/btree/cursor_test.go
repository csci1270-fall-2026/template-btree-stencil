package btree_test

import (
	"testing"

	"dinodb/pkg/btree"
	"dinodb/test/utils"
)

// =====================================================================
// HELPERS
// =====================================================================

// Number of entries inserted by the cursor tests, enough to span many leaf nodes
const cursorTestInserts int64 = 8 * btree.ENTRIES_PER_LEAF_NODE

// checkIteration walks a cursor from the start of the index to the end, checking that it
// visits exactly the keys in [0, numInserts) that aren't in the deleted range [delStart, delEnd).
// It uses Next and GetEntry directly so that failures point at the cursor rather than at Select
func checkIteration(t *testing.T, index *btree.BTreeIndex, numInserts int64, delStart int64, delEnd int64) {
	t.Helper()
	cur, err := index.CursorAtStart()
	if err != nil {
		t.Fatal("Failed to get a cursor at the start of the index:", err)
	}
	defer cur.Close()

	atEnd := false
	for key := range numInserts {
		if key >= delStart && key < delEnd {
			continue
		}
		if atEnd {
			t.Fatalf("Cursor reported the end of the B+Tree before reaching key %d", key)
		}
		entry, err := cur.GetEntry()
		if err != nil {
			t.Fatalf("GetEntry failed where the cursor should have been at key %d: %s", key, err)
		}
		// A cursor left in an empty or exhausted node reads a stale entry,
		// whose value won't match the key it claims to have
		utils.CheckEntry(t, entry, key, generateValue(key))
		atEnd = cur.Next()
	}
	if !atEnd {
		t.Error("Cursor did not report the end of the B+Tree after its last entry")
	}
}

// =====================================================================
// TESTS
// =====================================================================

func TestBTreeCursor(t *testing.T) {
	t.Run("SingleNode", testCursorSingleNode)
	t.Run("AcrossNodes", testCursorAcrossNodes)
}

func testCursorSingleNode(t *testing.T) {
	numInserts := btree.ENTRIES_PER_LEAF_NODE / 2
	index := standardBTreeSetup(t, numInserts)

	// All of the entries fit in the root leaf node
	checkIteration(t, index, numInserts, 0, 0)
	index.Close()
}

func testCursorAcrossNodes(t *testing.T) {
	index := standardBTreeSetup(t, cursorTestInserts)

	// The cursor has to follow right siblings from one node to the next
	checkIteration(t, index, cursorTestInserts, 0, 0)
	index.Close()
}
