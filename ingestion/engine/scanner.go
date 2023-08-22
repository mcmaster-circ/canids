// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// syncScanner prepares the scanner for consuming log entries. It will
// synchronize the existing database or generate a new database. All file(s) in
// the provided FilePath will be added to the local database. Any existing file
// that is no longer present will be deleted. An error will be returned if there
// is a file permission error.
func syncScanner(s *state) (*database, error) {
	s.DatabaseMutex.Lock()
	defer s.DatabaseMutex.Unlock()

	// check if database exists, create if doesn't exist
	db, err := dbLoad()
	if err != nil {
		if s.Debug {
			log.Println("[CanIDS DEBUG] local database does not exist, creating new database")
		}

		// create new entry
		db = &database{}

		// Generate new encryption key
		key, err := CreateKey()
		if err != nil {
			db.AssetID = ""
			db.Key = ""
			return db, err
		}

		db.Key = key
		s.EncryptionKey = key

		// Generate new assetID
		assetID, err := CreateAssetID()
		if err != nil {
			db.AssetID = ""
			return db, err
		}

		db.AssetID = assetID
		s.AssetID = assetID

	} else {
		if s.Debug {
			log.Println("[CanIDS DEBUG] local database exists, syncing database")
		}
		s.EncryptionKey = db.Key
		s.AssetID = db.AssetID
		// clear broken entries
		db.clean()
	}

	log.Printf("[CanIDS] info: s.FilePath %s, s.FileMode %s", s.FilePath, s.FileMode)

	switch s.FileMode {
	case fileRegular:
		// load single file in database (if not already)
		err = processRegularFile(s, s.FilePath, filepath.Base(s.FilePath), db)
		if err != nil {
			return db, err
		}
	case fileDirectory:
		// recursively load all files in database (if not already)
		err = processDirectory(s, s.FilePath, db)
		if len(db.Files) == 0 {
			log.Println("[CanIDS] warning: no files found in directory, nothing to send")
		}
	}

	// commit database changes
	err = db.commit(s)
	if err != nil {
		return db, err
	}
	return db, nil
}

// scannerGetFrame will generate the next frame to be sent over Websockets. If a
// frame cannot be generated, the scanner will sleep until a frame is available.
func scannerGetFrame(s *state, db *database, key []byte) (*UploadRequest, error) {
	s.DatabaseMutex.Lock()

	select {
	case <-s.ScannerAbort:
		s.DatabaseMutex.Unlock()
		// signalled to abort
		return nil, nil
	default:
	}

	// sleep if no files to upload
	if len(db.Files) == 0 {
		if s.Debug {
			log.Println("[CanIDS DEBUG] no files to upload, sleeping for", scannerSleep)
		}
		s.DatabaseMutex.Unlock()
		time.Sleep(scannerSleep)
		return scannerGetFrame(s, db, key)
	}

	// check if there is at least one file that has been modified
	fileIsModified := false
	for i := len(db.Files) - 1; i >= 0; i-- {
		file := db.Files[i]
		// get file info
		info, err := os.Stat(file.Path)
		if err != nil {
			// file cannot be read, remove from database
			if s.Debug {
				log.Println("[CanIDS DEBUG] can no longer read file, removing from local database", file.Path)
			}
			db.Files = removeFile(db.Files, i)
			// commit database
			err = db.commit(s)
			s.DatabaseMutex.Unlock()
			if err != nil {
				return nil, errSavingDatabase
			}
			return scannerGetFrame(s, db, key)
		}
		// check if file has gotten smaller (file rotation)
		if info.Size() < file.Size {
			// file rotated, remove from database
			if s.Debug {
				log.Println("[CanIDS DEBUG] file rotated, removing from local database", file.Path)
			}
			db.Files = removeFile(db.Files, i)
			// commit database
			err = db.commit(s)
			s.DatabaseMutex.Unlock()
			if err != nil {
				return nil, errSavingDatabase
			}
			return scannerGetFrame(s, db, key)
		}
		// file exists, see if file has been modified
		if info.Size() != file.Size {
			fileIsModified = true
			// dont break loop, need to remove all broken files with technique above
		}
	}
	// if no files are candidates for uploading, sleep and try again
	if !fileIsModified {
		if s.Debug {
			log.Println("[CanIDS DEBUG] no changes to upload, sleeping for", scannerSleep)
		}
		s.DatabaseMutex.Unlock()
		time.Sleep(scannerSleep)
		return scannerGetFrame(s, db, key)
	}

	// ensure local database is valid (sync)
	isSync := len(db.Files) > 0 && db.Next < len(db.Files)
	if !isSync {
		// state is not synchronized, must sync
		s.DatabaseMutex.Unlock()
		new, err := syncScanner(s)
		if err != nil {
			return nil, err
		}
		db.Next = 0 // start at zero for synchronization
		db.Files = new.Files
		return scannerGetFrame(s, db, key)
	}

	// get current file info
	file := db.Files[db.Next]
	info, err := os.Stat(file.Path)
	if err != nil {
		// cannot read file, remove from local database, get next frame
		if s.Debug {
			log.Println("[CanIDS DEBUG] cannot read current file, removing from local database", file.Path)
		}
		db.Files = removeFile(db.Files, db.Next)
		// decrement counter to get next file
		if db.Next == 0 {
			db.Next = len(db.Files) - 1
		} else {
			db.Next--
		}
		// commit database
		err = db.commit(s)
		s.DatabaseMutex.Unlock()
		if err != nil {
			return nil, errSavingDatabase
		}
		return scannerGetFrame(s, db, key)
	}

	// if file size hasn't changed, nothing to do, get next frame
	if info.Size() == file.Size {
		// decrement the counter to get next frame
		if db.Next == 0 {
			db.Next = len(db.Files) - 1
		} else {
			db.Next--
		}
		// commit database
		err = db.commit(s)
		s.DatabaseMutex.Unlock()
		if err != nil {
			return nil, errSavingDatabase
		}
		// get next frame
		return scannerGetFrame(s, db, key)
	}

	// generate frame (updated provided file)
	frame, frameErr := generateFrame(s, &file, info.Name(), key)

	// sync modified file with database and commit
	db.Files[db.Next] = file

	// decrement counter to get next file
	if db.Next == 0 {
		db.Next = len(db.Files) - 1
	} else {
		db.Next--
	}
	// commit database
	err = db.commit(s)
	s.DatabaseMutex.Unlock()
	if err != nil {
		return nil, errSavingDatabase
	}
	if frameErr != nil {
		return nil, frameErr
	}
	return frame, err
}
