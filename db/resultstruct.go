package db

type ContactList struct {
	Friends  []ContactListRow `json:"friends"`
	Tools    []ContactListRow `json:"tools"`
	Official []ContactListRow `json:"official"`
	Chatroom []ContactListRow `json:"chatroom"`
	Openim   []ContactListRow `json:"openim"`
	Block    []ContactListRow `json:"block"`
}

type ContactListRow struct {
	Wxid      string      `json:"wxid"`
	Initial   string      `json:"initial"`
	HeaderUrl string      `json:"headerUrl"`
	Nickname  string      `json:"nickname"`
	Sex       int         `json:"sex"`
	Remark    string      `json:"remark"`
	Signature string      `json:"signature"`
	Telphone  int         `json:"telphone"`
	Album     []AlbumRow  `json:"album"`
	Area      []string    `json:"area"`
	From      string      `json:"from"`
	Tag       string      `json:"tag"`
	Desc      ContactDesc `json:"desc"`
}

type AlbumRow struct {
	ImgSrc string `json:"imgSrc"`
	Date   int    `json:"date"`
}

type ContactDesc struct {
	Title  string `json:"title"`
	PicUrl string `json:"picUrl"`
}

type ChatList struct {
	Total int           `json:"total"`
	Rows  []ChatListRow `json:"rows"`
}

type ChatListRow struct {
	Talker      string `json:"talker"`
	Alias       string `json:"alias"`
	NickName    string `json:"nickname"`
	ConRemark   string `json:"conRemark"`
	Reserved1   string `json:"reserved1"`
	Reserved2   string `json:"reserved2"`
	LocalAvatar string `json:"localAvatar"`
	MsgCount    int    `json:"msgCount"`
	CreateTime  int64  `json:"createTime"`
	ChatType    int    `json:"chatType"`
	UserType    int    `json:"userType"`
}

type ChatDetailList struct {
	Total int                 `json:"total"`
	Rows  []ChatDetailListRow `json:"rows"`
}
type ChatDetailListRow struct {
	MsgId           string    `json:"msgId"`
	MsgSvrId        string    `json:"msgSvrId"`
	Type            int       `json:"type"`
	IsSend          int       `json:"isSend"`
	CreateTime      int64     `json:"createTime"`
	Talker          string    `json:"talker"`
	Content         string    `json:"content"`
	ImgPath         string    `json:"imgPath"`
	MediaPath       string    `json:"mediaPath"`
	MediaBCKPath    string    `json:"mediaBCKPath"`
	MediaSourcePath string    `json:"mediaSourcePath"`
	FileInfo        FileInfo  `json:"fileInfo"`
	EmojiInfo       EmojiInfo `json:"emojiInfo"`
	IsChatRoom      bool      `json:"isChatRoom"`
	UserInfo        UserInfo  `json:"userInfo"`
}

type UserInfo struct {
	UserName    string `json:"userName"`
	Alias       string `json:"alias"`
	ConRemark   string `json:"conRemark"`
	NickName    string `json:"nickName"`
	Reserved1   string `json:"reserved1"`
	Reserved2   string `json:"reserved2"`
	LocalAvatar string `json:"localAvatar"`
}

type FileInfo struct {
	FileName string `json:"fileName"`
	FileSize string `json:"fileSize"`
	FilePath string `json:"filePath"`
	FileExt  string `json:"fileExt"`
}

type EmojiInfo struct {
	Md5    string `json:"md5"`
	CDNUrl string `json:"cdnUrl"`
	W      int64  `json:"w"`
	H      int64  `json:"h"`
}
