package config

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// time.Duration forces you to specify in millis, and does not support days
// see https://stackoverflow.com/questions/48050945/how-to-unmarshal-json-into-durations
type TTL time.Duration

func (l TTL) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(l).String())
}

func (l *TTL) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		if value == "" {
			*l = 0
			return nil
		}
		if strings.HasSuffix(value, "d") {
			hours, err := strconv.Atoi(strings.TrimSuffix(value, "d"))
			*l = TTL(time.Duration(hours) * 24 * time.Hour)
			return err
		}
		d, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*l = TTL(d)
		return nil
	default:
		return errors.New("invalid TTL")
	}
}
