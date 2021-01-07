package jar

import (
	"net/http"
	"time"
)

// 기본 쿠키 관리자 인터페이스
type Jar interface {
	Set(cookie *http.Cookie)
	Remove(cookieName string, path string)
	Get(cookieName string) *http.Cookie
	Write()
}

// 쿠키 관리자 구현
type jar struct {
	request        map[string]*http.Cookie
	response       map[string]*http.Cookie
	responseWriter http.ResponseWriter
}

func (j *jar) Set(cookie *http.Cookie) {
	j.response[cookie.Name] = cookie
}

func (j *jar) Remove(cookieName string, path string) {
	if oldCookie, b := j.request[cookieName]; b {
		n := *oldCookie
		if n.MaxAge != 0 || len(n.RawExpires) > 0 {
			n.MaxAge = -1
			n.Expires = time.Now().Add(-time.Hour)
		}
		n.Value = ""
		n.Path = path
		j.response[cookieName] = &n
	} else {
		delete(j.response, cookieName)
	}
}

func (j *jar) Get(cookieName string) *http.Cookie {
	if ck, b := j.response[cookieName]; b {
		return ck
	}
	if ck, b := j.request[cookieName]; b {
		return ck
	}
	return nil
}

// response에 쓰기
func (j *jar) Write() {
	// 일단 지우고, 다른 데서 쿠키는 못 쓰게 만든다
	header := j.responseWriter.Header()
	header.Del("Set-Cookie")
	for _, ck := range j.response {
		if v := ck.String(); v != "" {
			header.Add("Set-Cookie", v)
		}
	}
}

func New(request *http.Request, response http.ResponseWriter) Jar {
	jar := jar{
		request:        make(map[string]*http.Cookie),
		response:       make(map[string]*http.Cookie),
		responseWriter: response,
	}
	for _, c := range request.Cookies() {
		jar.request[c.Name] = c
	}
	return &jar
}
