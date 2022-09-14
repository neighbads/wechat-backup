package db

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mutecomm/go-sqlcipher"
)

var MediaPathPrefix = "/media/"

const (
	VX_CONTACT_TYPE_IN_CONTACT   = 0x00000001 // 在通讯录里
	VX_CONTACT_TYPE_IN_MESSAGE   = 0x00000002 // 在聊天记录里
	VX_CONTACT_TYPE_CHATROOM_MEM = 0x00000004 // 群成员
	VX_CONTACT_TYPE_BLOCK        = 0x00000008 // 黑名单

	VX_CONTACT_TYPE_OFFICIAL_TOOLS = 0x00000020 // 官方工具号

	VX_CONTACT_TYPE_CANNOT_SEE_MY = 0x00000100 // 朋友圈权限, 不让他看
	VX_CONTACT_TYPE_MSG_MUTE      = 0x00000200 // 消息免打扰
)

var wxToolFriendStrs = []string{
	"qqsafe",      // QQ安全中心
	"mphelper",    // 公众平台安全助手
	"exmail_tool", // 腾讯企业邮箱
	"filehelper",  //文件传输助手
	"weixin",      // 微信团队
}

type EnMicroMsg struct {
	db          *sql.DB
	contactList *ContactList

	myInfo   UserInfo
	chatList *ChatList
}

func OpenEnMicroMsg(dbPath string, dbPassword string) *EnMicroMsg {
	em := &EnMicroMsg{}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("open db error: %v", err)
	}

	db.Exec(fmt.Sprintf("PRAGMA key = '%s';", dbPassword))
	db.Exec("PRAGMA cipher_use_hmac = OFF;")
	db.Exec("PRAGMA cipher_page_size = 1024;")
	db.Exec("PRAGMA kdf_iter = 4000;")

	// check db passwd
	_, err = db.Exec("select count(*) FROM sqlite_master;")
	if err != nil {
		log.Fatal(err)
	}

	em.db = db
	em.myInfo = em.GetMyInfo()

	em.contactList = em.GetContactList()
	return em
}

func (em *EnMicroMsg) Close() {
	em.db.Close()
}

func (em EnMicroMsg) GetContactList() *ContactList {
	contactList := &ContactList{}

	// get username,initial,conRemarkPYFull,nickname,conRemark
	rows, err := em.db.Query("select rc.type,rc.username,ifnull(rc.pyInitial,'') as initial,ifnull(rc.conRemarkPYFull,'') as conRemarkPYFull,ifnull(rc.nickname,'') as nickname,ifnull(rc.conRemark,'') as conRemark from rcontact rc")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var r ContactListRow
		var contactType int
		var conRemarkPYFull string

		err = rows.Scan(&contactType, &r.Wxid, &r.Initial, &conRemarkPYFull, &r.Nickname, &r.Remark)
		if err != nil {
			log.Printf("get Talker error,%s", err)
		}
		r.HeaderUrl = em.getLocalAvatar(r.Wxid)

		// get first letter
		if len(conRemarkPYFull) > 0 {
			r.Initial = conRemarkPYFull
		}
		if len(r.Initial) > 0 {
			r.Initial = strings.ToUpper(r.Initial[0:1])
			if r.Initial[0] < 65 || r.Initial[0] > 90 {
				r.Initial = "#"
			}
		} else {
			r.Initial = "#"
		}

		// 黑名单
		if (contactType & VX_CONTACT_TYPE_BLOCK) != 0 {
			contactList.Block = append(contactList.Block, r)
			continue
		}

		// 在通讯录里
		if (contactType & VX_CONTACT_TYPE_IN_CONTACT) != 0 {
			// 公共号
			if strings.HasPrefix(r.Wxid, "gh_") {
				contactList.Official = append(contactList.Official, r)
				continue
			}

			// 官方功能号
			if (contactType & VX_CONTACT_TYPE_OFFICIAL_TOOLS) != 0 {
				contactList.Tools = append(contactList.Tools, r)
				continue
			}

			// 群聊
			if strings.HasSuffix(r.Wxid, "@chatroom") {
				contactList.Chatroom = append(contactList.Chatroom, r)
				continue
			}

			// 企业联系人
			if strings.HasSuffix(r.Wxid, "@openim") {
				contactList.Openim = append(contactList.Openim, r)
				continue
			}

			// 其他功能号
			if inArray(r.Wxid, wxToolFriendStrs) {
				contactList.Tools = append(contactList.Tools, r)
			} else {
				// 普通联系人
				contactList.Friends = append(contactList.Friends, r)
			}
		}

	}
	rows.Close()

	log.Printf("get contact list success, friends(%d) tools(%d) chatroom(%d) openim(%d) block(%d)",
		len(contactList.Friends), len(contactList.Tools), len(contactList.Chatroom), len(contactList.Openim), len(contactList.Block))

	return contactList
}

