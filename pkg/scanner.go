package pkg

import (
	"bufio"
	"log"
	"math/rand"
	"time"
)

var (
	params      = make([]string, 0, 0)
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init(input ctx.Context) {
	rand.Seed(time.Now().UnixNano())

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatal(scanner.Err().Error())
		}
		params = append(params, scanner.Text())
	}
}

func GetParams() []string {
	return params
}

func RandStringRunes() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
