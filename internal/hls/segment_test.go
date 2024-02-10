package hls

import (
   "fmt"
   "os"
   "testing"
)

// gem.cbc.ca/downton-abbey/s01e04
const hls_encrypt = "https://cbcrcott-gem.akamaized.net/95bc1901-988d-400a-a7a3-624284880413/CBC_DOWNTON_ABBEY_S01E04.ism/QualityLevels(400047)/Manifest(video,format=m3u8-aapl)"

func TestBlock(t *testing.T) {
   res, err := http.Get(hls_encrypt)
   if err != nil {
      t.Fatal(err)
   }
   if res.StatusCode != http.StatusOK {
      t.Fatal(res.Status)
   }
   seg, err := NewScanner(res.Body).Segment()
   if err != nil {
      t.Fatal(err)
   }
   if err := res.Body.Close(); err != nil {
      t.Fatal(err)
   }
   key, err := get_key(seg.Key)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.ts")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   block, err := NewBlock(key)
   if err != nil {
      t.Fatal(err)
   }
   for i := 0; i <= 9; i++ {
      req, err := http.NewRequest("GET", seg.URI[i], nil)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(req.URL)
      req.URL = res.Request.URL.ResolveReference(req.URL)
      func() {
         res, err := new(http.Transport).RoundTrip(req)
         if err != nil {
            t.Fatal(err)
         }
         defer res.Body.Close()
         text, err := io.ReadAll(res.Body)
         if err != nil {
            t.Fatal(err)
         }
         text = block.DecryptKey(text)
         if _, err := file.Write(text); err != nil {
            t.Fatal(err)
         }
      }()
   }
}

var segment_names = []string{
   "audio_eng_aacl.m3u8",
   "video.m3u8",
}

func TestSegment(t *testing.T) {
   for _, name := range segment_names {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var segment MediaSegment
      segment.New(string(text))
      fmt.Printf("%+v\n", segment.Key)
      for _, uri := range segment.URI {
         fmt.Printf("%q\n", uri)
      }
   }
}

var raw_ivs = []string{
   "00000000000000000000000000000001",
   "0X00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func TestHex(t *testing.T) {
   for _, raw_iv := range raw_ivs {
      iv, err := Segment{RawIv: raw_iv}.IV()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(iv)
   }
}

func get_key(s string) ([]byte, error) {
   res, err := http.Get(s)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}
