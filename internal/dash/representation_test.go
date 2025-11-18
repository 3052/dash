package dash

import (
   "testing"
)

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
   // A mock period is now needed for the duration-based test.
   period := &Period{Duration: "PT60S"}

   rep := &Representation{
      ID:        "video-hd",
      Bandwidth: 5000000,
   }

   timelineTpl := &SegmentTemplate{
      Media:       "media-$RepresentationID$-t$Time$.m4s",
      StartNumber: 1,
      SegmentTimeline: &SegmentTimeline{
         Segments: []*S{
            {D: 100},
            {D: 100, R: 1},
         },
      },
   }

   numberTpl := &SegmentTemplate{
      Media:       "media-$Bandwidth$-$Number$.m4s",
      StartNumber: 1,
      EndNumber:   3,
   }

   // New template for duration-based calculation
   durationTpl := &SegmentTemplate{
      Media:     "media-$Number$.m4s",
      Timescale: 1000,
      Duration:  10000, // 10 seconds
   }

   testCases := []struct {
      name             string
      period           *Period
      rep              *Representation
      asTpl            *SegmentTemplate
      expectedLen      int
      expectedErr      bool
      expectedFirstURL string
      expectedLastURL  string
   }{
      {
         name:             "Timeline-based URLs",
         period:           period,
         rep:              rep,
         asTpl:            timelineTpl,
         expectedLen:      3,
         expectedErr:      false,
         expectedFirstURL: "media-video-hd-t0.m4s",
         expectedLastURL:  "media-video-hd-t200.m4s",
      },
      {
         name:             "Number-based URLs",
         period:           period,
         rep:              rep,
         asTpl:            numberTpl,
         expectedLen:      3,
         expectedErr:      false,
         expectedFirstURL: "media-5000000-1.m4s",
         expectedLastURL:  "media-5000000-3.m4s",
      },
      {
         name:        "No template available",
         period:      period,
         rep:         rep,
         asTpl:       nil,
         expectedErr: true,
      },
      {
         name:             "Duration-based calculation",
         period:           period, // PT60S
         rep:              rep,
         asTpl:            durationTpl, // 10s segments
         expectedLen:      6,           // 60s / 10s = 6 segments
         expectedErr:      false,
         expectedFirstURL: "media-1.m4s",
         expectedLastURL:  "media-6.m4s",
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         urls, err := tc.rep.ListMediaSegmentURLs(tc.period, tc.asTpl)

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
