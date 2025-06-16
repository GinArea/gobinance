package binapi

import (
	"fmt"
	"slices"

	"github.com/msw-x/moon/ujson"
)

type Error struct {
	Code ujson.Int64
	Text string
}

func (o *Error) Std() error {
	if o.Empty() {
		return nil
	} else {
		return o
	}
}

func (o *Error) Empty() bool {
	return o.Code == 0
}

func (o *Error) Error() string {
	return fmt.Sprintf("code[%d]: %s", o.Code.Value(), o.Text)
}

func (o *Error) ApiKeyInvalid() bool {
	codes := []ujson.Int64{
		-1, //TODO
	}
	return slices.Contains(codes, o.Code)
}
