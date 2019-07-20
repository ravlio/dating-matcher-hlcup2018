package dicts

import "errors"

// лежит тут, так как я устал бороться с циклическими зависимостями
type Sex uint8

const (
	SexMale   Sex = 1
	SexFemale Sex = 2
)

type Status uint8

const (
	StatusFree        Status = 1
	StatusOccupied    Status = 2
	StatusComplicated Status = 3
)

func StringToSex(k string) (Sex, error) {
	switch k {
	case "m":
		return SexMale, nil
	case "f":
		return SexFemale, nil
	default:
		return 0, errors.New("unknown sex")
	}
}

func SexToString(i Sex) (string, error) {
	switch i {
	case SexMale:
		return "m", nil
	case SexFemale:
		return "f", nil
	default:
		return "", errors.New("unknown sex")
	}
}

func StringToStatus(k string) (Status, error) {
	switch k {
	case `свободны`:
		return StatusFree, nil
	case `заняты`:
		return StatusOccupied, nil
	case `всё сложно`:
		return StatusComplicated, nil
	default:
		return 0, errors.New("unknown status")
	}
}

func StatusToString(i Status) (string, error) {
	switch i {
	case StatusFree:
		return `свободны`, nil
	case StatusOccupied:
		return `заняты`, nil
	case StatusComplicated:
		return `всё сложно`, nil
	default:
		return "", errors.New("unknown status")
	}
}
