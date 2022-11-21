package pkgui

import (
	"encoding/json"
	"fmt"
	"github.com/ying32/govcl/vcl"
	"log"
	"net/http"
	"os"
	"strings"
	"wuchieh.com/ilmstool/ilmsAPI"
)

func init() {
	api = ilmsAPI.New()
	file, err := os.ReadFile("settings.json")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(file, &setting)
	ilmsAPI.SetMaxConnect(setting.MaxConnect)
	if err != nil {
		log.Println(err)
		return
	}
}

func showDoWorkForm(w *ilmsAPI.Work) {
	Form2.SetCaption(w.Name)
	ok := false
	twc = TempWorkContent{
		title:    "",
		content:  []string{},
		fileData: []string{},
	}
	if w.Done {
		api.FileData, ok = api.GetWorkFiles(*w)
		title, content, _ := api.GetWorkContent(*w)
		Form2.EditTitle.SetText(title)
		Form2.MemoContent.SetText(content)
		twc.title = title
		for i := int32(0); i < Form2.MemoContent.Lines().Count(); i++ {
			twc.content = append(twc.content, Form2.MemoContent.Lines().S(i))
		}
		for _, v := range api.FileData {
			twc.fileData = append(twc.fileData, v.Name)
		}
	} else {
		ok = true
	}
	if ok {
		Form2.EditTitle.SetText("{{帳號}}")
		Form2.Show()
		workBeingDone = w
	}
}

func usernameListGenerator() {
	Form1.ComboBoxUserList.Items().Clear()
	Form3.ListBoxUserList.Items().Clear()
	if len(usernameList) == 0 {
		Form1.ComboBoxUserList.SetText("切換帳號")
		Form1.ComboBoxUserList.SetSelStart(-1)
		return
	}
	for i, _ := range usernameList {
		Form1.ComboBoxUserList.AddItem(i, nil)
		Form3.ListBoxUserList.Items().Add(i)
	}
}

//func createContent(api ilmsAPI.IlmsAPI) string {
//	note := ""
//	var StringList []string
//	for i := int32(0); i < Form2.MemoContent.Lines().Count(); i++ {
//		line := Form2.MemoContent.Lines().S(i)
//		StringList = append(StringList, line)
//	}
//	for _, v := range api.FileData {
//		for i, s := range StringList {
//			StringList[i] = strings.ReplaceAll(s, "{{"+v.Name+"}}", v.Js[len(v.Js)-1])
//		}
//	}
//	for i, v := range StringList {
//		if i == 0 {
//			note += v
//		} else {
//			note += "<div>" + v + "<br></div>"
//		}
//	}
//	return note
//}

func getContent() []string {
	var StringList []string
	for i := int32(0); i < Form2.MemoContent.Lines().Count(); i++ {
		line := Form2.MemoContent.Lines().S(i)
		StringList = append(StringList, line)
	}
	return StringList
}

func contentFormat(api ilmsAPI.IlmsAPI, content []string) string {
	note := ""
	for _, v := range api.FileData {
		if len(v.Js) == 0 || len(content) == 0 {
			break
		}
		for i, s := range content {
			fmt.Println(i, s)
			content[i] = strings.ReplaceAll(s, "{{"+v.Name+"}}", v.Js[len(v.Js)-1])
		}
	}
	for i, s := range content {
		content[i] = strings.ReplaceAll(s, "{{帳號}}", api.Username)
	}
	for i, v := range content {
		if i == 0 {
			note += v
		} else {
			note += "<div>" + v + "<br></div>"
		}
	}
	return note
}

func sendClick(title string, s []string) bool {
	//note := createContent(api)
	note := contentFormat(api, s)
	title = strings.ReplaceAll(title, "{{帳號}}", api.Username)
	// true = 未更動 false = 有更動
	if comparisonWork(title, s) {
		twc.title = title
		twc.content = s
		twc.fileData = []string{}

		for _, v := range api.FileData {
			twc.fileData = append(twc.fileData, v.Id)
		}
		vcl.ShowMessage("作業未有編輯")
		return false
	}
	twc.title = title
	twc.content = s
	twc.fileData = []string{}

	for _, v := range api.FileData {
		twc.fileData = append(twc.fileData, v.Id)
	}
	err := api.SendWork(workBeingDone, title, note)
	if err != nil {
		go vcl.ShowMessage(err.Error())
		return false
	}
	for ci, class := range api.ClassList {
		for wi, w := range class.Works {
			if w.Id == workBeingDone.Id {
				api.ClassList[ci].Works[wi] = *workBeingDone
			}
		}
	}
	return true
}

// 比對作業內容 如果沒有做出變更就會回傳true 如果有變更 就會回傳false
func comparisonWork(title string, s []string) bool {
	//比較標題
	if twc.title != title {
		return false
	}
	//比較內容
	if func(a, b []string) bool {
		if len(a) != len(b) {
			return true
		}
		for i := range a {
			if a[i] != b[i] {
				return true
			}
		}
		return false
	}(twc.content, s) {
		return false
	}
	//比較檔案
	if func() bool {
		if len(api.FileData) != len(twc.fileData) {
			return true
		}
		for _, v := range api.FileData {
			if !contains(twc.fileData, v.Id) {
				return false
			}
		}
		return false
	}() {
		return false
	}
	return true
}

// 查找List內是否有指定的值
func contains[T comparable](a []T, b T) bool {

	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// 儲存已登入的用戶訊息
func saveUserInfo() {
	userCookies := make(map[string][]http.Cookie)
	for i, v := range usernameList {
		for _, cookie := range v.Cookie {
			userCookies[i] = append(userCookies[i], *cookie)
		}
	}
	bytes, err := json.Marshal(userCookies)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.WriteFile("userinfo", bytes, 0666)
	if err != nil {
		log.Println(err)
		return
	}
}

func getUserInfoFile() {
	userCookies := make(map[string][]http.Cookie)
	file, err := os.ReadFile("userinfo")
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &userCookies)
	if err != nil {
		log.Println(err)
		return
	}
	for username, cookies := range userCookies {
		NewApi := ilmsAPI.New()
		NewApi.Username = username
		for _, cookieInfo := range cookies {
			cookie := &http.Cookie{
				Name:       cookieInfo.Name,
				Value:      cookieInfo.Value,
				Path:       cookieInfo.Path,
				Domain:     cookieInfo.Domain,
				Expires:    cookieInfo.Expires,
				RawExpires: cookieInfo.RawExpires,
				MaxAge:     cookieInfo.MaxAge,
				Secure:     cookieInfo.Secure,
				HttpOnly:   cookieInfo.HttpOnly,
				SameSite:   cookieInfo.SameSite,
				Raw:        cookieInfo.Raw,
				Unparsed:   cookieInfo.Unparsed,
			}
			NewApi.Cookie = append(NewApi.Cookie, cookie)
		}
		if NewApi.GetClassList() {
			usernameList[username] = NewApi
		}
	}
	if len(usernameList) != 0 {
		usernameListGenerator()
		Form1.ButtonUpdata.SetEnabled(true)
		Form1.ButtonLogout.SetEnabled(true)
	}
}
