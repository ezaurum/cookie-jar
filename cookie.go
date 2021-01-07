package jar

import (
	"net/http"
	"time"
)

// Jar 기본 쿠키 관리자 인터페이스
type Jar interface {
	Set(cookie *http.Cookie)
	Remove(cookieName string, path string)
	Get(cookieName string) *http.Cookie
	Extend(cookieName string, duration time.Duration) *http.Cookie
	Write()
}

// 쿠키 관리자 구현
type jar struct {
	request        map[string]*http.Cookie
	response       map[string]*http.Cookie
	responseWriter http.ResponseWriter
}

func (j *jar) Extend(cookieName string, duration time.Duration) {
	if get := j.Get(cookieName); nil == get {
		return
	} else {
		get.Expires.Add(duration)
	}
}

func (j *jar) Set(cookie *http.Cookie) {
	j.response[cookie.Name] = cookie
}

// Remove net/http 의 쿠키는 name, value 만 보내주므로 path 가 필요하다
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

// Write response 에 쓰기
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

func (j *jar) WriteTo(responseWriter http.ResponseWriter) {
	// 일단 지우고, 다른 데서 쿠키는 못 쓰게 만든다
	header := responseWriter.Header()
	header.Del("Set-Cookie")
	for _, ck := range j.response {
		if v := ck.String(); v != "" {
			header.Add("Set-Cookie", v)
		}
	}
}

// New Make cookie jar from http request
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

// MakeFrom http.Request를 가지고 쿠키를 만든다. 같은 이름이 있으면 나중 걸로 덮어씌워버림
func MakeFrom(request *http.Request) Jar {
	jar := jar{
		request:        make(map[string]*http.Cookie),
		response:       make(map[string]*http.Cookie),
	}
	for _, c := range request.Cookies() {
		jar.request[c.Name] = c
	}
	return &jar
}
