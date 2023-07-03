package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/rand"
	"time"
)

const charsetDown = "abcdefghijklmnopqrstuvwxyz"
const charsetUp = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charsetInt = "0123456789"
const charsetSpec = " !@#$%^&*()_+=-\\|/?.>,<\"';:`~"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano() * time.Now().UnixMilli()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func rndGenerateString(
	length int,
	csUp int,
	csDown int,
	csInt int,
	csSpec int) string {

	var stringTotal string
	srcString := make([]byte, len(charsetDown)*csDown+len(charsetUp)*csUp+len(charsetInt)*csInt+len(charsetSpec)*csSpec)

	for i := 0; i < csUp; i++ {
		srcString[i] = charsetUp[seededRand.Intn(len(charsetUp))]
	}
	stringTotal += string(srcString)
	for i := 0; i < csDown; i++ {
		srcString[i] = charsetDown[seededRand.Intn(len(charsetDown))]
	}
	stringTotal += string(srcString)
	for i := 0; i < csInt; i++ {
		srcString[i] = charsetInt[seededRand.Intn(len(charsetInt))]
	}
	stringTotal += string(srcString)
	for i := 0; i < csSpec; i++ {
		srcString[i] = charsetSpec[seededRand.Intn(len(charsetSpec))]
	}
	stringTotal += string(srcString)

	if *loggingLevel == "debug" {
		log.Println("string(srcString):", stringTotal)
	}
	return StringWithCharset(length, stringTotal)
}

func cryptoValue(s string, salt string) string {
	hash := sha256.New()
	encStr := salt + s + salt
	hash.Write([]byte(encStr))
	return hex.EncodeToString(hash.Sum(nil))
}
