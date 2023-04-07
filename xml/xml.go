package xml

import (
   "bytes"
   "encoding/xml"
   "errors"
   "io"
)

func Cut(data, sep []byte, v any) error {
   return decode(data, sep, v, false)
}

func Cut_Before(data, sep []byte, v any) error {
   return decode(data, sep, v, true)
}

func decode(data, sep []byte, v any, before bool) error {
   var found bool
   _, data, found = bytes.Cut(data, sep)
   if !found {
      return errors.New("sep not found")
   }
   if before {
      data = append(sep, data...)
   }
   dec := new_decoder(data)
   for {
      _, err := dec.Token()
      if err != nil {
         data = data[:dec.InputOffset()]
         return new_decoder(data).Decode(v)
      }
   }
}

func Indent(dst io.Writer, src io.Reader, prefix, indent string) error {
   dec := xml.NewDecoder(src)
   enc := xml.NewEncoder(dst)
   enc.Indent(prefix, indent)
   for {
      token, err := dec.Token()
      if err == io.EOF {
         return enc.Flush()
      }
      if err != nil {
         return err
      }
      data, ok := token.(xml.CharData)
      if ok {
         token = xml.CharData(bytes.TrimSpace(data))
      }
      if err := enc.EncodeToken(token); err != nil {
         return err
      }
   }
}

func new_decoder(data []byte) *xml.Decoder {
   dec := xml.NewDecoder(bytes.NewReader(data))
   dec.AutoClose = xml.HTMLAutoClose
   dec.Strict = false
   return dec
}
