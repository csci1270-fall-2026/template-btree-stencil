package btree_test

import (
	"fmt"
	"testing"

	"dinodb/pkg/btree"
	"dinodb/test/utils"
)

func TestBTreeSelect(t *testing.T) {
	t.Run("Increasing", TestSelectIncreasing)
}

func TestBTreeSelectRange(t *testing.T) {
	t.Run("Specific", TestSelectRangeSpecific)
	t.Run("Delete", TestSelectRangeDelete)
}

/*
Create and run basic test of inserting X entries + validating they are in the index with Select()
*/
func stageSelectIncreasingTest(numEntries int64) func(t *testing.T) {
	return func(t *testing.T) {
		index := standardBTreeSetup(t, numEntries)

		// Retrieve entries
		entries, err := index.Select()
		if err != nil {
			t.Error(err)
		}

		// check that size of entries slice is expected
		if int64(len(entries)) != numEntries {
			err = fmt.Errorf("Wrong number of entries returned by Select; len(entries) == %d; expected len(entries) is %d", int64(len(entries)), numEntries)
			t.Error(err)
		}
		for i, entry := range entries {
			key := int64(i)
			utils.CheckEntry(t, entry, key, generateValue(key))
		}
		index.Close()
	}
}

/*
Creates a BTree index, inserts entries with increasing keys,
and then retrieves all of the entries through Select
*/
func TestSelectIncreasing(t *testing.T) {
	// Define test cases, maps test name to number of entries inserted
	tests := map[string]int64{
		"Ten":     10,
		"Hundred": 100,
	}

	for name, numInserts := range tests {
		t.Run(name, stageSelectIncreasingTest(numInserts))
	}
}


/*
Creates a BTree index, inserts 1000 entries, and then retrieves some of the
entries through SelectRange
*/
func TestSelectRangeSpecific(t *testing.T) {
	index := standardBTreeSetup(t, 1000)

	// Retrieve entries
	start := int64(20)
	end := int64(100)
	entries, err := index.SelectRange(start, end)
	if err != nil {
		t.Error(err)
	}
	// check that size of entries slice is expected
	expectedLenEntries := (end - start)
	if int64(len(entries)) != expectedLenEntries {
		err = fmt.Errorf("Wrong number of entries returned by SelectRange; len(entries) == %d; expected len(entries) is %d", int64(len(entries)), expectedLenEntries)
		t.Error(err)
	}
	for i, entry := range entries {
		key := int64(i) + start
		utils.CheckEntry(t, entry, key, generateValue(key))
	}
	index.Close()
}

/*
Creates a BTree index, inserts 1000 entries, deletes some entries,
and makes sure deleted entries are not found in SelectRange
*/
func TestSelectRangeDelete(t *testing.T) {
	index := standardBTreeSetup(t, 1000)

	// Removes entries 200 to 499
	amountToDelete := int64(300)
	for i := range amountToDelete {
		err := index.Delete(i + 200)
		if err != nil {
			t.Error(err)
		}
	}
	// Retrieve all entries using SelectRange
	start := int64(0)
	end := int64(1000)
	entries, err := index.SelectRange(start, end)
	if err != nil {
		t.Error(err)
	}
	expectedLenEntries := ((end - start) - amountToDelete)
	// check that size of entries slice is expected
	if int64(len(entries)) != expectedLenEntries {
		err = fmt.Errorf("Wrong number of entries returned by SelectRange; len(entries) == %d; expected len(entries) is %d", int64(len(entries)), expectedLenEntries)
		t.Error(err)
	}
	// check that none of the entries are the deleted ones
	for _, entry := range entries {
		if entry.Key >= int64(200) && entry.Key < int64(500) {
			t.Error("Deleted entry found in slice returned from SelectRange")
			break
		}
	}
	index.Close()
}


