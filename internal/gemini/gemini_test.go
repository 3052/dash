package main

import (
   "log"
   "os/exec"
   "strings"
   "testing"
)

func output(name string, arg ...string) (string, error) {
   command := exec.Command(name, arg...)
   log.Print(command.Args)
   data, err := command.Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}

func TestGemini(t *testing.T) {
   log.SetFlags(log.Ltime)
   for _, testVar := range gemini_tests {
      arg := []string{"run", ".", testVar.name}
      if testVar.url != "" {
         arg = append(arg, testVar.url)
      }
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

var gemini_tests = []struct {
   name     string
   url      string
   contains map[string]string
   state    []string
}{
   {
      name: "../../testdata/criterion.mpd",
      url:  "https://vod-adaptive-ak.vimeocdn.com/exp=1752284211~acl=%2F15be2d09-cb01-46d4-9948-2667ba2e3907%2F%2A~hmac=6997e9aef9fd359a03a2b49a7a82db955064361a16ed4d875e1d927a62f2ca35/15be2d09-cb01-46d4-9948-2667ba2e3907/v2/playlist/drm/cenc,derived,325579370,e4576465a745213f336c1ef1bf5d513e/av/primary/sub/7433271-c-en/prot/bWF4X2hlaWdodD0xMDgw/playlist.mpd",
      contains: map[string]string{
         "subs-7433271": "",
         "audio-916e7eef-13ce-4a46-9bda-b2627ec04b4f": "",
         "video-888d2bc7-75b5-4264-bf57-08e3dc24ecbb": "",
      },
      state: []string{
         `Period.duration != "" (ignore)`,
         `Representation.SegmentList != nil`,
         `SegmentTemplate.SegmentTimeline == nil (endNumber or SegmentCount)`,
         `SegmentTemplate.duration >= 1 (SegmentCount)`,
         `SegmentTemplate.endNumber == 0 (SegmentTimeline or SegmentCount)`,
         `SegmentTemplate.startNumber == nil (startNumber = 1)`,
         `SegmentTemplate.timescale != nil (ignore)`,
         `URL.IsAbs == false`,
         `len(MPD.Period) == 1`,
      },
   },
   {
      name: "../../testdata/canal.mpd",
      url:  "https://cz-bks400-prod32-live.solocoo.tv:443/bpk-token/1ac@bwrqpnwcgc4vj01ychymvdb50uune2ltbkkz13ba/bpk-vod/playout01/default/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/index.mpd",
      contains: map[string]string{
         "audio_eng_1=576000": "https://cz-bks400-prod32-live.solocoo.tv:443/bpk-token/1ac@bwrqpnwcgc4vj01ychymvdb50uune2ltbkkz13ba/bpk-vod/playout01/default/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/dash/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD-audio_eng_1=576000-383904768.dash?serviceid=298f95e1bf91361258c44a2b1f4a2425",
         "video=3399914":      "https://cz-bks400-prod32-live.solocoo.tv:443/bpk-token/1ac@bwrqpnwcgc4vj01ychymvdb50uune2ltbkkz13ba/bpk-vod/playout01/default/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/dash/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD-video=3399914-4798800.dash?serviceid=298f95e1bf91361258c44a2b1f4a2425",
         "thumbnail":          "", // the MPD is actually invalid
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
