// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/urfave/cli.v1"
)

const (
	appName    = "canids-ingest"
	appUsage   = "realtime file uploader to CanIDS backend"
	appVersion = "2.0.0"
	appAuthor  = "Computing Infrastructure Research Centre, McMaster University"
)

var (
	// parameters are populated/modified by CLI app, defaults shown below (see
	// state for comments)
	valAssetID       = ""
	valHostname      = ""
	valInsecure      = false
	valDebug         = false
	valRetryDelay    = 5 * time.Second
	valFileMode      = zero
	valFileScan      = 5 * time.Second
	valFileChunkSize = 10
)

// Run executes the CLI app to begin ingestion. It will return an error upon
// user error.
func Run() error {
	app := cli.NewApp()

	app.Name = appName
	app.Usage = appUsage
	app.Version = appVersion
	app.Author = appAuthor

	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "asset, uid",
			Usage:       "unique asset (network tap) identifier",
			Destination: &valAssetID,
		},
		cli.StringFlag{
			Name:        "hostname, host",
			Usage:       "hostname and port of CanIDS gRPC backend",
			Destination: &valHostname,
		},
		cli.BoolFlag{
			Name:        "insecure",
			Usage:       "do not enforce secure connection with CanIDS backend",
			Destination: &valInsecure,
		},
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "enable verbose logging",
			Destination: &valDebug,
		},
		cli.DurationFlag{
			Name:        "delay",
			Usage:       "time delay before recovering connection",
			Value:       valRetryDelay,
			Destination: &valRetryDelay,
		},
		cli.DurationFlag{
			Name:        "scan",
			Usage:       "how often to scan file system for new files in directory",
			Value:       valFileScan,
			Destination: &valFileScan,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "upload",
			Aliases: []string{"u"},
			Usage:   "stream data to CanIDS backend",
			Action: func(c *cli.Context) error {
				return cmd(c)
			},
			Flags: flags,
		},
	}
	return app.Run(os.Args)
}

// cmd is called when the required parameters are provided to the CLI. It will
// validate parameters and attempt to start the client.
func cmd(c *cli.Context) error {
	// get + validate number of arguments
	args := c.Args()
	if len(args) == 0 {
		return errNoPath
	}
	if len(args) > 1 {
		return errMultiplePaths
	}
	// ensure hostname provided
	if valHostname == "" {
		return errHostname
	}
	// ensure asset id provided
	if valAssetID == "" || strings.ContainsAny(valAssetID, "`~!@#$%^&*()-_=+[]{}\\|;:'\",.<>/? ") {
		return errAssetID
	}
	// ensure directory/file exists
	valFilePath := args[0]
	info, err := os.Stat(valFilePath)
	if err != nil {
		return errNotFound
	}

	// update mode
	if info.Mode().IsDir() {
		valFileMode = fileDirectory
	} else if info.Mode().IsRegular() {
		valFileMode = fileRegular
	}

	// generate state
	config := &state{
		AssetID:       valAssetID,
		NetworkMutex:  &sync.Mutex{},
		DatabaseMutex: &sync.Mutex{},
		Session:       "",
		Insecure:      valInsecure,
		PollingAbort:  make(chan struct{}),
		ScannerAbort:  make(chan struct{}),
		Debug:         valDebug,
		RetryDelay:    valRetryDelay,
		FilePath:      valFilePath,
		FileMode:      valFileMode,
		FileScan:      valFileScan,
		FileChunkSize: valFileChunkSize,
	}

	// sync the scanner to retreive+update (or create) latest database
	db, err := syncScanner(config)
	if err != nil {
		log.Println("[CanIDS] local database error:", err)
		return nil
	}

	for {
		// initialize connection to gRPC and start
		err = ConnectWebsocketServer(config, db, valHostname)
		if config.Debug {
			log.Println("[CanIDS DEBUG]", err)
		}
		// reset config
		config = &state{
			AssetID:       valAssetID,
			NetworkMutex:  &sync.Mutex{},
			DatabaseMutex: &sync.Mutex{},
			Session:       "",
			Insecure:      valInsecure,
			PollingAbort:  make(chan struct{}),
			ScannerAbort:  make(chan struct{}),
			Debug:         valDebug,
			RetryDelay:    valRetryDelay,
			FilePath:      valFilePath,
			FileMode:      valFileMode,
			FileScan:      valFileScan,
			FileChunkSize: valFileChunkSize,
		}
		time.Sleep(config.RetryDelay)
	}
}
