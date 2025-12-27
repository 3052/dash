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

// ResolveURIs converts relative URLs to absolute URLs using the baseURL.
func (mp *MasterPlaylist) ResolveURIs(baseURL string) error {
   base, err := url.Parse(baseURL)
   if err != nil {
      return err
   }

   // Because we are iterating over pointers, we can modify items directly
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
   return nil
}

type Variant struct {
   URI              *url.URL
   Bandwidth        int
   AverageBandwidth int
   Codecs           string
   Resolution       string
   FrameRate        string
   Audio            string // ID linking to Media Group
   Subtitles        string // ID linking to Media Group
}

func (v *Variant) String() string {
   var builder strings.Builder
   builder.WriteString("Bandwidth: ")
   builder.WriteString(strconv.Itoa(v.Bandwidth))
   builder.WriteString("\nResolution: ")
   builder.WriteString(v.Resolution)
   builder.WriteString("\nCodecs: ")
   builder.WriteString(strconv.Quote(v.Codecs))
   return builder.String()
}

type Rendition struct {
   Type            string // AUDIO, VIDEO, SUBTITLES
   GroupID         string
   Name            string
   Language        string
   URI             *url.URL
   AutoSelect      bool
   Default         bool
   Forced          bool
   Channels        string
   Characteristics string
}

func parseMaster(lines []string) (*MasterPlaylist, error) {
   masterPlaylist := &MasterPlaylist{}

   for i := 0; i < len(lines); i++ {
      line := lines[i]

      if strings.HasPrefix(line, "#EXT-X-MEDIA:") {
         masterPlaylist.Medias = append(masterPlaylist.Medias, parseRendition(line))
      } else if strings.HasPrefix(line, "#EXT-X-SESSION-KEY:") {
         masterPlaylist.SessionKeys = append(masterPlaylist.SessionKeys, parseKey(line))
      } else if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
         // Parsing a Variant. The URI is on the *next* line.
         newVariant, err := parseVariant(line)
         if err != nil {
            return nil, err
         }

         // Check if next line is a URI (doesn't start with # or is empty)
         if i+1 < len(lines) {
            nextLine := lines[i+1]
            if !strings.HasPrefix(nextLine, "#") && nextLine != "" {
               if parsedURL, err := url.Parse(nextLine); err == nil {
                  newVariant.URI = parsedURL
               }
               i++ // Advance loop since we consumed the URI line
            }
         }
         masterPlaylist.Variants = append(masterPlaylist.Variants, newVariant)
      }
   }
   return masterPlaylist, nil
}

func parseVariant(line string) (*Variant, error) {
   attrs := parseAttributes(line, "#EXT-X-STREAM-INF:")

   // BANDWIDTH is required
   bwStr := attrs["BANDWIDTH"]
   if bwStr == "" {
      return nil, fmt.Errorf("missing required attribute BANDWIDTH")
   }
   bandwidth, err := strconv.Atoi(bwStr)
   if err != nil {
      return nil, fmt.Errorf("invalid BANDWIDTH %q: %w", bwStr, err)
   }

   // AVERAGE-BANDWIDTH is optional
   var averageBandwidth int
   if avgStr := attrs["AVERAGE-BANDWIDTH"]; avgStr != "" {
      val, err := strconv.Atoi(avgStr)
      if err != nil {
         return nil, fmt.Errorf("invalid AVERAGE-BANDWIDTH %q: %w", avgStr, err)
      }
      averageBandwidth = val
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
