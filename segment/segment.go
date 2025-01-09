package segment

import "41.neocities.org/dash/url"

type List struct {
   Initialization struct {
      SourceUrl url.Url `xml:"sourceURL,attr"`
   }
   SegmentUrl []struct {
      Media url.Url `xml:"media,attr"`
   }
}