func (em EnMicroMsg) GetChatList() *ChatList {
	chatList := &ChatList{}

	fmt.Print("Prepart chat list, Wait a moment...")

	// 原 sql 在 sqlcipher 会崩溃, 分步查询
	// get count,talker from message
	rows, err := em.db.Query("select count(*) as msgCount,msg.talker,msg.createtime from message msg left join rcontact rc on msg.talker=rc.username group by msg.talker order by msg.createTime desc")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var r ChatListRow

		err = rows.Scan(&r.MsgCount, &r.Talker, &r.CreateTime)
		if err != nil {
			log.Printf("get Talker error,%s", err)
		}

		r.LocalAvatar = em.getLocalAvatar(r.Talker)
		chatList.Rows = append(chatList.Rows, r)
	}
	rows.Close()

	// get Alias,NickName,ConRemark from rcontact
	for _, r := range chatList.Rows {
		queryChatInfoSql := fmt.Sprintf("select ifnull(rc.alias,'') as alias,ifnull(rc.nickname,'') as nickname,ifnull(rc.conRemark,'') as conRemark from rcontact rc where rc.username='%s'", r.Talker)
		err = em.db.QueryRow(queryChatInfoSql).Scan(&r.Alias, &r.NickName, &r.ConRemark)
		if err != nil {
			// log.Println("queryChatInfoSql", queryChatInfoSql)
			// log.Printf("get Alias,NickName,ConRemark error, %s ,%s", r.Talker, err)
			continue
		}
	}

	// get reserved1,reserved2 from img_flag
	for _, r := range chatList.Rows {
		queryChatInfoSql := fmt.Sprintf("select ifnull(imf.reserved1,'') as reserved1,ifnull(imf.reserved2,'') as reserved2 from img_flag imf where imf.username='%s'", r.Talker)
		err = em.db.QueryRow(queryChatInfoSql).Scan(&r.Reserved1, &r.Reserved2)
		if err != nil {
			// log.Println("queryChatInfoSql", queryChatInfoSql)
			// log.Printf("get reserved1,reserved2 error, %s ,%s", r.Talker, err)
			continue
		}
	}

	// update nickname from chatroom
	for _, r := range chatList.Rows {
		if len(strings.Split(r.Talker, "@")) == 2 && strings.Split(r.Talker, "@")[1] == "chatroom" {
			r.UserType = 1
			if r.NickName == "" {
				queryRoomSql := fmt.Sprintf("select displayname as nickname from chatroom where chatroomname='%s'", r.Talker)
				room, _ := em.db.Query(queryRoomSql)
				defer room.Close()
				for room.Next() {
					room.Scan(&r.NickName)
				}
			}
		} else if r.Talker[:3] == "gh_" {
			r.UserType = 2
		}
		// fmt.Printf("%+v\n", r)
	}

	chatList.Total = len(chatList.Rows)

	fmt.Println("Done, total:", chatList.Total)

	return chatList
}

func (em EnMicroMsg) ChatList(pageIndex int, pageSize int, all bool, name string) *ChatList {
	result := &ChatList{}

	result.Rows = append(result.Rows, em.chatList.Rows[pageIndex*pageSize:(pageIndex+1)*pageSize]...)
	result.Total = len(em.chatList.Rows)

	return result
}

