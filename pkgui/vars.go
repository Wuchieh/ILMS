package pkgui

import "wuchieh.com/ilmstool/ilmsAPI"

var (
	api          ilmsAPI.IlmsAPI
	usernameList = make(map[string]ilmsAPI.IlmsAPI)
	//fileData      []ilmsAPI.FileData
	workBeingDone *ilmsAPI.Work // 正在編輯的作業
	setting       Settings
	filePathList  []string
	twc           TempWorkContent // 作業內容緩存
)

type Settings struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	MaxConnect int    `json:"maxConnect"`
	Debug      bool   `json:"debug,omitempty"`
}

type TempWorkContent struct {
	title    string
	content  []string
	fileData []string //此處只有存檔案id
}
