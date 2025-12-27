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

func (k *Key) resolve(base *url.URL) {
   if k.URI != nil {
      k.URI = base.ResolveReference(k.URI)
   }
}

// DecodeData extracts and decodes the Base64 data directly from the URL Opaque field.
func (k *Key) DecodeData() ([]byte, error) {
   if k.URI == nil {
      return nil, errors.New("URI is nil")
   }

   if k.URI.Scheme != "data" {
      return nil, errors.New("URI is not a data URI")
   }

   // For data URIs, net/url stores the content (mime+encoding+data) in Opaque.
   // Format: [<mediatype>][;base64],<data>
   meta, dataString, found := strings.Cut(k.URI.Opaque, ",")
   if !found {
      return nil, errors.New("invalid data URI: missing comma separator")
   }

   // Verify base64 encoding is specified in the metadata (before the comma)
   if !strings.Contains(meta, ";base64") {
      return nil, errors.New("data URI does not contain base64 indicator")
   }

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

   newKey := &Key{
      Method:            attrs["METHOD"],
      KeyFormat:         attrs["KEYFORMAT"],
      KeyFormatVersions: attrs["KEYFORMATVERSIONS"],
      IV:                attrs["IV"],
      Characteristics:   attrs["CHARACTERISTICS"],
   }

   if value, ok := attrs["URI"]; ok && value != "" {
      if parsedURL, err := url.Parse(value); err == nil {
         newKey.URI = parsedURL
      }
   }
   return newKey
}

func parseMap(line string) *Map {
   attrs := parseAttributes(line, "#EXT-X-MAP:")
   segmentMap := &Map{}

   if value, ok := attrs["URI"]; ok && value != "" {
      if parsedURL, err := url.Parse(value); err == nil {
         segmentMap.URI = parsedURL
      }
   }
   return segmentMap
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
