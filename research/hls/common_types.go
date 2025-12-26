package hls

import (
   "strings"
   "time"
)

// Key represents encryption info (#EXT-X-KEY or #EXT-X-SESSION-KEY)
type Key struct {
   Method            string
   URI               string
   KeyFormat         string
   KeyFormatVersions string
   IV                string
   Characteristics   string // For session keys
}

// Map represents fMP4 initialization segments (#EXT-X-MAP)
type Map struct {
   URI string
}

// DateRange represents metadata time spans (#EXT-X-DATERANGE)
type DateRange struct {
   ID        string
   Class     string
   StartDate time.Time
   EndDate   time.Time
   Cue       string
   AssetList string
}

func parseKey(line string) *Key {
   prefix := "#EXT-X-KEY:"
   if strings.HasPrefix(line, "#EXT-X-SESSION-KEY:") {
      prefix = "#EXT-X-SESSION-KEY:"
   }
   attrs := parseAttributes(line, prefix)
   return &Key{
      Method:            attrs["METHOD"],
      URI:               attrs["URI"],
      KeyFormat:         attrs["KEYFORMAT"],
      KeyFormatVersions: attrs["KEYFORMATVERSIONS"],
      IV:                attrs["IV"],
      Characteristics:   attrs["CHARACTERISTICS"],
   }
}

func parseMap(line string) *Map {
   attrs := parseAttributes(line, "#EXT-X-MAP:")
   return &Map{
      URI: attrs["URI"],
   }
}

func parseDateRange(line string) *DateRange {
   attrs := parseAttributes(line, "#EXT-X-DATERANGE:")
   dateRange := &DateRange{
      ID:        attrs["ID"],
      Class:     attrs["CLASS"],
      Cue:       attrs["CUE"],
      AssetList: attrs["X-ASSET-LIST"],
   }
   // Swallow errors for simplicity
   dateRange.StartDate, _ = time.Parse(time.RFC3339, attrs["START-DATE"])
   if dateString, ok := attrs["END-DATE"]; ok {
      dateRange.EndDate, _ = time.Parse(time.RFC3339, dateString)
   }
   return dateRange
}
