package types

import (
	"bytes"
	"fmt"
	"strconv"
)

type Float64 float64

//goland:noinspection GoMixedReceiverTypes
func (f Float64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", f)), nil
}

//goland:noinspection GoMixedReceiverTypes
func (f *Float64) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	if len(data) == 0 || string(data) == "null" {
		*f = 0
		return nil
	}

	val, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return fmt.Errorf("invalid float value %q: %v", data, err)
	}
	*f = Float64(val)
	return nil
}
