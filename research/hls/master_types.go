package hls

import (
   "strconv"
   "strings"
)

type MasterPlaylist struct {
   Variants    []Variant
   Medias      []Rendition // #EXT-X-MEDIA
   SessionKeys []Key       // #EXT-X-SESSION-KEY
}

type Variant struct {
   URI              string
   Bandwidth        int
   AverageBandwidth int
   Codecs           string
   Resolution       string
   FrameRate        string
   Audio            string // ID linking to Media Group
   Subtitles        string // ID linking to Media Group
}

type Rendition struct {
   Type            string // AUDIO, VIDEO, SUBTITLES
   GroupID         string
   Name            string
   Language        string
   URI             string
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
         if i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "#") {
            variant.URI = lines[i+1]
            i++ // Advance loop since we consumed the URI line
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
   return Rendition{
      Type:            attrs["TYPE"],
      GroupID:         attrs["GROUP-ID"],
      Name:            attrs["NAME"],
      Language:        attrs["LANGUAGE"],
      URI:             attrs["URI"],
      Channels:        attrs["CHANNELS"],
      Characteristics: attrs["CHARACTERISTICS"],
      AutoSelect:      attrs["AUTOSELECT"] == "YES",
      Default:         attrs["DEFAULT"] == "YES",
      Forced:          attrs["FORCED"] == "YES",
   }
}
