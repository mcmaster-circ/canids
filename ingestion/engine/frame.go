package engine

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mcmaster-circ/canids-v2/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// generateFrame state and local database file. It will attempt to read
// unread lines in the file. For each line, the line will be parsed and generate
// a payload entry. If the line is not valid, it will be ignored. It
// also updates the provided file, updating how much if the file was read. It
// will return complete frame or an error.
func generateFrame(s *state, f *file, baseName string) (*protocol.UploadRequest, error) {
	// open file
	fs, err := os.Open(f.Path)
	if err != nil {
		return nil, errReadingFile
	}
	defer fs.Close()

	// create buffered reader, keep track of bytes read in current read
	buff := bufio.NewScanner(fs)
	byteCounter := scanByteCounter{}
	buff.Split(byteCounter.wrap(bufio.ScanLines))

	// extract header values (TSV)
	headerRaw := []string{}
	for i := 0; i <= 7; i++ {
		buff.Scan()
		// extract header fields
		headerRaw = append(headerRaw, buff.Text())
	}

	// generate header instance
	var h *header
	// if there's no separator, file is using JSON format
	if strings.Contains(headerRaw[0], "#separator") {
		// get seperator character
		sep := strings.Split(headerRaw[0], " ")[1]
		delimeter, err := strconv.ParseInt(sep[2:], 16, 64)
		if err == nil {
			h = &header{}
			h.separator = string(rune(delimeter))
			h.setSeperator = strings.Split(headerRaw[1], h.separator)[1]
			h.emptyField = strings.Split(headerRaw[2], h.separator)[1]
			h.unsetField = strings.Split(headerRaw[3], h.separator)[1]
			// ignore first column of fields and types
			h.fields = strings.Split(headerRaw[6], h.separator)[1:]
			h.types = strings.Split(headerRaw[7], h.separator)[1:]
		}
	}

	// reset buffer
	fs.Seek(0, 0)
	buff = bufio.NewScanner(fs)
	byteCounter = scanByteCounter{}
	buff.Split(byteCounter.wrap(bufio.ScanLines))

	// skip already read lines
	position := int64(0)
	for position < f.Lines {
		buff.Scan()
		position++
	}

	// append connection chunks
	chunks := [][]byte{}

	// read lines using specified chunk size
	newBytes := int64(0)
	for i := 0; i < s.FileChunkSize; i++ {
		// advance file pointer, stop of end of file
		newLine := buff.Scan()
		if !newLine {
			// End of file
			break
		}
		// parse the line
		line := buff.Text()
		// don't parse lines starting with # (header/comment)
		if line != "" && line[0:1] != "#" {
			payload, err := parseLine(line, h)
			// no error parsing, append to chunks
			if err == nil {
				chunks = append(chunks, payload)
			} else {
				// print error message
				log.Println(err)
			}
		}
		// valid update count of lines read and bytes read
		newBytes = byteCounter.BytesRead
		f.Lines++
	}

	// update total number of bytes read per chunk
	f.Size = newBytes

	// generate actual frame
	frame := &protocol.UploadRequest{
		Header: &protocol.Header{
			MsgUuid:      uuid.New().String(),
			MsgTimestamp: timestamppb.Now(),
			Status:       protocol.Status_REQUEST,
			ErrorMsg:     "",
			Session:      s.Session,
		},
		AssetId:  s.AssetID,
		FileName: baseName,
		Payload:  chunks,
	}

	return frame, nil
}
