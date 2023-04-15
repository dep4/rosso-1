package http

import (
   "2a.pages.dev/rosso/strconv"
   "fmt"
   "io"
   "net/http"
   "net/http/httputil"
   "os"
   "strings"
   "time"
)

const StatusFound = http.StatusFound

type Client struct {
   Log_Level int // this needs to work with flag.IntVar
   Status int
   http.Client
}

var Default_Client = Client{
   Client: http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   },
   Log_Level: 1,
   Status: http.StatusOK,
}

func (c Client) Do(req *Request) (*Response, error) {
   switch c.Log_Level {
   case 1:
      fmt.Println(req.Method, req.URL)
   case 2:
      dump, err := httputil.DumpRequest(req.Request, true)
      if err != nil {
         return nil, err
      }
      enc := strconv.Encode(dump)
      if strings.HasSuffix(enc, "\n") {
         fmt.Print(enc)
      } else {
         fmt.Println(enc)
      }
   }
   res, err := c.Client.Do(req.Request)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != c.Status {
      return nil, fmt.Errorf(res.Status)
   }
   return res, nil
}

func (c Client) Get(ref string) (*Response, error) {
   req, err := Get_URL(ref)
   if err != nil {
      return nil, err
   }
   return c.Do(req)
}

type Cookie = http.Cookie

type Header = http.Header

type Progress struct {
   bytes int64
   bytes_read int64
   bytes_written int
   chunks int
   chunks_read int64
   lap time.Time
   total time.Time
   w io.Writer
}

func Progress_Bytes(dst io.Writer, bytes int64) *Progress {
   return &Progress{w: dst, bytes: bytes}
}

func Progress_Chunks(dst io.Writer, chunks int) *Progress {
   return &Progress{w: dst, chunks: chunks}
}

func (p *Progress) Add_Chunk(bytes int64) {
   p.bytes_read += bytes
   p.chunks_read += 1
   p.bytes = int64(p.chunks) * p.bytes_read / p.chunks_read
}

func (p *Progress) Write(data []byte) (int, error) {
   if p.total.IsZero() {
      p.total = time.Now()
      p.lap = time.Now()
   }
   lap := time.Since(p.lap)
   if lap >= time.Second {
      total := time.Since(p.total).Seconds()
      var b []byte
      b = strconv.Ratio(p.bytes_written, p.bytes).Percent(b)
      b = append(b, "   "...)
      b = strconv.New_Number(p.bytes_written).Size(b)
      b = append(b, "   "...)
      b = strconv.Ratio(p.bytes_written, total).Rate(b)
      b = append(b, '\n')
      os.Stderr.Write(b)
      p.lap = p.lap.Add(lap)
   }
   write, err := p.w.Write(data)
   p.bytes_written += write
   return write, err
}

type Response = http.Response
