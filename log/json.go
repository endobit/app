package log

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

type niceJSON struct {
	Writer io.Writer
}

type message struct {
	Time  string         `json:"time,omitempty"`
	Level string         `json:"level,omitempty"`
	Msg   string         `json:"msg,omitempty"`
	Extra map[string]any `json:"extra,omitempty"`
}

// Write implements the io.Writer interface for n. It accepts a complete JSON
// strings, parses it, formats it nicely, and writes it's underlying io.Writer.
func (n *niceJSON) Write(buff []byte) (int, error) {
	var raw map[string]any

	buf := bytes.NewBuffer(buff)

	if err := json.NewDecoder(buf).Decode(&raw); err != nil {
		return n.Writer.Write(buff) // not a valid JSON, just write the raw message
	}

	var msg message

	if value, ok := raw["time"].(string); ok {
		t, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return n.Writer.Write(buff) // not a valid time, just write the raw message
		}

		msg.Time = t.Format(time.Kitchen)

		delete(raw, "time")
	}

	if l, ok := raw["level"].(string); ok {
		msg.Level = l

		delete(raw, "level")
	}

	if m, ok := raw["msg"].(string); ok {
		msg.Msg = m

		delete(raw, "msg")
	}

	msg.Extra = raw

	e := json.NewEncoder(n.Writer)
	e.SetIndent("", "    ")

	if err := e.Encode(msg); err != nil {
		return n.Writer.Write(buff) // if we cannot encode the message, just write the raw message
	}

	return len(buff), nil
}
