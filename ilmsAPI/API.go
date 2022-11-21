package ilmsAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	defer fmt.Println("wuchieh.com製作 ilmsAPI 以初始化完成")
}

func Ping() { // 測試用
	fmt.Println("PONG")
}

//func (a *IlmsAPI) Init() {
//	client := &http.Client{}
//	a.Client = client
//}

// New 創建新的 IlmsAPI
func New() IlmsAPI {
	client := &http.Client{}
	//a.Client = client
	return IlmsAPI{
		Username:      "",
		Client:        client,
		ClassList:     nil,
		Cookie:        nil,
		StringChannel: nil,
		FileData:      []FileData{},
	}
}

// SetMaxConnect 設定最大連接數
func SetMaxConnect(n int) {
	MaxConnect = n
	connect = make(chan int, MaxConnect)
}

// Login 登入
func (a *IlmsAPI) Login(Username, Password string) bool {
	loginData := fmt.Sprintf("account=%s&password=%s&secCode=na&stay=1", Username, Password)
	request, err := http.NewRequest("POST", LoginUrl, strings.NewReader(loginData))
	if err != nil {
		return printError(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := a.Client.Do(request)
	if err != nil {
		return printError(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	a.Cookie = response.Cookies()
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return printError(err)
	}
	err = json.Unmarshal(readAll, &loginStatus)
	if err != nil {
		return printError(err)
	}
	if loginStatus.Ret.Status == "true" {
		return true
	}
	return false
}

// Logout 登出
func (a *IlmsAPI) Logout() bool {
	request, err := http.NewRequest("POST", LogoutUrl, strings.NewReader("rk=0"))
	if err != nil {
		return printError(err)
	}
	request.Header.Set("Content-Type", "text/html")
	for _, v := range a.Cookie {
		request.Header.Set("Cookie", v.String())
	}
	//request.Header.Set("Cookie", a.Cookie[0].String())
	response, err := a.Client.Do(request)
	if err != nil {
		return printError(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	a.Cookie = response.Cookies()
	//readAll, err := io.ReadAll(response.Body)
	//if err != nil {
	//	return printError(err)
	//}
	return false
}

// GetClassList 取得課程清單
func (a *IlmsAPI) GetClassList() bool {
	request, err := http.NewRequest("GET", HomeUrl, nil)
	if err != nil {
		return printError(err)
	}
	a.setCookie(request)
	response, err := a.Client.Do(request)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err != nil {
		return printError(err)
	}
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return printError(err)
	}
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(readAll))
	if err != nil {
		return printError(err)
	}
	r, _ := regexp.Compile("\\d+")
	reader.Find(".mnuItem>a").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "成績查詢" {
			var class Class
			href, _ := s.Attr("href")
			class.Id = r.FindString(href)
			class.Name = s.Text()
			class.Url = HostUrl + href
			a.ClassList = append(a.ClassList, &class)
		}
	})
	if len(a.ClassList) == 0 {
		return false
	}
	return true
}

// GetWorkList 從課程中取得作業清單
func (a *IlmsAPI) GetWorkList(c *Class) bool {
	if len(c.Works) > 0 {
		return true
	}
	newUrl := WorkListUrl + c.Id + "&f=hwlist"
	request, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		return printError(err)
	}
	a.setCookie(request)
	response, err := a.Client.Do(request)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	if err != nil {
		return printError(err)
	}
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return printError(err)
	}
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(readAll))
	if err != nil {
		return printError(err)
	}
	r, _ := regexp.Compile("hw=\\d+")
	var wg sync.WaitGroup
	//if MaxConnect != 0 {
	//	for i := 0; i < MaxConnect; i++ {
	//		connect <- i
	//	}
	//}
	reader.Find("#main > div.tableBox > table > tbody > tr > td:nth-child(2) > a:nth-child(1)").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func() {
			if MaxConnect != 0 {
				connect <- 1
			}
			if s.Text() != "標題" {
				href, _ := s.Attr("href")
				workid := r.FindString(href)[3:]
				work := Work{
					Id:       workid,
					Name:     s.Text(),
					FolderID: c.Id,
					Done:     false,
					Url:      HostUrl + href,
					Cid:      "",
				}
				work.Done = a.GetWorkStatus(work)
				c.Works = append(c.Works, work)
			}
			defer func() {
				if MaxConnect != 0 {
					time.Sleep(time.Second)
					<-connect
				}
				wg.Done()
			}()
		}()
	})
	wg.Wait()
	for len(connect) > 0 {
		<-connect
	}
	sort.SliceStable(c.Works, func(i, j int) bool {
		atoi, _ := strconv.Atoi(c.Works[i].Id)
		atoj, _ := strconv.Atoi(c.Works[j].Id)
		return atoi < atoj
	})
	return true
}

