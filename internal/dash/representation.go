package dash

import "encoding/xml"

// Representation represents a Representation element in the MPD.
type Representation struct {
   XMLName                   xml.Name                     `xml:"Representation"`
   ID                        *string                      `xml:"id,attr,omitempty"`
   Bandwidth                 *uint64                      `xml:"bandwidth,attr,omitempty"`
   Codecs                    *string                      `xml:"codecs,attr,omitempty"`
   FrameRate                 *string                      `xml:"frameRate,attr,omitempty"`
   Height                    *uint64                      `xml:"height,attr,omitempty"`
   Width                     *uint64                      `xml:"width,attr,omitempty"`
   ScanType                  *string                      `xml:"scanType,attr,omitempty"`
   AudioSamplingRate         *string                      `xml:"audioSamplingRate,attr,omitempty"`
   BaseURL                   []*BaseURL                   `xml:"BaseURL,omitempty"`
   SegmentBase               *SegmentBase                 `xml:"SegmentBase,omitempty"`
   AudioChannelConfiguration []*AudioChannelConfiguration `xml:"AudioChannelConfiguration,omitempty"`
}

// AudioChannelConfiguration holds the value for the audio channel configuration.
type AudioChannelConfiguration struct {
   SchemeIDURI *string `xml:"schemeIdUri,attr,omitempty"`
   Value       *string `xml:"value,attr,omitempty"`
}
