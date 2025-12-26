package hls

import "strings"

// Playlist is a container that holds either a Master or Media playlist.
// Check IsMaster to determine which field to access.
type Playlist struct {
   IsMaster bool
   Master   *MasterPlaylist
   Media    *MediaPlaylist
}

// Decode parses the raw string content of an m3u8 file into a Playlist struct.
func Decode(content string) (*Playlist, error) {
   // Split the content by newline.
   // strings.Split is optimized and simpler than manual byte scanning.
   rawLines := strings.Split(content, "\n")

   // Allocate slice for cleaned lines
   lines := make([]string, 0, len(rawLines))

   isMaster := false

   for _, rawLine := range rawLines {
      // Trim whitespace (including \r from Windows line endings)
      line := strings.TrimSpace(rawLine)

      if line == "" {
         continue
      }

      lines = append(lines, line)

      // Detect if this is a Master Playlist based on specific tags
      if !isMaster {
         if strings.HasPrefix(line, "#EXT-X-STREAM-INF") ||
            strings.HasPrefix(line, "#EXT-X-MEDIA:") ||
            strings.HasPrefix(line, "#EXT-X-SESSION-KEY") ||
            strings.HasPrefix(line, "#EXT-X-I-FRAME-STREAM-INF") {
            isMaster = true
         }
      }
   }

   result := &Playlist{
      IsMaster: isMaster,
   }

   // Route to the specific parser functions (defined in master_types.go and media_types.go)
   if isMaster {
      result.Master = parseMaster(lines)
   } else {
      result.Media = parseMedia(lines)
   }

   return result, nil
}
