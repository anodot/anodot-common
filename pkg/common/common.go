package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AnodotTimestamp struct {
	time.Time
}

func (t AnodotTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(t.Unix())), nil
}

func (t AnodotTimestamp) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)

	i, err := strconv.ParseInt(strInput, 10, 64)
	if err != nil {
		panic(err)
	}

	t.Time = time.Unix(i, 0)
	return nil
}
