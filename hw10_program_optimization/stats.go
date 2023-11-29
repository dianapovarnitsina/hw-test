package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	re := regexp.MustCompile("\\." + domain)
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	p := fastjson.Parser{}

	for scanner.Scan() {
		line := scanner.Bytes()

		v, err := p.ParseBytes(line)
		if err != nil {
			return nil, err
		}

		email := string(v.GetStringBytes("Email"))
		domain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
		matched := re.MatchString(email)
		if matched {
			num := result[domain]
			num++
			result[domain] = num
		}
	}
	return result, nil
}
