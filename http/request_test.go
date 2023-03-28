package http

import (
   "net/http"
   "os"
   "testing"
)

func Test_URL(t *testing.T) {
   req, err := Get_Parse("http://httpbin.org/get")
   if err != nil {
      t.Fatal(err)
   }
   res, err := new(http.Transport).RoundTrip(req.Request)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
}

func Test_Request(t *testing.T) {
   req := Get()
   req.URL.Scheme = "http"
   req.URL.Host = "httpbin.org"
   req.URL.Path = "/get"
   Default_Client.Log_Level = 2
   res, err := Default_Client.Do(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
}

func Test_Body(t *testing.T) {
   req := Post("hello=world")
   req.URL.Scheme = "http"
   req.URL.Host = "httpbin.org"
   req.URL.Path = "/post"
   res, err := new(http.Transport).RoundTrip(req.Request)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
}
