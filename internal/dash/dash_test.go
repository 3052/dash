package dash

import (
   "os"
   "path/filepath"
   "strings"
   "testing"
)

func TestParse(t *testing.T) {
   // 1. Locate testdata directory
   testDir := "testdata"
   entries, err := os.ReadDir(testDir)
   if err != nil {
      t.Fatalf("Could not read testdata directory: %v. Make sure the 'testdata' folder exists.", err)
   }

   mpdFound := false

   // 2. Iterate over files
   for _, entry := range entries {
      if entry.IsDir() {
         continue
      }

      name := entry.Name()
      if strings.ToLower(filepath.Ext(name)) != ".mpd" {
         continue
      }

      mpdFound = true
      t.Run(name, func(t *testing.T) {
         path := filepath.Join(testDir, name)
         data, err := os.ReadFile(path)
         if err != nil {
            t.Fatalf("Failed to read file %s: %v", path, err)
         }

         // 3. Parse
         mpd, err := Parse(data)
         if err != nil {
            t.Fatalf("Parse failed for %s: %v", name, err)
         }

         // 4. Basic Validation (logging for visual verification)
         t.Logf("Parsed %s successfully", name)
         t.Logf("  Duration: %s", mpd.MediaPresentationDuration)
         t.Logf("  BaseURL: %s", mpd.BaseURL)
         t.Logf("  Periods: %d", len(mpd.Periods))

         for i, p := range mpd.Periods {
            t.Logf("    Period[%d] ID: %s, AdaptationSets: %d", i, p.ID, len(p.AdaptationSets))
            for j, as := range p.AdaptationSets {
               t.Logf("      AS[%d] Mime: %s, Reps: %d", j, as.MimeType, len(as.Representations))
               // Check for ContentProtection
               if len(as.ContentProtections) > 0 {
                  for _, cp := range as.ContentProtections {
                     if cp.Pssh != "" {
                        t.Logf("        Found PSSH in AdaptationSet (Scheme: %s)", cp.SchemeIdUri)
                     }
                  }
               }
            }
         }
      })
   }

   if !mpdFound {
      t.Log("No .mpd files found in testdata/. Skipping actual parsing.")
   }
}
