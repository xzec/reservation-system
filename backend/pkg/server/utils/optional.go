package utils

import (
	"encoding/json"
	"fmt"
)

type Optional[T any] struct {
	Defined bool
	Value   *T
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Defined = true
	return json.Unmarshal(data, &o.Value)
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	fmt.Println("o", o)
	if o.Defined {
		return json.Marshal(o.Value)
	}
	return []byte("null"), nil
}

func (o Optional[T]) String() string {
	if o.Defined == true {
		return fmt.Sprintf("{true %v}", o.Value)
	} else {
		return fmt.Sprintf("{false %v (undefined)}", o.Value)
	}
}
