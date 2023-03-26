package main

import (
   "2a.pages.dev/rosso/strconv"
   "bytes"
   "embed"
   "fmt"
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "strings"
   "text/template"
   "unicode/utf8"
)

func quote(s string) string {
   if can_backquote(s) {
      return "`" + s + "`"
   }
   return strconv.Quote(s)
}

// go.dev/ref/spec#String_literals
func can_backquote(s string) bool {
   for i := range s {
      b := s[i]
      if b == '\r' {
         return false
      }
      if b == '`' {
         return false
      }
      if strconv.Binary_Data(b) {
         return false
      }
   }
   return utf8.ValidString(s)
}

//go:embed _template.go
var content embed.FS

func write_request(req *http.Request, dst io.Writer) error {
   var v values
   if req.Body != nil && req.Method != "GET" {
      body, err := io.ReadAll(req.Body)
      if err != nil {
         return err
      }
      req.Body = io.NopCloser(bytes.NewReader(body))
      v.Raw_Req_Body = quote(string(body))
      v.Req_Body = "io.NopCloser(req_body)"
   } else {
      v.Raw_Req_Body = `""`
      v.Req_Body = "nil"
   }
   v.Query = req.URL.Query()
   v.Request = req
   temp, err := template.ParseFS(content, "_template.go")
   if err != nil {
      return err
   }
   return temp.Execute(dst, v)
}

type values struct {
   *http.Request
   Query url.Values
   Req_Body string
   Raw_Req_Body string
}

func write(req *http.Request, file io.Writer) error {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if file != nil {
      dump, err := httputil.DumpResponse(res, false)
      if err != nil {
         return err
      }
      fmt.Print(string(dump))
      if _, err := io.Copy(file, res.Body); err != nil {
         return err
      }
   } else {
      dump, err := httputil.DumpResponse(res, true)
      if err != nil {
         return err
      }
      enc := strconv.Encode(dump)
      if strings.HasSuffix(enc, "\n") {
         fmt.Print(enc)
      } else {
         fmt.Println(enc)
      }
   }
   return nil
}

type flags struct {
   golang bool
   https bool
   name string
   output string
}