// GetWorkStatus 獲取作業的狀態(是否已經繳交)
func (a *IlmsAPI) GetWorkStatus(w Work) bool {
	newUrl := WorkListUrl + w.FolderID + "&f=hw_doclist&hw=" + w.Id
	fmt.Println(newUrl)
	request, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		return printError(err)
	}
	a.setCookie(request)
	response, err := a.Client.Do(request)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err != nil {
		return printError(err)
	}
	//a.Cookie = response.Cookies()
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return printError(err)
	}
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(readAll))
	if err != nil {
		return printError(err)
	}
	if reader.Find("#main > div.toolWrapper > a:nth-child(4)").Text() == "交作業" {
		return false
	}
	return true
}

// GetWorkExpired 查詢作業是否已過期
func (a *IlmsAPI) GetWorkExpired(w *Work) bool { // 此處會編輯Work.Cid
	if w.Done {
		newUrl := w.Url
		err, readAll := a.httpGet(newUrl)
		if err != nil {
			return printError(err)
		}
		reader, err := goquery.NewDocumentFromReader(bytes.NewReader(readAll))
		if err != nil {
			return printError(err)
		}
		href, _ := reader.Find("#main > div.toolWrapper > a:nth-child(4)").Attr("href")
		newUrl = HostUrl + href
		//fmt.Println("newUrl:", newUrl)
		err, readAll = a.httpGet(newUrl)
		if err != nil {
			return printError(err)
		}
		reader, err = goquery.NewDocumentFromReader(bytes.NewReader(readAll))
		if err != nil {
			return printError(err)
		}
		if reader.Find("#tools > a:nth-child(1)").Text() == "編輯" || reader.Find("#tools > a").Text() == "編輯" {
			w.Expired = true
		} else {
			w.Expired = false
		}
		r, _ := regexp.Compile("cid=\\w+")
		w.Cid = r.FindString(newUrl)[4:]
		return w.Expired
	} else {
		newUrl := InsertWorkUrl + w.Id
		err, readAll := a.httpGet(newUrl)
		if err != nil {
			return printError(err)
		}
		reader, err := goquery.NewDocumentFromReader(bytes.NewReader(readAll))
		if err != nil {
			return printError(err)
		}
		if reader.Text() == "此作業已過繳交期限，無法繳交!" {
			w.Expired = false
		} else {
			w.Expired = true
			cid, _ := reader.Find("#id").Attr("value")
			w.Cid = cid
		}
		return w.Expired
	}
}

func (a *IlmsAPI) httpGet(newUrl string) (error, []byte) {
	request, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		return nil, nil
	}
	a.setCookie(request)
	response, err := a.Client.Do(request)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	//a.Cookie = response.Cookies()
	if err != nil {
		return nil, nil
	}
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil
	}
	return err, readAll
}

