package dash

import (
   "os"
   "testing"
)

func TestParse(t *testing.T) {
   filename := "rakuten.mpd"
   data, err := os.ReadFile(filename)
   if err != nil {
      t.Fatalf("failed to read %s: %v", filename, err)
   }

   mpd, err := Parse(data)
   if err != nil {
      t.Fatalf("failed to parse MPD: %v", err)
   }

   if mpd == nil {
      t.Fatal("returned MPD is nil")
   }

   t.Logf("MPD Type: %s", mpd.Type)
   t.Logf("MPD Duration: %s", mpd.MediaPresentationDuration)

   for i, period := range mpd.Periods {
      t.Logf("Period %d ID: %s", i, period.ID)
      for j, adaptationSet := range period.AdaptationSets {
         t.Logf("  AdaptationSet %d MimeType: %s", j, adaptationSet.MimeType)
         for k, representation := range adaptationSet.Representations {
            t.Logf("    Representation %d ID: %s, Bandwidth: %d", k, representation.ID, representation.Bandwidth)
         }
      }
   }
}
