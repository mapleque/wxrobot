package wxrobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

// Robot which instance use to send message to wx api
type Robot struct {
	webhook string

	debug bool
}

// New create the instance of Robot
func New(webhook string) *Robot {
	return &Robot{
		webhook: webhook,
	}
}

// Debug enable debug mode, only print message in log
func (r *Robot) Debug() {
	r.debug = true
}

// SendTextMessageWithTemplate send text format message building in go template
func (r *Robot) SendTextMessageWithTemplate(
	tpl string,
	v interface{},
) error {
	t, err := template.New("default").Parse(tpl)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}
	var buff bytes.Buffer
	if err := t.Execute(&buff, v); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}
	if r.debug {
		log.Printf("message will be send: \n----\n%s\n----\n", buff.String())
	}
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": buff.String(),
		},
	}
	bts, _ := json.Marshal(msg)
	resp, err := http.Post(
		r.webhook,
		"application/json",
		bytes.NewReader(bts),
	)
	if err != nil {
		return fmt.Errorf("request weixin api error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read weixin response error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"weixin response error code: %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}
	return nil
}