// GetWorkFiles 取得作業內已上傳的檔案
func (a *IlmsAPI) GetWorkFiles(w Work) ([]FileData, bool) {
	request, err := http.NewRequest("POST", GetWorkFilesUrl, strings.NewReader("id="+w.Cid))
	if err != nil {
		return []FileData{}, printError(err)
	}
	a.setCookie(request)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err := a.Client.Do(request)
	if err != nil {
		return []FileData{}, printError(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return []FileData{}, printError(err)
	}
	err = json.Unmarshal(readAll, &getFileStatus)
	if err != nil {
		return []FileData{}, printError(err)
	}
	return a.filesFormat(getFileStatus)
}

func (a *IlmsAPI) filesFormat(data GetFileStatus) ([]FileData, bool) {
	fileIds := data.Ret.AttachIDList
	fileIdList := strings.Split(fileIds, ",")
	var files []FileData
	var fileData FileData
	r, _ := regexp.Compile("<w*.*>")
	reader, err := goquery.NewDocumentFromReader(strings.NewReader(data.Ret.FileContent))
	if err != nil {
		return nil, false
	}
	reader.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "" {
			if fileData.Name != "" {
				files = append(files, fileData)
			}
			fileData.Name = s.Text()
		} else {
			js, _ := s.Attr("href")
			findString := r.FindString(js)
			fileData.Js = append(fileData.Js, findString)
		}
	})
	files = append(files, fileData)
	if len(fileIdList) == len(files) {
		for i := range files {
			//fmt.Printf("v.Id:%s, fileIdList[i]:%s\n", v.Id, fileIdList[i])
			files[i].Id = fileIdList[i]
		}
		//for _, v := range files {
		//	fmt.Println(v.Id == "")
		//}
		return files, true
	} else {
		return nil, false
	}
}

// GetWorkContent 取得作業的內容(此為已繳交的作業才有)
// 會返回 標題及內容以及是否成功取得
func (a *IlmsAPI) GetWorkContent(w Work) (string, string, bool) {
	newUrl := EditWorkUrl + w.Cid
	err, readAll := a.httpGet(newUrl)
	if err != nil {
		return "", "", printError(err)
	}
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(readAll))
	if err != nil {
		return "", "", printError(err)
	}
	content := reader.Find("#fmNote").Text()
	//fmt.Println(content)
	title, _ := reader.Find("#fmTitle").Attr("value")
	html := strings.ReplaceAll(content, "  ", "")
	html = strings.ReplaceAll(html, "<div>", "\n<div>")
	html = strings.ReplaceAll(html, "<div>", "")
	html = strings.ReplaceAll(html, "</div>", "")
	html = strings.ReplaceAll(html, "<br>", "")
	tagAHref, _ := regexp.Compile("/+.*id=\\d+")
	tagATarget, _ := regexp.Compile("=_\\w+")
	tagImg, _ := regexp.Compile("/sysdata.*?>")
	//tagBr, _ := regexp.Compile("<br.*?>")
	fmt.Println(a.FileData)
	for _, v := range a.FileData {
		if v.Name == "" {
			continue
		}
		js := v.Js[len(v.Js)-1]
		html = strings.ReplaceAll(html, js, "{{"+v.Name+"}}")
		if strings.HasPrefix(js, "<a") {
			href := tagAHref.FindString(js)
			js = tagAHref.ReplaceAllString(js, `"`+href+`"`)
			target := tagATarget.FindString(js)[1:]
			js = tagATarget.ReplaceAllString(js, `="`+target+`"`)
			//fmt.Println("a", js)
		} else if strings.HasPrefix(js, "<i") {
			src := tagImg.FindString(js)
			src = src[:len(src)-1]
			js = tagImg.ReplaceAllString(js, `"`+src+`">`)
			//fmt.Println("img", js)
		} else {
			//fmt.Println("未發生改變", js)
		}
		html = strings.ReplaceAll(html, js, "{{"+v.Name+"}}")
	}
	return title, html, true
}

