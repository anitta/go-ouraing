package ouraring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

)

const (
    defaultBaseURL = "https://api.ouraring.com/"
)

type Client struct {
    client *http.Client
    BaseURL *url.URL

    common service

    accessToken string

    // Services used for talking to different parts of the oura ring API.
    Sleep *SleepService
}

type service struct {
    client *Client
}

type ListOptions struct {
    StartTime string `url:"start"`
    EndTime string `url:"end"`
}

func NewClient(httpClient *http.Client, accessToken string) *Client {
    if httpClient == nil {
        httpClient = &http.Client{}
    }
    baseURL, _ := url.Parse(defaultBaseURL)

    c := &Client{client: httpClient, BaseURL: baseURL, accessToken: accessToken}

    c.common.client = c
    c.Sleep = (*SleepService)(&c.common)

    return c
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
    if ctx == nil {
        return nil, errors.New("context must be non-nil")
    }

    req = req.WithContext(ctx)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    err = CheckResponse(resp)
    if err != nil {
        return nil, err
    }

    if v != nil {
        if w, ok := v.(io.Writer); ok {
            io.Copy(w, req.Body)
        } else {
            decErr := json.NewDecoder(resp.Body).Decode(v)
            if decErr == io.EOF {
                decErr = nil
            }
            if decErr != nil {
                err = decErr
            }
        }
    }
    return resp, err
}

func CheckResponse(r *http.Response) error {
    if c := r.StatusCode; 200 <= c && c <= 299 {
        return nil
    }

    errorResponse := &http.Response{}
    data, err := ioutil.ReadAll(r.Body)
    if err == nil && data != nil {
        json.Unmarshal(data, errorResponse)
    }

    return errors.New(fmt.Sprintf("status code: %d, data: %#v",r.StatusCode,errorResponse.Body))
}

func (c *Client) NewRequest(urlStr string, body interface{}) (*http.Request, error) {
    u, err := c.BaseURL.Parse(urlStr)
    if err != nil {
        return nil, err
    }

    fmt.Println(c.accessToken)
    q := u.Query()
    q.Set("access_token", c.accessToken)

    u.RawQuery = q.Encode()

    fmt.Println(u.String())

    req, err := http.NewRequest("GET", u.String(), nil)
    if err != nil {
        return nil, err
    }

    return req, nil
}
