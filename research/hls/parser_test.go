package hls

import (
   "os"
   "path/filepath"
   "testing"
)

const (
   mediaFilename  = "8500_complete-95fe4117-98fe-4ab7-8895-b2eec69b2b63.m3u8"
   masterFilename = "ctr-all-fb600154-a5e0-4125-ab89-01d627163485-b123e16f-c381-4335-bf76-dcca65425460.m3u8"
)

func TestDecodeMedia(t *testing.T) {
   path := filepath.Join("testdata", mediaFilename)
   data, err := os.ReadFile(path)
   if err != nil {
      t.Fatalf("Failed to read file from %s: %v", path, err)
   }

   media, err := DecodeMedia(string(data))
   if err != nil {
      t.Fatalf("DecodeMedia failed: %v", err)
   }

   if media.TargetDuration != 9 {
      t.Errorf("Expected TargetDuration 9, got %d", media.TargetDuration)
   }

   // Test ResolveURIs
   err = media.ResolveURIs("https://example.com/video/")
   if err != nil {
      t.Errorf("ResolveURIs failed: %v", err)
   }
   expectedURI := "https://example.com/video/H264_1_CMAF_CENC_CTR_8500K/95fe4117-98fe-4ab7-8895-b2eec69b2b63/pts_0.mp4"

   if media.Segments[0].URI == nil {
      t.Fatal("Expected URI, got nil")
   }
   if media.Segments[0].URI.String() != expectedURI {
      t.Errorf("Expected Absolute URI %s, got %s", expectedURI, media.Segments[0].URI.String())
   }
}

func TestDecodeMaster(t *testing.T) {
   path := filepath.Join("testdata", masterFilename)
   data, err := os.ReadFile(path)
   if err != nil {
      t.Fatalf("Failed to read file from %s: %v", path, err)
   }

   master, err := DecodeMaster(string(data))
   if err != nil {
      t.Fatalf("DecodeMaster failed: %v", err)
   }

   if len(master.Variants) != 16 {
      t.Errorf("Expected 16 variants, got %d", len(master.Variants))
   }

   // Check URI of first variant
   if master.Variants[0].URI == nil {
      t.Error("Expected variant to have a valid URI")
   } else {
      // Just ensuring it parsed correctly (relative path)
      if master.Variants[0].URI.Path == "" {
         t.Error("Expected variant URI path to be populated")
      }
   }
}

func TestDecode_Mismatch(t *testing.T) {
   path := filepath.Join("testdata", mediaFilename)
   data, _ := os.ReadFile(path)

   _, err := DecodeMaster(string(data))
   if err != ErrNotMaster {
      t.Errorf("Expected ErrNotMaster, got %v", err)
   }
}
