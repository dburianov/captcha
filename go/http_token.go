package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type tokenResponse struct {
	ReturnTokenV string `json:"token"`
	ReturnHashV  string `json:"hash"`
	ReturnTtlV   int64  `json:"ttl"`
}

func tokenCheck(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	hash := mux.Vars(r)["hash"]
	clientIp := mux.Vars(r)["clientip"]
	ttl := mux.Vars(r)["ttl"]

	if len(token) < 1 || len(hash) < 1 || len(clientIp) < 7 || len(ttl) < 1 {
		http.Error(w, "No token or hash or clientip or ttl", http.StatusForbidden)
	} else {
		validate := checkTokenFunc(token, hash, clientIp, ttl)
		if validate == "Ok" {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "do not validate", http.StatusForbidden)
		}
	}
}

func checkTokenFunc(t string, h string, ip string, ttl string) string {
	soldHash := sha256.New()
	encString := t + seviceConfig.Salt + ip + ttl
	if *loggingLevel == "debug" {
		log.Println("encString: ", encString)
	}

	for i := 0; i < seviceConfig.Rand.RandCount; i++ {
		soldHash.Write([]byte(encString))
	}

	timeNow := time.Now().Add(seviceConfig.AuthExpire * time.Second).Unix()
	ttlInt64, _ := strconv.ParseInt(ttl, 10, 64)
	if *loggingLevel == "debug" {
		log.Println("ttlInt64: ", ttlInt64)
		log.Println("timeNow : ", timeNow)
	}

	if string(hex.EncodeToString(soldHash.Sum(nil))) == h && ttlInt64 <= timeNow {
		if *loggingLevel == "debug" {
			log.Println("Ok")
		}
		return "Ok"
	} else {
		log.Println("Nah")
		return "ERR"
	}
}

func newToken(w http.ResponseWriter, r *http.Request) {
	clientIp := mux.Vars(r)["clientip"]
	if len(clientIp) < 7 {
		http.Error(w, "No clientIp", http.StatusForbidden)
	} else {
		var str string
		for i := 0; i < seviceConfig.Rand.RandCount; i++ {
			rndString := rndGenerateString(
				seviceConfig.Rand.RandLenght,
				seviceConfig.Rand.CharacterUp,
				seviceConfig.Rand.CharacterDown,
				seviceConfig.Rand.CharacterInt,
				seviceConfig.Rand.CharacterSpec)
			str += rndString
		}

		ttl := time.Now().Unix()

		token := sha256.New()
		token.Write([]byte(str))
		returnToken := hex.EncodeToString(token.Sum(nil))

		hash := sha256.New()
		encStr := returnToken + seviceConfig.Salt + clientIp + strconv.FormatInt(ttl, 10)
		if *loggingLevel == "debug" {
			log.Println("encStr: ", encStr)
		}
		for i := 0; i < seviceConfig.Rand.RandCount; i++ {
			hash.Write([]byte(encStr))
		}
		returnHash := hex.EncodeToString(hash.Sum(nil))

		returnVal := tokenResponse{
			ReturnTokenV: returnToken,
			ReturnHashV:  returnHash,
			ReturnTtlV:   ttl,
		}
		jsonReturnVal, err := json.Marshal(returnVal)
		if err != nil {
			log.Println("error:", err)
		}

		str = ""
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonReturnVal))

		if *loggingLevel == "debug" {
			log.Println("prejson: ", returnVal)
			log.Print("json: ", string(jsonReturnVal))
		}
	}
}
