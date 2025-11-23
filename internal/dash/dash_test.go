package dash

import (
   "os"
   "path/filepath"
   "strings"
   "testing"
)

// TestSegmentGeneration reads all .mpd files in "testdata", parses them,
// and generates segment URLs for each Representation found.
func TestSegmentGeneration(t *testing.T) {
   testDataDir := "testdata"

   files, err := os.ReadDir(testDataDir)
   if err != nil {
      t.Fatalf("Failed to read testdata directory: %v", err)
   }

   for _, file := range files {
      if file.IsDir() || !strings.HasSuffix(file.Name(), ".mpd") {
         continue
      }

      t.Run(file.Name(), func(t *testing.T) {
         path := filepath.Join(testDataDir, file.Name())
         data, err := os.ReadFile(path)
         if err != nil {
            t.Fatalf("Failed to read file %s: %v", file.Name(), err)
         }

         mpd, err := Parse(data)
         if err != nil {
            t.Fatalf("Failed to parse MPD %s: %v", file.Name(), err)
         }

         // Set a dummy base URL to simulate the file being hosted,
         // ensuring relative BaseURLs resolve correctly.
         mpd.MPDURL = "http://localhost/dash/" + file.Name()

         count := 0
         for _, p := range mpd.Periods {
            for _, as := range p.AdaptationSets {
               for _, rep := range as.Representations {

                  // 1. Get the active SegmentTemplate (direct or inherited)
                  tmpl := rep.GetSegmentTemplate()
                  if tmpl == nil {
                     continue
                  }

                  // 2. Get the slice of replaced URLs
                  urls, err := tmpl.GetSegmentURLs(rep)
                  if err != nil {
                     t.Errorf("Failed to generate URLs for Representation %s: %v", rep.ID, err)
                     continue
                  }

                  count++
                  t.Logf("Representation ID: %s", rep.ID)

                  // 3. Print slice length
                  t.Logf("  Slice Length: %d", len(urls))

                  // 4. Print first and last URLs
                  if len(urls) > 0 {
                     t.Logf("  First URL: %s", urls[0].String())
                     if len(urls) > 1 {
                        t.Logf("  Last URL:  %s", urls[len(urls)-1].String())
                     }
                  }
               }
            }
         }

         if count == 0 {
            t.Log("No Representations with SegmentTemplates found in this MPD.")
         }
      })
   }
}
