package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	uuidFormat = "%s-%s"
)

func Md5(data string) string {
	d := []byte(data)
	m := md5.New()
	m.Write(d)

	return hex.EncodeToString(m.Sum(nil))
}

func SpanId() string {
	now := time.Now().UnixNano()
	rand.Seed(now)
	random := strconv.FormatInt(rand.Int63n(999999999), 10)
	return Md5(fmt.Sprintf(uuidFormat, random, strconv.FormatInt(now, 10)))
}
