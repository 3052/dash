package main

import (
   "encoding/json"
   "log"
   "os/exec"
   "testing"
)

var tests = []struct {
   name           string
   representation []representation
}{
   {
      name: "criterion.txt",
      representation: []representation{
         {
            content_type: type_text,
            id:           "subs-7433271",
            length:       1,
            url:          prefix + "texttrack/sub/7433271.vtt?pathsig=8c953e4f~UO056QMhmjVj394TCzXUSJJ4GI4BcpMoXktkwXsYSjw&r=dXMtY2VudHJhbDE%3D",
         },
         {
            content_type: type_video,
            id:           "video-888d2bc7-75b5-4264-bf57-08e3dc24ecbb",
            length: func() int {
               initialization := 1
               media := 1 + 1114 + 1
               return initialization + media
            }(),
            url: prefix + "drm/cenc,derived,325579370,e4576465a745213f336c1ef1bf5d513e/remux/avf/888d2bc7-75b5-4264-bf57-08e3dc24ecbb/segment.mp4?pathsig=8c953e4f~vEyD7FR7NMtgBhRbRGol6tYRL0pVp7AQxjE5pUlKliI&r=dXMtY2VudHJhbDE%3D&sid=1116&st=video",
         },
      },
   },
   {
      name: "molotov.txt",
      representation: []representation{
         {
            content_type: type_text,
            id:           "3=1000",
            length: func() int {
               initialization := 1
               media := 3339
               return initialization + media
            }(),
            url: prefix + "dash/32e3c47902de4911dca77b0ad73e9ac34965a1d8-3=1000-3339.m4s",
         },
         {
            content_type: type_video,
            id:           "video=4800000",
            length: func() int {
               initialization := 1
               media := 3555
               return initialization + media
            }(),
            url: prefix + "dash/32e3c47902de4911dca77b0ad73e9ac34965a1d8-video=4800000-3555.m4s",
         },
      },
   },
}

const prefix = "http://test.test/"

func Test(t *testing.T) {
   log.SetFlags(log.Ltime)
   for _, testVar := range tests {
      data, err := output("go", "run", ".", testVar.name)
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
         if representB[len(representB)-1] != representA.url {
            t.Fatal(
               "\npass", representA.url,
               "\nfail", representB[len(representB)-1],
            )
         }
      }
   }
}

func output(name string, arg ...string) ([]byte, error) {
   command := exec.Command(name, arg...)
   log.Print(command.Args)
   return command.Output()
}

type content_type int

const (
   type_image content_type = iota
   type_text
   type_video
)

type representation struct {
   id           string
   length       int
   url          string
   content_type content_type
}
