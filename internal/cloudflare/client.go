package cloudflare

import (
    "bytes"
    "net/http"
    "os"
)

type Client struct {
    HTTP   *http.Client
    Token  string
    ZoneID string
}

func New() *Client {
    return &Client{
        HTTP:   &http.Client{},
        Token:  os.Getenv("CF_API_TOKEN"),
        ZoneID: os.Getenv("CF_ZONE_ID"),
    }
}

func (c *Client) NewRequest(method, url string, body []byte) (*http.Request, error) {
    req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }	

    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    return req, nil
}
