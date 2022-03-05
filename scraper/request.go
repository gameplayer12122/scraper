package scraper

import (
	"net/http"
)

func Request(url string, client *http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "LunasScraper/1.0.0 (lunahatesgogle@gmail.com)")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
