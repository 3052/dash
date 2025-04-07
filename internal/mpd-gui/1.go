package main

import "github.com/tadvi/winc"

func main() {
   window := winc.NewForm(nil)
   left := winc.NewMultiEdit(window)
   window.Show()
   dock := winc.NewSimpleDock(window)
   dock.Dock(left, winc.Top)
   window.OnClose().Bind(func(*winc.Event) {
      winc.Exit()
   })
   dock.Update()
   winc.RunMainLoop()
}
