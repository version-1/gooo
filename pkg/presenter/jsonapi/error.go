package jsonapi

import (
	"fmt"
	"strings"
)

type Errors []Error

func (j Errors) JSONAPISerialize() (string, error) {
	str := "["
	for _, e := range j {
		json, err := e.JSONAPISerialize()
		if err != nil {
			return "", err
		}
		str += json + ","
	}
	str += "]"

	return str, nil
}

type Error struct {
	ID     string
	Status int
	Code   string
	Title  string
	Detail string
}

func (j Error) JSONAPISerialize() (string, error) {
	fields := []string{
		fmt.Sprintf("\"id\": %s", Stringify(j.ID)),
		fmt.Sprintf("\"status\": %s", Stringify(j.Status)),
		fmt.Sprintf("\"code\": %s", Stringify(j.Code)),
		fmt.Sprintf("\"title\": %s", Stringify(j.Title)),
		fmt.Sprintf("\"detail\": %s", Stringify(j.Detail)),
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, ", \n")), nil
}

type ErrorCompatible interface {
	ToJSONAPIError() Error
}
