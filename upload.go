package goyadisk

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
)

func (y *Yadisk) Upload(source string, destination string) error {
	parseBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("error parsing base URL: %s", err)
	}

	destinationPath := path.Join(y.appDir, destination)
	fmt.Println("Destination path: ", destinationPath)

	client := resty.New()

	var Response struct {
		Href string `json:"href"`
	}

	resp, err := client.R().
		SetHeader("Authorization", "OAuth "+y.token).
		ForceContentType("application/json").
		SetResult(&Response).
		SetQueryParams(map[string]string{
			"path":      destinationPath + ".mytmp",
			"overwrite": "true",
		}).
		Get(parseBaseURL.JoinPath("resources/upload").String())

	if err != nil {
		return fmt.Errorf("error getting upload URL: %s", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("response Status: %s .Error getting upload URL: %s", resp.Status(), resp)
	}
	fmt.Println("Upload Href: ", Response.Href)

	fileBytes, err := os.ReadFile(filepath.Clean(source))
	if err != nil {
		return fmt.Errorf("error reading file: %s", err)
	}
	start := time.Now()
	resp, err = client.R().
		SetHeader("Authorization", "OAuth "+y.token).
		SetBody(fileBytes).
		Put(Response.Href)

	if err != nil {
		return fmt.Errorf("error putting file: %s", err)
	}
	if resp.StatusCode() != 201 {
		return fmt.Errorf("response Status: %s .Error putting file: %s", resp.Status(), resp)
	}
	fmt.Println("PUT time: ", time.Since(start))

	resp, err = client.R().
		SetHeader("Authorization", "OAuth "+y.token).
		SetQueryParams(map[string]string{
			"from":      destinationPath + ".mytmp",
			"path":      destinationPath,
			"overwrite": "true",
		}).
		Post(parseBaseURL.JoinPath("resources/move").String())

	if err != nil {
		return fmt.Errorf("error renaming file: %s", err)
	}
	if resp.StatusCode() != 201 && resp.StatusCode() != 202 {
		return fmt.Errorf("response Status: %s .Error renaming file: %s", resp.Status(), resp)
	}
	return nil
}
