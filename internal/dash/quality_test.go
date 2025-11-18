package dash

import "testing"

func TestQuality_URLResolution(t *testing.T) {
   rep := &Representation{
      ID:      "video-hd",
      BaseURL: "video/",
      SegmentTemplate: &SegmentTemplate{
         Initialization: "init.mp4",
         Media:          "seg-$Number$.m4s",
         EndNumber:      1,
      },
   }

   period := &Period{
      ID:      "main_content",
      BaseURL: "period1/",
   }

   mpd := &MPD{
      BaseURL: "http://cdn.example.com/base/",
      Periods: []*Period{period},
   }

   ctx := &RepresentationContext{
      Period:        period,
      AdaptationSet: &AdaptationSet{},
   }

   quality := &Quality{
      Representation: rep,
      Contexts:       []*RepresentationContext{ctx},
      parentMPD:      mpd,
   }

   // Test Initialization URL
   initURL, err := quality.AbsoluteInitializationURL(ctx)
   if err != nil {
      t.Fatalf("AbsoluteInitializationURL() failed: %v", err)
   }
   expectedInitURL := "http://cdn.example.com/base/period1/video/init.mp4"
   if initURL != expectedInitURL {
      t.Errorf("expected init URL '%s', got '%s'", expectedInitURL, initURL)
   }

   // Test Media Segment URLs
   mediaURLs, err := quality.AbsoluteMediaSegmentURLs(ctx)
   if err != nil {
      t.Fatalf("AbsoluteMediaSegmentURLs() failed: %v", err)
   }
   if len(mediaURLs) != 1 {
      t.Fatalf("expected 1 media URL, got %d", len(mediaURLs))
   }
   expectedMediaURL := "http://cdn.example.com/base/period1/video/seg-1.m4s"
   if mediaURLs[0] != expectedMediaURL {
      t.Errorf("expected media URL '%s', got '%s'", expectedMediaURL, mediaURLs[0])
   }
}

func TestQuality_URLResolution_WithAbsoluteOverride(t *testing.T) {
   // Here, the Period BaseURL is absolute, so it should override the MPD's.
   rep := &Representation{
      ID: "video-hd",
      SegmentTemplate: &SegmentTemplate{
         Initialization: "init.mp4",
      },
   }
   period := &Period{
      ID:      "main_content",
      BaseURL: "http://other-cdn.com/special/",
   }
   mpd := &MPD{
      BaseURL: "http://cdn.example.com/base/",
      Periods: []*Period{period},
   }
   ctx := &RepresentationContext{
      Period:        period,
      AdaptationSet: &AdaptationSet{},
   }
   quality := &Quality{
      Representation: rep,
      Contexts:       []*RepresentationContext{ctx},
      parentMPD:      mpd,
   }

   initURL, err := quality.AbsoluteInitializationURL(ctx)
   if err != nil {
      t.Fatalf("AbsoluteInitializationURL() failed: %v", err)
   }
   expectedInitURL := "http://other-cdn.com/special/init.mp4"
   if initURL != expectedInitURL {
      t.Errorf("expected init URL '%s', got '%s'", expectedInitURL, initURL)
   }
}
