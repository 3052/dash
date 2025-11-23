package dash

import (
   "os"
   "path/filepath"
   "strings"
   "testing"
)

// TestSegmentGeneration reads all .mpd files in "testdata", parses them,
// and generates segment URLs. It processes only one Representation per unique MimeType.
func TestSegmentGeneration(t *testing.T) {
   testDataDir := "testdata"

   files, err := os.ReadDir(testDataDir)
   if err != nil {
      if os.IsNotExist(err) {
         t.Logf("testdata folder not found, skipping tests")
         return
      }
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

         // Set the MPDURL as requested
         mpd.MPDURL = "http://hello.test"

         // Track processed mimeTypes to ensure we only get one slice per type
         processedMimes := make(map[string]bool)
         count := 0

         for _, p := range mpd.Periods {
            for _, as := range p.AdaptationSets {
               for _, rep := range as.Representations {
                  mime := rep.GetMimeType()

                  // Skip if we have already processed this mimeType for this file
                  if processedMimes[mime] {
                     continue
                  }

                  // 1. Get the active SegmentTemplate (direct or inherited)
                  tmpl := rep.GetSegmentTemplate()
                  if tmpl == nil {
                     continue
                  }

                  // 2. Get the slice of replaced URLs
                  urls, err := tmpl.GetSegmentURLs(rep)
                  if err != nil {
                     t.Errorf("Failed to generate URLs for Representation %s (mime: %s): %v", rep.ID, mime, err)
                     continue
                  }

                  // Mark this mimeType as processed so we don't repeat it
                  processedMimes[mime] = true
                  count++

                  t.Logf("MimeType: %s (Rep ID: %s)", mime, rep.ID)

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
