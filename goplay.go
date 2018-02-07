package goplay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Result is a go playground compilation result
type Result struct {
	// Compilation errors
	Errors string
	// Aggregated stdout
	Stdout string
	// Aggregated stderr
	Stderr string
}

const (
	playCompileURL = "https://play.golang.org/compile"

	kindStdout = "stdout"
	kindStderr = "stderr"
)

type playResponse struct {
	Errors string
	Events []struct {
		Message string
		Kind    string
		Delay   int
	}
}

// Fetch returns code from go playground by link
func Fetch(link string) (string, error) {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		return "", fmt.Errorf("unable to fetch link %v", err)
	}
	code := doc.Find("#code").Text()
	if code == "" {
		return "", fmt.Errorf("unable to find code block on page")
	}
	return code, nil
}

// Compile compiles code and returns aggregated result with compilation errors, stdout and stderr
func Compile(code string) (*Result, error) {
	form := url.Values{}
	form.Add("version", "2")
	form.Add("body", code)
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest(http.MethodPost, playCompileURL, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %v", err)
	}

	playResp := &playResponse{}
	err = json.NewDecoder(resp.Body).Decode(playResp)
	if err != nil {
		return nil, fmt.Errorf("unable to parse response: %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close response body: %v", err)
	}

	if playResp.Errors != "" {
		return &Result{Errors: playResp.Errors}, nil
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	for _, event := range playResp.Events {
		switch event.Kind {
		case kindStdout:
			_, err = stdout.WriteString(event.Message)
		case kindStderr:
			_, err = stderr.WriteString(event.Message)
		}
		if err != nil {
			return nil, fmt.Errorf("unable to write to buffer: %v", err)
		}
	}
	return &Result{
		Errors: "",
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, nil
}
