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

   if len(mpd.Periods) == 0 {
      t.Fatal("expected at least one Period")
   }

   // Flags to ensure we test the parsing logic at least once if the elements exist.
   foundAdaptationSet := false
   foundContentProtection := false
   foundSegmentList := false
   var totalRepresentationsWithID int

   // Iterate through all periods to check their contents.
   for _, period := range mpd.Periods {
      if len(period.AdaptationSets) > 0 {
         foundAdaptationSet = true
      }

      for _, as := range period.AdaptationSets {
         if len(as.ContentProtections) > 0 {
            foundContentProtection = true
            if as.ContentProtections[0].SchemeIDURI == "" {
               t.Error("expected ContentProtection to have a non-empty schemeIdUri")
            }
         }

         for _, rep := range as.Representations {
            if rep.ID != "" {
               totalRepresentationsWithID++
            }
            if rep.SegmentList != nil {
               foundSegmentList = true
               if len(rep.SegmentList.SegmentURLs) == 0 {
                  t.Error("expected SegmentList to have at least one SegmentURL")
               }
               if rep.SegmentList.SegmentURLs[0].Media == "" {
                  t.Error("expected SegmentURL to have a non-empty media attribute")
               }
            }
         }
      }
   }

   if !foundAdaptationSet {
      t.Error("expected at least one AdaptationSet in at least one Period")
   }

   if !foundContentProtection {
      t.Log("Warning: No ContentProtection elements found in the provided MPD to test against.")
   }
   if !foundSegmentList {
      t.Log("Warning: No SegmentList elements found in the provided MPD to test against.")
   }

   // Test the RepresentationsByID method
   repsByID := mpd.RepresentationsByID()
   if repsByID == nil {
      t.Fatal("RepresentationsByID() returned a nil map")
   }

   // Count the total number of representations stored in the map.
   var mapRepCount int
   for _, reps := range repsByID {
      mapRepCount += len(reps)
   }

   if totalRepresentationsWithID > 0 && mapRepCount == 0 {
      t.Error("Representations with IDs were found, but the map is empty")
   }

   if totalRepresentationsWithID != mapRepCount {
      t.Errorf("Mismatch in representation count: found %d with IDs, but map contains %d", totalRepresentationsWithID, mapRepCount)
   }

   // Sanity check one of the entries.
   for id, reps := range repsByID {
      if len(reps) > 0 {
         if reps[0].ID != id {
            t.Errorf("Representation with ID '%s' is stored under the wrong key '%s'", reps[0].ID, id)
         }
      }
      break // just check the first one
   }
}
