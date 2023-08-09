// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import (
	"errors"
	"sync"
	"time"
)

const (
	// local database filename
	dbFileName = ".canids-ingestion-v1.0.0.db"

	// scannerSleep indicates how long to sleep for if there is no new frames to
	// generate (used to avoid busy wait and heavy I/O activity)
	scannerSleep = 5 * time.Second

	logConn     = "conn.log"
	logDHCP     = "dhcp.log"
	logDNS      = "dns.log"
	logFTP      = "ftp.log"
	logHTTP     = "http.log"
	logIRC      = "irc.log"
	logKerberos = "kerberos.log"
	logModbus   = "modbus.log"
	logMySQL    = "mysql.log"
	logNTP      = "ntp.log"
	logRadius   = "radius.log"
	logRDP      = "rdp.log"
	logSIP      = "sip.log"
	logSMTP     = "smtp.log"
	logSNMP     = "snmp.log"
	logSocks    = "socks.log"
	logSSH      = "ssh.log"
	logSSL      = "ssl.log"
	logSyslog   = "syslog.log"
	logStats    = "stats.log"
	logTunnel   = "tunnel.log"
	logWeird    = "weird.log"
	logNotice   = "notice.log"
)

var (
	errNoPath         = errors.New("[CanIDS] error: must provide path of file or directory containing Zeek log(s)")
	errMultiplePaths  = errors.New("[CanIDS] error: must provide single path of file or directory containing Zeek log(s)")
	errNotFound       = errors.New("[CanIDS] error: provided file or directory not found or insufficient permissions")
	errReadingFile    = errors.New("[CanIDS] error: failed to read file system, please check permissions")
	errSavingDatabase = errors.New("[CanIDS] error: failed to save local database, please check permissions")
	errEmptyLine      = errors.New("[CanIDS] error: empty line")
	errInvalidLine    = errors.New("[CanIDS] error: invalid line")
	errHostname       = errors.New("[CanIDS] error: must provide hostname of backend")
	errAssetID        = errors.New("[CanIDS] error: must provide unique asset (network tap) identifier, only alphanumeric characters, no spaces")
	errBadJSON        = errors.New("[CanIDS] error: malformed JSON")
	errBadTSV         = errors.New("[CanIDS] error: malformed TSV")
)

// fileMode indicates if a single regular file or directory was passed
type fileMode int

const (
	zero fileMode = 0
	// fileRegular is regular file
	fileRegular fileMode = 1
	// fileDirectory is directory
	fileDirectory fileMode = 2
)

// state represents client state.
type state struct {
	AssetID       string        // AssetID identifies the data in the database
	NetworkMutex  *sync.Mutex   // NetworkMutex is for preventing concurrent writing to websocket
	DatabaseMutex *sync.Mutex   // DatabaseMutex is for preventing concurrent operations to local database
	Session       string        // Session is the session identifier
	PollingAbort  chan struct{} // PollingAbort is for signalling the file system polling loop to terminate
	ScannerAbort  chan struct{} // ScannerAbort is for signalling the recursive scanner to terminate
	Debug         bool          // Debug indicates if debugging logging should be used
	RetryDelay    time.Duration // RetryDelay is delay before attempting reconnect
	FilePath      string        // FilePath is the file or directory to upload form
	FileMode      fileMode      // FileMode indicates type of file mode being used (regular file or directory provided)
	FileScan      time.Duration // FileScan indicates how often to scan for new files on the file system
	FileChunkSize int           // FileChunkSize indicates number of lines to send in frame
}
