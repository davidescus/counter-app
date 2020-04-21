package syncsrv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type client struct{}

func (c *client) send(data []Occurrence, addr string) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", addr, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("response Statuscode: %d from: %s", resp.StatusCode, addr))
	}

	return nil
}
