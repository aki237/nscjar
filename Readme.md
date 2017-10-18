# nscjar

nscjar is a small golang package used to parse and output Netscape/Mozilla's old-school cookie files.
It also implements a simple cookie jar struct to manage the cookies added to the cookie jar.

## Usage

### Parser
The usage is simple. The cookie data file has a simple format.

```
DOMAIN    FLAGS    PATH     SECURE     EXPIRY_TIMESTAMP     NAME      VALUE
some.com   TRUE     /        TRUE         15333423          lang      en_US
```

Parsing this in go is simple.

```go
    jar := nscjar.Parser{}
    cookies, err := jar.Unmarshal(f)            // f : io.Reader
    ...handle error...
```

And writing back to file is very simple too :

```go
    jar := nscjar.Parser{}
    jar.Marshal(w, c)          // w : io.Writer, c : http.Cookie
```

See the [test.go](example/test.go) file for the usage.

### Cookie Jar

The cookie jar is used to manage the cookie data added to the struct.
Usage is simple. Say you have a HTTP Response object (`r := http.Response{}`).

```go
    jar := nscjar.NewCookieJar()
    jar.AddCookies(r.Cookies())          // To add multiple *http.Cookie at the same time.

    // Unmarshalling it to a Netscape Cookie file format is easy.
    f, _ := os.OpenFile("cookie.txt", os.O_CREATE|os.O_RDWR, 0666)
    // please do the error checks.
    jar.Marshal(f) // any io.Writer will do. This returns an error.
```
