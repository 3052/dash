package main

import (
   "log"
   "os/exec"
   "strings"
   "testing"
)

var kimi_tests = []struct {
   name     string
   url      string
   contains map[string]string
   state    []string
}{
   {
      name: "../../testdata/canal.mpd",
      url:  "https://cz-bks400-prod31-live.solocoo.tv:443/bpk-token/1ac@xbve3bnlusuhuoq2iaob0kj0dkjifjpix3nnjrca/bpk-vod/playout01/default/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/index.mpd",
      contains: map[string]string{
         "thumbnail":          "", // the MPD is actually invalid
         "video=3399914":      "dash/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD-video=3399914-4798800.dash?serviceid=298f95e1bf91361258c44a2b1f4a2425",
         "audio_eng_1=576000": "dash/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD-audio_eng_1=576000-383904768.dash?serviceid=298f95e1bf91361258c44a2b1f4a2425",
      },
      state: []string{
         `Period.duration != "" (ignore)`,
         `Representation.SegmentTemplate != nil`,
         `SegmentTemplate.SegmentTimeline != nil`,
         `SegmentTemplate.timescale != nil (ignore)`,
         `len(MPD.Period) == 1`,
         `URL.IsAbs == false`,
         `strings.Contains(SegmentTemplate.media, "$Time$")`,
         `SegmentTemplate.startNumber == nil (startNumber = 1)`,
         `SegmentTemplate.duration == 0 (SegmentTimeline or endNumber)`,
         `SegmentTemplate.endNumber == 0 (SegmentTimeline or SegmentCount)`,
      },
   },
}

func Test(t *testing.T) {
   log.SetFlags(log.Ltime)
   for _, testVar := range kimi_tests {
      arg := []string{"run", ".", testVar.name}
      data, err := output("go", arg...)
      if err != nil {
         t.Fatal(data)
      }
      for _, value := range testVar.contains {
         if !strings.Contains(data, value) {
            t.Fatal(value)
         }
      }
   }
}

func output(name string, arg ...string) (string, error) {
   command := exec.Command(name, arg...)
   log.Print(command.Args)
   data, err := command.Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}
