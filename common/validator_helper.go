package common

import (
	"strings"
	forks "todo-list/forks"
)

var commonDomains = []string{"gmail.com", "yahoo.com", "hotmail.com", "outlook.com"}

func IsTypoDomain(email string) (bool, string) {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false, ""
    }

    domain := parts[1]
    for _, d := range commonDomains {
        distance := forks.ComputeDistance(domain, d)
		
        if distance > 0 && distance <= 2 { // cho phép sai tối đa 2 ký tự
            return true, d
        }
    }

    return false, ""
}