package raftboltdb

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func testBoltStore(t *testing.T) *BoltStore {
	fh, err := ioutil.TempFile("", "bolt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	os.Remove(fh.Name())

	// Successfully creates and returns a store
	store, err := NewBoltStore(fh.Name())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return store
}

func TestNewBoltStore(t *testing.T) {
	fh, err := ioutil.TempFile("", "bolt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	os.Remove(fh.Name())
	defer os.Remove(fh.Name())

	// Successfully creates and returns a store
	store, err := NewBoltStore(fh.Name())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure the file was created
	if store.path != fh.Name() {
		t.Fatalf("unexpected file path %q", store.path)
	}
	if _, err := os.Stat(fh.Name()); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Close the store so we can open again
	if err := store.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure our tables were created
	db, err := bolt.Open(fh.Name(), dbFileMode, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tx, err := db.Begin(true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if _, err := tx.CreateBucket([]byte(dbLogs)); err != bolt.ErrBucketExists {
		t.Fatalf("bad: %v", err)
	}
	if _, err := tx.CreateBucket([]byte(dbConf)); err != bolt.ErrBucketExists {
		t.Fatalf("bad: %v", err)
	}
}

func TestBoltStore_FirstIndex(t *testing.T) {
	store := testBoltStore(t)
	defer store.Close()
	defer os.Remove(store.path)

	// Should get 0 index on empty log
	idx, err := store.FirstIndex()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if idx != 0 {
		t.Fatalf("bad: %v", idx)
	}

	// Set a mock raft log
	err = store.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dbLogs))
		bucket.Put(uint64ToBytes(1), []byte("log1"))
		bucket.Put(uint64ToBytes(2), []byte("log2"))
		return nil
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Fetch the first Raft index
	idx, err = store.FirstIndex()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if idx != 1 {
		t.Fatalf("bad: %d", idx)
	}
}

func TestBoltStore_LastIndex(t *testing.T) {
	store := testBoltStore(t)
	defer store.Close()
	defer os.Remove(store.path)

	// Should get 0 index on empty log
	idx, err := store.LastIndex()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if idx != 0 {
		t.Fatalf("bad: %v", idx)
	}

	// Set a mock raft log
	err = store.conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dbLogs))
		bucket.Put(uint64ToBytes(1), []byte("log1"))
		bucket.Put(uint64ToBytes(2), []byte("log2"))
		return nil
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Fetch the last Raft index
	idx, err = store.LastIndex()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if idx != 2 {
		t.Fatalf("bad: %d", idx)
	}
}