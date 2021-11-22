package types

import (
	"encoding/json"
	"fmt"
)

type Input map[string]interface{}

func (i Input) ToString() (string, error) {
	bb, err := json.Marshal(i)
	if err != nil {
		fmt.Println(fmt.Errorf("json marshal input err: %v", err))
		return "", err
	}
	return string(bb), nil
}
