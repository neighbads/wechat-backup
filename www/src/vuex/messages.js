/**
 * mid-消息的id 唯一标识，重要
 * type
 * group_name
 * group_qrCode
 * read-true；已读 false：未读
 * newMsgCount
 * quiet-true：消息免打扰 false：提示此好友/群的新消息
 * msg-对话框的聊天记录 新消息 push 进
 *  msg
 *    text-消息
 *    date-时间
 *    name-发送者
 *    headerUrl-发送者头像
 * 
 *  {   //普通消息列表
        "mid": 1, //消息的id 唯一标识，重要
        "type": "friend",
        "group_name": "",
        "group_qrCode": "",
        "read": true, //true；已读 false：未读
        "newMsgCount": 1,
        "quiet": false, // true：消息免打扰 false：提示此好友/群的新消息
        "msg": [{ //对话框的聊天记录 新消息 push 进
            "text": "长按这些白色框消息，唤醒消息操作菜单，长按这些白色框消息，唤醒消息操作菜单",
            "date": 1488117964495,
            "name": "阿荡",
            "headerUrl": "/images/header/header01.png"
        }, {
            "text": '点击空白处，操作菜单消失',
            "date": 1488117964495,
            "name": "阿荡",
            "headerUrl": "/images/header/header01.png"
        }, {
            "text": '来呀 快活啊',
            "date": 1488117964495,
            "name": "阿荡",
            "headerUrl": "/images/header/header01.png"
        }],
        "user": [] // 此消息的用户数组 长度为1则为私聊 长度大于1则为群聊
    },
 */

const messages = []

const message = {
    messages
}

export default message
