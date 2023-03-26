package protobuf

import (
   "2a.pages.dev/rosso/strconv"
   "google.golang.org/protobuf/encoding/protowire"
)

// If you need fmt.GoStringer with indent, just use `go fmt`.
type Message map[Number]Encoder

func (m Message) Add(num Number, value Message) error {
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = value
   case Message:
      m[num] = Slice[Message]{lvalue, value}
   case Slice[Message]:
      m[num] = append(lvalue, value)
   default:
      return type_error{num, lvalue, value}
   }
   return nil
}

func (m Message) Add_Fixed32(num Number, value uint32) error {
   rvalue := Fixed32(value)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Fixed32:
      m[num] = Slice[Fixed32]{lvalue, rvalue}
   case Slice[Fixed32]:
      m[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) Add_Fixed64(num Number, value uint64) error {
   rvalue := Fixed64(value)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Fixed64:
      m[num] = Slice[Fixed64]{lvalue, rvalue}
   case Slice[Fixed64]:
      m[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) Add_String(num Number, value string) error {
   rvalue := String(value)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case String:
      m[num] = Slice[String]{lvalue, rvalue}
   case Slice[String]:
      m[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) Add_Varint(num Number, value uint64) error {
   rvalue := Varint(value)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Varint:
      m[num] = Slice[Varint]{lvalue, rvalue}
   case Slice[Varint]:
      m[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) Get(num Number) Message {
   switch rvalue := m[num].(type) {
   case Message:
      return rvalue
   case Raw:
      return rvalue.Message
   }
   return nil
}

func (m Message) Get_Bytes(num Number) ([]byte, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return nil, type_error{num, lvalue, rvalue}
   }
   return rvalue.Bytes, nil
}

func (m Message) Get_Fixed64(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Fixed64)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) Get_Messages(num Number) []Message {
   switch rvalue := m[num].(type) {
   case Message:
      return []Message{rvalue}
   case Slice[Message]:
      return rvalue
   case Raw:
      return []Message{rvalue.Message}
   case Slice[Raw]:
      var mes []Message
      for _, raw := range rvalue {
         mes = append(mes, raw.Message)
      }
      return mes
   }
   return nil
}

func (m Message) Get_String(num Number) (string, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return "", type_error{num, lvalue, rvalue}
   }
   return rvalue.String, nil
}

func (m Message) Get_Varint(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Varint)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) consume_fixed32(num Number, buf []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Fixed32(num, val); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func (m Message) consume_fixed64(num Number, buf []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Fixed64(num, val); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func (m Message) consume_raw(num Number, buf []byte) ([]byte, error) {
   var (
      length int
      rvalue Raw
   )
   rvalue.Bytes, length = protowire.ConsumeBytes(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   rvalue.String = strconv.Encode(rvalue.Bytes)
   rvalue.Message, _ = Unmarshal(rvalue.Bytes)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Raw:
      m[num] = Slice[Raw]{lvalue, rvalue}
   case Slice[Raw]:
      m[num] = append(lvalue, rvalue)
   default:
      return nil, type_error{num, lvalue, rvalue}
   }
   return buf[length:], nil
}

func (m Message) consume_varint(num Number, buf []byte) ([]byte, error) {
   val, length := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Varint(num, val); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func (m Message) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

func (Message) get_type() string { return "Message" }

func (t type_error) Error() string {
   get_type := func(enc Encoder) string {
      if enc == nil {
         return "nil"
      }
      return enc.get_type()
   }
   var b []byte
   b = append(b, "field "...)
   b = strconv.AppendInt(b, int64(t.Number), 10)
   b = append(b, " is "...)
   b = append(b, get_type(t.lvalue)...)
   b = append(b, ", not "...)
   b = append(b, get_type(t.rvalue)...)
   return string(b)
}

type Raw struct {
   Bytes []byte
   String string
   Message Message
}

type type_error struct {
   Number
   lvalue Encoder
   rvalue Encoder
}

type Bytes []byte

func (b Bytes) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b)
}

func (Bytes) get_type() string { return "Bytes" }

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

type Fixed32 uint32

func (f Fixed32) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

func (Fixed32) get_type() string { return "Fixed32" }

type Fixed64 uint64

func (f Fixed64) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

func (Fixed64) get_type() string { return "Fixed64" }

type Number = protowire.Number

func (r Raw) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, r.Bytes)
}

func (Raw) get_type() string { return "Raw" }

type Slice[T Encoder] []T

func (s Slice[T]) encode(buf []byte, num Number) []byte {
   for _, value := range s {
      buf = value.encode(buf, num)
   }
   return buf
}

func (Slice[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

type String string

func (s String) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendString(buf, string(s))
}

func (String) get_type() string { return "String" }

type Varint uint64

func (v Varint) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
}

func (Varint) get_type() string { return "Varint" }
