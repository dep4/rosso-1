package http

import (
   "2a.pages.dev/rosso/strconv"
   "fmt"
   "net/http"
   "net/http/httputil"
   "strings"
)

func (c Client) Do(req Request) (*Response, error) {
   switch c.Log_Level {
   case 1:
      fmt.Println(req.Method, req.URL)
   case 2:
      dump, err := httputil.DumpRequest(req.Request, true)
      if err != nil {
         return nil, err
      }
      enc := strconv.Encode(dump)
      if strings.HasSuffix(enc, "\n") {
         fmt.Print(enc)
      } else {
         fmt.Println(enc)
      }
   }
   res, err := c.Client.Do(req.Request)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != c.Status {
      return nil, fmt.Errorf(res.Status)
   }
   return res, nil
}

type Client struct {
   Log_Level int // this needs to work with flag.IntVar
   Status int
   http.Client
}

var Default_Client = Client{
   Client: http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   },
   Log_Level: 1,
   Status: http.StatusOK,
}

func (c Client) Clone() Client {
   return c
}
