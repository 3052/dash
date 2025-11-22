package dash

import (
   "os"
   "path/filepath"
   "strings"
   "testing"
)

// TestParseMpdFiles reads all .mpd files in the testdata folder and verifies
// that they can be parsed without error.
func TestParseMpdFiles(t *testing.T) {
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
      if file.IsDir() {
         continue
      }

      if strings.HasSuffix(file.Name(), ".mpd") {
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

            // Basic Verification of Navigation Links
            verifyLinks(t, mpd)
         })
      }
   }
}

func verifyLinks(t *testing.T, m *MPD) {
   for _, p := range m.Periods {
      if p.Parent != m {
         t.Error("Navigation 10.3 failed: Period -> MPD link missing")
      }
      for _, as := range p.AdaptationSets {
         if as.Parent != p {
            t.Error("Navigation 10.1 failed: AdaptationSet -> Period link missing")
         }

         if as.SegmentTemplate != nil {
            if as.SegmentTemplate.ParentAdaptationSet != as {
               t.Error("Navigation 10.6 failed: SegmentTemplate -> AdaptationSet link missing")
            }
         }

         for _, rep := range as.Representations {
            if rep.Parent != as {
               t.Error("Navigation 10.4 failed: Representation -> AdaptationSet link missing")
            }

            if rep.SegmentTemplate != nil {
               if rep.SegmentTemplate.ParentRepresentation != rep {
                  t.Error("Navigation 10.7 failed: SegmentTemplate -> Representation link missing")
               }
            }

            if rep.SegmentList != nil {
               if rep.SegmentList.Parent != rep {
                  t.Error("Navigation 10.5 failed: SegmentList -> Representation link missing")
               }

               if rep.SegmentList.Initialization != nil {
                  if rep.SegmentList.Initialization.Parent != rep.SegmentList {
                     t.Error("Navigation 10.2 failed: Initialization -> SegmentList link missing")
                  }
               }

               for _, url := range rep.SegmentList.SegmentURLs {
                  if url.Parent != rep.SegmentList {
                     t.Error("Navigation 10.8 failed: SegmentURL -> SegmentList link missing")
                  }
               }
            }
         }
      }
   }
}
