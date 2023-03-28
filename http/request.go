package http

import (
   "bytes"
   "io"
   "net/http"
   "net/url"
)

type Body interface {
   []byte | string
}

type Request struct {
   *http.Request
}

func Get() *Request {
   return New_Request(http.MethodGet, new(url.URL))
}

func New_Request(method string, ref *url.URL) *Request {
   return &Request{
      &http.Request{
         Header: make(http.Header),
         Method: method,
         ProtoMajor: 1,
         ProtoMinor: 1,
         URL: ref,
      },
   }
}

func Parse_Get(ref string) (*Request, error) {
   href, err := url.Parse(ref)
   if err != nil {
      return nil, err
   }
   return New_Request(http.MethodGet, href), nil
}

func Parse_Post[T Body](ref string, body T) (*Request, error) {
   href, err := url.Parse(ref)
   if err != nil {
      return nil, err
   }
   read := bytes.NewReader([]byte(body))
   req := New_Request(http.MethodPost, href)
   req.Body = io.NopCloser(read)
   return req, nil
}

func Post[T Body](body T) *Request {
   read := bytes.NewReader([]byte(body))
   req := New_Request(http.MethodPost, new(url.URL))
   req.Body = io.NopCloser(read)
   return req
}
