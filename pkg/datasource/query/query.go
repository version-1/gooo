package query

import (
	"errors"
	"fmt"
	"strings"
)

func Select(table string, fields []string, where *string) string {
	cond := "1 = 1"
	if where != nil {
		cond = *where
	}

	return fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(fields, ","),
		table,
		cond,
	)
}

func Insert(table string, fields []string, returnings *[]string) string {
	rtns := []string{"id", "created_at", "updated_at"}
	if returnings != nil {
		rtns = *returnings
	}

	placeholders := BuildPlaceholders(len(fields))

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		table,
		strings.Join(fields, ","),
		placeholders,
		strings.Join(rtns, ","),
	)
}

func Update(table string, fields []string, condition string) (string, error) {
	if condition == "" {
		return "", errors.New("where clause is required")
	}

	set := ""
	for i, f := range fields {
		if i == len(fields)-1 {
			set += fmt.Sprintf("%s = $%d", f, i+1)
		} else {
			set += fmt.Sprintf("%s = $%d, ", f, i+1)
		}
	}

	return fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(fields, ","),
		condition,
	), nil
}

func BuildPlaceholders(n int) string {
	list := []string{}
	for i := 1; i <= n; i++ {
		list = append(list, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf("%s", strings.Join(list, ", "))
}
