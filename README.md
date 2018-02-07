# goplay

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/yanzay/goplay)

Usage:

```go
code, err := goplay.Fetch("https://play.golang.org/p/HmnNoBf0p1z")
if err != nil {
    // ...
}
result, err := goplay.Compile(code)
if err != nil {
    // ...
}
fmt.Println(result.Errors)
fmt.Println(result.Stdout)
fmt.Println(result.Stderr)
```
