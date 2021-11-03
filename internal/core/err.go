package core

import "fmt"

// not real error
type ReturnErr struct {
	value interface{}
}

func (r ReturnErr) Error() string {
	return fmt.Sprintf("%+v", r.value)
}
