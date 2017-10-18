package nscjar

import (
	"fmt"
	"io"
	"net/http"
)

// CookieJar is a simple struct used to manage the cookies
type CookieJar struct {
	cookies []*http.Cookie
}

// NewCookieJar function returns a newly initialized CookieJar struct
func NewCookieJar() *CookieJar {
	return &CookieJar{make([]*http.Cookie, 0)}
}

// AddCookie method is used to add a cookie to the jar
func (j *CookieJar) AddCookie(c *http.Cookie) {
	if c == nil {
		return
	}
	changed := false
	for i, val := range j.cookies {
		if val.Domain == c.Domain && val.Path == c.Path && val.Name == c.Name {
			j.cookies[i] = c
			fmt.Println("Cookie added : ", c)
			changed = true
			break
		}
	}
	if !changed {
		j.cookies = append(j.cookies, c)
	}
}

// AddCookies method is used to add multiple cookies to the jar at once
func (j *CookieJar) AddCookies(c ...*http.Cookie) {
	for _, val := range c {
		j.AddCookie(val)
	}
}

// Marshal method is used to write the cookie data in the passed io.Writer in netscape cookie data format
func (j CookieJar) Marshal(wr io.Writer) error {
	p := Parser{}
	for _, val := range j.cookies {
		err := p.Marshal(wr, val)
		if err != nil {
			return err
		}
	}
	return nil
}
