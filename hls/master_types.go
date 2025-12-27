package hls

import (
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

type MasterPlaylist struct {
   Variants    []*Variant
   Medias      []*Rendition
   SessionKeys []*Key
}

// ResolveURIs converts relative URLs to absolute URLs using the base URL.
func (mp *MasterPlaylist) ResolveURIs(base *url.URL) {
   for _, variantItem := range mp.Variants {
      if variantItem.URI != nil {
         variantItem.URI = base.ResolveReference(variantItem.URI)
      }
   }
   for _, renditionItem := range mp.Medias {
      if renditionItem.URI != nil {
         renditionItem.URI = base.ResolveReference(renditionItem.URI)
      }
   }
   for _, keyItem := range mp.SessionKeys {
      keyItem.resolve(base)
   }
}

type Variant struct {
   URI              *url.URL
   Bandwidth        int
   AverageBandwidth int
   Codecs           string
   Resolution       string
   FrameRate        string
   Audio            string
   Subtitles        string
   ID               int
}

// String returns a multi-line summary of the Variant.
func (v *Variant) String() string {
   var builder strings.Builder
   builder.WriteString("bandwidth = ")
   builder.WriteString(strconv.Itoa(v.Bandwidth))
   if v.Resolution != "" {
      builder.WriteString("\nresolution = ")
      builder.WriteString(v.Resolution)
   }
   if v.Codecs != "" {
      builder.WriteString("\ncodecs = ")
      builder.WriteString(strconv.Quote(v.Codecs))
   }
   builder.WriteString("\nid = ")
   builder.WriteString(strconv.Itoa(v.ID))
   return builder.String()
}

type Rendition struct {
   Type            string
   GroupID         string
   Name            string
   Language        string
   URI             *url.URL
   AutoSelect      bool
   Default         bool
   Forced          bool
   Channels        string
   Characteristics string
   ID              int
}

// String returns a multi-line summary of the Rendition.
func (r *Rendition) String() string {
   var builder strings.Builder
   builder.WriteString("type = ")
   builder.WriteString(r.Type)
   if r.Name != "" {
      builder.WriteString("\nname = ")
      builder.WriteString(strconv.Quote(r.Name))
   }
   if r.Language != "" {
      builder.WriteString("\nlang = ")
      builder.WriteString(r.Language)
   }
   if r.GroupID != "" {
      builder.WriteString("\ngroup = ")
      builder.WriteString(r.GroupID)
   }
   builder.WriteString("\nid = ")
   builder.WriteString(strconv.Itoa(r.ID))
   return builder.String()
}

func parseMaster(lines []string) (*MasterPlaylist, error) {
   masterPlaylist := &MasterPlaylist{}
   streamCounter := 0 // Unified counter for both variants and renditions.
   for i := 0; i < len(lines); i++ {
      line := lines[i]
      if strings.HasPrefix(line, "#EXT-X-MEDIA:") {
         rendition := parseRendition(line)
         rendition.ID = streamCounter
         streamCounter++
         masterPlaylist.Medias = append(masterPlaylist.Medias, rendition)
      } else if strings.HasPrefix(line, "#EXT-X-SESSION-KEY:") {
         masterPlaylist.SessionKeys = append(masterPlaylist.SessionKeys, parseKey(line))
      } else if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
         newVariant, err := parseVariant(line)
         if err != nil {
            return nil, err
         }
         if i+1 < len(lines) {
            nextLine := lines[i+1]
            if !strings.HasPrefix(nextLine, "#") && nextLine != "" {
               if parsedURL, err := url.Parse(nextLine); err == nil {
                  newVariant.URI = parsedURL
               }
               i++
            }
         }
         newVariant.ID = streamCounter
         streamCounter++
         masterPlaylist.Variants = append(masterPlaylist.Variants, newVariant)
      }
   }
   return masterPlaylist, nil
}

func parseVariant(line string) (*Variant, error) {
   attrs := parseAttributes(line, "#EXT-X-STREAM-INF:")
   bwStr := attrs["BANDWIDTH"]
   if bwStr == "" {
      return nil, fmt.Errorf("missing required attribute BANDWIDTH")
   }
   bandwidth, err := strconv.Atoi(bwStr)
   if err != nil {
      return nil, fmt.Errorf("invalid BANDWIDTH %q: %w", bwStr, err)
   }
   var averageBandwidth int
   if avgStr := attrs["AVERAGE-BANDWIDTH"]; avgStr != "" {
      average, err := strconv.Atoi(avgStr)
      if err != nil {
         return nil, fmt.Errorf("invalid AVERAGE-BANDWIDTH %q: %w", avgStr, err)
      }
      averageBandwidth = average
   }
   return &Variant{
      Bandwidth:        bandwidth,
      AverageBandwidth: averageBandwidth,
      Codecs:           attrs["CODECS"],
      Resolution:       attrs["RESOLUTION"],
      FrameRate:        attrs["FRAME-RATE"],
      Audio:            attrs["AUDIO"],
      Subtitles:        attrs["SUBTITLES"],
   }, nil
}

func parseRendition(line string) *Rendition {
   attrs := parseAttributes(line, "#EXT-X-MEDIA:")
   newRendition := &Rendition{
      Type:            attrs["TYPE"],
      GroupID:         attrs["GROUP-ID"],
      Name:            attrs["NAME"],
      Language:        attrs["LANGUAGE"],
      Channels:        attrs["CHANNELS"],
      Characteristics: attrs["CHARACTERISTICS"],
      AutoSelect:      attrs["AUTOSELECT"] == "YES",
      Default:         attrs["DEFAULT"] == "YES",
      Forced:          attrs["FORCED"] == "YES",
   }
   if value, ok := attrs["URI"]; ok && value != "" {
      if parsedURL, err := url.Parse(value); err == nil {
         newRendition.URI = parsedURL
      }
   }
   return newRendition
}
