import Vue from 'vue'
import Vuex from 'vuex'
import contact from './contacts' //存放所有联系人的数据
import message from './messages' //存放消息列表
import mutations from "./mutations"
import actions from "./actions"
Vue.use(Vuex)
    // 统一管理接口域名 
let apiPublicDomain = window.location.host;

if(location.hostname == 'localhost'){
    apiPublicDomain = '/';
}else{
    apiPublicDomain = '/';
}

const state = {
    currentLang: "zh", //当前使用的语言 zh：简体中文 en:英文 后期需要
    newMsgCount: 0, //新消息数量
    allContacts: contact, //所有联系人
    currentPageName: "微信", //用于在wx-header组件中显示当前页标题
    //backPageName: "", //用于在返回按钮出 显示前一页名字 已遗弃
    headerStatus: true, //显示（true）/隐藏（false）wx-header组件
    tipsStatus: false, //控制首页右上角菜单的显示(true)/隐藏(false)
    // 所有接口地址 后期需要
    apiUrl:  apiPublicDomain + "api/",
    msgList: {
        stickMsg: [], //置顶消息列表 后期需要
        baseMsg: message.messages
    }
}
export default new Vuex.Store({
    state,
    mutations,
    actions,
})