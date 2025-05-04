package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type GObj struct {
}

// Scan Scanner
func (o *GObj) Scan(value any, obj any) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}
	return json.Unmarshal(b, obj)
}

// Value Valuer
func (o GObj) Value(obj any) (driver.Value, error) {
	if obj == nil {
		return nil, nil
	}
	return json.Marshal(obj)
}

func (o *GObj) GormDataType() string {
	return "obj"
}
