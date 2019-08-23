package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type RunnerClient struct {
	baseURL string
}

func NewRunnerClient() (*RunnerClient, error) {
	baseURL := os.Getenv("SERVICE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("No service url found. Please provide service url in SERVICE_URL")
	}
	return &RunnerClient{baseURL: baseURL}, nil
}

type request struct {
	Image   string `json:"image"`
	Tag     string `json:"tag"`
	Timeout int64  `json:"timeout"`
}

type callback struct {
	CallbackURL string `json:"callbackUrl"`
}

func (client *RunnerClient) SubmitJob(image, tag string, timeout int64) (string, error) {
	jobRequest := request{Image: image, Tag: tag, Timeout: timeout}
	b, _ := json.Marshal(jobRequest)
	response, err := http.Post(fmt.Sprintf("%v/execute", client.baseURL), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return "", err
	}

	var callbackURL callback
	err = json.Unmarshal(responseBytes, &callbackURL)
	if err != nil {
		return "", err
	}

	return callbackURL.CallbackURL, nil
}

type statusDetails struct {
	Status string
}

func (client *RunnerClient) GetJobStatus(jobURL string) (string, error) {
	url := fmt.Sprintf("%v%v", client.baseURL, jobURL)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return "", err
	}

	var status statusDetails
	err = json.Unmarshal(responseBytes, &status)

	if err != nil {
		return "", err
	}

	return status.Status, nil
}
