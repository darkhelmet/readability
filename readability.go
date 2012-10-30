package readability

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const (
    Parser = "https://readability.com/api/content/v1/parser"
)

var (
    ErrTransient = errors.New("readability: transient error, probably a 5xx, maybe try again")
)

type Response struct {
    Domain        string  `json:"domain"`
    Author        *string `json:"author"`
    URL           URL     `json:"url"`
    ShortURL      URL     `json:"short_url"`
    Title         string  `json:"title"`
    TotalPages    int     `json:"total_pages"`
    WordCount     int     `json:"word_count"`
    Content       string  `json:"content"`
    DatePublished *Time   `json:"date_published"`
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

func parseResponse(uri string, r io.Reader) (*Response, error) {
    var rresp Response
    decoder := json.NewDecoder(r)
    err := decoder.Decode(&rresp)
    if err != nil {
        return nil, fmt.Errorf("readability: JSON error (%s): %s", uri, err)
    }
    return &rresp, nil
}

func (e *Endpoint) Extract(uri string) (*Response, error) {
    resp, err := http.Get(e.buildUrl(uri))
    if err != nil {
        return nil, fmt.Errorf("readability: HTTP error (%s): %s", uri, err)
    }
    defer resp.Body.Close()

    switch {
    case resp.StatusCode >= 500:
        // Eat and throw away the body
        io.Copy(ioutil.Discard, resp.Body)
        return nil, ErrTransient
    case resp.StatusCode == 200:
        // All is well
    default:
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("readability: HTTP error (%s): %d, %s", uri, resp.StatusCode, body)
    }

    return parseResponse(uri, resp.Body)
}

func (e *Endpoint) ExtractWithContent(uri, content string) (*Response, error) {
    resp, err := http.Post(e.buildUrl(uri), "application/x-www-form-urlencoded", strings.NewReader(url.QueryEscape(content)))
    if err != nil {
        return nil, fmt.Errorf("readability: HTTP error (%s): %s", uri, err)
    }
    defer resp.Body.Close()

    switch {
    case resp.StatusCode >= 500:
        io.Copy(ioutil.Discard, resp.Body)
        return nil, ErrTransient
    case resp.StatusCode == 200:
        // All is well
    default:
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("readability: HTTP error (%s): %d, %s", uri, resp.StatusCode, body)
    }

    return parseResponse(uri, resp.Body)
}
