package dash

import (
   "net/url"
   "os"
   "testing"
)

func TestUrl(t *testing.T) {
   data, err := os.ReadFile("testdata/criterion.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   media.BaseUrl[0] = &url.URL{Path: "/0/1/2/3/4/5/6/7/8/9/10"}
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   if media.BaseUrl.UnmarshalText([]byte{'\n'}) == nil {
      t.Fatal("BaseUrl.UnmarshalText")
   }
}

func TestDuration(t *testing.T) {
   var d Duration
   if d.UnmarshalText(nil) == nil {
      t.Fatal("Duration.UnmarshalText")
   }
}

func TestMpd(t *testing.T) {
   data, err := os.ReadFile("testdata/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
}
