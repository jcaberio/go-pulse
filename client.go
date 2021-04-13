package pulse

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/jcaberio/go-pulse/internal"
)

// Client is wrapper for Feedzai Pulse API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient returns a client, error will be non-nil if the authentication failed.
func New(options *Options) (*Client, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Jar:     cookieJar,
		Timeout: options.Timeout,
	}

	client := &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(options.BaseURL, "/"),
	}

	creds := newCredentials(options.Username, options.Password)
	credsPayload, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}

	loginUrl := client.baseURL + "/pulseviews/api/sessions"

	loginResp, err := client.post(loginUrl, credsPayload)
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

	return client, nil
}

func (c *Client) post(url string, payload []byte) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	request.Header.Set("content-type", "application/json")
	return c.httpClient.Do(request)
}

func (c *Client) put(url string, payload []byte) (*http.Response, error) {
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	request.Header.Set("content-type", "application/json")
	return c.httpClient.Do(request)
}

func (c *Client) get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	request.Header.Set("content-type", "application/json")
	return c.httpClient.Do(request)
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

	uploadResp, err := c.upload(filename, uploadUrl)
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

func (c *Client) DownloadList(filename string, listID string) error {
	downloadUrl := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/managedlists/%s/managedlistitems/csv",
		c.baseURL, listID)
	return c.download(filename, downloadUrl)
}

func (c *Client) upload(filename string, url string) (*http.Response, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := w.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	uploadReq, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, err
	}
	uploadReq.Header.Set("content-type", w.FormDataContentType())
	uploadReq.Header.Set("x-requested-with", "XMLHttpRequest")

	return c.httpClient.Do(uploadReq)
}

func (c *Client) download(filename, url string) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func (c *Client) ExportApp(filename string) error {
	exportUrl := c.baseURL + "/pulseviews/api/apps/paymaya/export"
	return c.download(filename, exportUrl)
}

func (c *Client) ImportResource(zipFile, workflowName, workflowElement string) error {
	workflowElementID, err := c.getWorkflowElementID(workflowName, workflowElement)
	if err != nil {
		return err
	}

	partialResp, err := c.partialImportPrepare(zipFile)
	if err != nil {
		return err
	}

	schemaRequestRulesProjects := make([]internal.RulesProject, len(partialResp.RulesProjects))
	for i, rulesProject := range partialResp.RulesProjects {
		snapshots := make([]internal.Snapshot, len(rulesProject.Snapshots))

		for j, snapshot := range rulesProject.Snapshots {
			snapshots[j] = internal.Snapshot{
				Desc: snapshot.Desc,
				ID:   snapshot.ID,
				WorkflowMappings: []internal.WorkflowMapping{
					{WorkflowElementId: workflowElementID, WorkflowId: "workflow"},
				},
			}
		}

		schemaRequestRulesProjects[i] = internal.RulesProject{
			ID:              rulesProject.ID,
			DestinationDesc: rulesProject.Desc,
			DestinationId:   rulesProject.ID,
			Snapshots:       snapshots,
		}
	}

	schemaRequest := &internal.PartialImportSchemasRequest{
		Lists:         []internal.List{},
		Models:        partialResp.Models,
		Plans:         partialResp.Plans,
		RulesProjects: schemaRequestRulesProjects,
	}

	schemaPayload, err := json.Marshal(schemaRequest)
	if err != nil {
		return err
	}

	checkSchemaURL := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/partialImportCheckSchemas/%s",
		c.baseURL, partialResp.ImportID)
	checkSchemaResp, err := c.post(checkSchemaURL, schemaPayload)
	if err != nil {
		return err
	}
	defer checkSchemaResp.Body.Close()

	if checkSchemaResp.StatusCode != http.StatusOK {
		return errors.New("pulse: failed schema validation")
	}

	commitURL := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/partialImportCommit/%s",
		c.baseURL, partialResp.ImportID)
	commitResp, err := c.post(commitURL, schemaPayload)
	if err != nil {
		return err
	}
	defer commitResp.Body.Close()

	if commitResp.StatusCode != http.StatusNoContent {
		errors.New("pulse: failed partial commit")
	}

	return c.update()
}

