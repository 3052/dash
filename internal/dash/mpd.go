package dash

import "encoding/xml"

// MPD represents the Media Presentation Description element.
type MPD struct {
   XMLName                   xml.Name   `xml:"MPD"`
   Type                      *string    `xml:"type,attr,omitempty"`
   MediaPresentationDuration *string    `xml:"mediaPresentationDuration,attr,omitempty"`
   MinBufferTime             *string    `xml:"minBufferTime,attr,omitempty"`
   Profiles                  *string    `xml:"profiles,attr,omitempty"`
   BaseURL                   []*BaseURL `xml:"BaseURL,omitempty"`
   Periods                   []*Period  `xml:"Period,omitempty"`
}

// BaseURL represents the BaseURL element.
type BaseURL struct {
   XMLName xml.Name `xml:"BaseURL"`
   URL     string   `xml:",innerxml"`
}
