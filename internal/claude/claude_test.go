package main

import (
   "encoding/json"
   "log"
   "os/exec"
   "slices"
   "testing"
)

const prefix = "http://test.test/"

func output(name string, arg ...string) ([]byte, error) {
   command := exec.Command(name, arg...)
   log.Print(command.Args)
   return command.Output()
}

var tests = []struct {
   name           string
   representation []representation
}{
   {
      name: "paramount.txt",
      representation: []representation{
         {
            id: "5", // avc1.640028
            url: "TPIR_0722_100824_2997DF_1920x1080_178_2CH_PRORESHQ_2CH_2939373_4500/seg_571.m4s",
            length: func() int {
               initialization := 1
               media := 539 + 1 + 1 + 29 + 1
               return initialization + media
            }(),
         },
         {
            id: "8", // wvtt
            length: func() int {
               initialization := 1
               media := 540 + 1 + 22
               return initialization + media
            }(),
            url: "TPIR_0722_2997_2CH_DF_1728406422/seg_563.m4s",
         },
         {
            id: "thumb_320x180",
            length: 11,
            url: "thumb_320x180/tile_11.jpg",
         },
      },
   },
   {
      name: "molotov.txt",
      representation: []representation{
         {
            id: "video=4800000",
            length: func() int {
               initialization := 1
               media := 3555
               return initialization + media
            }(),
            url: "dash/32e3c47902de4911dca77b0ad73e9ac34965a1d8-video=4800000-3555.m4s",
         },
         {
            id: "3=1000",
            length: func() int {
               initialization := 1
               media := 3339
               return initialization + media
            }(),
            url: "dash/32e3c47902de4911dca77b0ad73e9ac34965a1d8-3=1000-3339.m4s",
         },
      },
   },
   {
      name: "criterion.txt",
      representation: []representation{
         {
            id: "video-888d2bc7-75b5-4264-bf57-08e3dc24ecbb",
            length: func() int {
               initialization := 1
               media := 1 + 1114 + 1
               return initialization + media
            }(),
            url: "drm/cenc,derived,325579370,e4576465a745213f336c1ef1bf5d513e/remux/avf/888d2bc7-75b5-4264-bf57-08e3dc24ecbb/segment.mp4?pathsig=8c953e4f~vEyD7FR7NMtgBhRbRGol6tYRL0pVp7AQxjE5pUlKliI&r=dXMtY2VudHJhbDE%3D&sid=1116&st=video",
         },
         {
            id:     "subs-7433271",
            length: 1,
            url:    "texttrack/sub/7433271.vtt?pathsig=8c953e4f~UO056QMhmjVj394TCzXUSJJ4GI4BcpMoXktkwXsYSjw&r=dXMtY2VudHJhbDE%3D",
         },
      },
   },
}

func Test(t *testing.T) {
   log.SetFlags(log.Ltime)
   type representation struct {
      Id   string
      Urls []string
   }
   for _, testVar := range tests {
      data, err := output("go", "run", ".", testVar.name)
      if err != nil {
         t.Fatal(string(data))
      }
      var representsB struct {
         Representations []*representation
      }
      err = json.Unmarshal(data, &representsB)
      if err != nil {
         t.Fatal(err)
      }
      for _, representA := range testVar.representation {
         index := slices.IndexFunc(representsB.Representations,
            func(r *representation) bool { return r.Id == representA.id },
         )
         representB := representsB.Representations[index].Urls
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

type representation struct {
   id     string
   length int
   url    string
}
