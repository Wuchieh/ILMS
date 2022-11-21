package pkgui

import (
	"fmt"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"runtime"
	"strings"
	"sync"
)

// ::private::
type TForm2Fields struct {
}

func (f *TForm2) OnButtonInsertClick(sender vcl.IObject) {
	file := api.FileData[f.ComboBoxFiles.ItemIndex()]
	//f.MemoContent.Lines().Add(file.Js[len(file.Js)-1])
	sprintf := fmt.Sprintf("{{%s}}", file.Name)
	f.MemoContent.Lines().Add(sprintf)
}

func (f *TForm2) OnButtonDeleteClick(sender vcl.IObject) {
	if f.ComboBoxFiles.ItemIndex() == -1 {
		return
	}
	file := api.FileData[f.ComboBoxFiles.ItemIndex()]

	if api.WorkDeleteFile(*workBeingDone, file) {
		ok := false
		api.FileData, ok = api.GetWorkFiles(*workBeingDone)
		if !ok {
			return
		}
		f.OnFormShow(sender)
	}
}

func (f *TForm2) OnButtonSendClick(sender vcl.IObject) {
	//note := createContent(api)
	//
	//err := api.SendWork(workBeingDone, f.EditTitle.Text(), note)
	//if err != nil {
	//	vcl.ShowMessage(err.Error())
	//	return
	//}
	//for ci, class := range api.ClassList {
	//	for wi, w := range class.Works {
	//		if w.Id == workBeingDone.Id {
	//			api.ClassList[ci].Works[wi] = *workBeingDone
	//		}
	//	}
	//}
	//ft := Form1.TreeViewClass.Items()
	//for i := int32(0); i < ft.Count(); i++ {
	//	s := ft.Item(i).Text()
	//	s1 := strings.Split(s, " - ")
	//	if len(s1) > 1 {
	//		if s1[1] == workBeingDone.Name {
	//			ft.Item(i).SetText("已完成 - " + workBeingDone.Name)
	//		}
	//	}
	//}
	s := getContent()
	title := Form2.EditTitle.Text()
	if sendClick(title, s) {
		f.ButtonMultipleSendWorks.SetEnabled(false)
		ft := Form1.TreeViewClass.Items()
		for i := int32(0); i < ft.Count(); i++ {
			s := ft.Item(i).Text()
			s1 := strings.Split(s, " - ")
			if len(s1) > 1 {
				if s1[1] == workBeingDone.Name {
					ft.Item(i).SetText("已完成 - " + workBeingDone.Name)
				}
			}
			Form2.ButtonMultipleSendWorks.SetEnabled(false)
		}
		vcl.ShowMessage("作業繳交完成")
	}
}

func (f *TForm2) OnButtonUploadClick(sender vcl.IObject) {
	f.OpenDialogFiles.Execute()
	if f.OpenDialogFiles.Files().Count() == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(int(f.OpenDialogFiles.Files().Count()))
	for i := int32(0); i < f.OpenDialogFiles.Files().Count(); i++ {
		file := f.OpenDialogFiles.Files().Strings(i)

		filePathList = append(filePathList, file)
		go api.WorkUploadFile(file, *workBeingDone, &wg)
	}
	go f.uploadFileWait(&wg)
	f.OpenDialogFiles.Files().Clear()
	f.OpenDialogFiles.SetFileName("")
}

func (f *TForm2) uploadFileWait(wg *sync.WaitGroup) {
	wg.Wait()
	api.FileData, _ = api.GetWorkFiles(*workBeingDone)
	for i, v := range api.FileData {
		for _, file := range filePathList {
			func() {
				var fileNameSplit []string
				sysType := runtime.GOOS
				if sysType == "windows" {
					fileNameSplit = strings.Split(file, "\\")
				} else {
					fileNameSplit = strings.Split(file, "/")
				}
				filename := fileNameSplit[len(fileNameSplit)-1]
				if filename == v.Name {
					api.FileData[i].Path = file
				}
			}()
		}
	}
	go vcl.ShowMessage("檔案上傳完成")
	f.OnFormShow(nil)
}

func (f *TForm2) OnButtonDebugClick(sender vcl.IObject) {
	fmt.Println(twc.title)
}

func (f *TForm2) OnFormCreate(sender vcl.IObject) {

}

func (f *TForm2) OnComboBoxFilesClick(sender vcl.IObject) {

}

func (f *TForm2) OnComboBoxFilesEnter(sender vcl.IObject) {

}

func (f *TForm2) OnFormShow(sender vcl.IObject) {
	if !setting.Debug {
		f.ButtonDebug.Hide()
		//f.ButtonMultipleSendWorks.SetEnabled(false)
	}
	//Form1.Top() + Form1.Height()/2
	//Form1.Left() + Form1.Width()/2
	f.SetTop(Form1.Top() + Form1.Height()/2 - f.Height()/2)
	f.SetLeft(Form1.Left() + Form1.Width()/2 - f.Width()/2)
	f.ComboBoxFiles.Clear()
	for _, v := range api.FileData {
		f.ComboBoxFiles.AddItem(v.Name, nil)
	}
	if f.ComboBoxFiles.Items().Count() > 0 {
		f.ComboBoxFiles.SetItemIndex(0)
	}
}

func (f *TForm2) OnFormCloseQuery(sender vcl.IObject, canClose *bool) {
	f.SetCaption("作業")
	f.ComboBoxFiles.Items().Clear()
	api.FileData = nil
	f.ComboBoxFiles.SetText("")
	f.EditTitle.SetText("")
	f.MemoContent.Lines().Clear()
}

func (f *TForm2) OnFormClose(sender vcl.IObject, closeAction *types.TCloseAction) {
	Form1.Show()
}

func (f *TForm2) OnFormResize(sender vcl.IObject) {

}

func (f *TForm2) OnFormConstrainedResize(sender vcl.IObject, minWidth *int32, minHeight *int32, maxWidth *int32, maxHeight *int32) {
	f.MemoContent.SetWidth(f.Width() - 10 - f.MemoContent.Left())
	f.MemoContent.SetHeight(f.Height() - 17 - f.MemoContent.Top())
}

func (f *TForm2) OnButtonMultipleSendWorksClick(sender vcl.IObject) {
	Form3.Show()
}
