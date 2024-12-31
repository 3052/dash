package dash

func (s *SegmentTemplate) set() {
   // dashif.org/Guidelines-TimingModel#addressing-simple
   if s.StartNumber == nil {
      value := 1
      s.StartNumber = &value
   }
   // dashif.org/Guidelines-TimingModel#timing-sampletimeline
   if s.Timescale == nil {
      value := 1
      s.Timescale = &value
   }
}

type SegmentTemplate struct {
   Media                  Media          `xml:"media,attr"`
   Initialization         Initialization `xml:"initialization,attr"`
   Duration               int            `xml:"duration,attr"`
   PresentationTimeOffset int            `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   Timescale   *int `xml:"timescale,attr"`
}
