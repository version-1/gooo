package helper

import (
	"fmt"
	"strings"
)

type ConnInfo struct {
	Value    string
	Server   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func ParseConnstr(s string) (*ConnInfo, error) {
	server := ""
	sections := strings.Split(s, "://")
	if len(sections) < 2 {
		return nil, fmt.Errorf("invalid connection string")
	} else {
		server = sections[0]
	}

	username := make([]int, 2)
	password := make([]int, 2)
	host := make([]int, 2)
	port := make([]int, 2)
	database := make([]int, 2)

	target := sections[1]

	for i, char := range target {
		if username[1] == 0 && string(char) == ":" {
			username[1] = i
			password[0] = i + 1
		}

		if string(char) == "@" {
			password[1] = i
			host[0] = i + 1
		}

		if host[0] != 0 && string(char) == ":" {
			host[1] = i
			port[0] = i + 1
		}

		if username[1] != 0 && string(char) == "/" {
			if port[0] != 0 {
				port[1] = i
			}
			if host[1] == 0 {
				host[1] = i
			}

			database[0] = i + 1
		}

		if database[0] != 0 && string(char) == "?" {
			database[1] = i
		}
	}

	if database[1] == 0 {
		database[1] = len(target)
	}
	fmt.Printf("host: %#v port: %#v\n", host, port)

	return &ConnInfo{
		Value:    s,
		Server:   server,
		Host:     substr(target, host[0], host[1]),
		Port:     substr(target, port[0], port[1]),
		Username: substr(target, username[0], username[1]),
		Password: substr(target, password[0], password[1]),
		Database: substr(target, database[0], database[1]),
	}, nil
}

func substr(str string, a, b int) string {
	if a < 0 || len(str) < b || a > b {
		return ""
	}

	return str[a:b]
}
