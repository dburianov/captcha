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
	charsetDownCount := len(charsetDown) * csDown
	charsetUpCount := len(charsetUp) * csUp
	charsetIntCount := len(charsetInt) * csInt
	charsetSpecCount := len(charsetSpec) * csSpec
	srcStringLen := charsetDownCount + charsetUpCount + charsetIntCount + charsetSpecCount
	srcStringDown := make([]byte, charsetDownCount)
	srcStringUp := make([]byte, charsetUpCount)
	srcStringInt := make([]byte, charsetIntCount)
	srcStringSpec := make([]byte, charsetSpecCount)

	if *loggingLevel == "debug" {
		log.Println("srcStringLen", srcStringLen)
	}

	for i := 0; i < charsetUpCount; i++ {
		srcStringUp[i] = charsetUp[seededRand.Intn(len(charsetUp))]
	}
	stringTotal += string(srcStringUp)

	for i := 0; i < charsetDownCount; i++ {
		srcStringDown[i] = charsetDown[seededRand.Intn(len(charsetDown))]
	}
	stringTotal += string(srcStringDown)

	for i := 0; i < charsetIntCount; i++ {
		srcStringInt[i] = charsetInt[seededRand.Intn(len(charsetInt))]
	}
	stringTotal += string(srcStringInt)

	for i := 0; i < charsetSpecCount; i++ {
		srcStringSpec[i] = charsetSpec[seededRand.Intn(len(charsetSpec))]
	}
	stringTotal += string(srcStringSpec)

	if *loggingLevel == "debug" {
		log.Println("string(srcString):", stringTotal)
		log.Println(len(stringTotal))
	}
	return StringWithCharset(length, stringTotal)
}

func cryptoValue(s string, salt string) string {
	hash := sha256.New()
	encStr := salt + s + salt
	hash.Write([]byte(encStr))
	return hex.EncodeToString(hash.Sum(nil))
}
