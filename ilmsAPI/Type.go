package ilmsAPI

import (
	"net/http"
)

const (
	LoginUrl          = "https://elearning.uch.edu.tw/sys/lib/ajax/login_submit.php"
	LogoutUrl         = "https://elearning.uch.edu.tw/sys/lib/ajax/logout.php"
	HomeUrl           = "https://elearning.uch.edu.tw/home.php"
	HostUrl           = "https://elearning.uch.edu.tw"
	WorkListUrl       = "https://elearning.uch.edu.tw/course.php?courseID="                // + Class.Id &f=hwlist"
	InsertWorkUrl     = "https://elearning.uch.edu.tw/course/doc_insert.php?folderID="     // + Work.Id
	EditWorkUrl       = "https://elearning.uch.edu.tw/course/doc_edit.php?id="             // + Work.Cid
	GetWorkFilesUrl   = "https://elearning.uch.edu.tw/course/http_get_doc_attach_file.php" // POST id=Work.Cid
	WorkDeleteFileUrl = "https://elearning.uch.edu.tw/course/http_doc_attach_delete.php"   // POST id=Work.Cid&attachIDs=FileData.Id
	UploadFileUrl     = "https://elearning.uch.edu.tw/course/upload_attach.php?id="        // +Work.Cid POST
	SendEditWorkUrl   = "https://elearning.uch.edu.tw/course/doc_edit.php"                 // POST
	DeleteWorkUrl     = "https://elearning.uch.edu.tw/course/doc_delete.php?id="           //+ Work.Cid POST fmSubmit=yes
)

var (
	MaxConnect        = 0 // 設定最大連接數(讀取作業列表/上傳檔案)[0 無限制]
	connect           = make(chan int, MaxConnect)
	loginStatus       LoginStatus
	getFileStatus     GetFileStatus
	deleteFilesStatus DeleteFilesStatus
)

type Class struct {
	Id    string // courseID 課程ID
	Name  string
	Url   string // https://elearning.uch.edu.tw/course/+ Class.Id
	Works []Work
}

type Work struct {
	Id       string // hw 作業ID
	Name     string
	FolderID string // 課程ID Class.Id
	Done     bool   // 是否已經完成
	Expired  bool
	Url      string // 此為 已交名單頁面
	Cid      string // 專屬ID 只有準備繳交 或 已經繳交的作業會有
}

type IlmsAPI struct {
	Username      string // 學號
	Client        *http.Client
	ClassList     []*Class // 課程清單
	Cookie        []*http.Cookie
	StringChannel chan string
	FileData      []FileData
}

type FileData struct {
	Id   string
	Name string
	Js   []string
	Path string
}

type LoginStatus struct {
	Ret struct {
		Status  string `json:"status"`
		Email   string `json:"email,omitempty"`
		Name    string `json:"name,omitempty"`
		Phone   string `json:"phone,omitempty"`
		Info    string `json:"info,omitempty"`
		Unknow1 string `json:"unknow1,omitempty"`
		Unknow2 string `json:"unknow2,omitempty"`
		DivName string `json:"divName,omitempty"`
		DivCode string `json:"divCode,omitempty"`
		Focus   string `json:"focus,omitempty"`
		Msg     string `json:"msg,omitempty"`
	} `json:"ret"`
}

type GetFileStatus struct {
	Ret struct {
		Status       string `json:"status"`
		AttachIDList string `json:"attachID_list"`
		FileContent  string `json:"fileContent"`
		Msg          string `json:"msg"`
	} `json:"ret"`
}

type DeleteFilesStatus struct {
	Ret struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	} `json:"ret"`
}