// WorkDeleteFile 刪除作業內的檔案
func (a *IlmsAPI) WorkDeleteFile(w Work, f FileData) bool {
	s := fmt.Sprintf("id=%s&attachIDs=%s", w.Cid, f.Id)
	request, err := http.NewRequest("POST", WorkDeleteFileUrl, strings.NewReader(s))
	if err != nil {
		return printError(err)
	}
	a.setCookie(request)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err := a.Client.Do(request)
	if err != nil {
		return printError(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return printError(err)
	}
	err = json.Unmarshal(readAll, &deleteFilesStatus)
	if err != nil {
		return printError(err)
	}
	if deleteFilesStatus.Ret.Status == "true" {
		return true
	} else {
		return false
	}
}

// WorkUploadFile 上傳檔案至作業
func (a *IlmsAPI) WorkUploadFile(file string, w Work, options ...*sync.WaitGroup) bool {
	if MaxConnect != 0 {
		connect <- 1
	}
	defer func() {
		if MaxConnect != 0 {
			<-connect
		}
	}()
	newUrl := UploadFileUrl + w.Cid
	oF := openFile(file)
	values := map[string]io.Reader{
		"Filedata":        oF,
		"fmSubmit":        strings.NewReader("yes"),
		"id":              strings.NewReader(w.Cid), // Work.Cid
		"from":            strings.NewReader(""),
		"resType":         strings.NewReader("1"),
		"fmTitle":         strings.NewReader(""), // 作業的標題
		"fmAttachDisplay": strings.NewReader("1"),
		"fmNote":          strings.NewReader(""), // 作業的內容

	}
	err := upload(a, newUrl, values)
	defer func() {
		for _, wg := range options {
			wg.Done()
		}
	}()
	if err != nil {
		return printError(err)
	}
	return true
}

// SendWork 繳交作業
func (a *IlmsAPI) SendWork(w *Work, title, content string) error {
	if title == "" {
		return errors.New("標題不得為空")
	}
	if w.Done {
		ss := fmt.Sprintf("fmSubmit=yes&id=%s&from=&resType=1&fmTitle=%s&Filedata=&fmAttachDisplay=1&fmNote=%s", w.Cid, title, content)
		return a.editWork(w, ss)
	} else {
		ss := fmt.Sprintf("fmSubmit=yes&id=%s&from=&resType=1&fmTitle=%s&Filedata=&fmAttachDisplay=1&fmNote=%s&folderID=%s&addWiki=0&actID=0", w.Cid, title, content, w.FolderID)
		return a.sendNewWork(w, ss)
	}
}

func (a *IlmsAPI) sendNewWork(w *Work, ss string) (err error) {
	header := map[string]string{
		`Accept`:                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		`Accept-Encoding`:           `gzip, deflate, br`,
		`Accept-Language`:           `zh-TW,zh;q=0.9`,
		`Cache-Control`:             `max-age=0`,
		`Connection`:                `keep-alive`,
		`Content-Length`:            `126`,
		`Content-Type`:              `application/x-www-form-urlencoded`,
		`Host`:                      `elearning.uch.edu.tw`,
		`Origin`:                    `https://elearning.uch.edu.tw`,
		`Referer`:                   `https://elearning.uch.edu.tw/course/doc_edit.php?id=2073482`,
		`sec-ch-ua`:                 `"Microsoft Edge";v="107", "Chromium";v="107", "Not=A?Brand";v="24"`,
		`sec-ch-ua-mobile`:          `?0`,
		`sec-ch-ua-platform`:        `"Windows"`,
		`Sec-Fetch-Dest`:            `document`,
		`Sec-Fetch-Mode`:            `navigate`,
		`Sec-Fetch-Site`:            `same-origin`,
		`Sec-Fetch-User`:            `?1`,
		`Upgrade-Insecure-Requests`: `1`,
		//`User-Agent`:                `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.26`,
	}
	req, err := http.NewRequest("POST", SendEditWorkUrl, strings.NewReader(ss))
	if err != nil {
		return err
	}
	a.setCookie(req)
	for i, v := range header {
		req.Header.Set(i, v)
	}
	res, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	w.Done = a.GetWorkStatus(*w)
	return
}

func (a *IlmsAPI) editWork(i *Work, ss string) (err error) {
	header := map[string]string{
		`Accept`:                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		`Accept-Encoding`:           `gzip, deflate, br`,
		`Accept-Language`:           `zh-TW,zh;q=0.9`,
		`Cache-Control`:             `max-age=0`,
		`Connection`:                `keep-alive`,
		`Content-Length`:            `126`,
		`Content-Type`:              `application/x-www-form-urlencoded`,
		`Host`:                      `elearning.uch.edu.tw`,
		`Origin`:                    `https://elearning.uch.edu.tw`,
		`Referer`:                   `https://elearning.uch.edu.tw/course/doc_edit.php?id=2073482`,
		`sec-ch-ua`:                 `"Microsoft Edge";v="107", "Chromium";v="107", "Not=A?Brand";v="24"`,
		`sec-ch-ua-mobile`:          `?0`,
		`sec-ch-ua-platform`:        `"Windows"`,
		`Sec-Fetch-Dest`:            `document`,
		`Sec-Fetch-Mode`:            `navigate`,
		`Sec-Fetch-Site`:            `same-origin`,
		`Sec-Fetch-User`:            `?1`,
		`Upgrade-Insecure-Requests`: `1`,
		`User-Agent`:                `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.26`,
	}
	req, err := http.NewRequest("POST", SendEditWorkUrl, strings.NewReader(ss))
	if err != nil {
		return err
	}
	a.setCookie(req)
	for i, v := range header {
		req.Header.Set(i, v)
	}
	res, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

// DeleteWork 刪除作業
func (a *IlmsAPI) DeleteWork(w *Work) bool {
	header := map[string]string{`Accept`: `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		`Accept-Encoding`: `gzip, deflate, br`,
		`Accept-Language`: `zh-TW,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6`,
		`Cache-Control`:   `max-age=0`,
		`Connection`:      `keep-alive`,
		`Content-Length`:  `12`,
		`Content-Type`:    `application/x-www-form-urlencoded`,
		//`Cookie`:                    `_ga=GA1.3.817377558.1666625227; cookie_locale=zh-tw; cookie_account=b10813141; cookie_passwd=ed6ee117e65c72ccabb0c69e6576cda2; PHPSESSID=912drvnbs9bqkc8c1fmvrb35m0; cookie_reload2080592=1`,
		`Host`:                      `elearning.uch.edu.tw`,
		`Origin`:                    `https://elearning.uch.edu.tw`,
		`Referer`:                   `https://elearning.uch.edu.tw/course/doc_delete.php?id=2080592`,
		`sec-ch-ua`:                 `"Microsoft Edge";v="107", "Chromium";v="107", "Not=A?Brand";v="24"`,
		`sec-ch-ua-mobile`:          `?0`,
		`sec-ch-ua-platform`:        `"Windows"`,
		`Sec-Fetch-Dest`:            `iframe`,
		`Sec-Fetch-Mode`:            `navigate`,
		`Sec-Fetch-Site`:            `same-origin`,
		`Sec-Fetch-User`:            `?1`,
		`Upgrade-Insecure-Requests`: `1`,
		`User-Agent`:                `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.35`,
	}
	newUrl := DeleteWorkUrl + w.Cid
	request, err := http.NewRequest("POST", newUrl, strings.NewReader("fmSubmit=yes"))
	if err != nil {
		return !printError(err)
	}

	for i, v := range header {
		request.Header.Set(i, v)
	}
	a.setCookie(request)
	_, err = a.Client.Do(request)
	if err != nil {
		return !printError(err)
	}
	w.Done = a.GetWorkStatus(*w)
	w.Cid = ""
	return !w.Done
}

// ClearWorkFiles 清除作業內的所有檔案
func (a *IlmsAPI) ClearWorkFiles(w Work) bool {
	var getFileStatus GetFileStatus
	request, err := http.NewRequest("POST", GetWorkFilesUrl, strings.NewReader("id="+w.Cid))
	if err != nil {
		return printError(err)
	}
	a.setCookie(request)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err := a.Client.Do(request)
	if err != nil {
		return printError(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	readAll, err := io.ReadAll(response.Body)
	err = json.Unmarshal(readAll, &getFileStatus)
	return func() bool {
		if getFileStatus.Ret.AttachIDList == "" {
			return true
		}
		s := fmt.Sprintf("id=%s&attachIDs=%s", w.Cid, getFileStatus.Ret.AttachIDList)
		request, err := http.NewRequest("POST", WorkDeleteFileUrl, strings.NewReader(s))
		if err != nil {
			return printError(err)
		}
		a.setCookie(request)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		response, err := a.Client.Do(request)
		if err != nil {
			return printError(err)
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)
		return true
	}()
}

func (a *IlmsAPI) setCookie(request *http.Request) {
	cookie := ""
	for _, v := range a.Cookie {
		cookie += v.String()
	}
	request.Header.Set("cookie", cookie)
}
