package main

import (
   "os"
   "os/exec"
   "testing"
)

var tests = []struct {
   name     string
   contains string
   url      string
}{
   {
      name:     "../../testdata/canal.mpd",
      contains: "",
      url:      "https://cz-bks400-prod32-live.solocoo.tv:443/bpk-token/1ac@bwrqpnwcgc4vj01ychymvdb50uune2ltbkkz13ba/bpk-vod/playout01/default/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/appletvcz_A007300100102_2464C3BF9652075492E7CF48A400F243_HD/index.mpd",
   },
}

func Test(t *testing.T) {
   for _, testVar := range tests {
      arg := []string{"run", ".", testVar.name}
      if testVar.url != "" {
         arg = append(arg, testVar.url)
      }
      data, err := exec.Command("go", arg...).Output()
      if err != nil {
         t.Fatal(string(data))
      }
      os.Stdout.Write(data)
   }
}
