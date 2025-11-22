package dash

import (
   "fmt"
   "os"
   "path/filepath"
   "testing"
)

func TestParse(t *testing.T) {
   // Locate mpd files in testdata folder
   pattern := filepath.Join("testdata", "*.mpd")
   files, err := filepath.Glob(pattern)
   if err != nil {
      t.Fatalf("Failed to glob testdata: %v", err)
   }

   if len(files) == 0 {
      t.Log("No .mpd files found in testdata/ to test.")
      return
   }

   for _, file := range files {
      t.Run(filepath.Base(file), func(t *testing.T) {
         content, err := os.ReadFile(file)
         if err != nil {
            t.Fatalf("Failed to read file %s: %v", file, err)
         }

         mpd, err := Parse(content)
         if err != nil {
            t.Errorf("Failed to parse %s: %v", file, err)
            return
         }

         // Basic Validation to ensure structs are populating
         if mpd == nil {
            t.Errorf("Resulting MPD struct is nil for file %s", file)
            return // Added return to prevent nil pointer dereference
         }

         // Print debug info
         t.Logf("Parsed %s successfully. Periods: %d", file, len(mpd.Period))
         for _, p := range mpd.Period {
            t.Logf("  Period ID: %s, Duration: %s, AdaptationSets: %d", p.ID, p.Duration, len(p.AdaptationSet))
         }
      })
   }
}

func ExampleParse() {
   // Example usage
   xmlData := []byte(`
      <MPD mediaPresentationDuration="PT1H30M">
         <Period duration="PT10M">
            <AdaptationSet mimeType="video/mp4">
               <Representation id="1" bandwidth="1000000" width="1920" height="1080" />
            </AdaptationSet>
         </Period>
      </MPD>
   `)

   mpd, err := Parse(xmlData)
   if err != nil {
      panic(err)
   }

   fmt.Printf("MPD Duration: %s\n", mpd.MediaPresentationDuration)
   if len(mpd.Period) > 0 {
      fmt.Printf("First Period Duration: %s\n", mpd.Period[0].Duration)
   }

   // Output:
   // MPD Duration: PT1H30M
   // First Period Duration: PT10M
}
