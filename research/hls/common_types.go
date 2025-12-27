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
   URI               *url.URL
   KeyFormat         string
   KeyFormatVersions string
   IV                string
   Characteristics   string // For session keys
}

func (key *Key) resolve(base *url.URL) {
   if key.URI != nil {
      key.URI = base.ResolveReference(key.URI)
   }
}

// DecodeData extracts and decodes the Base64 data from the URI if it is a Data URI.
func (key *Key) DecodeData() ([]byte, error) {
   if key.URI == nil {
      return nil, errors.New("URI is nil")
   }

   // Reconstruct the string to parse it manually as a data URI string.
   // net/url puts the content in Opaque for opaque schemes like "data",
   // but using String() provides a consistent view.
   uriString := key.URI.String()

   if !strings.HasPrefix(uriString, "data:") {
      return nil, errors.New("URI is not a data URI")
   }

   // The format is data:[<mediatype>][;base64],<data>
   commaIndex := strings.IndexByte(uriString, ',')
   if commaIndex == -1 {
      return nil, errors.New("invalid data URI: missing comma separator")
   }

   // Ensure it specifies base64 encoding in the metadata prefix
   meta := uriString[:commaIndex]
   if !strings.Contains(meta, ";base64") {
      return nil, errors.New("data URI does not contain base64 indicator")
   }

   dataString := uriString[commaIndex+1:]
   return base64.StdEncoding.DecodeString(dataString)
}

// Map represents fMP4 initialization segments (#EXT-X-MAP)
type Map struct {
   URI *url.URL
}

func (m *Map) resolve(base *url.URL) {
   if m.URI != nil {
      m.URI = base.ResolveReference(m.URI)
   }
}

// DateRange represents metadata time spans (#EXT-X-DATERANGE)
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

   k := &Key{
      Method:            attrs["METHOD"],
      KeyFormat:         attrs["KEYFORMAT"],
      KeyFormatVersions: attrs["KEYFORMATVERSIONS"],
      IV:                attrs["IV"],
      Characteristics:   attrs["CHARACTERISTICS"],
   }

   if val, ok := attrs["URI"]; ok && val != "" {
      // Ignore error, keep nil if invalid
      if u, err := url.Parse(val); err == nil {
         k.URI = u
      }
   }
   return k
}

func parseMap(line string) *Map {
   attrs := parseAttributes(line, "#EXT-X-MAP:")
   m := &Map{}

   if val, ok := attrs["URI"]; ok && val != "" {
      if u, err := url.Parse(val); err == nil {
         m.URI = u
      }
   }
   return m
}

func parseDateRange(line string) *DateRange {
   attrs := parseAttributes(line, "#EXT-X-DATERANGE:")
   return &DateRange{
      ID:        attrs["ID"],
      Class:     attrs["CLASS"],
      StartDate: attrs["START-DATE"],
      EndDate:   attrs["END-DATE"],
      Cue:       attrs["CUE"],
      AssetList: attrs["X-ASSET-LIST"],
   }
}
