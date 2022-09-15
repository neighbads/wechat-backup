<template>
    <!--公众号组件-->
    <div :class="{'search-open-contact':!$store.state.headerStatus}" class="official-account">
        <header id="wx-header">
            <div class="center">
                <router-link to="/contact" tag="div" class="iconfont icon-return-arrow">
                    <span>通讯录</span>
                </router-link>
                <span>公众号</span>
            </div>
        </header>
        <!--这里的 search 组件的样式也需要修改一下-->
        <search></search>
        <!--公众号集合-->
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
                pageName: ""
            }
        },
        mounted() {
            // 获取列表
            this.getList();
        },
        computed: {
            contactsInitialList() {
                return contact.getInitialList(contact.official)
            },
            contactsList() {
                return contact.getContactsListGroupByInitial(contact.official, this.contactsInitialList)
            }
        },
        methods: {
            getList(){
                this.$http({
                    url:this.$store.state.apiUrl+'contact/list',
                    method:'GET',
                    params:{
                        type: "official"
                    }
                }).then(function (res) {
                    // console.log(res);
                    contact.official = res.data;
                });
            },
        }

    }
</script>
<style>
    .official-account {
        padding-bottom: 20px;
    }

    .official-account .weui-cell_access:active {
        background-color: rgba(177, 177, 177, 0.53)
    }
</style>