package pkgui

import (
	"fmt"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"strings"
	"wuchieh.com/ilmstool/ilmsAPI"
)

// ::private::
type TForm1Fields struct {
}

func (f *TForm1) OnButtonDebugClick(sender vcl.IObject) {
	//Form2.Show()
	for _, v := range usernameList {
		for _, cookie := range v.Cookie {
			fmt.Println(cookie)
		}
	}
}

func (f *TForm1) OnFormCreate(sender vcl.IObject) {
	f.SetTop(vcl.Screen.Height()/2 - f.Height()/2)
	f.SetLeft(vcl.Screen.Width()/2 - f.Width()/2)
	f.EditUsername.SetText(setting.Username)
	f.EditPassword.SetText(setting.Password)
	if !setting.Debug {
		f.ButtonDebug.Hide()
		//f.ComboBoxUserList.SetEnabled(false)
	}
	f.ButtonLogout.SetEnabled(false)
}

func (f *TForm1) OnButtonLoginClick(sender vcl.IObject) {
	username := f.EditUsername.Text()
	password := f.EditPassword.Text()
	username = strings.ToUpper(username)
	if username == "" || password == "" {
		vcl.ShowMessage("帳號及密碼不得為空")
		return
	}
	f.MemoConsole.Lines().Add("登入中...")
	NewApi := ilmsAPI.New()
	if !NewApi.Login(username, password) {
		go vcl.ShowMessage("登入失敗")
		f.MemoConsole.Lines().Add("登入失敗")
		return
	}
	//f.ButtonLogin.SetEnabled(false)
	f.ButtonLogout.SetEnabled(true)
	f.MemoConsole.Lines().Add(username + " 登入成功")
	f.ButtonUpdata.SetEnabled(true)
	//vcl.ShowMessage(username + "登入成功")
	if !NewApi.GetClassList() {
		go vcl.ShowMessage("課程列表取得失敗")
		f.MemoConsole.Lines().Add("課程列表取得失敗")
	}
	f.MemoConsole.Lines().Add("課程列表取得成功")
	f.TreeViewClass.Items().Clear()
	for _, v := range NewApi.ClassList {
		add := f.TreeViewClass.Items().Add(nil, v.Name)
		f.TreeViewClass.Items().AddChild(add, "1")
	}
	NewApi.Username = username
	usernameList[username] = NewApi
	api = NewApi
	f.ComboBoxUserList.SetText(username)
	usernameListGenerator()
}

func (f *TForm1) OnButtonLogoutClick(sender vcl.IObject) {
	Form2.Close()
	Form3.Close()
	index := f.ComboBoxUserList.ItemIndex()
	if index == -1 {
		return
	}
	username := f.ComboBoxUserList.Items().S(index)
	i := usernameList[username]
	if i.Client != nil {
		i.Logout()
		delete(usernameList, username)
	}
	f.MemoConsole.Lines().Add(username + " 登出成功")
	if len(usernameList) == 0 {
		api = ilmsAPI.New()
		f.ButtonUpdata.SetEnabled(false)
	}
	usernameListGenerator()
	f.TreeViewClass.Items().Clear()
	if f.ComboBoxUserList.Items().Count() == 0 {
		f.ButtonLogout.SetEnabled(false)
	}
	text := ""
	if f.ComboBoxUserList.Items().Count() == 0 {
		return
	}
	text = f.ComboBoxUserList.Items().S(0)
	f.ComboBoxUserList.SetText(text)
	f.OnComboBoxUserListChange(nil)
}

func (f *TForm1) OnTreeViewClassExpanded(sender vcl.IObject, node *vcl.TTreeNode) {
	if node.Item(0).Text() != "1" {
		return
	}
	f.MemoConsole.Lines().Add("正在讀取 " + node.Text() + " 作業清單...")
	class := api.ClassList[node.Index()]
	if !api.GetWorkList(class) {
		vcl.ShowMessage("發生錯誤")
	}
	for _, v := range class.Works {
		if v.Done {
			f.TreeViewClass.Items().AddChild(node, "已完成 - "+v.Name)
		} else {
			f.TreeViewClass.Items().AddChild(node, "未完成 - "+v.Name)
		}
	}
	f.MemoConsole.Lines().Add("讀取完成作業清單...")
	node.Items(0).Delete()
}

func (f *TForm1) OnTreeViewClassDblClick(sender vcl.IObject) {
	node := f.TreeViewClass.Selected()
	if node == nil {
		//vcl.ShowMessage("請先選擇作業")
		return
	}
	for node.GetPrevSibling() != nil {
		node = node.GetPrevSibling()
	}
	if node.GetPrev() != nil {
		//fmt.Println(node.GetPrev().Text())
		//fmt.Println(f.TreeViewClass.Selected().Text())
		Form2.Close()
		workIndex := f.TreeViewClass.Selected().Index()
		classIndex := node.GetPrev().Index()
		selectWork := api.ClassList[classIndex].Works[workIndex]
		fmt.Println(selectWork.Url)
		if api.GetWorkExpired(&selectWork) {
			//fmt.Println("可編輯", selectWork)
			if selectWork.Done {
				Form2.ButtonMultipleSendWorks.SetEnabled(false)
			} else {
				Form2.ButtonMultipleSendWorks.SetEnabled(true)
			}
			showDoWorkForm(&selectWork)
		} else {
			vcl.ShowMessage(selectWork.Name + " 不可編輯")
		}
	}
}

