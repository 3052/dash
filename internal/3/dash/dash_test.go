package dash

import (
   "os"
   "testing"
)

func TestParse(t *testing.T) {
   filename := "testdata/rakuten.mpd"
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

   // Verify manual traversal for scope sanity check
   for i, period := range mpd.Periods {
      pScope := PeriodScope{Period: &mpd.Periods[i], MPD: mpd}
      if pScope.MPD != mpd {
         t.Error("PeriodScope MPD pointer mismatch")
      }
      for j := range period.AdaptationSets {
         asScope := AdaptationSetScope{AdaptationSet: &period.AdaptationSets[j], Scope: pScope}
         if asScope.Scope.Period != &mpd.Periods[i] {
            t.Error("AdaptationSetScope parent Period pointer mismatch")
         }
      }
   }

   // Test GetRepresentations
   t.Log("Testing GetRepresentations...")
   groupedReps := mpd.GetRepresentations()

   if len(groupedReps) == 0 {
      // We expect at least some representations if the file is valid DASH
      t.Log("Warning: GetRepresentations returned empty map (input file might be empty or have no representations)")
   }

   totalReps := 0
   for id, scopes := range groupedReps {
      t.Logf("Representation ID: %s, Count: %d", id, len(scopes))
      for _, scope := range scopes {
         totalReps++

         // Validate ID consistency
         if scope.Representation.ID != id {
            t.Errorf("Scope ID mismatch: map key %s vs representation ID %s", id, scope.Representation.ID)
         }

         // Validate Scope Chain Pointers
         if scope.Scope.AdaptationSet == nil {
            t.Error("Scope chain broken: AdaptationSet is nil")
         }
         if scope.Scope.Scope.Period == nil {
            t.Error("Scope chain broken: Period is nil")
         }
         if scope.Scope.Scope.MPD != mpd {
            t.Error("Scope chain broken: MPD does not match original object")
         }

         // Test GetSegmentTemplateScope
         stScope := scope.GetSegmentTemplateScope()
         if stScope != nil {
            // Verify consistency
            if stScope.Scope.Representation != scope.Representation {
               t.Error("SegmentTemplateScope Representation pointer mismatch via Scope")
            }

            // Manual check to verify correct inheritance logic
            expectedSt := scope.Representation.SegmentTemplate
            if expectedSt == nil {
               expectedSt = scope.Scope.AdaptationSet.SegmentTemplate
            }

            if stScope.SegmentTemplate != expectedSt {
               t.Error("SegmentTemplateScope SegmentTemplate pointer mismatch (inheritance logic failed)")
            }
         } else {
            t.Logf("  No SegmentTemplate found for Representation %s", id)
         }
      }
   }
   t.Logf("Total Representations found: %d", totalReps)
}
