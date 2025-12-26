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

func TestDecode_MediaPlaylist(t *testing.T) {
   // 1. Read from disk (testdata folder)
   path := filepath.Join("testdata", mediaFilename)
   data, err := os.ReadFile(path)
   if err != nil {
      t.Fatalf("Failed to read file from %s: %v", path, err)
   }

   // 2. Decode
   playlist, err := Decode(string(data))
   if err != nil {
      t.Fatalf("Decode failed: %v", err)
   }

   // 3. Verify it is NOT a master playlist
   if playlist.IsMaster {
      t.Fatal("Expected Media Playlist, got Master Playlist")
   }

   // 4. Validate specific fields based on your provided file
   media := playlist.Media

   if media.TargetDuration != 9 {
      t.Errorf("Expected TargetDuration 9, got %d", media.TargetDuration)
   }

   if media.MediaSequence != 0 {
      t.Errorf("Expected MediaSequence 0, got %d", media.MediaSequence)
   }

   if media.PlaylistType != "VOD" {
      t.Errorf("Expected PlaylistType VOD, got %s", media.PlaylistType)
   }

   // Check Segments (File has 2 segments)
   if len(media.Segments) != 2 {
      t.Errorf("Expected 2 segments, got %d", len(media.Segments))
   }

   // Check details of the first segment
   firstSegment := media.Segments[0]
   if firstSegment.Duration != 8.0 {
      t.Errorf("Expected first segment duration 8.0, got %f", firstSegment.Duration)
   }
   expectedURI := "H264_1_CMAF_CENC_CTR_8500K/95fe4117-98fe-4ab7-8895-b2eec69b2b63/pts_0.mp4"
   if firstSegment.URI != expectedURI {
      t.Errorf("Expected URI %s, got %s", expectedURI, firstSegment.URI)
   }

   // Check DRM Keys
   if len(media.Keys) < 2 {
      t.Errorf("Expected at least 2 Keys, got %d", len(media.Keys))
   }
   if media.Keys[0].Method != "SAMPLE-AES-CTR" {
      t.Errorf("Expected Key Method SAMPLE-AES-CTR, got %s", media.Keys[0].Method)
   }

   // Check DateRanges
   if len(media.DateRanges) != 3 {
      t.Errorf("Expected 3 DateRanges, got %d", len(media.DateRanges))
   }
}

func TestDecode_MasterPlaylist(t *testing.T) {
   // 1. Read from disk (testdata folder)
   path := filepath.Join("testdata", masterFilename)
   data, err := os.ReadFile(path)
   if err != nil {
      t.Fatalf("Failed to read file from %s: %v", path, err)
   }

   // 2. Decode
   playlist, err := Decode(string(data))
   if err != nil {
      t.Fatalf("Decode failed: %v", err)
   }

   // 3. Verify it IS a master playlist
   if !playlist.IsMaster {
      t.Fatal("Expected Master Playlist, got Media Playlist")
   }

   master := playlist.Master

   // 4. Validate Variants (16 streams in your file: 8 AAC + 8 EAC-3)
   if len(master.Variants) != 16 {
      t.Errorf("Expected 16 variants, got %d", len(master.Variants))
   }

   // Check the first variant
   // #EXT-X-STREAM-INF:BANDWIDTH=3893429,...
   firstVariant := master.Variants[0]
   if firstVariant.Bandwidth != 3893429 {
      t.Errorf("Expected first variant bandwidth 3893429, got %d", firstVariant.Bandwidth)
   }
   if firstVariant.Resolution != "1280x720" {
      t.Errorf("Expected first variant resolution 1280x720, got %s", firstVariant.Resolution)
   }
   if firstVariant.Codecs != "avc1.64001f,mp4a.40.2" {
      t.Errorf("Expected codecs avc1.64001f,mp4a.40.2, got %s", firstVariant.Codecs)
   }

   // 5. Validate Media Groups (Audio/Subtitles) - 8 entries in your file
   // (English + Audio Desc for eac-3, aac-128, aac-64) + (2 subtitles) = 8
   if len(master.Medias) != 8 {
      t.Errorf("Expected 8 media entries, got %d", len(master.Medias))
   }

   // Check first Audio group
   firstMedia := master.Medias[0]
   if firstMedia.Type != "AUDIO" {
      t.Errorf("Expected Type AUDIO, got %s", firstMedia.Type)
   }
   if firstMedia.GroupID != "eac-3" {
      t.Errorf("Expected GroupID eac-3, got %s", firstMedia.GroupID)
   }

   // 6. Validate Session Keys
   if len(master.SessionKeys) < 1 {
      t.Error("Expected Session Keys, found none")
   }
}
