package hls

import (
   "strconv"
   "strings"
   "text/scanner"
)

func (s Scanner) Segment() (*Segment, error) {
   var seg Segment
   for s.y.Scan() != scanner.EOF {
      line := s.y.TokenText()
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
         s.x.Init(strings.NewReader(line))
         for s.x.Scan() != scanner.EOF {
            switch s.x.TokenText() {
            case "IV":
               s.x.Scan()
               s.x.Scan()
               seg.Raw_IV = s.x.TokenText()
            case "URI":
               s.x.Scan()
               s.x.Scan()
               seg.Key, err = strconv.Unquote(s.x.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         s.x.Init(strings.NewReader(line))
         for s.x.Scan() != scanner.EOF {
            switch s.x.TokenText() {
            case "URI":
               s.x.Scan()
               s.x.Scan()
               seg.Map, err = strconv.Unquote(s.x.TokenText())
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
