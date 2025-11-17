package dash

import (
   "os"
   "testing"
)

func TestParse(t *testing.T) {
   // The user will provide the rakuten.mpd file.
   // For testing purposes, we assume the file is in the same directory.
   mpdBytes, err := os.ReadFile("rakuten.mpd")
   if err != nil {
      t.Fatalf("reading rakuten.mpd: %v", err)
   }

   mpd, err := Parse(mpdBytes)
   if err != nil {
      t.Fatalf("Parse() error = %v", err)
   }

   if mpd.Type != "static" {
      t.Errorf("expected type 'static', got '%s'", mpd.Type)
   }

   if mpd.Period == nil {
      t.Fatal("expected Period not to be nil")
   }

   if len(mpd.Period.AdaptationSets) == 0 {
      t.Error("expected at least one AdaptationSet")
   }
}
