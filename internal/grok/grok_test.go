package main

import (
   "encoding/json"
   "log"
   "os/exec"
   "testing"
)

var tests = []struct {
   name           string
   url            string
   representation []representation
}{
   {
      name: "../../testdata/criterion.mpd",
      url:  "https://vod-adaptive-ak.vimeocdn.com/exp=1752284211~acl=%2F15be2d09-cb01-46d4-9948-2667ba2e3907%2F%2A~hmac=6997e9aef9fd359a03a2b49a7a82db955064361a16ed4d875e1d927a62f2ca35/15be2d09-cb01-46d4-9948-2667ba2e3907/v2/playlist/drm/cenc,derived,325579370,e4576465a745213f336c1ef1bf5d513e/av/primary/sub/7433271-c-en/prot/bWF4X2hlaWdodD0xMDgw/playlist.mpd",
      representation: []representation{
         {
            id:     "video-888d2bc7-75b5-4264-bf57-08e3dc24ecbb",
            length: 1 + 1 + 1114 + 1,
            url:    "drm/cenc,derived,325579370,e4576465a745213f336c1ef1bf5d513e/remux/avf/888d2bc7-75b5-4264-bf57-08e3dc24ecbb/segment.mp4?pathsig=8c953e4f~vEyD7FR7NMtgBhRbRGol6tYRL0pVp7AQxjE5pUlKliI&r=dXMtY2VudHJhbDE%3D&sid=1116&st=video",
         },
         {
            id:     "subs-7433271",
            length: 1,
            url:    "texttrack/sub/7433271.vtt?pathsig=8c953e4f~UO056QMhmjVj394TCzXUSJJ4GI4BcpMoXktkwXsYSjw&r=dXMtY2VudHJhbDE%3D",
         },
      },
   },
   {
      name: "../../testdata/molotov.mpd",
      url:  "https://vod-molotov.akamaized.net/output/v2/d8/a1/65/32e3c47902de4911dca77b0ad73e9ac34965a1d8/32e3c47902de4911dca77b0ad73e9ac34965a1d8.ism/fhdready.mpd",
      representation: []representation{
         {
            id:     "video=4800000",
            length: 1 + 3555,
            url:    "dash/32e3c47902de4911dca77b0ad73e9ac34965a1d8-video=4800000-3555.m4s",
         },
         {
            id: "3=1000",
            length: 1 + 3339,
            url: "dash/32e3c47902de4911dca77b0ad73e9ac34965a1d8-3=1000-3339.m4s",
         },
      },
   },
   {
      name: "../../testdata/paramount.mpd",
      url:  "https://vod-gcs-cedexis.cbsaavideo.com/intl_vms/2024/10/01/2376943683811/2939404_cenc_precon_dash/stream.mpd",
      representation: []representation{
         {
            id: "5",
            length: 1 + 539 + 1 + 1 + 29 + 1,
            url: "TPIR_0722_100824_2997DF_1920x1080_178_2CH_PRORESHQ_2CH_2939373_4500/seg_571.m4s",
         },
         {
            id: "8",
            length: 1 + 540 + 1 + 22,
            url: "TPIR_0722_2997_2CH_DF_1728406422/seg_563.m4s",
         },
         {
            id: "thumb_320x180",
            length: 11,
            url: "thumb_320x180/tile_11.jpg",
         },
      },
   },
}

func output(name string, arg ...string) ([]byte, error) {
   command := exec.Command(name, arg...)
   log.Print(command.Args)
   return command.Output()
}

type representation struct {
   id     string
   length int
   url    string
}

func Test(t *testing.T) {
   log.SetFlags(log.Ltime)
   for _, testVar := range tests {
      data, err := output("go", "run", ".", "-input", testVar.name)
      if err != nil {
         t.Fatal(string(data))
      }
      var representsB map[string][]string
      err = json.Unmarshal(data, &representsB)
      if err != nil {
         t.Fatal(err)
      }
      for _, representA := range testVar.representation {
         representB := representsB[representA.id]
         if len(representB) != representA.length {
            t.Fatal(
               representA.id,
               "pass", representA.length,
               "fail", len(representB),
            )
         }
         if representB[len(representB)-1] != prefix+representA.url {
            t.Fatal(
               "\npass", prefix+representA.url,
               "\nfail", representB[len(representB)-1],
            )
         }
      }
   }
}

const prefix = "http://test.test/"
