// 由res2go IDE插件自动生成。
package main

import (
	"github.com/ying32/govcl/vcl"
	"wuchieh.com/ilmstool/pkgui"
)

func main() {
    vcl.Application.SetScaled(true)
    vcl.Application.SetTitle("ILMS作業繳交工具")
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
    vcl.Application.CreateForm(&pkgui.Form1)
    vcl.Application.CreateForm(&pkgui.Form2)
    vcl.Application.CreateForm(&pkgui.Form3)
	vcl.Application.Run()
}
