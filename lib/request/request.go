package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Get(url string, params map[string]string) (resp *http.Response, err error) {
	query := ""
	for key, value := range params {
		query = fmt.Sprintf("%s&%s=%s", query, key, value)
	}

	fullUrl := fmt.Sprintf("%s?%s", url, query)
	resp, err = http.Get(fullUrl)
	if err != nil {
		currentTime := time.Now()
		fmt.Printf("[%s] [error] fail sending request: url => %s, params => %v\n", currentTime.Format("2006/01/02 15:04:05"), url, params)
	}
	defer resp.Body.Close()

	return
}

func Post(url string, params interface{}) (resp *http.Response, err error) {
	jsonParams, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonParams))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		currentTime := time.Now()
		fmt.Printf(" ==> [%s] [error] fail sending request: url => %s, params => %v\n", currentTime.Format("2006/01/02 15:04:05"), url, params)
	}
	defer resp.Body.Close()

	return
}
