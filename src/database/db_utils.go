package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Uint64Slice []uint64

func (slice Uint64Slice) Value() (driver.Value, error) {
	jsonSlice, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}
	return string(jsonSlice), nil
}

func (slice *Uint64Slice) Scan(src interface{}) (err error) {
	var s []uint64

	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &s)
	case []byte:
		err = json.Unmarshal(src.([]byte), &s)
	default:
		err = fmt.Errorf("invalid type for []uint64: %T", src)
	}

	if err != nil {
		return err
	}

	*slice = s
	return nil
}
