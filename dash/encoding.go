package dash

import "fmt"

type Range struct {
   Start uint64
   End uint64
}

type RawRange string

func (r RawRange) Scan() (*Range, error) {
   var v Range
   _, err := fmt.Sscanf(string(r), "%v-%v", &v.Start, &v.End)
   if err != nil {
      return nil, err
   }
   return &v, nil
}

type Representation struct {
   adaptation_set *adaptation_set
   Bandwidth int64 `xml:"bandwidth,attr"`
   ID string `xml:"id,attr"`
   // this might not exist
   BaseURL string
   // this might not exist, or might be under AdaptationSet
   Codecs string `xml:"codecs,attr"`
   // this might be under AdaptationSet
   ContentProtection []ContentProtection
   // this might not exist
   Height int64 `xml:"height,attr"`
   // this might be under AdaptationSet
   MimeType string `xml:"mimeType,attr"`
   // this might not exist
   SegmentBase *struct {
      Initialization struct {
         Range RawRange `xml:"range,attr"`
      }
      IndexRange RawRange `xml:"indexRange,attr"`
   }
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width int64 `xml:"width,attr"`
}

type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // this might not exist
   DefaultKid string `xml:"default_KID,attr"`
   // this might not exist
   PSSH string `xml:"pssh"`
}

