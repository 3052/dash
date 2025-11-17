package dash

import "encoding/xml"

// Parse takes a byte slice of an MPD file and returns a parsed MPD struct.
func Parse(data []byte) (*MPD, error) {
   var mpd MPD
   err := xml.Unmarshal(data, &mpd)
   if err != nil {
      return nil, err
   }
   return &mpd, nil
}
