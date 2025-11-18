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

   // Test for ContentProtection. This assumes the provided rakuten.mpd
   // contains at least one AdaptationSet with a ContentProtection element.
   foundContentProtection := false
   for _, as := range mpd.Period.AdaptationSets {
      if len(as.ContentProtections) > 0 {
         foundContentProtection = true
         cp := as.ContentProtections[0]
         if cp.SchemeIDURI == "" {
            t.Error("expected ContentProtection to have a non-empty schemeIdUri")
         }
         // You could add more specific checks here if you know the URI, for example:
         // if cp.SchemeIDURI != "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         //    t.Errorf("unexpected schemeIdUri: got %s", cp.SchemeIDURI)
         // }
         break
      }
   }

   // This check is to ensure that the test is actually validating the parsing.
   // If your MPD file is not expected to have content protection, you can remove this.
   if !foundContentProtection {
      t.Log("Warning: No ContentProtection elements found in the provided MPD to test against.")
   }
}
