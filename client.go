package pulse

import (
	"bytes"
	"encoding/base64"
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
	"time"

	"github.com/jcaberio/go-pulse/internal"
)

// Client is wrapper for Feedzai Pulse API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	appName    string
}

// New returns a client, error will be non-nil if the authentication failed.
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
		appName:    options.AppName,
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

func (c *Client) do(url string, method string, payload []byte) (*http.Response, error) {
	var body io.Reader
	if method == http.MethodGet || method == http.MethodDelete {
		body = nil
	} else {
		body = bytes.NewBuffer(payload)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	request.Header.Set("content-type", "application/json")
	return c.httpClient.Do(request)
}

func (c *Client) post(url string, payload []byte) (*http.Response, error) {
	return c.do(url, http.MethodPost, payload)
}

func (c *Client) put(url string, payload []byte) (*http.Response, error) {
	return c.do(url, http.MethodPut, payload)
}

func (c *Client) get(url string) (*http.Response, error) {
	return c.do(url, http.MethodGet, nil)
}

func (c *Client) delete(url string) (*http.Response, error) {
	return c.do(url, http.MethodDelete, nil)
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

	uploadUrl := fmt.Sprintf("%s/pulseviews/api/apps/%s/managedlists/%s/managedlistitems",
		c.baseURL, c.appName, listID)

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
	downloadUrl := fmt.Sprintf("%s/pulseviews/api/apps/%s/managedlists/%s/managedlistitems/csv",
		c.baseURL, c.appName, listID)
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
	exportUrl := fmt.Sprintf("%s/pulseviews/api/apps/%s/export", c.baseURL, c.appName)
	return c.download(filename, exportUrl)
}

func (c *Client) importPlan(zipFile string) error {

	// partialImportPrepare
	partialResp, err := c.partialImportPrepare(zipFile)
	if err != nil {
		return err
	}

	// partialImportCheckSchemas
	schemaRequestPlans := make([]internal.Plans, len(partialResp.Plans))
	for i, plan := range partialResp.Plans {
		execs := make([]internal.Executions, len(plan.Executions))

		for j, exec := range plan.Executions {
			execs[j] = internal.Executions{
				ID: exec.ID,
			}
		}

		schemaRequestPlans[i] = internal.Plans{
			Executions:      execs,
			ID:              plan.ID,
			DestinationId:   plan.ID,
			DestinationDesc: plan.Desc,
		}
	}

	schemaRequest := &internal.PartialImportSchemasRequest{
		Lists:         []internal.List{},
		Models:        partialResp.Models,
		Plans:         schemaRequestPlans,
		RulesProjects: partialResp.RulesProjects,
	}

	schemaPayload, err := json.Marshal(schemaRequest)
	if err != nil {
		return err
	}

	checkSchemaURL := fmt.Sprintf("%s/pulseviews/api/apps/%s/partialImportCheckSchemas/%s",
		c.baseURL, c.appName, partialResp.ImportID)
	checkSchemaResp, err := c.post(checkSchemaURL, schemaPayload)
	if err != nil {
		return err
	}
	defer checkSchemaResp.Body.Close()

	if checkSchemaResp.StatusCode != http.StatusOK {
		return errors.New("pulse: failed schema validation")
	}

	// partialImportCommit
	commitURL := fmt.Sprintf("%s/pulseviews/api/apps/%s/partialImportCommit/%s",
		c.baseURL, c.appName, partialResp.ImportID)
	commitResp, err := c.post(commitURL, schemaPayload)
	if err != nil {
		return err
	}
	defer commitResp.Body.Close()

	if commitResp.StatusCode != http.StatusNoContent {
		return errors.New("pulse: failed partial commit")
	}
	// update
	return c.update()

}

func (c *Client) ImportRule(zipFile, workflowName, workflowElement string) error {
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

	checkSchemaURL := fmt.Sprintf("%s/pulseviews/api/apps/%s/partialImportCheckSchemas/%s",
		c.baseURL, c.appName, partialResp.ImportID)
	checkSchemaResp, err := c.post(checkSchemaURL, schemaPayload)
	if err != nil {
		return err
	}
	defer checkSchemaResp.Body.Close()

	if checkSchemaResp.StatusCode != http.StatusOK {
		return errors.New("pulse: failed schema validation")
	}

	commitURL := fmt.Sprintf("%s/pulseviews/api/apps/%s/partialImportCommit/%s",
		c.baseURL, c.appName, partialResp.ImportID)
	commitResp, err := c.post(commitURL, schemaPayload)
	if err != nil {
		return err
	}
	defer commitResp.Body.Close()

	if commitResp.StatusCode != http.StatusNoContent {
		return errors.New("pulse: failed partial commit")
	}

	return c.update()
}

func (c *Client) partialImportPrepare(zipFile string) (*internal.PartialImportPrepareResponse, error) {
	partialImportURL := fmt.Sprintf("%s/pulseviews/api/apps/%s/partialImportPrepare",
		c.baseURL, c.appName)
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
	rteURL := fmt.Sprintf("%s/pulseviews/api/apps/%s/rte_workflows/paged?limit=5&sort_by=desc&order=ASC&_=%d",
		c.baseURL, c.appName, time.Now().UnixNano()/int64(time.Millisecond))

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

func (c *Client) DeleteApp() error {
	url := fmt.Sprintf("%s/pulseviews/api/apps/%s", c.baseURL, c.appName)
	_, err := c.delete(url)
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

	url := fmt.Sprintf("%s/pulseviews/api/apps/import", c.baseURL)

	err = c.submit(url, http.MethodPost, importReqPayload, http.StatusOK, "pulse: failed to import app")
	if err != nil {
		return err
	}

	return c.start()
}

func (c *Client) lifecycle(cycle string) error {
	url := fmt.Sprintf("%s/pulseviews/api/apps/%s/lifecycle/%s",
		c.baseURL, c.appName, cycle)

	req := &internal.PublishRequest{
		Async:        true,
		FullReload:   true,
		Rolling:      true,
		SkipRecovery: false,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = c.submit(url, http.MethodPost, payload, http.StatusOK, "pulse: failed to publish")
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) start() error {
	return c.lifecycle("start")
}

func (c *Client) update() error {
	return c.lifecycle("update")
}

func (c *Client) Restart() error {
	body, item, err := c.getWorkflowState()
	if err != nil {
		return err
	}

	payload, err := c.validate(body, item)
	if err != nil {
		return err
	}

	err = c.validateRestoreState(payload)
	if err != nil {
		return err
	}

	err = c.saveWorkflow(body)
	if err != nil {
		return err
	}

	return c.update()
}

func (c *Client) submit(url string, method string, body []byte, statusCode int, errMsg string) error {
	var resp *http.Response
	var err error
	if method == http.MethodPost {
		resp, err = c.post(url, body)
	}
	if method == http.MethodPut {
		resp, err = c.put(url, body)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != statusCode {
		return errors.New(errMsg)
	}
	return nil
}

func (c *Client) saveWorkflow(body []byte) error {
	url := fmt.Sprintf("%s/pulseviews/api/apps/%s/rte_workflows/workflow", c.baseURL, c.appName)
	return c.submit(url, http.MethodPut, body, http.StatusOK, "pulse: failed saving workflow")
}

func (c *Client) validateRestoreState(payload []byte) error {
	url := fmt.Sprintf("%s/pulseviews/api/apps/%s/rte_workflows/validaterestorestate",
		c.baseURL, c.appName)

	return c.submit(url, http.MethodPost, payload, http.StatusNoContent, "pulse: failed validating restore state")
}

func (c *Client) validate(body []byte, item internal.Item) ([]byte, error) {
	url := fmt.Sprintf("%s/pulseviews/api/apps/%s/rte_workflows/validate",
		c.baseURL, c.appName)

	err := c.submit(url, http.MethodPost, body, http.StatusOK, "pulse: failed validating workflow")
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(internal.ValidateRestoreState{
		item.Config.RecoveryExpression})
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (c *Client) getWorkflowState() ([]byte, internal.Item, error) {
	rteURL := fmt.Sprintf("%s/pulseviews/api/apps/%s/rte_workflows/workflow?_=%d",
		c.baseURL, c.appName, time.Now().UnixNano()/int64(time.Millisecond))

	resp, err := c.get(rteURL)
	if err != nil {
		return nil, internal.Item{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, internal.Item{}, err
	}
	var item internal.Item
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, internal.Item{}, err
	}
	return body, item, nil
}

func (c *Client) IsPublishInProgress() bool {
	url := fmt.Sprintf("%s/pulseviews/api/apps/%s/lifecycle/currentOperationProgress?_=%d",
		c.baseURL, c.appName, time.Now().UnixNano()/int64(time.Millisecond))
	resp, _ := c.get(url)
	return resp.StatusCode == http.StatusOK
}

func (c *Client) Abort() error {
	url := fmt.Sprintf("%s/pulseviews/api/apps/%s/lifecycle/currentOperationProgress?_=%d",
		c.baseURL, c.appName, time.Now().UnixNano()/int64(time.Millisecond))
	resp, err := c.get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		var progress internal.ProgressResponse
		if err := json.Unmarshal(body, &progress); err != nil {
			return err
		}

		url := fmt.Sprintf("%s/pulseviews/api/apps/%s/lifecycle/cancel/%s",
			c.baseURL, c.appName, progress.OperationId)

		resp, err = c.post(url, []byte{})
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New("pulse: failed aborting publish")
		}
	}
	return nil
}
