<template>
    <!--我的群聊组件-->
    <div :class="{'search-open-contact':!$store.state.headerStatus}">
        <header id="wx-header">
            <div class="center">
                <router-link to="/contact" tag="div" class="iconfont icon-return-arrow">
                    <span>通讯录</span>
                </router-link>
                <span>群聊</span>
            </div>
        </header>
        <!--这里的 search 组件的样式需要修改一下-->
        <search></search>
        <!--群聊集合-->
        <template v-for="(value,key) in contactsList">
            <div class="weui-cells__title" :key="key">{{key}}</div>
            <div class="weui-cells" :key="key+1">
                <router-link :key="item.wxid" :to="{path:'/contact/details/details',query:{wxid:item.wxid}}"
                        class="weui-cell weui-cell_access" v-for="item in value" tag="div">
                    <div class="weui-cell__hd">
                        <img :src="item.headerUrl" class="home__mini-avatar___1nSrW">
                    </div>
                    <div class="weui-cell__bd">
                        {{item.remark?item.remark:item.nickname}}
                    </div>
                </router-link>
            </div>
        </template>
    </div>
</template>
<script>
    import search from "../common/search"
    import contact from "../../vuex/contacts"
    export default {
        components: {
            search
        },
        data() {
            return {

            }
        },
        mounted() {
            // 获取列表
            this.getList();
        },
        computed: {
            contactsInitialList() {
                return contact.getInitialList(contact.chatroom)
            },
            contactsList() {
                return contact.getContactsListGroupByInitial(contact.chatroom, this.contactsInitialList)
            }
        },
        methods: {
            getList(){
                this.$http({
                    url:this.$store.state.apiUrl+'contact/list',
                    method:'GET',
                    params:{
                        type: "chatroom"
                    }
                }).then(function (res) {
                    // console.log(res);
                    contact.chatroom = res.data;
                });
            },
        }
    }
</script>
<style>
    .header-box {
        position: relative;
        float: left;
        width: 38px;
        height: 38px;
        margin-right: 10px;
    }

    .header-box .header {
        height: 100%;
        display: flex;
        display: -webkit-flex;
        flex-direction: row;
        flex-wrap: wrap;
        align-items: flex-start;
        overflow: hidden;
        background: #dddbdb;
    }

    .header-box .header img {
        width: 10%;
        height: auto;
        flex-grow: 2;
        border: 0;
    }

    .multi-header img {
        margin: 1px;
    }
</style>