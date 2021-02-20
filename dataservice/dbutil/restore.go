package dbutil

import (
	"encoding/json"
	"io"
	"net/http"
)

// RestoreJSON restores v in form of JSON via http POST to url.
func RestoreJSON(v interface{}, url string, body io.Reader) error {
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
