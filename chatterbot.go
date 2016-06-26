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
	Headers		map[string]string
	Vars		map[string]string
	Conversation	string
	Cookies		*cookiejar.Jar
}

func NewCleverbot(lang string) *CleverbotSession {

	if len(lang) < 2 {
		lang = "en"
	}

	session := &CleverbotSession{}

	session.Headers = make(map[string]string)
	session.Headers["Accept-Language"] = lang + ";q=1.0"

	session.Vars = make(map[string]string)
	session.Vars["islearning"] = "1";
	session.Vars["icognoid"] = "wsf";

	session.Cookies, _ = cookiejar.New(nil)

	Request(cleverbotBaseUrl, session.Cookies, "", nil)

	return session
}

func (session *CleverbotSession) ThinkThrough(though string) string {

	session.Vars["stimulus"] = though

	// our digest will always contain of this string, but cut at the cleverbotEndIndex
	digestString := url.QueryEscape(though) + "&islearning=" + session.Vars["islearning"] + "&icognoid=";
	digest := digestString[0:cleverbotEndIndex]
	session.Vars["icognocheck"] = fmt.Sprintf("%x", md5.Sum([]byte(digest)))

	// our default params; always the same for the first request
	params := "stimulus=" + url.QueryEscape(session.Vars["stimulus"]) +
	"&islearning=" + url.QueryEscape(session.Vars["islearning"]) +
	"&icognoid=" + url.QueryEscape(session.Vars["icognoid"]) +
	"&icognocheck=" + url.QueryEscape(session.Vars["icognocheck"])

	// if there is already a conversation going, we need to add the vTexts, even if empty
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

	response, err := Request(cleverbotServiceUrl, session.Cookies, params, session.Headers)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	session.Conversation += "," + response
	return response
}

func Request(uri string, cookies *cookiejar.Jar, params string, headers map[string]string) (string, error) {

	var request *http.Request

	if len(params) == 0 {
		request, _ = http.NewRequest(http.MethodGet, uri, nil)
	} else {
		request, _ = http.NewRequest(http.MethodPost, uri, bytes.NewBufferString(params))
		request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	}

	if len(headers) > 0 {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}

	httpClient := &http.Client{
		Jar: cookies,
	}
	resp, err := httpClient.Do(request)
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