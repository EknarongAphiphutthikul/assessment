package common

import (
	"strconv"
	"strings"
)

type Error struct {
	Code          int
	Desc          string
	OriginalError error
}

func (e Error) Error() string {
	if e.Code == 0 {
		return e.Desc
	}
	return strings.Join([]string{strconv.Itoa(e.Code), e.Desc}, ":")
}
