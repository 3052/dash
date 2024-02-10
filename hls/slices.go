package hls

import (
   "strconv"
   "strings"
)

type Media struct {
   GroupId string
   Type string
   Name string
   Characteristics string
   RawUri string
}

func (m Media) String() string {
   var b strings.Builder
   b.WriteString("group ID: ")
   b.WriteString(m.GroupId)
   b.WriteString("\ntype: ")
   b.WriteString(m.Type)
   b.WriteString("\nname: ")
   b.WriteString(m.Name)
   if m.Characteristics != "" {
      b.WriteString("\ncharacteristics: ")
      b.WriteString(m.Characteristics)
   }
   return b.String()
}

type Stream struct {
   Bandwidth int64
   RawUri string
   Audio string
   Codecs string
   Resolution string
}

func (m Stream) String() string {
   var b []byte
   if m.Resolution != "" {
      b = append(b, "resolution: "...)
      b = append(b, m.Resolution...)
      b = append(b, '\n')
   }
   b = append(b, "bandwidth: "...)
   b = strconv.AppendInt(b, m.Bandwidth, 10)
   if m.Codecs != "" {
      b = append(b, "\ncodecs: "...)
      b = append(b, m.Codecs...)
   }
   if m.Audio != "" {
      b = append(b, "\naudio: "...)
      b = append(b, m.Audio...)
   }
   return string(b)
}
