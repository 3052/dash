package hls

import (
   "strings"
)

// DecodeMaster parses a Master Playlist.
func DecodeMaster(content string) (*MasterPlaylist, error) {
   lines := splitLines(content)
   return parseMaster(lines)
}

// DecodeMedia parses a Media Playlist.
func DecodeMedia(content string) (*MediaPlaylist, error) {
   lines := splitLines(content)
   return parseMedia(lines)
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
