package main

import "github.com/tadvi/winc"

func main() {
   winc.DefaultFont = winc.NewFont("", 12, 0)
   form := winc.NewForm(nil)
   {
      panel := winc.NewPanel(form)
      {
         button := winc.NewPushButton(panel)
         dock := winc.NewSimpleDock(panel)
         input := winc.NewEdit(panel)

         output := winc.NewMultiEdit(panel)
         output.SetSize(100, 100)
         output.AddLine("hello")
         output.AddLine("world")

         dock.Dock(input, winc.Top)
         dock.Dock(button, winc.Top)
         dock.Dock(output, winc.Top)
      }
   }
   form.Show()
   form.OnClose().Bind(func(*winc.Event) {
      winc.Exit()
   })
   winc.RunMainLoop()
}