func (em EnMicroMsg) ChatList_bak(pageIndex int, pageSize int, all bool, name string) *ChatList {
	result := &ChatList{}
	result.Total = 10
	result.Rows = make([]ChatListRow, 0)
	var queryRowsSqlTmp string
	var queryRowsSql string
	queryRowsSqlTmp = "select count(*) as msgCount,ifnull(rc.alias,'') as alias,msg.talker,ifnull(rc.nickname,'') as nickname,ifnull(rc.conRemark,'') as conRemark,ifnull(imf.reserved1,'') as reserved1,ifnull(imf.reserved2,'') as reserved2,msg.createtime from message msg left join rcontact rc on msg.talker=rc.username  left join img_flag imf on msg.talker=imf.username "
	if name != "" {
		queryRowsSqlTmp = queryRowsSqlTmp + "where nickname like '%" + name + "%'  or conRemark like '%" + name + "%'"
	}
	queryRowsSqlTmp = queryRowsSqlTmp + " group by msg.talker order by msg.createTime desc "
	queryRowsSql = queryRowsSqlTmp
	if !all {
		queryRowsSql = queryRowsSql + fmt.Sprintf("limit %d,%d", pageIndex*pageSize, pageSize)
	}
	rows, err := em.db.Query(queryRowsSql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r ChatListRow
		r.UserType = 0
		err = rows.Scan(&r.MsgCount, &r.Alias, &r.Talker, &r.NickName, &r.ConRemark, &r.Reserved1, &r.Reserved2, &r.CreateTime)
		// 判断是否是群聊
		if len(strings.Split(r.Talker, "@")) == 2 && strings.Split(r.Talker, "@")[1] == "chatroom" {
			r.UserType = 1
			if r.NickName == "" {
				queryRoomSql := fmt.Sprintf("select displayname as nickname from chatroom where chatroomname='%s'", r.Talker)
				room, _ := em.db.Query(queryRoomSql)
				defer room.Close()
				for room.Next() {
					room.Scan(&r.NickName)
				}
			}
		} else if r.Talker[:3] == "gh_" {
			r.UserType = 2
		}

		if err != nil {
			log.Printf("未查询到聊天列表,%s", err)
		}
		r.LocalAvatar = em.getLocalAvatar(r.Talker)
		result.Rows = append(result.Rows, r)
	}

	queryTotalSql := "select count(*) as total FROM (" + queryRowsSqlTmp + ") as d"
	var total int
	err = em.db.QueryRow(queryTotalSql).Scan(&total)
	if err != nil {
		log.Printf("未查询到聊天列表总数,%s", err)
	} else {
		result.Total = total
	}

	return result
}

func (em EnMicroMsg) ChatDetailList(talker string, pageIndex int, pageSize int) *ChatDetailList {
	result := &ChatDetailList{}
	result.Total = 10
	result.Rows = make([]ChatDetailListRow, 0)
	queryRowsSql := fmt.Sprintf("SELECT ifnull(msgId,'') as msgId,ifnull(msgSvrId,'') as msgSvrId,type,isSend,createTime,talker,ifnull(content,'') as content,ifnull(imgPath,'') as imgPath FROM message WHERE talker='%s' order by createtime desc limit %d,%d", talker, pageIndex*pageSize, pageSize)
	rows, err := em.db.Query(queryRowsSql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r ChatDetailListRow
		err = rows.Scan(&r.MsgId, &r.MsgSvrId, &r.Type, &r.IsSend, &r.CreateTime, &r.Talker, &r.Content, &r.ImgPath)
		if err != nil {
			log.Printf("未查询到聊天历史记录,%s", err)
		}
		// 表情图片
		if r.Type == 47 {
			r.EmojiInfo = em.GetEmojiInfo(r.ImgPath)
		}
		// em.getMediaPath(&r, wxfileindex)
		result.Rows = append(result.Rows, r)
	}
	return result
}

