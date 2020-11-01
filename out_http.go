package polldance

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/k0kubun/pp"
)

type OutHTTP struct {
	URLTmpl string
	Method  string
	Debug   bool
}

func (o *OutHTTP) Handler(ev *EventData) error {
	url, err := o.buildURL(ev.Source)
	if err != nil {
		return fmt.Errorf("failed to handle event: %w", err)
	}
	pp.Println(ev, url)

	client := &http.Client{}
	req, err := http.NewRequest(o.Method, url, bytes.NewBuffer([]byte(ev.Data)))
	if err != nil {
		return fmt.Errorf("request setup error: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request %s: %w", url, err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("http status=%d, body=%s", resp.StatusCode, b)
	}
	return nil
}

type urlParam struct {
	Path     string
	Filename string
}

func (o *OutHTTP) buildURL(p string) (string, error) {
	tmpl, err := template.New("url").Parse(o.URLTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse url template: %w", err)
	}
	b := &bytes.Buffer{}
	param := &urlParam{
		Path:     p,
		Filename: path.Base(p),
	}
	if err := tmpl.Execute(b, param); err != nil {
		return "", fmt.Errorf("failed to execute url template: %w", err)
	}
	return b.String(), nil
}
