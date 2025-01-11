package models

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type Float64 float64

// UnmarshalJSON read a JSON string into a custom Float64 type
func (f *Float64) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	parsed, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.Wrap(err, "invalid float value: %v")
	}

	*f = Float64(parsed)

	return nil
}

// Float64 returns the float64 value of the custom Float64 type
func (f Float64) Float64() float64 {
	return float64(f)
}
