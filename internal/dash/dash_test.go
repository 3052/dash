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

   // Verify manual traversal for node sanity check
   for i, period := range mpd.Periods {
      pNode := PeriodNode{Period: &mpd.Periods[i], MPD: mpd}
      if pNode.MPD != mpd {
         t.Error("PeriodNode MPD pointer mismatch")
      }
      for j := range period.AdaptationSets {
         asNode := AdaptationSetNode{AdaptationSet: &period.AdaptationSets[j], Node: pNode}
         if asNode.Node.Period != &mpd.Periods[i] {
            t.Error("AdaptationSetNode parent Period pointer mismatch")
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
   for id, nodes := range groupedReps {
      t.Logf("Representation ID: %s, Count: %d", id, len(nodes))
      for _, node := range nodes {
         totalReps++

         // Validate ID consistency
         if node.Representation.ID != id {
            t.Errorf("Node ID mismatch: map key %s vs representation ID %s", id, node.Representation.ID)
         }

         // Validate Node Chain Pointers
         if node.Node.AdaptationSet == nil {
            t.Error("Node chain broken: AdaptationSet is nil")
         }
         if node.Node.Node.Period == nil {
            t.Error("Node chain broken: Period is nil")
         }
         if node.Node.Node.MPD != mpd {
            t.Error("Node chain broken: MPD does not match original object")
         }

         // Test GetSegmentTemplateNode
         stNode := node.GetSegmentTemplateNode()
         if stNode != nil {
            // Verify consistency
            if stNode.Node.Representation != node.Representation {
               t.Error("SegmentTemplateNode Representation pointer mismatch via Node")
            }

            // Manual check to verify correct inheritance logic
            expectedSt := node.Representation.SegmentTemplate
            if expectedSt == nil {
               expectedSt = node.Node.AdaptationSet.SegmentTemplate
            }

            if stNode.SegmentTemplate != expectedSt {
               t.Error("SegmentTemplateNode SegmentTemplate pointer mismatch (inheritance logic failed)")
            }
         } else {
            t.Logf("  No SegmentTemplate found for Representation %s", id)
         }
      }
   }
   t.Logf("Total Representations found: %d", totalReps)
}
