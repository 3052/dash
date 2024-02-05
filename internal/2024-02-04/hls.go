package hls

import (
   "strconv"
   "strings"
   "text/scanner"
)

type Media struct {
   Group_ID string
   Type string
   Name string
   Characteristics string
   Raw_URI string
}

type Stream struct {
   Bandwidth int64
   Raw_URI string
   Audio string
   Codecs string
   Resolution string
}

func (s Scanner) Segment() (*Segment, error) {
   var seg Segment
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
      var err error
      switch {
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         seg.URI = append(seg.URI, line)
      case line == "#EXT-X-DISCONTINUITY":
         if seg.Key != "" {
            return &seg, nil
         }
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         seg.URI = nil
         s.Init(strings.NewReader(line))
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               seg.Raw_IV = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               seg.Key, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         s.Init(strings.NewReader(line))
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "URI":
               s.Scan()
               s.Scan()
               seg.Map, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      }
   }
   return &seg, nil
}

type Segment struct {
   Key string
   Map string
   Raw_IV string
   URI []string
}
