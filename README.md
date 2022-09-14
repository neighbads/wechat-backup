## 介绍
基于 [wechat-backup](https://github.com/greycodee/wechat-backup) 和 [vue-WeChat](https://github.com/zhaohaodang/vue-WeChat) 修改的 wechat 备份文件展示 Web


## 使用流程

// TODO: 备份导出流程

## 计划

- [ ] 集成 vue-WeChat
- [ x ] 1. 联系人列表
- [ ] 2. 消息列表
- [ ] 3. 个人信息
- [ ] 4. 朋友圈
- [ ] 5. 删除消息

## 开发测试

### Sqlcipher 测试
```bash
./sqlcipher ./data_backup/20220907/MicroMsg/849ffe37857c38427b0ca7d1b62b2dbe/EnMicroMsg.db

PRAGMA key = '375421e';
PRAGMA cipher_use_hmac = off;
PRAGMA kdf_iter = 4000;
PRAGMA cipher_page_size = 1024;
PRAGMA cipher_hmac_algorithm = HMAC_SHA1;
PRAGMA cipher_kdf_algorithm = PBKDF2_HMAC_SHA1;

```
