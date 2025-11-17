package dash

import "encoding/xml"

// MPD represents the Media Presentation Description (MPD) element.
type MPD struct {
   XMLName                   xml.Name `xml:"MPD"`
   Type                      string   `xml:"type,attr"`
   MinBufferTime             string   `xml:"minBufferTime,attr"`
   MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr"`
   Profiles                  string   `xml:"profiles,attr"`
   Period                    *Period  `xml:"Period"`
}
