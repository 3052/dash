package segment

type List struct {
   Initialization struct {
      SourceUrl Url `xml:"sourceURL,attr"`
   }
   SegmentUrl []struct {
      Media Url `xml:"media,attr"`
   }
}

type Url string
