package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/steambap/captcha"
)

func captchaIndexHandle(w http.ResponseWriter, _ *http.Request) {
	doc, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	doc.Execute(w, nil)
}

func captchaTextHandle(w http.ResponseWriter, _ *http.Request) {
	img, err := captcha.New(seviceConfig.Captcha.Text.Width, seviceConfig.Captcha.Text.Hight, func(options *captcha.Options) {
		options.TextLength = seviceConfig.Captcha.TextLength
		options.CharPreset = seviceConfig.Captcha.CharPreset
	})
	if err != nil {
		fmt.Fprint(w, nil)
		fmt.Println(err.Error())
		return
	}

	expiration := time.Now().Add(seviceConfig.Captcha.Cookie.Expire * time.Minute)
	cookie := http.Cookie{Name: "captchatoken", Value: cryptoValue(img.Text, seviceConfig.Captcha.Salt), Path: "/", Secure: true, Expires: expiration}
	http.SetCookie(w, &cookie)
	img.WriteImage(w)
	if *loggingLevel == "debug" {
		log.Println(img.Text)
	}
}

func captchaMathHandle(w http.ResponseWriter, _ *http.Request) {
	img, err := captcha.NewMathExpr(seviceConfig.Captcha.Math.Width, seviceConfig.Captcha.Math.Hight)
	if err != nil {
		fmt.Fprint(w, nil)
		fmt.Println(err.Error())
		return
	}
	expiration := time.Now().Add(seviceConfig.Captcha.Cookie.Expire * time.Second)
	cookie := http.Cookie{Name: "captchamath", Value: cryptoValue(img.Text, seviceConfig.Captcha.Salt), Path: "/", Secure: true, Expires: expiration}
	http.SetCookie(w, &cookie)
	img.WriteImage(w)
}

func captchaCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	tokenCookie, err := r.Cookie("captchatoken")
	if err != nil {
		log.Fatalf("Error occured while reading cookie")
	}
	if *loggingLevel == "debug" {
		log.Println(tokenCookie)
	}

	tokenMath, err := r.Cookie("captchamath")
	if err != nil {
		log.Fatalf("Error occured while reading cookie")
	}
	captcha := r.FormValue("captcha")
	math := r.FormValue("math")

	ips, err := r.Cookie("cookieIP")
	if err != nil {
		log.Fatalf("Error occured while reading cookie")
	}
	if *loggingLevel == "debug" {
		log.Println(tokenCookie)
	}
	if cryptoValue(captcha, seviceConfig.Captcha.Salt) == tokenCookie.Value && cryptoValue(math, seviceConfig.Captcha.Salt) == tokenMath.Value {

		apiAddress := "http://127.0.0.1:" + strconv.FormatUint(uint64(seviceConfig.Bind.Port), 10) + "/api/v0/token/new/" + ips.Value

		req, err := http.NewRequest(http.MethodGet, apiAddress, bytes.NewReader([]byte("")))
		if err != nil {
			log.Println(err)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		bodyReaderResponse := bytes.NewReader(body)

		var jsontokenResponse tokenResponse
		jsonParser := json.NewDecoder(bodyReaderResponse)
		jsonParser.Decode(&jsontokenResponse)

		if *loggingLevel == "debug" {
			log.Println("jsontokenResponse:", jsontokenResponse)
			log.Println("str ttl: ", jsontokenResponse.ReturnTtlV)
		}

		expiration := time.Now().Add(seviceConfig.AuthExpire * time.Second)

		cookieToken := http.Cookie{Name: "cookieToken", Value: jsontokenResponse.ReturnTokenV, Path: "/", Secure: true, Expires: expiration, HttpOnly: true}
		http.SetCookie(w, &cookieToken)

		cookieHash := http.Cookie{Name: "cookieHash", Value: jsontokenResponse.ReturnHashV, Path: "/", Secure: true, Expires: expiration, HttpOnly: true}
		http.SetCookie(w, &cookieHash)

		cookieTtl := http.Cookie{Name: "cookieTtl", Value: strconv.FormatInt(int64(jsontokenResponse.ReturnTtlV), 10), Path: "/", Secure: true, Expires: expiration, HttpOnly: true}
		http.SetCookie(w, &cookieTtl)

		xRedirect, err := r.Cookie("cookieUrl")
		if err != nil {
			log.Fatalf("Error occured while reading cookie")
		}
		if *loggingLevel == "debug" {
			log.Println(xRedirect.Value)
		}
		http.Redirect(w, r, xRedirect.Value, 302)
	}
}
