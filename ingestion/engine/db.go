// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import (
	"encoding/gob"
	"os"
)

// database is a small database used for tracking file progress.
type database struct {
	Files   []file // Files is a list of files
	Next    int    // Next indicates the index of the next file to scan
	AssetID string
	Key     string
}

// file is a file and it's progress.
type file struct {
	Path  string // Path is the location of the file path
	Lines int64  // Lines is the number of lines read and already uploaded
	Size  int64  // Size indicates file size last time file was read
}

// commit stores the database. It will return an error if the database cannot be
// stored.
func (db *database) commit(s *state) error {
	db.AssetID = s.AssetID
	db.Key = s.EncryptionKey
	file, err := os.Create(dbFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return encoder.Encode(db)
}

// dbLoad loads the database. It will return an error if the database cannot be
// loaded.
func dbLoad() (*database, error) {
	var db *database
	file, err := os.Open(dbFileName)
	if err != nil {
		return db, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&db)
	return db, err
}

// fileExists returns if the file exists in the database.
func (db *database) fileExists(path string) bool {
	for _, file := range db.Files {
		if file.Path == path {
			return true
		}
	}
	return false
}

// clean will remove all non existent files from the database.
func (db *database) clean() {
	// list of broken indexes to remove
	broken := []int{}
	for i, file := range db.Files {
		_, err := os.Stat(file.Path)
		if err != nil {
			// file path not existent, delete index
			broken = append(broken, i)
		}

	}
	// remove broken (delete backwards due to index shifts)
	for i := len(broken) - 1; i >= 0; i-- {
		db.Files = removeFile(db.Files, broken[i])
	}
}

// removeFile removes an element from slice, returning a new slice.
func removeFile(slice []file, s int) []file {
	return append(slice[:s], slice[s+1:]...)
}
