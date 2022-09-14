/**
 * wxid-微信id
 * initial-姓名首字母
 * headerUrl-头像地址
 * nickname-昵称
 * sex-性别 男1女0
 * remark-备注
 * signature-个性签名
 * telphone-电话号码
 * album-相册
 * area-地区
 * from-来源
 * desc-描述
 * {
        "wxid": "wxid_zhaohd",
        "initial": 'z',
        "headerUrl": "/images/header/header01.png",
        "nickname": "阿荡",
        "sex": 1,
        "remark": "阿荡",
        "signature": "填坑小能手",
        "telphone": 18896586152,
        "album": [{
            imgSrc: ""
        }],
        "area": ["中国", "北京", "海淀"],
        "from": "",
        "tag": "",
        "desc": {
        }
    }
 */

const contact = {
    friends: [],    // 普通好友
    tools: [],      // 工具号
    chatroom: [],   // 群聊
    openim: [],     // 企业好友
    official: [],   // 公众号
    block: [],      // 黑名单

    getInitialList: function(list) {
        var initialList = []
        for (var i = 0; i < list.length; i++) {
            var initial = list[i].initial.toUpperCase()
            if (initialList.indexOf(initial) == -1 && initial != "#" && initial != "") {
                initialList.push(list[i].initial.toUpperCase())
            }
        }
        initialList = initialList.sort()
        initialList.push("#")
        // console.log(initialList);
        return initialList
    },
    getContactsListGroupByInitial: function(list, initialList) {
        var contactsList = {}
        for (var i = 0; i < initialList.length; i++) {
            var protoTypeName = initialList[i]
            contactsList[protoTypeName] = []
            for (var j = 0; j < list.length; j++) {
                if (list[j].initial.toUpperCase() === protoTypeName) {
                    contactsList[protoTypeName].push(list[j])
                }
            }
        }
        return contactsList
    },

    getUserInfo: function(wxid) {
        if (!wxid) {
            return;
        } else {
            for (var index0 in contact.friends) {
                if (contact.friends[index0].wxid === wxid) {
                    return contact.friends[index0]
                }
            }
            for (var index1 in contact.tools) {
                if (contact.tools[index1].wxid === wxid) {
                    return contact.tools[index1]
                }
            }
            for (var index2 in contact.chatroom) {
                if (contact.chatroom[index2].wxid === wxid) {
                    return contact.chatroom[index2]
                }
            }
            for (var index3 in contact.openim) {
                if (contact.openim[index3].wxid === wxid) {
                    return contact.openim[index3]
                }
            }
            for (var index4 in contact.official) {
                if (contact.official[index4].wxid === wxid) {
                    return contact.official[index4]
                }
            }
            for (var index5 in contact.block) {
                if (contact.block[index5].wxid === wxid) {
                    return contact.block[index5]
                }
            }
        }
    }
}

export default contact