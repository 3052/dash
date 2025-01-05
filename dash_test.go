package dash

import (
   "encoding/xml"
   "log"
   "net/http"
   "net/url"
   "os"
   "testing"
)

var range_tests = []struct{
   data string
   ok bool
}{
   {"-3", true},
   {"2-", false},
   {"2-3", true},
}

func TestRange(t *testing.T) {
   for _, test := range range_tests {
      var r Range
      ok := r.UnmarshalText([]byte(test.data)) == nil
      if ok != test.ok {
         t.Fatal(test)
      }
   }
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Print(req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}

func TestInitialization(t *testing.T) {
   data, err := os.ReadFile("ignore/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   base, err := url.Parse(pluto_mpd)
   if err != nil {
      t.Fatal(err)
   }
   media.Set(base)
   var represent Representation
   for represent = range media.Representation() {
      if *represent.MimeType == "video/mp4" {
         break
      }
   }
   initial, err := represent.SegmentTemplate.Initialization.Url(&represent)
   if err != nil {
      t.Fatal(err)
   }
   http.DefaultClient.Transport = transport{}
   resp, err := http.Head(initial.String())
   if err != nil {
      t.Fatal(err)
   }
   if resp.StatusCode != http.StatusOK {
      t.Fatal(resp.Status)
   }
}

func TestPssh(t *testing.T) {
   var p Pssh
   if p.UnmarshalText([]byte{0}) == nil {
      t.Fatal("Pssh.UnmarshalText")
   }
}

func Test2Initialization(t *testing.T) {
   _, err := Initialization{"\n"}.Url(&Representation{})
   if err == nil {
      t.Fatal("Initialization.Url")
   }
}

func TestDuration(t *testing.T) {
   var d Duration
   if d.UnmarshalText(nil) == nil {
      t.Fatal("Duration.UnmarshalText")
   }
}
