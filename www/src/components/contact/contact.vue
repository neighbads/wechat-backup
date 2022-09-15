<template>
    <div id="contact">
        <section>
            <div class="weui-cells_contact-head weui-cells weui-cells_access" style="margin-top:-1px">
                <router-link to="/contact/new-friends" class="weui-cell">
                    <div class="weui-cell_hd"> <img class="img-obj-cover"
                            src="/images/contact_top-friend-notify.png"> </div>
                    <div class="weui-cell_bd weui-cell_primary">
                        <p>新的朋友</p>
                    </div>
                </router-link>
                <router-link to="/contact/group-list" class="weui-cell">
                    <div class="weui-cell_hd"> <img class="img-obj-cover"
                            src="/images/contact_top-addgroup.png"> </div>
                    <div class="weui-cell_bd weui-cell_primary">
                        <p>群聊</p>
                    </div>
                </router-link>
                <router-link to="/contact/tags" class="weui-cell">
                    <div class="weui-cell_hd"> <img class="img-obj-cover" src="/images/contact_top-tag.png">
                    </div>
                    <div class="weui-cell_bd weui-cell_primary">
                        <p>标签</p>
                    </div>
                </router-link>
                <router-link to="/contact/official-accounts" class="weui-cell">
                    <div class="weui-cell_hd"><img class="img-obj-cover"
                            src="/images/contact_top-offical.png"></div>
                    <div class="weui-cell_bd weui-cell_primary">
                        <p>公众号</p>
                    </div>
                </router-link>
                <router-link to="/contact/official-tools" class="weui-cell">
                    <div class="weui-cell_hd"><img class="img-obj-cover"
                            src="/images/contact_top-offical.png"></div>
                    <div class="weui-cell_bd weui-cell_primary">
                        <p>工具号</p>
                    </div>
                </router-link>
                <router-link to="/contact/block-friends" class="weui-cell">
                    <div class="weui-cell_hd"><img class="img-obj-cover"
                            src="/images/contact_top-offical.png"></div>
                    <div class="weui-cell_bd weui-cell_primary">
                        <p>黑名单</p>
                    </div>
                </router-link>
            </div>
            <!--联系人集合-->
            <div >
            <template v-for="(value,key) in contactsList">
                <!--首字母-->
                <div :ref="`key_${key}`" :key="key" class="weui-cells__title">{{key}}</div>
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
        </section>
        <!--检索-->
        <div class="initial-bar"><span @click="toPs(i)" :key="i+1" v-for="i in contactsInitialList">{{i}}</span></div>
    </div>
</template>
<script>
    import contact from "../../vuex/contacts"
    export default {
        mixins: [window.mixin],
        data() {
            return {
                "pageName": "通讯录"
            }
        },
        mounted() {
            // mutations.js中有介绍
            this.$store.commit("toggleTipsStatus", -1)

            // 获取列表
            this.getList();
        },
        activated() {
            this.$store.commit("toggleTipsStatus", -1)
        },
        computed: {
            contactsInitialList() {
                return contact.getInitialList(contact.friends)
            },
            contactsList() {
                return contact.getContactsListGroupByInitial(contact.friends, this.contactsInitialList)
            }
        },
        methods: {
            toPs(i) {
                window.scrollTo(0, this.$refs['key_' + i][0].offsetTop)
            },
            getList(){
                this.$http({
                    url:this.$store.state.apiUrl+'contact/list',
                    method:'GET',
                    params:{
                        type: "friend"
                    }
                }).then(function (res) {
                    // console.log(res);
                    contact.friends = res.data;
                });
            },
        }
    }
</script>
<style lang="less">
    @import "../../assets/less/contact.less";
</style>