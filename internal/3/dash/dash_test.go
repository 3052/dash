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

   // Verify manual traversal for context sanity check
   for i, period := range mpd.Periods {
      pCtx := PeriodContext{Period: &mpd.Periods[i], MPD: mpd}
      if pCtx.MPD != mpd {
         t.Error("PeriodContext MPD pointer mismatch")
      }
      for j := range period.AdaptationSets {
         asCtx := AdaptationSetContext{AdaptationSet: &period.AdaptationSets[j], Context: pCtx}
         if asCtx.Context.Period != &mpd.Periods[i] {
            t.Error("AdaptationSetContext parent Period pointer mismatch")
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
   for id, contexts := range groupedReps {
      t.Logf("Representation ID: %s, Count: %d", id, len(contexts))
      for _, ctx := range contexts {
         totalReps++

         // Validate ID consistency
         if ctx.Representation.ID != id {
            t.Errorf("Context ID mismatch: map key %s vs representation ID %s", id, ctx.Representation.ID)
         }

         // Validate Context Chain Pointers
         if ctx.Context.AdaptationSet == nil {
            t.Error("Context chain broken: AdaptationSet is nil")
         }
         if ctx.Context.Context.Period == nil {
            t.Error("Context chain broken: Period is nil")
         }
         if ctx.Context.Context.MPD != mpd {
            t.Error("Context chain broken: MPD does not match original object")
         }
      }
   }
   t.Logf("Total Representations found: %d", totalReps)
}
