package dash

import (
   "os"
   "testing"
)

// The TestParse function remains unchanged from the previous version.
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

func TestRepresentation_ResolveURL(t *testing.T) {
   testCases := []struct {
      name        string
      rep         *Representation
      template    string
      expectedURL string
   }{
      {
         name: "ID and Bandwidth replacement",
         rep: &Representation{
            ID:        "video-1080p",
            Bandwidth: 4000000,
         },
         template:    "init-$RepresentationID$-$Bandwidth$.m4s",
         expectedURL: "init-video-1080p-4000000.m4s",
      },
      {
         name:        "Nil representation",
         rep:         nil,
         template:    "init-$RepresentationID$.m4s",
         expectedURL: "init-$RepresentationID$.m4s",
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         resolvedURL := tc.rep.ResolveURL(tc.template)
         if resolvedURL != tc.expectedURL {
            t.Errorf("expected URL: %s, got: %s", tc.expectedURL, resolvedURL)
         }
      })
   }
}

func TestSegmentTimeline_GetSegments(t *testing.T) {
   st := &SegmentTimeline{
      Segments: []*S{
         {D: 100},        // time: 0, duration: 100
         {D: 100, R: 1},  // time: 100, 200. duration 100 (2 segments total)
         {T: 500, D: 50}, // time: 500, duration: 50
         {D: 50},         // time: 550, duration: 50
      },
   }
   segments := st.GetSegments()
   expectedCount := 5
   if len(segments) != expectedCount {
      t.Fatalf("expected %d segments, got %d", expectedCount, len(segments))
   }

   expectedTimes := []uint64{0, 100, 200, 500, 550}
   for i, s := range segments {
      if s.StartTime != expectedTimes[i] {
         t.Errorf("segment %d: expected start time %d, got %d", i, expectedTimes[i], s.StartTime)
      }
   }
}

func TestRepresentation_ListMediaSegmentURLs(t *testing.T) {
   rep := &Representation{
      ID:        "video-hd",
      Bandwidth: 5000000,
   }

   timelineTpl := &SegmentTemplate{
      Media:       "media-$RepresentationID$-t$Time$.m4s",
      StartNumber: 1,
      SegmentTimeline: &SegmentTimeline{
         Segments: []*S{
            {D: 100},       // t=0
            {D: 100, R: 1}, // t=100, t=200
         },
      },
   }

   numberTpl := &SegmentTemplate{
      Media:       "media-$Bandwidth$-$Number$.m4s",
      StartNumber: 1,
      EndNumber:   3,
   }

   testCases := []struct {
      name             string
      rep              *Representation
      asTpl            *SegmentTemplate
      expectedLen      int
      expectedErr      bool
      expectedFirstURL string
      expectedLastURL  string
   }{
      {
         name:             "Timeline-based URLs",
         rep:              rep,
         asTpl:            timelineTpl,
         expectedLen:      3,
         expectedErr:      false,
         expectedFirstURL: "media-video-hd-t0.m4s",
         expectedLastURL:  "media-video-hd-t200.m4s",
      },
      {
         name:             "Number-based URLs",
         rep:              rep,
         asTpl:            numberTpl,
         expectedLen:      3,
         expectedErr:      false,
         expectedFirstURL: "media-5000000-1.m4s",
         expectedLastURL:  "media-5000000-3.m4s",
      },
      {
         name:        "No template available",
         rep:         rep,
         asTpl:       nil,
         expectedLen: 0,
         expectedErr: true,
      },
      {
         name: "Template with no timeline or endNumber",
         rep:  rep,
         asTpl: &SegmentTemplate{
            Media: "foo.m4s",
         },
         expectedLen: 0,
         expectedErr: true,
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         urls, err := tc.rep.ListMediaSegmentURLs(tc.asTpl)

         if tc.expectedErr {
            if err == nil {
               t.Errorf("expected an error but got none")
            }
            return
         }
         if err != nil {
            t.Fatalf("expected no error but got: %v", err)
         }
         if len(urls) != tc.expectedLen {
            t.Fatalf("expected %d URLs, but got %d", tc.expectedLen, len(urls))
         }
         if tc.expectedLen > 0 {
            if urls[0] != tc.expectedFirstURL {
               t.Errorf("expected first URL to be '%s', but got '%s'", tc.expectedFirstURL, urls[0])
            }
            if urls[len(urls)-1] != tc.expectedLastURL {
               t.Errorf("expected last URL to be '%s', but got '%s'", tc.expectedLastURL, urls[len(urls)-1])
            }
         }
      })
   }
}
