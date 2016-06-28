package chatterbot

import (
	"net/http"
	"fmt"
	"crypto/md5"
	"log"
	"net/http/cookiejar"
	"bytes"
	"net/url"
	"strings"
	"strconv"
)

const (
	cleverbotBaseUrl = "http://www.cleverbot.com"
	cleverbotServiceUrl = "http://www.cleverbot.com/webservicemin?uc=165"
	cleverbotEndIndex = 26
)

type CleverbotSession struct {
	Conversation	string
	Cookies		*cookiejar.Jar
	HttpClient	*http.Client
}

func NewCleverbot() *CleverbotSession {

	session := &CleverbotSession{}
	session.Cookies, _ = cookiejar.New(nil)
	session.HttpClient = &http.Client{
		Jar: session.Cookies,
	}

	// initial request (get cookies for session)
	session.Request(cleverbotBaseUrl, "")

	return session
}

func (session *CleverbotSession) ThinkThrough(though string) string {

	// our digest will always contain this string, but cut at the cleverbotEndIndex
	digestString := url.QueryEscape(though) + "&islearning=1&icognoid=wsf";
	digest := digestString[0:cleverbotEndIndex]
	icognocheck := fmt.Sprintf("%x", md5.Sum([]byte(digest)))

	// our default params; always the same for the first request
	params := "stimulus=" + url.QueryEscape(though) +
	"&islearning=1" +
	"&icognoid=wsf" +
	"&icognocheck=" + url.QueryEscape(icognocheck)

	// if there is already a conversation going, we need to add the vTexts, even empty ones
	if len(session.Conversation) > 0 {

		vTexts := strings.Split(session.Conversation, ",")
		if len(vTexts) > 0 {

			params += "logurl=&"
			for index, vText := range vTexts {

				// we need to be the last char a dot, if not something else
				lastChar := string(vText[len(vText)-1])
				if strings.ContainsAny(lastChar, "! ?") {
					vText += "."
				}
				params += "vText" + strconv.Itoa(index + 2) + "=" + url.QueryEscape(vText) + "&"
			}
			for i := len(vTexts)-1; i <= 6; i++ {
				params += "vText" + strconv.Itoa(i + 2) + "=&"
			}
			params += "prevref="
		}
	}

	// keep our whole conversation comma-separated for vTexts
	if session.Conversation != "" {
		session.Conversation += ","
	}
	session.Conversation += though

	response, err := session.Request(cleverbotServiceUrl, params)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	if response != "" {
		session.Conversation += "," + response
		return response
	}

	fmt.Println("missing response")
	return ""
}

func (session *CleverbotSession) Request(uri string, params string) (string, error) {

	var request *http.Request

	if len(params) == 0 {
		request, _ = http.NewRequest(http.MethodGet, uri, nil)
	} else {
		request, _ = http.NewRequest(http.MethodPost, uri, bytes.NewBufferString(params))
		request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	}

	// set language statically, because it doesn't seem to change anything
	request.Header.Add("Accept-Language", "en;q=1.0")

	session.HttpClient.Jar = session.Cookies
	resp, err := session.HttpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	response := resp.Header.Get("Cboutput")
	if response == "" {
		return "", nil
	}

	decodedResponse, err := url.QueryUnescape(response)
	if err != nil {
		return "", err
	}

	return decodedResponse, nil
}