func (f *TForm1) OnButtonUpdataClick(sender vcl.IObject) {
	f.TreeViewClass.Items().Clear()
	for _, v := range api.ClassList {
		add := f.TreeViewClass.Items().Add(nil, v.Name)
		f.TreeViewClass.Items().AddChild(add, "1")
	}
	for _, class := range api.ClassList {
		class.Works = []ilmsAPI.Work{}
	}
}

func (f *TForm1) OnFormResize(sender vcl.IObject) {

}

func (f *TForm1) OnFormConstrainedResize(sender vcl.IObject, minWidth *int32, minHeight *int32, maxWidth *int32, maxHeight *int32) {
	f.MemoConsole.SetHeight(f.Height() - f.MemoConsole.Top() - 7)
}

func (f *TForm1) OnMenuItem1Click(sender vcl.IObject) {
	f.OnTreeViewClassDblClick(sender)
}

func (f *TForm1) OnMenuItem2Click(sender vcl.IObject) {
	node := f.TreeViewClass.Selected()
	if node == nil {
		vcl.ShowMessage("請先選擇作業")
		return
	}
	for node.GetPrevSibling() != nil {
		node = node.GetPrevSibling()
	}
	if node.GetPrev() != nil {
		workIndex := f.TreeViewClass.Selected().Index()
		classIndex := node.GetPrev().Index()
		selectWork := api.ClassList[classIndex].Works[workIndex]
		if api.GetWorkExpired(&selectWork) && api.GetWorkStatus(selectWork) {
			if api.DeleteWork(&selectWork) {
				vcl.ShowMessage("刪除成功")
				Form2.Close()
				Form3.Close()
				f.TreeViewClass.Selected().SetText("未完成 - " + selectWork.Name)
				api.ClassList[classIndex].Works[workIndex] = selectWork
			}
		} else {
			vcl.ShowMessage(selectWork.Name + " 不可編輯")
		}
	}
}

func (f *TForm1) OnEditUsernameKeyPress(sender vcl.IObject, key *types.Char) {
	if *key == 13 {
		f.OnButtonLoginClick(sender)
	}
}

func (f *TForm1) OnEditPasswordKeyPress(sender vcl.IObject, key *types.Char) {
	if *key == 13 {
		f.OnButtonLoginClick(sender)
	}
}

func (f *TForm1) OnTreeViewClassMouseDown(sender vcl.IObject, button types.TMouseButton, shift types.TShiftState, x int32, y int32) {
	//if button == 1 {
	//	fmt.Println(x, y)
	//	for i := int32(0); i < f.TreeViewClass.Items().Count(); i++ {
	//		item := f.TreeViewClass.Items().Item(i)
	//		if item.Top() != 0 && y >= item.Top() && y < item.Bottom() && strings.Contains(item.Text(), " - ") {
	//			f.TreeViewClass.SetSelected(item)
	//			f.PopupMenu1.Popup2()
	//		}
	//	}
	//}
}

func (f *TForm1) OnComboBoxUserListChange(sender vcl.IObject) {
	index := f.ComboBoxUserList.ItemIndex()
	item := f.ComboBoxUserList.Items().Strings(index)
	if api.Username == usernameList[item].Username && f.TreeViewClass.Items().Count() != 0 {
		return
	}
	api = usernameList[item]
	f.TreeViewClass.Items().Clear()
	for _, v := range api.ClassList {
		add := f.TreeViewClass.Items().Add(nil, v.Name)
		f.TreeViewClass.Items().AddChild(add, "1")
	}
	Form2.Close()
	Form3.Close()
}

func (f *TForm1) OnButtonDeleteWorkClick(sender vcl.IObject) {
	item := f.TreeViewClass.Selected()
	node := item
	if strings.Contains(item.Text(), " - ") {
		for node.GetPrevSibling() != nil {
			node = node.GetPrevSibling()
		}
		classIndex := node.GetPrev().Index()
		work := api.ClassList[classIndex].Works[item.Index()]
		if api.GetWorkExpired(&work) && api.GetWorkStatus(work) {
			if !api.DeleteWork(&work) {
				vcl.ShowMessage("作業刪除失敗")
				return
			}
			vcl.ShowMessage("刪除成功")
			Form2.Close()
			Form3.Close()
			f.TreeViewClass.Selected().SetText("未完成 - " + work.Name)
			api.ClassList[classIndex].Works[item.Index()] = work
		} else {
			vcl.ShowMessage("作業簿可編輯")
		}
	} else {
		vcl.ShowMessage("不可刪除課程")
	}
}

func (f *TForm1) OnFormClose(sender vcl.IObject, closeAction *types.TCloseAction) {
}

func (f *TForm1) OnFormCloseQuery(sender vcl.IObject, canClose *bool) {
	if vcl.MessageDlg("是否退出?", types.MtCustom, types.MbYes, types.MbNo) == types.MrYes {
		if len(usernameList) != 0 && vcl.MessageDlg("是否儲存已登入的資訊?", types.MtCustom, types.MbYes, types.MbNo) == types.MrYes {
			saveUserInfo()
		}
	} else {
		*canClose = false
	}
}

func (f *TForm1) OnFormShow(sender vcl.IObject) {
	getUserInfoFile()
	//for i, v := range usernameList {
	//	fmt.Println(i)
	//	for _, v2 := range v.ClassList {
	//		fmt.Println(v2.Name)
	//	}
	//}
}
