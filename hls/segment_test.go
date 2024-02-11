package hls

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "testing"
)

func TestSegment(t *testing.T) {
   for _, name := range segment_names {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var segment MediaSegment
      segment.New(string(text))
      fmt.Printf("%+v\n", segment.Key)
      for _, uri := range segment.Uri {
         fmt.Printf("%q\n", uri)
      }
   }
}

func get(s string) (*url.URL, []byte, error) {
   res, err := http.Get(s)
   if err != nil {
      return nil, nil, err
   }
   defer res.Body.Close()
   data, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, nil, err
   }
   return res.Request.URL, data, nil
}

// gem.cbc.ca/downton-abbey/s01e04
const hls_encrypt = "https://cbcrcott-gem.akamaized.net/95bc1901-988d-400a-a7a3-624284880413/CBC_DOWNTON_ABBEY_S01E04.ism/QualityLevels(400047)/Manifest(video,format=m3u8-aapl)"

func TestIv(t *testing.T) {
   var media MediaSegment
   media.Key.Iv = "0X000102030405060708090A0B0C0D0E0F"
   iv, err := media.Iv()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(iv)
}

func TestDecrypt(t *testing.T) {
   var segment MediaSegment
   base, text, err := get(hls_encrypt)
   if err != nil {
      t.Fatal(err)
   }
   segment.New(string(text))
   _, key, err := get(segment.Key.Uri)
   if err != nil {
      t.Fatal(err)
   }
   block, err := NewCipher(key)
   if err != nil {
      t.Fatal(err)
   }
   iv, err := segment.Iv()
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.ts")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   for i := range 9 {
      uri, err := base.Parse(segment.Uri[i])
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(uri)
      _, data, err := get(uri.String())
      if err != nil {
         t.Fatal(err)
      }
      data = Decrypt(block, iv, data)
      if _, err := file.Write(data); err != nil {
         t.Fatal(err)
      }
   }
}

var segment_names = []string{
   "m3u8/audio_eng_aacl.m3u8",
   "m3u8/video.m3u8",
}
