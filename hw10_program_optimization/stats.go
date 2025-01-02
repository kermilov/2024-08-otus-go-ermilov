package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	// ID       int
	// Name     string
	// Username string
	Email string
	// Phone    string
	// Password string
	// Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (users, error) {
	scanner := bufio.NewScanner(r)
	var result users
	var user User
	for scanner.Scan() {
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.HasSuffix(user.Email, "."+domain)

		if matched {
			x := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[x]
			num++
			result[x] = num
		}
	}
	return result, nil
}
