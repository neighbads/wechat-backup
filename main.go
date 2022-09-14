package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/greycodee/wechat-backup/db"
)

var wcdb *db.WCDB

//go:embed www/dist
var htmlFile embed.FS

var serverAddr = flag.String("h", "", "server listen port")
var serverPort = flag.String("l", "8081", "server listen port")
var basePath = flag.String("d", "./data_backup", "wechat bak folder")
var dbPasswd = flag.String("p", "", "sqlite db password")

func init() {
	flag.Parse()

	// select bak folder
	dirList := make([]string, 0)
	files, _ := ioutil.ReadDir(*basePath)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		MicroMsgFiles, _ := ioutil.ReadDir(*basePath + "/" + f.Name() + "/MicroMsg")
		for _, MicroMsgFile := range MicroMsgFiles {
			if !MicroMsgFile.IsDir() {
				continue
			}

			DbPath := *basePath + "/" + f.Name() + "/MicroMsg/" + MicroMsgFile.Name()
			// fmt.Println("DbPath", DbPath)

			enMicroMsgPath := DbPath + "/EnMicroMsg.db"
			wxFileIndexPath := DbPath + "/WxFileIndex.db"
			// fmt.Println("enMicroMsgPath", enMicroMsgPath)

			// path exist
			_, err := os.Stat(enMicroMsgPath)
			if err != nil {
				continue
			}
			_, err = os.Stat(wxFileIndexPath)
			if err != nil {
				continue
			}
			fmt.Println(len(dirList), ": ", DbPath)
			dirList = append(dirList, DbPath)
		}
	}

	// not found
	if len(dirList) == 0 {
		log.Fatal("not found EnMicroMsg.db")
	}

	// get input
	if len(dirList) > 1 {
		var input int
		fmt.Print("Please select backup folder[0]: ")
		fmt.Scanln(&input)
		if input >= len(dirList) {
			log.Fatal("Invalid Input")
		}
		*basePath = dirList[input]
	} else {
		*basePath = dirList[0]
	}
	fmt.Println("")

	fmt.Println("EnMicroMsg.db path: \t", *basePath+"/EnMicroMsg.db")
	fmt.Println("WxFileIndex.db path: \t", *basePath+"/WxFileIndex.db")

	// read db passwd
	if *dbPasswd == "" {
		pwdPath := filepath.Clean(*basePath + "/../../passwd.txt")
		pwdPath = "./" + strings.Replace(pwdPath, "\\", "/", -1)
		fmt.Println("Passwd.txt path: \t", pwdPath)
		dbPasswdFileData, err := ioutil.ReadFile(pwdPath)
		if err != nil {
			log.Fatal("please set sqlite db password, or create passwd.txt file in bak folder")
		}

		if len(dbPasswdFileData) < 7 {
			log.Fatal("passwd.txt file content error, need 7 char")
		}
		*dbPasswd = string(dbPasswdFileData[:7])
	}

	fmt.Println("sqlite db password: \t[ ", *dbPasswd, " ]")
	fmt.Println("")
}

func main() {

	wcdb = db.InitWCDB(*basePath, *dbPasswd)

	fsys, _ := fs.Sub(htmlFile, "www/dist")
	staticHandle := http.FileServer(http.FS(fsys))

	// 文件路由
	fs := http.FileServer(http.Dir(*basePath))
	http.Handle(db.MediaPathPrefix, http.StripPrefix(db.MediaPathPrefix, fs))

	http.Handle("/", staticHandle)
	http.Handle("/api/", route())

	listenAddr := *serverAddr + ":" + *serverPort
	log.Println("Server start at", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func route() http.Handler {
	return &API{}
}

type API struct {
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	apiMap[path](w, r)
}

var apiMap = map[string]func(w http.ResponseWriter, r *http.Request){
	"/api/contact/list": func(w http.ResponseWriter, r *http.Request) {
		// 联系人列表
		params := r.URL.Query()
		contactType := params.Get("type")

		result, err := json.Marshal(wcdb.GetContactList(contactType))
		if err != nil {
			log.Fatalf("json marshal error: %v", err)
		}
		w.Write(result)
	},
	"/api/chat/list": func(w http.ResponseWriter, r *http.Request) {
		// 聊天列表
		params := r.URL.Query()
		pageIndex, _ := strconv.Atoi(params["pageIndex"][0])
		pageSize, _ := strconv.Atoi(params["pageSize"][0])
		names, _ := params["name"]
		name := ""
		if len(names) > 0 {
			name = names[0]
		}
		all, _ := strconv.ParseBool(params["all"][0])
		result, err := json.Marshal(wcdb.ChatList(pageIndex-1, pageSize, all, name))
		if err != nil {
			log.Fatalf("json marshal error: %v", err)
		}
		w.Write(result)
	},
	"/api/chat/detail": func(w http.ResponseWriter, r *http.Request) {
		//聊天记录
		params := r.URL.Query()
		talker := params["talker"][0]
		pageIndex, _ := strconv.Atoi(params["pageIndex"][0])
		pageSize, _ := strconv.Atoi(params["pageSize"][0])

		result, err := json.Marshal(wcdb.ChatDetailList(talker, pageIndex-1, pageSize))
		if err != nil {
			log.Fatalf("json marshal error: %v", err)
		}
		w.Write(result)
	},
	"/api/user/info": func(w http.ResponseWriter, r *http.Request) {
		// 用户信息
		params := r.URL.Query()
		username := params["username"][0]
		result, err := json.Marshal(wcdb.GetUserInfo(username))
		if err != nil {
			log.Fatalf("json marshal error: %v", err)
		}
		w.Write(result)
	},
	"/api/user/myinfo": func(w http.ResponseWriter, r *http.Request) {
		// 自己的信息
		result, err := json.Marshal(wcdb.GetMyInfo())
		if err != nil {
			log.Fatalf("json marshal error: %v", err)
		}
		w.Write(result)
	},
	"/api/media/img": func(w http.ResponseWriter, r *http.Request) {
		// 图片
		params := r.URL.Query()
		msgId := params["msgId"][0]
		w.Write([]byte(wcdb.GetImgPath(msgId)))
	},
	"/api/media/video": func(w http.ResponseWriter, r *http.Request) {
		// 视频
		params := r.URL.Query()
		msgId := params["msgId"][0]
		w.Write([]byte(wcdb.GetVideoPath(msgId)))
	},
	"/api/media/voice": func(w http.ResponseWriter, r *http.Request) {
		// 语音
		params := r.URL.Query()
		msgId := params["msgId"][0]
		w.Write([]byte(wcdb.GetVoicePath(msgId)))
	},
}
