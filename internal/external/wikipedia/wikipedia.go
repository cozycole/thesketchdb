package wikipedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var URL_TEMPLATE = "https://en.wikipedia.org/wiki/%s"

// WikipediaResponse models the structure of the response we care about.
type WikipediaResponse struct {
	Query struct {
		Pages map[string]struct {
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

// GetExtract fetches the intro extract for a given Wikipedia page.
func GetExtract(pageName string) (string, error) {
	endpoint := "https://en.wikipedia.org/w/api.php"

	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("titles", pageName)
	params.Add("prop", "extracts")
	params.Add("exintro", "")

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("User-Agent", "theSketchDb/1.0 (https://thesketchdb.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body failed: %w", err)
	}

	var wikiResp WikipediaResponse
	if err := json.Unmarshal(body, &wikiResp); err != nil {
		return "", fmt.Errorf("json unmarshal failed: %w", err)
	}

	for _, page := range wikiResp.Query.Pages {
		return page.Extract, nil
	}

	return "", fmt.Errorf("no extract found for page %q", pageName)
}

// ExtractPageName takes a Wikipedia URL and returns the page name (e.g. "Shane_Gillis").
func ExtractPageName(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Wikipedia pages live under /wiki/{Page_Name}
	if !strings.HasPrefix(u.Path, "/wiki/") {
		return "", fmt.Errorf("not a valid Wikipedia page URL: %s", rawURL)
	}

	page := strings.TrimPrefix(u.Path, "/wiki/")
	page = path.Clean(page)

	// Remove fragment (#Section) if present
	if u.Fragment != "" {
		page = strings.Split(page, "#")[0]
	}

	return page, nil
}
