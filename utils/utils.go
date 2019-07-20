package utils

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"math/rand"
	"strconv"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ParseEmailDomain(email string) (string, error) {
	ec := strings.Split(email, "@")

	if len(ec) == 1 {
		return "", errors.New("wrong email")
	}

	return ec[1], nil
}

func GetPhone(country uint16, city uint16, phone uint32) int64 {
	return int64(country)*100000000000 + int64(city)*100000000 + int64(phone)
}

// TODO переделать на builder!!
func PhoneToString(phone int64) string {
	p := strconv.Itoa(int(phone))
	return fmt.Sprintf("%s(%s)%s", p[:1], p[1:4], p[4:])
}

type argKV struct {
	Key   []byte
	Value []byte
}

func GetURIArgs(ctx *fasthttp.RequestCtx) []argKV {
	args := make([]argKV, 0, ctx.Request.URI().QueryArgs().Len())

	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		args = append(args, argKV{Key: key, Value: value})
	})

	return args
}
