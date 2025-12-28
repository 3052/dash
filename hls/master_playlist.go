package hls

import (
   "fmt"
   "net/url"
   "sort"
   "strconv"
   "strings"
)

// VariantInfo holds the metadata from a single #EXT-X-STREAM-INF tag.
// It describes one possible way to play the content from its parent Stream's URI.
type VariantInfo struct {
   Bandwidth        int
   AverageBandwidth int
   Codecs           string
   Resolution       string
   FrameRate        string
   Audio            string // Refers to a Rendition GROUP-ID for audio
   Subtitles        string // Refers to a Rendition GROUP-ID for subtitles
}

// String returns a multi-line summary of the VariantInfo.
func (vi *VariantInfo) String() string {
   var builder strings.Builder
   builder.WriteString("  bandwidth = ")
   builder.WriteString(strconv.Itoa(vi.Bandwidth))
   if vi.Codecs != "" {
      builder.WriteString("\n  codecs = ")
      builder.WriteString(vi.Codecs)
   }
   if vi.Audio != "" {
      builder.WriteString("\n  audio_group = ")
      builder.WriteString(vi.Audio)
   }
   return builder.String()
}

// Stream represents a single media playlist (URI) and all the different
// variant renditions (#EXT-X-STREAM-INF tags) that point to it.
type Stream struct {
   URI      *url.URL
   Variants []*VariantInfo
   ID       int // Unique ID within the manifest
}

// String returns a multi-line summary of the Stream and its nested variants.
func (s *Stream) String() string {
   var builder strings.Builder
   builder.WriteString(fmt.Sprintf("Stream (ID: %d)", s.ID))
   if s.URI != nil {
      builder.WriteString("\nURI = ")
      builder.WriteString(s.URI.String())
   }
   for i, v := range s.Variants {
      builder.WriteString(fmt.Sprintf("\n- Variant %d:\n", i+1))
      builder.WriteString(v.String())
   }
   return builder.String()
}

// MaxBandwidth returns the highest bandwidth found among its variants.
func (s *Stream) MaxBandwidth() int {
   max := 0
   for _, v := range s.Variants {
      if v.Bandwidth > max {
         max = v.Bandwidth
      }
   }
   return max
}

type MasterPlaylist struct {
   Streams     []*Stream    // Changed from Variants
   Medias      []*Rendition // Audio, Subtitle, etc. renditions
   SessionKeys []*Key
}

// ResolveURIs converts relative URLs to absolute URLs using the base URL.
func (mp *MasterPlaylist) ResolveURIs(base *url.URL) {
   for _, streamItem := range mp.Streams {
      if streamItem.URI != nil {
         streamItem.URI = base.ResolveReference(streamItem.URI)
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

// Sort sorts the Streams and Medias slices in place.
// Streams are sorted by their maximum bandwidth (ascending).
// Renditions (Medias) are sorted by GroupID.
func (mp *MasterPlaylist) Sort() {
   sort.Slice(mp.Streams, func(i, j int) bool {
      return mp.Streams[i].MaxBandwidth() < mp.Streams[j].MaxBandwidth()
   })
   sort.Slice(mp.Medias, func(i, j int) bool {
      return mp.Medias[i].GroupID < mp.Medias[j].GroupID
   })
}

// Rendition represents an #EXT-X-MEDIA tag.
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
      builder.WriteString(r.Name)
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
   streamCounter := 0
   streamMap := make(map[string]*Stream) // Map URL to Stream to handle grouping

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
         variantInfo, err := parseVariantInfo(line)
         if err != nil {
            return nil, err
         }

         // The URI MUST be on the next line
         if i+1 >= len(lines) {
            continue // Malformed, missing URI
         }
         i++
         uriLine := lines[i]

         // Check if we've already seen this URI
         if stream, ok := streamMap[uriLine]; ok {
            // URI exists, append this new variant info to the existing stream
            stream.Variants = append(stream.Variants, variantInfo)
         } else {
            // First time seeing this URI, create a new Stream
            newStream := &Stream{ID: streamCounter}
            streamCounter++
            if parsedURL, err := url.Parse(uriLine); err == nil {
               newStream.URI = parsedURL
            }
            newStream.Variants = append(newStream.Variants, variantInfo)
            streamMap[uriLine] = newStream
            masterPlaylist.Streams = append(masterPlaylist.Streams, newStream)
         }
      }
   }
   return masterPlaylist, nil
}

func parseVariantInfo(line string) (*VariantInfo, error) {
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
   return &VariantInfo{
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
