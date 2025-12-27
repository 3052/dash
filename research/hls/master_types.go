package hls

import (
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

type MasterPlaylist struct {
   Variants    []Variant
   Medias      []Rendition // #EXT-X-MEDIA
   SessionKeys []Key       // #EXT-X-SESSION-KEY
}

// ResolveURIs converts relative URLs to absolute URLs using the baseURL.
func (mp *MasterPlaylist) ResolveURIs(baseURL string) error {
   base, err := url.Parse(baseURL)
   if err != nil {
      return err
   }

   for i := range mp.Variants {
      if mp.Variants[i].URI != nil {
         mp.Variants[i].URI = base.ResolveReference(mp.Variants[i].URI)
      }
   }
   for i := range mp.Medias {
      if mp.Medias[i].URI != nil {
         mp.Medias[i].URI = base.ResolveReference(mp.Medias[i].URI)
      }
   }
   for i := range mp.SessionKeys {
      mp.SessionKeys[i].resolve(base)
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

func (v Variant) String() string {
   uriStr := "<nil>"
   if v.URI != nil {
      uriStr = v.URI.String()
   }
   return fmt.Sprintf("Variant{Bandwidth: %d, Resolution: %s, Codecs: %q, URI: %s}", v.Bandwidth, v.Resolution, v.Codecs, uriStr)
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

func parseMaster(lines []string) *MasterPlaylist {
   masterPlaylist := &MasterPlaylist{}

   for i := 0; i < len(lines); i++ {
      line := lines[i]

      if strings.HasPrefix(line, "#EXT-X-MEDIA:") {
         masterPlaylist.Medias = append(masterPlaylist.Medias, parseRendition(line))
      } else if strings.HasPrefix(line, "#EXT-X-SESSION-KEY:") {
         masterPlaylist.SessionKeys = append(masterPlaylist.SessionKeys, *parseKey(line))
      } else if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
         // Parsing a Variant. The URI is on the *next* line.
         variant := parseVariant(line)

         // Check if next line is a URI (doesn't start with # or is empty)
         if i+1 < len(lines) {
            nextLine := lines[i+1]
            if !strings.HasPrefix(nextLine, "#") && nextLine != "" {
               if u, err := url.Parse(nextLine); err == nil {
                  variant.URI = u
               }
               i++ // Advance loop since we consumed the URI line
            }
         }
         masterPlaylist.Variants = append(masterPlaylist.Variants, variant)
      }
   }
   return masterPlaylist
}

func parseVariant(line string) Variant {
   attrs := parseAttributes(line, "#EXT-X-STREAM-INF:")
   bandwidth, _ := strconv.Atoi(attrs["BANDWIDTH"])
   averageBandwidth, _ := strconv.Atoi(attrs["AVERAGE-BANDWIDTH"])

   return Variant{
      Bandwidth:        bandwidth,
      AverageBandwidth: averageBandwidth,
      Codecs:           attrs["CODECS"],
      Resolution:       attrs["RESOLUTION"],
      FrameRate:        attrs["FRAME-RATE"],
      Audio:            attrs["AUDIO"],
      Subtitles:        attrs["SUBTITLES"],
   }
}

func parseRendition(line string) Rendition {
   attrs := parseAttributes(line, "#EXT-X-MEDIA:")

   rendition := Rendition{
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

   if val, ok := attrs["URI"]; ok && val != "" {
      if u, err := url.Parse(val); err == nil {
         rendition.URI = u
      }
   }
   return rendition
}
