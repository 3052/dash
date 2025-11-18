package dash

import "encoding/xml"

// Representation represents the Representation element.
type Representation struct {
   XMLName           xml.Name     `xml:"Representation"`
   ID                string       `xml:"id,attr"`
   Bandwidth         int          `xml:"bandwidth,attr"`
   Codecs            string       `xml:"codecs,attr"`
   Width             int          `xml:"width,attr,omitempty"`
   Height            int          `xml:"height,attr,omitempty"`
   FrameRate         string       `xml:"frameRate,attr,omitempty"`
   AudioSamplingRate string       `xml:"audioSamplingRate,attr,omitempty"`
   BaseURL           string       `xml:"BaseURL,omitempty"`
   SegmentList       *SegmentList `xml:"SegmentList"`
}
