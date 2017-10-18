/*
Package nscjar is used to parse Netscape/Mozilla's old-school cookie files

These cookie files are also used in cURL for storing cookies.
*/
package nscjar

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Parser is a type meant for parsing the cookie file
type Parser struct {
}

// Marshal method is used to reform the Netscape/Mozilla's old school cookie file for the given cookie
func (p Parser) Marshal(wr io.Writer, c *http.Cookie) error {
	if c.Value == "" || c.Name == "" {
		return errors.New("not a valid cookie")
	}

	cookieLine := ""

	if c.HttpOnly {
		cookieLine += "#HttpOnly_"
	}

	cookieLine += c.Domain + "\tTRUE\t"

	if c.Path == "" {
		c.Path = "/"
	}

	cookieLine += c.Path + "\t"

	if c.Secure {
		cookieLine += "TRUE\t"
	} else {
		cookieLine += "FALSE\t"
	}

	cookieLine += fmt.Sprintf("%d\t", c.Expires.Unix())

	cookieLine += c.Name + "\t"

	addend := ""

	if strings.Contains(c.Value, " ") {
		addend = "\""
	}

	cookieLine += addend + c.Value + addend + "\n"

	_, err := wr.Write([]byte(cookieLine))

	return err
}

// Unmarshal method is used to parse the contents from io.Reader and return the corresponding cookies
func (p Parser) Unmarshal(rd io.Reader) ([]*http.Cookie, error) {
	b := bufio.NewReader(rd)
	err := error(nil)
	line := ""
	cs := make([]*http.Cookie, 0)
	for err == nil {
		line, err = b.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") || line == "" {
			continue
		}
		c, err := getCookieFromString(line)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	if err != io.EOF {
		return cs, err
	}
	return cs, nil
}

// getCookieFromString is used to parse the string passed and generate the corresponding http.Cookie
func getCookieFromString(line string) (*http.Cookie, error) {
	if line == "" {
		return nil, errors.New("empty string passed")
	}

	//pos := 1
	tokens := make([]string, 0)
	inQuotes := false
	token := ""

	for _, ch := range line {
		switch ch {
		case '\t':
			if inQuotes {
				token += string(ch)
				break
			}
			tokens = append(tokens, token)
			token = ""
		case '"':
			inQuotes = !inQuotes
		default:
			token += string(ch)
		}
	}
	tokens = append(tokens, token)

	c := &http.Cookie{}
	for i, val := range tokens {
		err := setCookieField(i, val, c)
		if err != nil {
			return nil, err
		}
	}

	if c.Value == "" || c.Name == "" {
		return nil, errors.New("not a valid cookie data")
	}

	return c, nil
}

// setCookieField function is used to parse the val and assign it to the cookie according
// to the index of the string. Generally used after tokenizing the cookie data line from
// Netscape cookie data. Each ith field have to assigned to a particular struct field.
func setCookieField(i int, val string, c *http.Cookie) error {
	switch i {
	case 0:
		c.Domain = val
		if strings.HasPrefix("#HttpOnly_", val) {
			c.Domain = val[10:]
			c.HttpOnly = true
		}
	case 2:
		c.Path = val
	case 3:
		if val == "TRUE" {
			c.Secure = true
			break
		}
		if val == "FALSE" {
			c.Secure = false
		} else {
			return errors.New("unexpected boolean value found for secure flag")
		}
	case 4:
		timestamp, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		c.Expires = time.Unix(timestamp, 0)
	case 5:
		c.Name = val
	case 6:
		c.Value = val
	}
	return nil
}
