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
            url:    "/" + "drm/cenc,derived,325579370,e4576465a745213f336c1ef1bf5d513e/remux/avf/888d2bc7-75b5-4264-bf57-08e3dc24ecbb/segment.mp4?pathsig=8c953e4f~vEyD7FR7NMtgBhRbRGol6tYRL0pVp7AQxjE5pUlKliI&r=dXMtY2VudHJhbDE%3D&sid=1116&st=video",
         },
         {
            id:     "audio-916e7eef-13ce-4a46-9bda-b2627ec04b4f",
            length: 1 + 1047 + 69,
            url:    "/" + "drm/cenc,derived,325579370,e4576465a745213f336c1ef1bf5d513e/remux/avf/916e7eef-13ce-4a46-9bda-b2627ec04b4f/segment.mp4?pathsig=8c953e4f~ta8gBIdEHUUP_p39sTXzQwDgjmCoZMymCeDBI6DL2H4&r=dXMtY2VudHJhbDE%3D&sid=1116&st=audio",
         },
         {
            id:     "subs-7433271",
            length: 1,
            url:    "/" + "texttrack/sub/7433271.vtt?pathsig=8c953e4f~UO056QMhmjVj394TCzXUSJJ4GI4BcpMoXktkwXsYSjw&r=dXMtY2VudHJhbDE%3D",
         },
      },
   },
   {
      name: "../../testdata/canal.mpd",
      url:  "https://cz-bks400-prod16-live.solocoo.tv:443/bpk-token/1ac@315gh2fp412qxkiyzybvo5noyt44xzbp31cbtmca/bpk-vod/playout01/default/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/index.mpd",
      representation: []representation{
         {
            id:     "video=3399914",
            length: 1 + 1 + 1332 + 1,
            url:    "/" + "dash/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD-video=3399914-4798800.dash?serviceid=298f95e1bf91361258c44a2b1f4a2425",
         },
         {
            id:     "audio_eng_1=576000",
            length: 1 + 1334,
            url:    "/" + "dash/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD-audio_eng_1=576000-383904768.dash?serviceid=298f95e1bf91361258c44a2b1f4a2425",
         },
         {
            id:     "thumbnail", // the MPD is actually invalid
            length: 80,
            url:    "/" + "dash/thumbnail/tile_80.jpeg?serviceid=298f95e1bf91361258c44a2b1f4a2425",
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
      data, err := output("go", "run", ".", testVar.name)
      if err != nil {
         t.Fatal(data)
      }
      var representsB struct {
         RepresentationUrls map[string][]string `json:"representation_urls"`
      }
      err = json.Unmarshal(data, &representsB)
      if err != nil {
         t.Fatal(err)
      }
      for _, representA := range testVar.representation {
         representB := representsB.RepresentationUrls[representA.id]
         if len(representB) != representA.length {
            t.Fatal(
               representA.id,
               "pass", representA.length,
               "fail", len(representB),
            )
         }
         if representB[len(representB)-1] != representA.url {
            t.Fatal(
               "\npass", representA.url,
               "\nfail", representB[len(representB)-1],
            )
         }
      }
   }
}
