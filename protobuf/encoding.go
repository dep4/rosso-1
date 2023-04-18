package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "errors"
   "sort"
)

// If you need indent, just use this with `go fmt`.
func (m Message) GoString() string {
   b := []byte("protobuf.Message{")
   for num, tok := range m {
      b = fmt.Append(b, num, ":")
      switch tok.(type) {
      case Fixed32:
         b = fmt.Appendf(b, "protobuf.Fixed32(%v)", tok)
      case Fixed64:
         b = fmt.Appendf(b, "protobuf.Fixed64(%v)", tok)
      case String:
         b = fmt.Appendf(b, "protobuf.String(%q)", tok)
      case Varint:
         b = fmt.Appendf(b, "protobuf.Varint(%v)", tok)
      default:
         b = fmt.Appendf(b, "%#v", tok)
      }
      b = append(b, ',')
   }
   b = append(b, "\n}"...)
   return string(b)
}

func (m Raw) GoString() string {
   b := []byte("protobuf.Raw{")
   b = fmt.Appendf(b, "Bytes:%#v", m.Bytes)
   b = fmt.Appendf(b, ",String:%q", m.String)
   b = fmt.Appendf(b, ",Message:%#v", m.Message)
   b = append(b, '}')
   return string(b)
}

func (m Message) Marshal() []byte {
   var nums []Number
   for num := range m {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   var bufs []byte
   for _, num := range nums {
      bufs = m[num].encode(bufs, num)
   }
   return bufs
}

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, length := protowire.ConsumeTag(buf)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      buf = buf[length:]
      switch typ {
      case protowire.VarintType:
         buf, err = mes.consume_varint(num, buf)
      case protowire.Fixed64Type:
         buf, err = mes.consume_fixed64(num, buf)
      case protowire.Fixed32Type:
         buf, err = mes.consume_fixed32(num, buf)
      case protowire.BytesType:
         buf, err = mes.consume_raw(num, buf)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}
