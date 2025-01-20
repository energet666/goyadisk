package goyadisk

import (
	"fmt"
	"path"

	"github.com/go-resty/resty/v2"
)

type YaItem struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (y *Yadisk) GetItems(relPath string) (items []YaItem, err error) {
	client := resty.New()

	var Response struct {
		Embedded struct {
			Sort  string   `json:"sort"`
			Total int      `json:"total"`
			Limit int      `json:"limit"`
			Items []YaItem `json:"items"`
		} `json:"_embedded"`
	}

	targetPath := path.Join(y.appDir, relPath)

	resp, err := client.R().
		EnableTrace().
		SetHeader("Authorization", "OAuth "+y.token).
		ForceContentType("application/json").
		SetResult(&Response).
		SetQueryParams(map[string]string{
			"path":  targetPath,
			"sort":  "-modified",
			"limit": "9999",
		}).
		Get(baseURL + "/resources")
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %s", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error getting resources: %s", resp.Status())
	}

	fmt.Println("Total request time: ", resp.Request.TraceInfo().TotalTime)
	fmt.Println("Sort: ", Response.Embedded.Sort)
	fmt.Println("Total: ", Response.Embedded.Total)
	fmt.Println("Limit: ", Response.Embedded.Limit)
	fmt.Println("Response Status: ", resp.Status())

	return Response.Embedded.Items, nil
}
