package hls

import (
   "errors"
   "strings"
)

var (
   ErrNotMaster = errors.New("content appears to be a media playlist, not a master playlist")
   ErrNotMedia  = errors.New("content appears to be a master playlist, not a media playlist")
)

// Decode detects the playlist type and parses it.
// It returns a MasterPlaylist OR a MediaPlaylist (one will be nil).
func Decode(content string) (*MasterPlaylist, *MediaPlaylist, error) {
   lines := splitLines(content)
   if isMaster(lines) {
      return parseMaster(lines), nil, nil
   }
   return nil, parseMedia(lines), nil
}

// DecodeMaster parses a Master Playlist.
// Returns an error if the content looks like a Media Playlist.
func DecodeMaster(content string) (*MasterPlaylist, error) {
   lines := splitLines(content)
   if !isMaster(lines) {
      return nil, ErrNotMaster
   }
   return parseMaster(lines), nil
}

// DecodeMedia parses a Media Playlist.
// Returns an error if the content looks like a Master Playlist.
func DecodeMedia(content string) (*MediaPlaylist, error) {
   lines := splitLines(content)
   if isMaster(lines) {
      return nil, ErrNotMedia
   }
   return parseMedia(lines), nil
}

// Helper to split and trim lines
func splitLines(content string) []string {
   rawLines := strings.Split(content, "\n")
   lines := make([]string, 0, len(rawLines))
   for _, raw := range rawLines {
      line := strings.TrimSpace(raw)
      if line != "" {
         lines = append(lines, line)
      }
   }
   return lines
}

// Heuristic to check for Master Playlist tags
func isMaster(lines []string) bool {
   for _, line := range lines {
      if strings.HasPrefix(line, "#EXT-X-STREAM-INF") ||
         strings.HasPrefix(line, "#EXT-X-MEDIA:") ||
         strings.HasPrefix(line, "#EXT-X-SESSION-KEY") {
         return true
      }
   }
   return false
}
