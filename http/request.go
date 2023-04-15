package http

import (
   "bytes"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (r Request) Body_Bytes(b []byte) {
   body := bytes.NewReader(b)
   r.Body = io.NopCloser(body)
   r.ContentLength = body.Size()
}

func (r Request) Body_String(s string) {
   body := strings.NewReader(s)
   r.Body = io.NopCloser(body)
   r.ContentLength = body.Size()
}

type Request struct {
   *http.Request
}

func Get() *Request {
   return New_Request(http.MethodGet, new(url.URL))
}

func Get_URL(ref string) (*Request, error) {
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

func Patch() *Request {
   return New_Request(http.MethodPatch, new(url.URL))
}

func Post() *Request {
   return New_Request(http.MethodPost, new(url.URL))
}

func Post_URL(ref string) (*Request, error) {
   href, err := url.Parse(ref)
   if err != nil {
      return nil, err
   }
   return New_Request(http.MethodPost, href), nil
}

func Put() *Request {
   return New_Request(http.MethodPut, new(url.URL))
}
