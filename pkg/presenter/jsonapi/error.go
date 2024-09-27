package jsonapi

import (
	"encoding/json"
	"fmt"
	"strings"

	goooerrors "github.com/version-1/gooo/pkg/errors"
)

type Errors []Error

func (j Errors) Error() string {
	list := []string{}
	for _, e := range j {
		list = append(list, e.Error())
	}

	return fmt.Sprintf("[%s]", strings.Join(list, ", "))
}

func (j Errors) JSONAPISerialize() (string, error) {
	str := "["
	for i, e := range j {
		json, err := e.JSONAPISerialize()
		if err != nil {
			return "", err
		}

		comma := ""
		if i != len(j)-1 {
			comma = ","
		}
		str += json + comma
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

func (j Error) Error() string {
	return j.Detail
}

func (j Error) JSONAPISerialize() (string, error) {
	id, err := escape(j.ID)
	if err != nil {
		return "", err
	}

	status, err := escape(j.Status)
	if err != nil {
		return "", err
	}

	code, err := escape(j.Code)
	if err != nil {
		return "", err
	}

	title, err := escape(j.Title)
	if err != nil {
		return "", err
	}

	detail, err := escape(j.Detail)
	if err != nil {
		return "", err
	}

	fields := []string{
		fmt.Sprintf("\"id\": %s", id),
		fmt.Sprintf("\"status\": %s", status),
		fmt.Sprintf("\"code\": %s", code),
		fmt.Sprintf("\"title\": %s", title),
		fmt.Sprintf("\"detail\": %s", detail),
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, ", \n")), nil
}

func escape(i any) (string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return "", goooerrors.Wrap(err)
	}

	return string(b), nil
}

type Errable interface {
	ToJSONAPIError() Error
	Error() string
}
