package pulse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

// Client is wrapper for Feedzai Pulse API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient returns a client, error will be non-nil if the authentication failed.
func NewClient(options *Options) (*Client, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Jar:     cookieJar,
		Timeout: options.Timeout,
	}

	creds := newCredentials(options.Username, options.Password)
	credsPayload, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}

	baseURL := strings.TrimRight(options.BaseURL, "/")
	loginUrl := baseURL + "/pulseviews/api/sessions"

	loginReq, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(credsPayload))
	if err != nil {
		return nil, err
	}
	loginReq.Header.Set("x-requested-with", "XMLHttpRequest")

	loginResp, err := httpClient.Do(loginReq)
	if err != nil {
		return nil, err
	}
	defer loginResp.Body.Close()
	if loginResp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(loginResp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("pulse: %s\n", body)
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}, nil
}

// UploadList uploads the contents of a CSV file named filename to a Pulse list identified with listID.
// error will be non-nil of there are network issues, duplicate entries or the HTTP status
// returned by Feedzai's API is not StatusNoContent.
func (c *Client) UploadList(filename string, listID string) error {

	if filename == "" {
		return errors.New("pulse: empty filename")
	}

	if listID == "" {
		return errors.New("pulse: empty listID")
	}

	uploadUrl := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/managedlists/%s/managedlistitems",
		c.baseURL, listID)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	csvfile, err := os.Open(filename)
	if err != nil {
		return err
	}

	part, err := w.CreateFormFile("file", csvfile.Name())
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, csvfile); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	uploadReq, err := http.NewRequest("POST", uploadUrl, &buf)
	if err != nil {
		return err
	}
	uploadReq.Header.Set("content-type", w.FormDataContentType())
	uploadReq.Header.Set("x-requested-with", "XMLHttpRequest")

	uploadResp, err := c.httpClient.Do(uploadReq)
	if err != nil {
		return err
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(uploadResp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("pulse: %s\n", body)
	}

	return nil
}
