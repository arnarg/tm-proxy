package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"

	"github.com/gin-gonic/gin"
)

type fastGPTResults struct {
	Meta struct {
		ID         string  `json:"id"`
		Node       string  `json:"node"`
		MS         int64   `json:"ms"`
		APIBalance float64 `json:"api_balance"`
	} `json:"meta"`
	Data fastGPTData `json:"data"`
}

type fastGPTData struct {
	Output     string              `json:"output"`
	Tokens     int64               `json:"tokens"`
	References []fastGPTReferences `json:"references"`
}

type fastGPTReferences struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

var tmpl = `-----
Results for search query "{{.Query}}"
-----
{{.Results}}

References:
{{- range $i, $r := .References}}
[{{inc $i}}]: {{$r.Title}} ({{$r.URL}})
{{- end}}
`

type FastGPTResponse struct {
	Content string `json:"content"`
}

func fetchFastGPT(c *gin.Context) {
	qq := c.Query("q")
	if qq == "" {
		c.JSON(
			http.StatusBadRequest,
			&ServiceResponse{
				Success:        false,
				Message:        "`url` query parameter missing",
				StatusCode:     http.StatusBadRequest,
				ResponseObject: nil,
			},
		)
		return
	}

	// Unescape query
	query, err := url.QueryUnescape(qq)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			&ServiceResponse{
				Success:        false,
				Message:        err.Error(),
				StatusCode:     http.StatusBadRequest,
				ResponseObject: nil,
			},
		)
		return
	}

	// Get API key from header
	key := c.Request.Header.Get("Kagi-API-Key")
	if key == "" {
		c.JSON(
			http.StatusBadRequest,
			&ServiceResponse{
				Success:        false,
				Message:        "`Kagi-API-Key` header missing",
				StatusCode:     http.StatusBadRequest,
				ResponseObject: nil,
			},
		)
		return
	}

	// Make request
	results, err := fastGPTRequest(query, key)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			&ServiceResponse{
				Success:        false,
				Message:        err.Error(),
				StatusCode:     http.StatusInternalServerError,
				ResponseObject: nil,
			},
		)
		return
	}

	// Template output
	output, err := templateFastGPTResults(query, results)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			&ServiceResponse{
				Success:        false,
				Message:        err.Error(),
				StatusCode:     http.StatusInternalServerError,
				ResponseObject: nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		&ServiceResponse{
			Success:    true,
			Message:    "Content fetched successfully",
			StatusCode: http.StatusOK,
			ResponseObject: &FastGPTResponse{
				Content: output,
			},
		},
	)
}

func fastGPTRequest(query, key string) (*fastGPTResults, error) {
	// API endpoint
	const apiURL = "https://kagi.com/api/v0/fastgpt"

	// Prepare request body
	requestData := map[string]string{
		"query": query,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set required headers
	req.Header.Set("Authorization", "Bot "+key)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status: %d %s", resp.StatusCode, string(body))
	}

	// Decode the successful JSON response
	results := &fastGPTResults{}
	if err := json.NewDecoder(resp.Body).Decode(results); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return results, nil
}

func templateFastGPTResults(query string, data *fastGPTResults) (string, error) {
	// Prepare data for the template
	templateData := map[string]any{
		"Query":      query,
		"Results":    data.Data.Output,
		"References": data.Data.References,
	}

	// Parse and execute the template
	t, err := template.
		New("fastgpt_result").
		Funcs(template.FuncMap{
			"inc": func(i int) int {
				return i + 1
			},
		}).
		Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
