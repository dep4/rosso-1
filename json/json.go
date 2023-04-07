package json

import (
   "bytes"
   "encoding/json"
   "io"
)

// github.com/golang/go/blob/go1.20.3/src/encoding/xml/xml.go
func unmarshal(data, sep []byte, v any, before bool) error {
   _, data, found := bytes.Cut(data, sep)
   if !found {
      return io.EOF
   }
   if before {
      data = append(sep, data...)
   }
   dec := NewDecoder(bytes.NewReader(data))
   for {
      _, err := dec.Token()
      if err != nil {
         data = data[:dec.InputOffset()]
         return json.Unmarshal(data, v)
      }
   }
}

var (
   MarshalIndent = json.MarshalIndent
   NewDecoder = json.NewDecoder
)

func Cut(data, sep []byte, v any) error {
   return unmarshal(data, sep, v, false)
}

func Cut_Before(data, sep []byte, v any) error {
   return unmarshal(data, sep, v, true)
}
