package dash

import (
   "os"
   "path/filepath"
   "strings"
   "testing"
)

// TestPrintPSSH reads all .mpd files in "testdata" and prints the
// ContentProtection information (SchemeIdUri and PSSH) for each Representation.
func TestPrintPSSH(t *testing.T) {
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

         foundPssh := false

         for _, p := range mpd.Periods {
            for _, as := range p.AdaptationSets {
               for _, rep := range as.Representations {

                  // Get ContentProtection elements (Rep overrides AS)
                  cps := rep.GetContentProtection()

                  if len(cps) > 0 {
                     t.Logf("Representation: %s (Mime: %s)", rep.ID, rep.GetMimeType())

                     for _, cp := range cps {
                        if cp.Pssh != "" {
                           foundPssh = true
                           t.Logf("  Scheme: %s", cp.SchemeIdUri)
                           t.Logf("  PSSH:   %s", cp.Pssh)
                        }
                     }
                  }
               }
            }
         }

         if !foundPssh {
            t.Log("No PSSH data found in this MPD.")
         }
      })
   }
}
