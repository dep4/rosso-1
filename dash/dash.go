package dash

import (
   "strconv"
   "strings"
)

func (r Representation) String() string {
   b := r.Marshal_Indent("\t")
   return string(b)
}

func (r Representation) Marshal_Indent(indent string) []byte {
   var b []byte
   b = append(b, "id: "...)
   b = append(b, r.Id...)
   if r.Width >= 1 {
      b = append(b, '\n')
      b = append(b, indent...)
      b = append(b, "width: "...)
      b = strconv.AppendInt(b, r.Width, 10)
   }
   if r.Height >= 1 {
      b = append(b, '\n')
      b = append(b, indent...)
      b = append(b, "height: "...)
      b = strconv.AppendInt(b, r.Height, 10)
   }
   if r.Bandwidth >= 1 {
      b = append(b, '\n')
      b = append(b, indent...)
      b = append(b, "bandwidth: "...)
      b = strconv.AppendInt(b, r.Bandwidth, 10)
   }
   if r.Codecs != "" {
      b = append(b, '\n')
      b = append(b, indent...)
      b = append(b, "codecs: "...)
      b = append(b, r.Codecs...)
   }
   b = append(b, '\n')
   b = append(b, indent...)
   b = append(b, "MIME type: "...)
   b = append(b, r.MIME_Type...)
   if r.Adaptation.Role != nil {
      b = append(b, '\n')
      b = append(b, indent...)
      b = append(b, "role: "...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   if r.Adaptation.Lang != "" {
      b = append(b, '\n')
      b = append(b, indent...)
      b = append(b, "lang: "...)
      b = append(b, r.Adaptation.Lang...)
   }
   return b
}

type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Role *struct {
      Value string `xml:"value,attr"`
   }
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   Representation []Representation
}

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Height int64 `xml:"height,attr"`
   Id string `xml:"id,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   Width int64 `xml:"width,attr"`
}

func (r Representation) Widevine() *Content_Protection {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

type Content_Protection struct {
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
}

func (p Presentation) Representation() Representations {
   var reps []Representation
   for i, ada := range p.Period.Adaptation_Set {
      for _, rep := range ada.Representation {
         rep.Adaptation = &p.Period.Adaptation_Set[i]
         if rep.Codecs == "" {
            rep.Codecs = ada.Codecs
         }
         if rep.Content_Protection == nil {
            rep.Content_Protection = ada.Content_Protection
         }
         if rep.MIME_Type == "" {
            rep.MIME_Type = ada.MIME_Type
         }
         if rep.Segment_Template == nil {
            rep.Segment_Template = ada.Segment_Template
         }
         reps = append(reps, rep)
      }
   }
   return reps
}

func (r Representation) Initialization() string {
   return r.replace_ID(r.Segment_Template.Initialization)
}

func (r Representation) Media() []string {
   var start int
   if r.Segment_Template.Start_Number != nil {
      start = *r.Segment_Template.Start_Number
   }
   var refs []string
   for _, seg := range r.Segment_Template.Segment_Timeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := r.replace_ID(r.Segment_Template.Media)
         if r.Segment_Template.Start_Number != nil {
            ref = strings.Replace(ref, "$Number$", seg.Time(), 1)
            seg.T++
            start++
         } else {
            ref = strings.Replace(ref, "$Time$", seg.Time(), 1)
            seg.T += seg.D
            start += seg.D
         }
         refs = append(refs, ref)
      }
   }
   return refs
}

func (r Representations) Filter(f func(Representation) bool) Representations {
   var carry []Representation
   for _, item := range r {
      if f(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (r Representations) Video() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MIME_Type == "video/mp4"
   })
}

func (r Representations) Audio() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MIME_Type == "audio/mp4"
   })
}

func (r Representations) Index(f func(a, b Representation) bool) int {
   carry := -1
   for i, item := range r {
      if carry == -1 || f(r[carry], item) {
         carry = i
      }
   }
   return carry
}

func (r Representations) Bandwidth(v int64) int {
   distance := func(a Representation) int64 {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return r.Index(func(carry, item Representation) bool {
      return distance(item) < distance(carry)
   })
}

type Segment struct {
   D int `xml:"d,attr"` // duration
   R int `xml:"r,attr"` // repeat
   T int `xml:"t,attr"` // time
}

func (s Segment) Time() string {
   return strconv.Itoa(s.T)
}

type Segment_Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   Segment_Timeline struct {
      S []Segment
   } `xml:"SegmentTimeline"`
   Start_Number *int `xml:"startNumber,attr"`
}

type Representations []Representation

type Presentation struct {
   Period struct {
      Adaptation_Set []Adaptation `xml:"AdaptationSet"`
   }
}

func (r Representation) Ext() string {
   switch r.MIME_Type {
   case "video/mp4":
      return ".m4v"
   case "audio/mp4":
      return ".m4a"
   }
   return ""
}

func (r Representation) Role() string {
   if r.Adaptation.Role == nil {
      return ""
   }
   return r.Adaptation.Role.Value
}

func (r Representation) replace_ID(s string) string {
   return strings.Replace(s, "$RepresentationID$", r.Id, 1)
}
