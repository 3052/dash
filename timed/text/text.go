package text

import (
   "html"
   "encoding/xml"
   "strings"
)

const WebVtt = "WEBVTT"

type Markup struct {
   Body struct {
      Div  struct {
         P []Paragraph `xml:"p"`
      } `xml:"div"`
   } `xml:"body"`
}

func (m *Markup) Unmarshal(b []byte) error {
   return xml.Unmarshal(b, m)
}

type Paragraph struct {
   Begin  string `xml:"begin,attr"`
   End    string `xml:"end,attr"`
   Text   string `xml:",innerxml"`
}

func (p Paragraph) String() string {
   p.Text = html.UnescapeString(p.Text)
   p.Text = strings.ReplaceAll(p.Text, "<br />", "\n")
   var b strings.Builder
   b.WriteByte('\n')
   b.WriteString(p.Begin)
   b.WriteString(" --> ")
   b.WriteString(p.End)
   b.WriteByte('\n')
   b.WriteString(p.Text)
   return b.String()
}
