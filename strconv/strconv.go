package strconv

import (
   "encoding/hex"
   "errors"
   "unicode/utf8"
   "strconv"
)

// mimesniff.spec.whatwg.org#binary-data-byte
func Binary_Data(b byte) bool {
   if b <= 0x08 {
      return true
   }
   if b == 0x0B {
      return true
   }
   if b >= 0x0E && b <= 0x1A {
      return true
   }
   if b >= 0x1C && b <= 0x1F {
      return true
   }
   return false
}

var (
   AppendInt = strconv.AppendInt
   AppendUint = strconv.AppendUint
   Quote = strconv.Quote
)

func Encode(src []byte) string {
   var dst []byte
   for len(src) >= 1 {
      r, size := decode_rune(src)
      s := src[:size]
      if r == utf8.RuneError && size == 1 {
         var d [2]byte
         hex.Encode(d[:], s)
         dst = append(dst, escape_character)
         dst = append(dst, d[:]...)
      } else {
         dst = append(dst, s...)
      }
      src = src[size:]
   }
   return string(dst)
}

const escape_character = '~'

var error_escape = errors.New("invalid printable escape")

func decode(src string) ([]byte, error) {
   var dst []byte
   for len(src) >= 1 {
      s := src[0]
      if s == escape_character {
         if len(src) <= 2 {
            return nil, error_escape
         }
         d, err := hex.DecodeString(src[1:3])
         if err != nil {
            return nil, err
         }
         dst = append(dst, d...)
         src = src[3:]
      } else {
         dst = append(dst, s)
         src = src[1:]
      }
   }
   return dst, nil
}

func decode_rune(p []byte) (rune, int) {
   if len(p) >= 1 {
      b := p[0]
      if b == escape_character || Binary_Data(b) {
         return utf8.RuneError, 1
      }
   }
   return utf8.DecodeRune(p)
}

type unit_measure struct {
   factor float64
   name string
}

type Ordered interface {
   ~float32 | ~float64 |
   ~int | ~int8 | ~int16 | ~int32 | ~int64 |
   ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
