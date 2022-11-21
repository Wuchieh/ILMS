package pkgui

import (
	"fmt"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"log"
	"strings"
	"sync"
	"wuchieh.com/ilmstool/ilmsAPI"
)

// ::private::
type TForm3Fields struct {
}

func (f *TForm3) OnButtonEnterClick(sender vcl.IObject) {
	title := Form2.EditTitle.Text()
	//for i := int32(0); i < Form2.MemoContent.Lines().Count(); i++ {
	//	if i == 0 {
	//		note += Form2.MemoContent.Lines().S(i)
	//	} else {
	//		note += "<div>" + Form2.MemoContent.Lines().S(i) + "<br></div>"
	//	}
	//}
	content := getContent()
	var wg sync.WaitGroup
	for i := int32(0); i < f.ListBoxUserList.Items().Count(); i++ {
		if f.ListBoxUserList.Selected(i) {
			wg.Add(1)
			username := f.ListBoxUserList.Items().Strings(i)
			go func() {
				sendWork(username, title, content)
				defer wg.Done()
			}()
		}
	}
	wg.Wait()
	Form2.ButtonMultipleSendWorks.SetEnabled(false)
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
	vcl.ShowMessage("批量繳交完成")
}

func sendWork(u, title string, content []string) {
	if usernameList[u].Username == api.Username {
		fmt.Println("節流")
		sendClick(title, content)
		return
	}
	fmt.Println("無節流")
	api := usernameList[u]
	work := ilmsAPI.Work{
		Id:       workBeingDone.Id,
		Name:     workBeingDone.Name,
		FolderID: workBeingDone.FolderID,
		Done:     false,
		Expired:  false,
		Url:      workBeingDone.Url,
		Cid:      "",
	}
	title = strings.ReplaceAll(title, "{{帳號}}", api.Username)
	api.GetWorkExpired(&work)
	work.Done = api.GetWorkStatus(work)
	if !api.ClearWorkFiles(work) {
		fmt.Println(u, "作業繳交失敗, 清除檔案出錯")
	}
	var fileWg sync.WaitGroup
	for _, file := range filePathList {
		fileWg.Add(1)
		go func(file string) {
			api.WorkUploadFile(file, work)

			defer fileWg.Done()
		}(file)
	}
	fileWg.Wait()
	ok := false
	api.FileData, ok = api.GetWorkFiles(work)
	if !ok {
		fmt.Println("已上傳文件 檔案讀取錯誤")
		return
	}
	note := contentFormat(api, content)
	err := api.SendWork(&work, title, note)
	if err != nil {
		log.Println(err)
		return
	}
}

func (f *TForm3) OnButtonCancelClick(sender vcl.IObject) {

}

func (f *TForm3) OnFormCreate(sender vcl.IObject) {

}

func (f *TForm3) OnFormShow(sender vcl.IObject) {
	f.SetTop(Form1.Top() + Form1.Height()/2 - f.Height()/2)
	f.SetLeft(Form1.Left() + Form1.Width()/2 - f.Width()/2)
}

func (f *TForm3) OnFormClose(sender vcl.IObject, closeAction *types.TCloseAction) {
	Form1.Show()
}

func (f *TForm3) OnFormCloseQuery(sender vcl.IObject, canClose *bool) {

}
