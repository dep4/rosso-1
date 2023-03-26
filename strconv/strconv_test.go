package strconv

import (
   "bytes"
   "net/http"
   "net/http/httputil"
   "testing"
)

func Test_Binary(t *testing.T) {
   src, dst, err := round_trip(binary)
   if err != nil {
      t.Fatal(err)
   }
   if !bytes.Equal(src, dst) {
      t.Fatal(dst)
   }
}

func round_trip(s string) ([]byte, []byte, error) {
   res, err := http.Get(s)
   if err != nil {
      return nil, nil, err
   }
   src, err := httputil.DumpResponse(res, true)
   if err != nil {
      return nil, nil, err
   }
   if err := res.Body.Close(); err != nil {
      return nil, nil, err
   }
   dst, err := decode(Encode(src))
   if err != nil {
      return nil, nil, err
   }
   return src, dst, nil
}

const (
   binary = "https://picsum.photos/1"
   text = "http://httpbin.org/get"
)

func Test_Text(t *testing.T) {
   src, dst, err := round_trip(text)
   if err != nil {
      t.Fatal(err)
   }
   if !bytes.Equal(src, dst) {
      t.Fatal(dst)
   }
}
