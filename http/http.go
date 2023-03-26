package http

import (
   "2a.pages.dev/rosso/strconv"
   "bufio"
   "bytes"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "os"
   "strings"
   "time"
)

func Get() Request {
   var r Request
   r.Request = new(http.Request) // .Request
   r.Header = make(http.Header) // .Request.Header
   r.Method = "GET" // .Request.Method
   r.URL = new(url.URL) // .Request.URL
   return r
}

func Post() Request {
   var r Request
   r.Request = new(http.Request) // .Request
   r.Header = make(http.Header) // .Request.Header
   r.Method = "POST" // .Request.Method
   r.URL = new(url.URL) // .Request.URL
   return r
}

type Request struct {
   *http.Request
}

func (r Request) Set_Body(body io.Reader) {
   var ok bool
   r.Body, ok = body.(io.ReadCloser)
   if !ok {
      r.Body = io.NopCloser(body)
   }
}

func (r Request) Set_URL(ref string) error {
   var err error
   r.URL, err = url.Parse(ref)
   if err != nil {
      return err
   }
   return nil
}

func Read_Request(r *bufio.Reader) (*http.Request, error) {
   var req http.Request
   text := textproto.NewReader(r)
   // .Method
   raw_method_path, err := text.ReadLine()
   if err != nil {
      return nil, err
   }
   method_path := strings.Fields(raw_method_path)
   req.Method = method_path[0]
   // .URL
   ref, err := url.Parse(method_path[1])
   if err != nil {
      return nil, err
   }
   req.URL = ref
   // .URL.Host
   head, err := text.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   if req.URL.Host == "" {
      req.URL.Host = head.Get("Host")
   }
   // .Header
   req.Header = http.Header(head)
   // .Body
   buf := new(bytes.Buffer)
   length, err := text.R.WriteTo(buf)
   if err != nil {
      return nil, err
   }
   if length >= 1 {
      req.Body = io.NopCloser(buf)
   }
   req.ContentLength = length
   return &req, nil
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
