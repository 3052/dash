package hls

import (
   "net/url"
   "os"
   "path/filepath"
   "strings"
   "testing"
)

const (
   mediaFilename  = "8500_complete-95fe4117-98fe-4ab7-8895-b2eec69b2b63.m3u8"
   masterFilename = "ctr-all-fb600154-a5e0-4125-ab89-01d627163485-b123e16f-c381-4335-bf76-dcca65425460.m3u8"
)

func TestDecodeMedia(t *testing.T) {
   path := filepath.Join("../testdata", mediaFilename)
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
   baseURL, err := url.Parse("https://example.com/video/")
   if err != nil {
      t.Fatalf("Failed to parse base URL: %v", err)
   }

   media.ResolveURIs(baseURL)

   expectedURI := "https://example.com/video/H264_1_CMAF_CENC_CTR_8500K/95fe4117-98fe-4ab7-8895-b2eec69b2b63/pts_0.mp4"

   if media.Segments[0].URI == nil {
      t.Fatal("Expected URI, got nil")
   }
   if media.Segments[0].URI.String() != expectedURI {
      t.Errorf("Expected Absolute URI %s, got %s", expectedURI, media.Segments[0].URI.String())
   }
}

func TestDecodeMaster(t *testing.T) {
   path := filepath.Join("../testdata", masterFilename)
   data, err := os.ReadFile(path)
   if err != nil {
      t.Fatalf("Failed to read file from %s: %v", path, err)
   }

   master, err := DecodeMaster(string(data))
   if err != nil {
      t.Fatalf("DecodeMaster failed: %v", err)
   }

   // The sample manifest has 8 unique video stream URIs, not 16 variants.
   if len(master.Streams) != 8 {
      t.Errorf("Expected 8 unique streams, got %d", len(master.Streams))
   }

   // Check URI of first stream before sorting
   if master.Streams[0].URI == nil {
      t.Error("Expected stream to have a valid URI")
   } else {
      if master.Streams[0].URI.Path == "" {
         t.Error("Expected stream URI path to be populated")
      }
   }

   // Find a specific stream to verify grouping
   var foundStream *Stream
   for _, stream := range master.Streams {
      if strings.Contains(stream.URI.Path, "8500_complete") {
         foundStream = stream
         break
      }
   }
   if foundStream == nil {
      t.Fatal("Could not find expected stream '8500_complete' to test grouping")
   }
   // This specific stream has two #EXT-X-STREAM-INF tags pointing to it.
   if len(foundStream.Variants) != 2 {
      t.Errorf("Expected stream to have 2 variants, got %d", len(foundStream.Variants))
   }

   // Sort the renditions and streams
   master.Sort()

   // Print all renditions (Medias) first
   t.Log("--- Renditions (sorted by GroupID) ---")
   for _, rendition := range master.Medias {
      t.Logf("Rendition:\n%s\n---", rendition)
   }

   // Print all streams and their grouped variants
   t.Log("\n--- Streams (sorted by Max Bandwidth) ---")
   for _, stream := range master.Streams {
      t.Logf("%s\n---", stream.String())
   }
}
