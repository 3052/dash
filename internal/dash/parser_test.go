package dash

import (
   "os"
   "testing"
)

func TestParse(t *testing.T) {
   // os.ReadFile is the modern replacement for opening and reading an entire file.
   // It handles opening, reading, and closing the file automatically.
   byteValue, err := os.ReadFile("rakuten.mpd.txt")
   if err != nil {
      t.Fatalf("Error reading test file: %v", err)
   }

   // Call the Parse function
   mpd, err := Parse(byteValue)
   if err != nil {
      t.Fatalf("Parse() returned an error: %v", err)
   }

   // Basic checks to ensure parsing was successful
   if mpd == nil {
      t.Fatal("Parse() returned a nil MPD struct")
   }

   expectedDuration := "PT47M48.000S"
   if mpd.MediaPresentationDuration == nil || *mpd.MediaPresentationDuration != expectedDuration {
      t.Errorf("Expected MediaPresentationDuration to be %s, but got %v", expectedDuration, mpd.MediaPresentationDuration)
   }

   if len(mpd.Periods) != 1 {
      t.Errorf("Expected 1 Period, but got %d", len(mpd.Periods))
   }

   if len(mpd.Periods[0].AdaptationSets) != 2 {
      t.Errorf("Expected 2 AdaptationSets, but got %d", len(mpd.Periods[0].AdaptationSets))
   }

   videoSet := mpd.Periods[0].AdaptationSets[0]
   if len(videoSet.Representations) != 6 {
      t.Errorf("Expected 6 video representations, but got %d", len(videoSet.Representations))
   }

   audioSet := mpd.Periods[0].AdaptationSets[1]
   expectedLang := "spa"
   if audioSet.Lang == nil || *audioSet.Lang != expectedLang {
      t.Errorf("Expected audio language to be %s, but got %v", expectedLang, audioSet.Lang)
   }
}
