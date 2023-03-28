package http

import (
   "bytes"
   "io"
   "net/http"
   "net/url"
   "strings"
)

type Request struct {
   *http.Request
}

func Get() *Request {
   return New_Request(http.MethodGet, new(url.URL))
}

func Get_Parse(ref string) (*Request, error) {
   href, err := url.Parse(ref)
   if err != nil {
      return nil, err
   }
   return New_Request(http.MethodGet, href), nil
}

func New_Request(method string, ref *url.URL) *Request {
   req := http.Request{
      Header: make(http.Header),
      Method: method,
      ProtoMajor: 1,
      ProtoMinor: 1,
      URL: ref,
   }
   return &Request{&req}
}

func Post(body string) *Request {
   req := New_Request(http.MethodPost, new(url.URL))
   req.Body = io.NopCloser(strings.NewReader(body))
   return req
}

func Post_Bytes(body []byte) *Request {
   req := New_Request(http.MethodPost, new(url.URL))
   req.Body = io.NopCloser(bytes.NewReader(body))
   return req
}

func Post_Parse(ref string, body []byte) (*Request, error) {
   href, err := url.Parse(ref)
   if err != nil {
      return nil, err
   }
   req := New_Request(http.MethodPost, href)
   req.Body = io.NopCloser(bytes.NewReader(body))
   return req, nil
}
