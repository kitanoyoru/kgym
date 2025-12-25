package file

import "github.com/pkg/errors"

type State string

const (
	StatePending   State = "pending"
	StateCompleted State = "completed"
	StateFailed    State = "failed"
)

func StateFromString(s string) (State, error) {
	switch s {
	case "pending":
		return StatePending, nil
	case "completed":
		return StateCompleted, nil
	case "failed":
		return StateFailed, nil
	}

	return "", errors.New("invalid state")
}

func (s State) String() string {
	return string(s)
}
