package hls

import (
   "encoding/base64"
   "errors"
   "net/url"
   "strings"
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

func (key *Key) resolve(base *url.URL) {
   if key.URI != "" {
      key.URI = resolveReference(base, key.URI)
   }
}

// DecodeData extracts and decodes the Base64 data from the URI if it is a Data URI.
func (key *Key) DecodeData() ([]byte, error) {
   if !strings.HasPrefix(key.URI, "data:") {
      return nil, errors.New("URI is not a data URI")
   }

   // The format is data:[<mediatype>][;base64],<data>
   commaIndex := strings.IndexByte(key.URI, ',')
   if commaIndex == -1 {
      return nil, errors.New("invalid data URI: missing comma separator")
   }

   // Ensure it specifies base64 encoding in the metadata prefix
   meta := key.URI[:commaIndex]
   if !strings.Contains(meta, ";base64") {
      return nil, errors.New("data URI does not contain base64 indicator")
   }

   dataString := key.URI[commaIndex+1:]
   return base64.StdEncoding.DecodeString(dataString)
}

// Map represents fMP4 initialization segments (#EXT-X-MAP)
type Map struct {
   URI string
}

func (m *Map) resolve(base *url.URL) {
   if m.URI != "" {
      m.URI = resolveReference(base, m.URI)
   }
}

// DateRange represents metadata time spans (#EXT-X-DATERANGE)
// StartDate and EndDate are stored as raw strings to avoid "time" package dependency.
type DateRange struct {
   ID        string
   Class     string
   StartDate string
   EndDate   string
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
   return &DateRange{
      ID:        attrs["ID"],
      Class:     attrs["CLASS"],
      StartDate: attrs["START-DATE"], // Kept as string
      EndDate:   attrs["END-DATE"],   // Kept as string
      Cue:       attrs["CUE"],
      AssetList: attrs["X-ASSET-LIST"],
   }
}