func (em EnMicroMsg) GetUserInfo(username string) UserInfo {
	r := UserInfo{}
	querySql := fmt.Sprintf("select rc.username,rc.alias,rc.conRemark,rc.nickname,ifnull(imf.reserved1,'') as reserved1,ifnull(imf.reserved2,'') as reserved2 from rcontact rc LEFT JOIN img_flag imf on rc.username=imf.username where rc.username='%s';", username)
	err := em.db.QueryRow(querySql).Scan(&r.UserName, &r.Alias, &r.ConRemark, &r.NickName, &r.Reserved1, &r.Reserved2)
	if err != nil {
		log.Printf("未查询到用户信息,%s", err)
	} else {
		r.LocalAvatar = em.getLocalAvatar(r.UserName)
	}
	r.LocalAvatar = em.getLocalAvatar(r.UserName)
	return r
}

func (em EnMicroMsg) GetMyInfo() UserInfo {
	r := UserInfo{}
	querySql := "select rc.username,rc.alias,ifnull(rc.conRemark,''),rc.nickname,ifnull(imf.reserved1,'') as reserved1,ifnull(imf.reserved2,'') as reserved2 from rcontact rc left join img_flag imf on rc.username=imf.username where rc.username=(select value from userinfo WHERE id = 2)"
	err := em.db.QueryRow(querySql).Scan(&r.UserName, &r.Alias, &r.ConRemark, &r.NickName, &r.Reserved1, &r.Reserved2)
	if err != nil {
		log.Fatal("未查询到个人信息", err)
	} else {
		r.LocalAvatar = em.getLocalAvatar(r.UserName)
	}
	return r
}

func (em EnMicroMsg) getLocalAvatar(username string) string {
	md5Str := fmt.Sprintf("%x", md5.Sum([]byte(username)))
	filePath := fmt.Sprintf("%savatar/%s/%s/user_%s.png", MediaPathPrefix, md5Str[0:2], md5Str[2:4], md5Str)
	return filePath
}

func (em EnMicroMsg) formatImagePath(path string) string {
	imgFileName := strings.Split(path, "://")[1]
	return fmt.Sprintf("%simage2/%s/%s/%s", MediaPathPrefix, imgFileName[3:5], imgFileName[5:7], imgFileName)
}
func (em EnMicroMsg) formatImageBCKPath(chat ChatDetailListRow) string {
	var imgFileName string
	if chat.IsSend == 0 {
		// 接收
		imgFileName = fmt.Sprintf("%s_%s_%s_backup", chat.Talker, em.myInfo.UserName, chat.MsgSvrId)
	} else {
		// 发送
		imgFileName = fmt.Sprintf("%s_%s_%s_backup", em.myInfo.UserName, chat.Talker, chat.MsgSvrId)
	}
	return fmt.Sprintf("%simage2/%s/%s/%s", MediaPathPrefix, imgFileName[0:2], imgFileName[2:4], imgFileName)
}
func (em EnMicroMsg) formatVoicePath(path string) string {
	p := md5.Sum([]byte(path))
	md5Str := fmt.Sprintf("%x", p)
	// 这边原本后缀为amr格式，由于网页不能播放amr格式，所以要用转换工具转换格式，转换后为mp3格式，所以后缀为mp3
	return fmt.Sprintf("%svoice2/%s/%s/msg_%s.mp3", MediaPathPrefix, md5Str[0:2], md5Str[2:4], path)
}
func (em EnMicroMsg) formatVideoPath(path string) string {
	return fmt.Sprintf("%svideo/%s.mp4", MediaPathPrefix, path)
}

func (em EnMicroMsg) GetEmojiInfo(imgPath string) EmojiInfo {
	emojiInfo := EmojiInfo{}
	querySql := fmt.Sprintf("select md5, cdnUrl,width,height from EmojiInfo where md5='%s'", imgPath)

	err := em.db.QueryRow(querySql).Scan(&emojiInfo.Md5, &emojiInfo.CDNUrl, &emojiInfo.W, &emojiInfo.H)
	if err != nil {
		log.Printf("未查询到Emoji,%s", err)
	}
	return emojiInfo
}

func inArray(s string, arr []string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
