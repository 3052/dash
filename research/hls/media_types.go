package hls

import (
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

type MediaPlaylist struct {
   TargetDuration int
   MediaSequence  int
   Version        int
   PlaylistType   string
   Segments       []*Segment
   Keys           []*Key       // Global keys or keys rotating
   DateRanges     []*DateRange // Interstitials/Ads
   EndList        bool
}

// ResolveURIs converts relative URLs to absolute URLs using the baseURL.
func (mp *MediaPlaylist) ResolveURIs(baseURL string) error {
   base, err := url.Parse(baseURL)
   if err != nil {
      return err
   }

   for _, keyItem := range mp.Keys {
      keyItem.resolve(base)
   }

   for _, segmentItem := range mp.Segments {
      segmentItem.resolve(base)
   }
   return nil
}

type Segment struct {
   URI      *url.URL
   Duration float64
   Title    string
   Key      *Key // Encrypt key specific to this segment (if any)
   Map      *Map // Init segment specific to this segment (if any)
}

// resolve updates the Segment's URI and its nested Key/Map URIs to be absolute.
func (s *Segment) resolve(base *url.URL) {
   if s.URI != nil {
      s.URI = base.ResolveReference(s.URI)
   }
   if s.Key != nil {
      s.Key.resolve(base)
   }
   if s.Map != nil {
      s.Map.resolve(base)
   }
}

func parseMedia(lines []string) (*MediaPlaylist, error) {
   mediaPlaylist := &MediaPlaylist{}

   // State trackers
   var currentKey *Key
   var currentMap *Map

   for i := 0; i < len(lines); i++ {
      line := lines[i]

      switch {
      case strings.HasPrefix(line, "#EXT-X-VERSION:"):
         val, err := strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-VERSION:"))
         if err != nil {
            return nil, fmt.Errorf("invalid EXT-X-VERSION: %w", err)
         }
         mediaPlaylist.Version = val

      case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
         val, err := strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-TARGETDURATION:"))
         if err != nil {
            return nil, fmt.Errorf("invalid EXT-X-TARGETDURATION: %w", err)
         }
         mediaPlaylist.TargetDuration = val

      case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
         val, err := strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"))
         if err != nil {
            return nil, fmt.Errorf("invalid EXT-X-MEDIA-SEQUENCE: %w", err)
         }
         mediaPlaylist.MediaSequence = val

      case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
         mediaPlaylist.PlaylistType = strings.TrimPrefix(line, "#EXT-X-PLAYLIST-TYPE:")

      case strings.HasPrefix(line, "#EXT-X-ENDLIST"):
         mediaPlaylist.EndList = true

      case strings.HasPrefix(line, "#EXT-X-DATERANGE:"):
         mediaPlaylist.DateRanges = append(mediaPlaylist.DateRanges, parseDateRange(line))

      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         newKey := parseKey(line)
         mediaPlaylist.Keys = append(mediaPlaylist.Keys, newKey)
         currentKey = newKey // Apply this key to subsequent segments

      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         currentMap = parseMap(line)

      case strings.HasPrefix(line, "#EXTINF:"):
         // Parse duration and title
         // Format: #EXTINF:duration,[title]
         raw := strings.TrimPrefix(line, "#EXTINF:")
         durationStr, title, _ := strings.Cut(raw, ",")

         duration, err := strconv.ParseFloat(durationStr, 64)
         if err != nil {
            return nil, fmt.Errorf("invalid EXTINF duration: %w", err)
         }

         newSegment := &Segment{
            Key:      currentKey,
            Map:      currentMap,
            Duration: duration,
            Title:    strings.TrimSpace(title),
         }

         // The URI is on the next line
         if i+1 < len(lines) {
            nextLine := lines[i+1]
            if !strings.HasPrefix(nextLine, "#") && nextLine != "" {
               if parsedURL, err := url.Parse(nextLine); err == nil {
                  newSegment.URI = parsedURL
               }
               i++
            }
         }
         mediaPlaylist.Segments = append(mediaPlaylist.Segments, newSegment)
      }
   }

   return mediaPlaylist, nil
}
