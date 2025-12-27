package hls

import (
   "net/url"
   "strconv"
   "strings"
)

type MediaPlaylist struct {
   TargetDuration int
   MediaSequence  int
   Version        int
   PlaylistType   string
   Segments       []Segment
   Keys           []Key       // Global keys or keys rotating
   DateRanges     []DateRange // Interstitials/Ads
   EndList        bool
}

// ResolveURIs converts relative URLs to absolute URLs using the baseURL.
func (mp *MediaPlaylist) ResolveURIs(baseURL string) error {
   base, err := url.Parse(baseURL)
   if err != nil {
      return err
   }

   for i := range mp.Keys {
      mp.Keys[i].resolve(base)
   }
   for i := range mp.Segments {
      if mp.Segments[i].URI != "" {
         mp.Segments[i].URI = resolveReference(base, mp.Segments[i].URI)
      }
      if mp.Segments[i].Key != nil {
         mp.Segments[i].Key.resolve(base)
      }
      if mp.Segments[i].Map != nil {
         mp.Segments[i].Map.resolve(base)
      }
   }
   return nil
}

type Segment struct {
   URI      string
   Duration float64
   Title    string
   Key      *Key // Encrypt key specific to this segment (if any)
   Map      *Map // Init segment specific to this segment (if any)
}

func parseMedia(lines []string) *MediaPlaylist {
   mediaPlaylist := &MediaPlaylist{}

   // State trackers
   var currentKey *Key
   var currentMap *Map

   for i := 0; i < len(lines); i++ {
      line := lines[i]

      switch {
      case strings.HasPrefix(line, "#EXT-X-VERSION:"):
         mediaPlaylist.Version, _ = strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-VERSION:"))

      case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
         mediaPlaylist.TargetDuration, _ = strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-TARGETDURATION:"))

      case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
         mediaPlaylist.MediaSequence, _ = strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"))

      case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
         mediaPlaylist.PlaylistType = strings.TrimPrefix(line, "#EXT-X-PLAYLIST-TYPE:")

      case strings.HasPrefix(line, "#EXT-X-ENDLIST"):
         mediaPlaylist.EndList = true

      case strings.HasPrefix(line, "#EXT-X-DATERANGE:"):
         mediaPlaylist.DateRanges = append(mediaPlaylist.DateRanges, *parseDateRange(line))

      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         key := parseKey(line)
         mediaPlaylist.Keys = append(mediaPlaylist.Keys, *key)
         currentKey = key // Apply this key to subsequent segments

      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         currentMap = parseMap(line)

      case strings.HasPrefix(line, "#EXTINF:"):
         // Parse duration and title
         raw := strings.TrimPrefix(line, "#EXTINF:")
         parts := strings.SplitN(raw, ",", 2)

         segment := Segment{
            Key: currentKey,
            Map: currentMap,
         }

         if len(parts) > 0 {
            segment.Duration, _ = strconv.ParseFloat(parts[0], 64)
         }
         if len(parts) > 1 {
            segment.Title = strings.TrimSpace(parts[1])
         }

         // The URI is on the next line
         if i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "#") {
            segment.URI = lines[i+1]
            i++
         }
         mediaPlaylist.Segments = append(mediaPlaylist.Segments, segment)
      }
   }

   return mediaPlaylist
}
