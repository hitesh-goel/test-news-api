package apireq

import (
	"encoding/json"
	"net/http"
	"net/url"
)

//GetAPIRequests calls GET API to fetch results
func GetAPIRequests(apiURL string, resp interface{}) error {
	u, err := url.ParseRequestURI(apiURL)
	if err != nil {
		return err
	}

	urlStr := u.String()
	httpResp, err := http.Get(urlStr)
	if err != nil {
		return err
	}

	err = json.NewDecoder(httpResp.Body).Decode(resp)

	if err != nil {
		return err
	}

	defer httpResp.Body.Close()

	return nil
}
