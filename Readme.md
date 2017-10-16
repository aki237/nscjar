# nscjar

nscjar is a small golang package used to parse and output Netscape/Mozilla's old-school cookie files.

## Usage

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
