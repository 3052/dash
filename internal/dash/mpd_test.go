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
}

func TestMPD_QualityOptions(t *testing.T) {
   repVideo := &Representation{ID: "video-hd", Bandwidth: 2500}
   repAudioEn := &Representation{ID: "audio-en-stereo", Bandwidth: 128}
   repAudioEs := &Representation{ID: "audio-es-stereo", Bandwidth: 128}

   // A representation with a colliding ID but different context
   repAudioEn_period2 := &Representation{ID: "audio-en-stereo", Bandwidth: 192}

   mpd := &MPD{
      Periods: []*Period{
         {
            AdaptationSets: []*AdaptationSet{
               {ContentType: "video", Representations: []*Representation{repVideo}},
               {ContentType: "audio", Lang: "en", Representations: []*Representation{repAudioEn}},
               {ContentType: "audio", Lang: "es", Representations: []*Representation{repAudioEs}},
            },
         },
         {
            AdaptationSets: []*AdaptationSet{
               {ContentType: "audio", Lang: "en", Representations: []*Representation{repAudioEn_period2}},
            },
         },
      },
   }

   options := mpd.QualityOptions()

   if len(options) != 3 {
      t.Fatalf("expected 3 unique quality option IDs, but got %d", len(options))
   }

   // Test a simple, non-colliding ID
   videoQuals, ok := options["video-hd"]
   if !ok {
      t.Fatal("expected to find key 'video-hd'")
   }
   if len(videoQuals) != 1 {
      t.Fatalf("expected 1 quality for 'video-hd', got %d", len(videoQuals))
   }
   if videoQuals[0].ContentType != "video" {
      t.Errorf("expected content type 'video', got '%s'", videoQuals[0].ContentType)
   }
   if videoQuals[0].Bandwidth != 2500 {
      t.Errorf("incorrect bandwidth for 'video-hd'")
   }

   // Test the colliding ID
   audioEnQuals, ok := options["audio-en-stereo"]
   if !ok {
      t.Fatal("expected to find key 'audio-en-stereo'")
   }
   if len(audioEnQuals) != 2 {
      t.Fatalf("expected 2 qualities for 'audio-en-stereo' due to collision, got %d", len(audioEnQuals))
   }

   // Check context of the first one
   if audioEnQuals[0].Lang != "en" {
      t.Errorf("expected lang 'en', got '%s'", audioEnQuals[0].Lang)
   }
   if audioEnQuals[0].Bandwidth != 128 {
      t.Errorf("wrong bandwidth for first 'audio-en-stereo'")
   }

   // Check context of the second one
   if audioEnQuals[1].Lang != "en" {
      t.Errorf("expected lang 'en', got '%s'", audioEnQuals[1].Lang)
   }
   if audioEnQuals[1].Bandwidth != 192 {
      t.Errorf("wrong bandwidth for second 'audio-en-stereo'")
   }
}
