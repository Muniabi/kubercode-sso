package values

import (
	"encoding/json"
	"fmt"
)

type IsMentor struct {
	IsMentor bool `json:"is_mentor"`
}

func NewIsMentor(isMentor bool) (*IsMentor, error) {
	return &IsMentor{
		IsMentor: isMentor,
	}, nil
}

func (i *IsMentor) UnmarshalJSON(data []byte) error {
	var temp struct {
		IsMentor bool `json:"is_mentor"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.IsMentor = temp.IsMentor
	return nil
}

func (i *IsMentor) ToString() string {
	return fmt.Sprintf("%t", i.IsMentor)
}

func (i *IsMentor) GetIsMentor() bool {
	return i.IsMentor
}
