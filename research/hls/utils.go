package hls

import (
   "strings"
)

// parseAttributes parses HLS attribute lists (e.g., KEY="VAL",KEY2=VAL).
// It handles quoted strings containing commas.
func parseAttributes(line string, tagPrefix string) map[string]string {
   line = strings.TrimPrefix(line, tagPrefix)
   attributes := make(map[string]string)

   var keyBuilder, valBuilder strings.Builder
   inKey := true
   inQuote := false

   for i := 0; i < len(line); i++ {
      char := line[i]

      if inKey {
         if char == '=' {
            inKey = false
         } else {
            keyBuilder.WriteByte(char)
         }
      } else {
         // Inside the value part
         if char == '"' {
            inQuote = !inQuote
            continue // Skip the actual quote character
         }

         // If we hit a comma and we are NOT in a quote, it's the end of the pair
         if char == ',' && !inQuote {
            keyString := strings.TrimSpace(keyBuilder.String())
            valueString := valBuilder.String()
            attributes[keyString] = valueString

            // Reset
            keyBuilder.Reset()
            valBuilder.Reset()
            inKey = true
         } else {
            valBuilder.WriteByte(char)
         }
      }
   }

   // Flush the final pair
   if keyBuilder.Len() > 0 {
      keyString := strings.TrimSpace(keyBuilder.String())
      valueString := valBuilder.String()
      attributes[keyString] = valueString
   }

   return attributes
}
