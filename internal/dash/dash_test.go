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
         name: "Only ID replacement",
         rep: &Representation{
            ID:        "audio-main",
            Bandwidth: 128000,
         },
         template:    "init-$RepresentationID$.m4s",
         expectedURL: "init-audio-main.m4s",
      },
      {
         name: "Only Bandwidth replacement",
         rep: &Representation{
            ID:        "video-low",
            Bandwidth: 500000,
         },
         template:    "init-$Bandwidth$.m4s",
         expectedURL: "init-500000.m4s",
      },
      {
         name: "No identifiers",
         rep: &Representation{
            ID:        "video-high",
            Bandwidth: 8000000,
         },
         template:    "initialization.m4s",
         expectedURL: "initialization.m4s",
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

func TestRepresentation_ListMediaSegmentURLs(t *testing.T) {
   rep := &Representation{
      ID:        "video-hd",
      Bandwidth: 5000000,
   }

   // Template defined at the AdaptationSet level
   asTpl := &SegmentTemplate{
      Media:       "media-$RepresentationID$-$Number$.m4s",
      StartNumber: 1,
      EndNumber:   3,
   }

   // Template defined at the Representation level (should take precedence)
   repTpl := &SegmentTemplate{
      Media:       "media-rep-$Bandwidth$-$Number$.m4s",
      StartNumber: 10,
      EndNumber:   12,
   }

   testCases := []struct {
      name             string
      rep              *Representation
      asTpl            *SegmentTemplate // AdaptationSet template
      expectedLen      int
      expectedErr      bool
      expectedFirstURL string
      expectedLastURL  string
   }{
      {
         name:             "Uses AdaptationSet template",
         rep:              rep, // This rep has no template of its own
         asTpl:            asTpl,
         expectedLen:      3,
         expectedErr:      false,
         expectedFirstURL: "media-video-hd-1.m4s",
         expectedLastURL:  "media-video-hd-3.m4s",
      },
      {
         name: "Uses Representation template",
         rep: &Representation{
            ID:              "video-sd",
            Bandwidth:       1000000,
            SegmentTemplate: repTpl,
         },
         asTpl:            asTpl, // This should be ignored
         expectedLen:      3,
         expectedErr:      false,
         expectedFirstURL: "media-rep-1000000-10.m4s",
         expectedLastURL:  "media-rep-1000000-12.m4s",
      },
      {
         name:        "No template available",
         rep:         rep,
         asTpl:       nil,
         expectedLen: 0,
         expectedErr: true,
      },
      {
         name: "Template with no end number",
         rep:  rep,
         asTpl: &SegmentTemplate{
            Media:       "foo-$Number$.m4s",
            StartNumber: 1,
         },
         expectedLen: 0,
         expectedErr: true,
      },
      {
         name: "Template with default start number",
         rep:  rep,
         asTpl: &SegmentTemplate{
            Media:     "bar-$Number$.m4s",
            EndNumber: 2,
         },
         expectedLen:      2,
         expectedErr:      false,
         expectedFirstURL: "bar-1.m4s",
         expectedLastURL:  "bar-2.m4s",
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         urls, err := tc.rep.ListMediaSegmentURLs(tc.asTpl)

         if tc.expectedErr {
            if err == nil {
               t.Errorf("expected an error but got none")
            }
            return // Test ends here for error cases
         }

         if err != nil {
            t.Fatalf("expected no error but got: %v", err)
         }

         if len(urls) != tc.expectedLen {
            t.Fatalf("expected %d URLs, but got %d", tc.expectedLen, len(urls))
         }

         if tc.expectedLen > 0 && urls[0] != tc.expectedFirstURL {
            t.Errorf("expected first URL to be '%s', but got '%s'", tc.expectedFirstURL, urls[0])
         }

         if tc.expectedLen > 0 && urls[len(urls)-1] != tc.expectedLastURL {
            t.Errorf("expected last URL to be '%s', but got '%s'", tc.expectedLastURL, urls[len(urls)-1])
         }
      })
   }
}
