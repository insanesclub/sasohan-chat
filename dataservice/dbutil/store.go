package dbutil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// StoreJSON stores v in form of JSON via http POST request to url.
func StoreJSON(v interface{}, url string) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
