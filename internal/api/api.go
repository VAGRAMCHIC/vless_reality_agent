package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Server string // например "xray:10085"
	Tag    string // inbound tag в Xray, например "vless-reality"
}

// New создает нового клиента
func New(server, tag string) *Client {
	return &Client{
		Server: server,
		Tag:    tag,
	}
}

// структура запроса для Xray API
type xrayCommand struct {
	Command   string         `json:"command"`
	Arguments map[string]any `json:"arguments"`
}

// runCommand отправляет команду на Xray API
func (c *Client) runCommand(cmd string, args map[string]any) error {
	url := fmt.Sprintf("http://%s/", c.Server)

	payload := xrayCommand{
		Command:   cmd,
		Arguments: args,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("xray api error: %s", string(body))
	}

	return nil
}

// AddUser добавляет пользователя через Xray API
func (c *Client) AddUser(uuid, email string) error {
	args := map[string]any{
		"tag":   c.Tag,
		"email": email,
		"id":    uuid,
	}
	return c.runCommand("addUser", args)
}

// RemoveUser удаляет пользователя через Xray API
func (c *Client) RemoveUser(email string) error {
	args := map[string]any{
		"tag":   c.Tag,
		"email": email,
	}
	return c.runCommand("removeUser", args)
}
