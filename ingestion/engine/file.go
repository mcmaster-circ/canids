// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

// processRegularFile will add the file path to the database if it's not already
// present. It will return an error if the file cannot be read.
func processRegularFile(s *state, filePath string, fileName string, db *database) error {
	// get absolute path of single file
	abs, err := filepath.Abs(filePath)
	if err != nil {
		return errReadingFile
	}
	// do nothing if file exists
	if db.fileExists(abs) {
		if s.Debug {
			log.Println("[CanIDS DEBUG]", "file in local database", abs)
		}
		return nil
	}
	// create new file (no lines/bytes read), add to database, commit
	f := file{
		Path:  abs,
		Lines: 0,
		Size:  0,
	}
	// conn.log dns.log http.log notice.log sip.log ssl.log stats.log weird.log telemetry.log capture_loss.log
	whitelist := []string{"conn.log", "dns.log", "http.log", "sip.log", "ssl.log", "stats.log", "weird.log", "telemetry.log"}
	if !slices.Contains(whitelist, fileName) {
		if s.Debug {
			log.Println("[CanIDS DEBUG]", "Ignoring non-whitelisted file", abs)
		}
		return nil
	}
	db.Files = append(db.Files, f)
	if s.Debug {
		log.Println("[CanIDS DEBUG]", "file not in local database, adding file", abs)
	}
	return nil
}

// processDiretory will recursively add all files in the file path to the
// database if it's not already present. It will return an error if the file
// cannot be read.
func processDirectory(s *state, filePath string, db *database) error {
	// recursively walk file path
	err := filepath.Walk(filePath,
		func(path string, info os.FileInfo, err error) error {
			// if regular file, process it
			if info.Mode().IsRegular() {
				return processRegularFile(s, path, info.Name(), db)
			}
			return nil
		})
	return err
}
