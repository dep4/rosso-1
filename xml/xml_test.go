package xml

import (
   "fmt"
   "os"
   "strings"
   "testing"
   
   "encoding/xml"
)

func Test_Cut(t *testing.T) {
   
   var a int
   if err := xml.Unmarshal(nil, &a); err != nil {
      t.Fatal(err)
   }
   
   
   return
   
   data, sep := []byte(dirty), []byte(" world\n")
   var rating regional_rating
   err := Cut(data, sep, &rating)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", rating)
}

func Test_Before(t *testing.T) {
   data, sep := []byte(dirty), []byte("<regionalRating>")
   var rating regional_rating
   err := Cut_Before(data, sep, &rating)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", rating)
}

const clean = `
<regionalRating>
   <rating>TV-PG</rating>
   <region>CA</region>
</regionalRating>
`

func Test_Indent(t *testing.T) {
   err := Indent(os.Stdout, strings.NewReader(clean), "", " ")
   if err != nil {
      t.Fatal(err)
   }
}

const dirty = `
hello world
<regionalRating>
   <rating>TV-PG</rating>
   <region>CA</region>
</regionalRating>
`

type regional_rating struct {
   Rating string `xml:"rating"`
   Region string `xml:"region"`
} 
