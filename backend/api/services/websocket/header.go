package websocket

import "errors"

func Validate(header *Header) error {
	if header == nil {
		return errors.New("invalid header provided")
	}
	if header.MsgUuid == "" {
		return errors.New("invalid header UUID")
	}
	if header.MsgTimestamp.IsZero() {
		return errors.New("invalid header timestamp")
	}
	return nil
}
