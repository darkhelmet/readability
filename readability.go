package readability

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

const (
    Parser = "https://readability.com/api/content/v1/parser"
)

type Response struct {
    Domain        string  `json:"domain"`
    Author        *string  `json:"author"`
    URL           URL     `json:"url"`
    ShortURL      URL     `json:"short_url"`
    Title         string  `json:"title"`
    TotalPages    int     `json:"total_pages"`
    WordCount     int     `json:"word_count"`
    Content       string  `json:"content"`
    DatePublished Time    `json:"date_published"`
    NextPageId    *string `json:"next_page_id"`
    RenderedPages int     `json:"rendered_pages"`
}

type Endpoint struct {
    token string
}

func New(token string) *Endpoint {
    return &Endpoint{token}
}

func (e *Endpoint) buildUrl(uri string) string {
    return fmt.Sprintf("%s?url=%s&token=%s", Parser, url.QueryEscape(uri), url.QueryEscape(e.token))
}

func (e *Endpoint) Extract(uri string) (*Response, error) {
    resp, err := http.Get(e.buildUrl(uri))
    if err != nil {
        return nil, fmt.Errorf("readability: HTTP error (%s): %s", uri, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("readability: HTTP error (%s): %d", uri, resp.StatusCode)
    }

    var rresp Response
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&rresp)
    if err != nil {
        return nil, fmt.Errorf("readability: JSON error (%s): %s", uri, err)
    }
    return &rresp, nil
}
