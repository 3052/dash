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
}

func TestMPD_QualityOptions(t *testing.T) {
   // Representation data is defined once.
   repHD := &Representation{
      ID:        "video-hd",
      Bandwidth: 2000,
   }

   // This manifest has two periods. The 'video-hd' Representation appears in both,
   // but with a different SegmentTemplate in each context.
   mpd := &MPD{
      Periods: []*Period{
         { // Main content period
            ID: "main_content",
            AdaptationSets: []*AdaptationSet{
               {
                  ContentType: "video", Lang: "en",
                  Representations: []*Representation{repHD},
                  SegmentTemplate: &SegmentTemplate{
                     Media:     "content/segment-$Number$.m4s",
                     EndNumber: 10,
                  },
               },
            },
         },
         { // Ad break period
            ID: "ad_break",
            AdaptationSets: []*AdaptationSet{
               {
                  ContentType: "video", Lang: "en",
                  Representations: []*Representation{repHD}, // Same Representation data
                  SegmentTemplate: &SegmentTemplate{
                     Media:     "ads/segment-$Number$.m4s",
                     EndNumber: 2,
                  },
               },
            },
         },
      },
   }

   options := mpd.QualityOptions()

   // There should only be one key for "video-hd"
   if len(options) != 1 {
      t.Fatalf("expected 1 unique quality option, but got %d", len(options))
   }

   hdQuality, ok := options["video-hd"]
   if !ok {
      t.Fatal("expected to find key 'video-hd'")
   }

   // The single Quality object should have a Bandwidth of 2000
   if hdQuality.Bandwidth != 2000 {
      t.Errorf("expected bandwidth of 2000, got %d", hdQuality.Bandwidth)
   }

   // And it should have two different contexts
   if len(hdQuality.Contexts) != 2 {
      t.Fatalf("expected 2 contexts for 'video-hd', but got %d", len(hdQuality.Contexts))
   }

   // --- Generate URLs for the first context (main content) ---
   mainContext := hdQuality.Contexts[0]
   if mainContext.Period.ID != "main_content" {
      t.Errorf("expected first context to be for 'main_content'")
   }
   mainURLs, err := hdQuality.ListMediaSegmentURLs(mainContext)
   if err != nil {
      t.Fatalf("error getting main content URLs: %v", err)
   }
   if len(mainURLs) != 10 {
      t.Errorf("expected 10 segments for main content, got %d", len(mainURLs))
   }
   if mainURLs[0] != "content/segment-1.m4s" {
      t.Errorf("unexpected URL for main content: %s", mainURLs[0])
   }

   // --- Generate URLs for the second context (ad break) ---
   adContext := hdQuality.Contexts[1]
   if adContext.Period.ID != "ad_break" {
      t.Errorf("expected second context to be for 'ad_break'")
   }
   adURLs, err := hdQuality.ListMediaSegmentURLs(adContext)
   if err != nil {
      t.Fatalf("error getting ad URLs: %v", err)
   }
   if len(adURLs) != 2 {
      t.Errorf("expected 2 segments for ad, got %d", len(adURLs))
   }
   if adURLs[0] != "ads/segment-1.m4s" {
      t.Errorf("unexpected URL for ad content: %s", adURLs[0])
   }
}