func (c *Client) partialImportPrepare(zipFile string) (*internal.PartialImportPrepareResponse, error) {
	partialImportURL := c.baseURL + "/pulseviews/api/apps/paymaya/partialImportPrepare"
	resp, err := c.upload(zipFile, partialImportURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var partialImportResp internal.PartialImportPrepareResponse

	if err := json.Unmarshal(body, &partialImportResp); err != nil {
		return nil, err
	}
	return &partialImportResp, nil
}

func (c *Client) getWorkflowElementID(workflowName, workflowElementName string) (string, error) {
	rteURL := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/rte_workflows/paged?limit=5&sort_by=desc&order=ASC&_=%d",
		c.baseURL, (time.Now().UnixNano() / int64(time.Millisecond)))

	resp, err := c.get(rteURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var rteWorkflow internal.RteWorkflow
	if err := json.Unmarshal(body, &rteWorkflow); err != nil {
		return "", err
	}

	for _, item := range rteWorkflow.Items {
		if item.Desc == workflowName {
			for _, elem := range item.Config.Elements {
				if elem.Desc == workflowElementName {
					return elem.ID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("pulse: workflowElementId for %s not found", workflowElementName)
}

func (c *Client) DeleteApp(appName string) error {
	deleteAppURL := fmt.Sprintf("%s/pulseviews/api/apps/%s", c.baseURL, appName)
	request, err := http.NewRequest("DELETE", deleteAppURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	_, err = c.httpClient.Do(request)
	return err
}

func (c *Client) ImportApp(filename string) error {
	prepareImportURL := fmt.Sprintf("%s/pulseviews/api/apps/prepareImport", c.baseURL)
	resp, err := c.upload(filename, prepareImportURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return err
	}

	var prepareImportResp internal.PrepareImportResponse
	if err := json.Unmarshal(decoded, &prepareImportResp); err != nil {
		return err
	}

	importReq := &internal.ImportRequest{
		ImportID: prepareImportResp.ImportID,
		OwnershipGroupsMatching: struct {
			ImportID string `json:"importId"`
		}{
			prepareImportResp.ImportID,
		},
	}

	importReqPayload, err := json.Marshal(importReq)
	if err != nil {
		return err
	}

	importURL := fmt.Sprintf("%s/pulseviews/api/apps/import", c.baseURL)
	resp, err = c.post(importURL, importReqPayload)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("pulse: failed to import app")
	}

	startURL := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/lifecycle/start", c.baseURL)

	startReq := &internal.PublishRequest{
		Async:        true,
		FullReload:   true,
		Rolling:      true,
		SkipRecovery: false,
	}
	startReqPayload, err := json.Marshal(startReq)
	if err != nil {
		return err
	}

	startResp, err := c.post(startURL, startReqPayload)
	if err != nil {
		return err
	}

	if startResp.StatusCode != http.StatusOK {
		return errors.New("pulse: failed to start publish")
	}
	return nil
}

func (c *Client) update() error {

	updateRequest := &internal.PublishRequest{
		Async:        true,
		FullReload:   true,
		Rolling:      true,
		SkipRecovery: false,
	}
	updatePayload, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}

	updateURL := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/lifecycle/update", c.baseURL)

	updateResp, err := c.post(updateURL, updatePayload)
	if err != nil {
		return err
	}
	defer updateResp.Body.Close()

	body, err := ioutil.ReadAll(updateResp.Body)
	if err != nil {
		return err
	}
	var update internal.UpdateResponse

	if err := json.Unmarshal(body, &update); err != nil {
		return err
	}

	return nil
}

func (c *Client) Restart() error {
    rteURL := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/rte_workflows/paged?limit=1&_=%d",
    		c.baseURL, (time.Now().UnixNano() / int64(time.Millisecond)))

    resp, err := c.get(rteURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    var rteWorkflow internal.RteWorkflow
    if err := json.Unmarshal(body, &rteWorkflow); err != nil {
        return err
    }
    log.Printf("%s\n", body)

    items := rteWorkflow.Items
    if len(items) != 0 {
        item := items[0]
        payload, err := json.Marshal(item)
        if err != nil {
        	return err
        }
		log.Printf("%s\n", payload)

        workflowURL := fmt.Sprintf("%s/pulseviews/api/apps/paymaya/rte_workflows/workflow", c.baseURL)
        resp, err := c.put(workflowURL, payload)
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
            return errors.New("pulse: failed saving workflow")
        }
    }

    return c.update()
}