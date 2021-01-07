package jar

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestGet(t *testing.T) {

	// 테스트 서버
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j := New(r, w)
		get := j.Get("GetCookie")
		assert.Equal(t, "GetCookieValue", get.Value)
	}))

	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	assert.NoError(t, err)
	client := &http.Client{}

	// 리퀘스트에 쿠키가 있으면

	getCookie := http.Cookie{
		Name:   "GetCookie",
		Value:  "GetCookieValue",
		Path:   u.Path,
		Domain: u.Hostname(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Cookie", getCookie.String())
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	// Get 으로 가져올 수 있다

	// response 에는 없다
	get := res.Header.Get("Set-Cookie")
	assert.Equal(t, 0, len(get))
}

func TestSet(t *testing.T) {

	// 테스트 서버
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j := New(r, w)
		// 리퀘스트에 쿠키가 있으면
		get := j.Get("GetCookie")
		assert.Equal(t, "GetCookieValue", get.Value)

		j.Set(&http.Cookie{
			Name:  "SetCookie",
			Value: "SetCookieValue",
			Path:  r.RequestURI,
		})

		j.Set(&http.Cookie{
			Name:  "GetCookie",
			Value: "GetCookieValue2",
			Path:  r.RequestURI,
		})
	}))

	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	assert.NoError(t, err)
	client := &http.Client{}

	getCookie := http.Cookie{
		Name:   "GetCookie",
		Value:  "GetCookieValue",
		Path:   u.Path,
		Domain: u.Hostname(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Cookie", getCookie.String())
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	// Get 으로 가져올 수 있다
	// Set으로 쿠키값을 지정하면
	// response 에 해당값이 있다
	get := res.Header["Set-Cookie"]
	assert.Regexp(t, "SetCookie=SetCookieValue", get)
	assert.Regexp(t, "GetCookie=GetCookieValue2", get)
}

func TestRemove(t *testing.T) {

	// 테스트 서버
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j := New(r, w)
		// 리퀘스트에 쿠키가 있으면
		get := j.Get("GetCookie")
		assert.Equal(t, "GetCookieValue", get.Value)
		j.Remove("GetCookie", r.RequestURI)

		j.Set(&http.Cookie{
			Name:  "SetCookie",
			Value: "SetCookieValue",
			Path:  r.RequestURI,
		})
	}))

	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	assert.NoError(t, err)
	client := &http.Client{}

	getCookie := http.Cookie{
		Name:   "GetCookie",
		Value:  "GetCookieValue",
		Path:   u.Path,
		Domain: u.Hostname(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Cookie", getCookie.String())
	assert.NoError(t, err)

	res, err := client.Do(req)
	assert.NoError(t, err)
	// Get 으로 가져올 수 있다
	// Set으로 쿠키값을 지정하면
	// response 에 해당값이 있다
	get := res.Header["Set-Cookie"]
	// 쿠키는 2개
	assert.Equal(t, 2, len(get))
	setCookie := false
	setReg := regexp.MustCompile("SetCookie=SetCookieValue")
	getCookie2 := false
	getReg := regexp.MustCompile(`GetCookie=(;|$)`)
	for _, s := range get {
		if setReg.MatchString(s) {
			setCookie = true
		}
		if getReg.MatchString(s) {
			getCookie2 = true
		}
	}
	assert.Equal(t, true, setCookie, "SetCookie is not set")
	assert.Equal(t, true, getCookie2, "GetCookie is not changed")
}
