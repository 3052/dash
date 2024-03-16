package ttml

import (
   "154.pages.dev/sofia"
   "io"
   "encoding/xml"
)

type timed_text struct {
   Body struct {
      Div  struct {
         P    []struct {
            Begin  string `xml:"begin,attr"`
            End    string `xml:"end,attr"`
            Text   string `xml:",chardata"`
         } `xml:"p"`
      } `xml:"div"`
   } `xml:"body"`
}

func (t *timed_text) decode(r io.Reader) error {
   var file sofia.File
   err := file.Decode(r)
   if err != nil {
      return err
   }
   switch length := len(file.MediaData.Data); length {
   case 0:
      panic(0)
   case 1:
   default:
      panic(length)
   }
   return xml.Unmarshal(file.MediaData.Data[0], t)
}
