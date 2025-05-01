---
title: N8n
emoji: ⚡
colorFrom: blue
colorTo: gray
sdk: docker
app_port: 5678
pinned: false
license: mit
short_description: n8n free hosted with supebase
---

Space variables (Public)说明:   
GENERIC_TIMEZONE=Asia/Shanghai  :时区   
TZ=Asia/Shanghai   :时区   
DB_TABLE_PREFIX=n8n_    :数据库表前缀   
N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS=true   :重置文件权限   
WEBHOOK_URL=https://用户名-space名.hf.space/ !!!!!不要漏了最后的反斜杠      


Space secrets(Private)说明:   
N8N_ENCRYPTION_KEY:  加密密钥【需要保存，如不保存以后重置数据无法恢复】   
DB_TYPE: 数据库类型:postgresdb,sqlite【不填默认sqlite】如不用postgresql下面可不填   
DB_POSTGRESDB_SCHEMA: 数据库协议（public）   
DB_POSTGRESDB_HOST: 数据库主机地址   
DB_POSTGRESDB_PORT: 数据库端口   
DB_POSTGRESDB_DATABASE: 数据库名称   
DB_POSTGRESDB_USER: 数据库用户名   
DB_POSTGRESDB_PASSWORD: 数据库密码   
DB_POSTGRESDB_SSL_CA: 数据库链接SSL证书内容   

## 如使用Rclone方法需要设置:

Space variables (Public)说明:   
GENERIC_TIMEZONE=Asia/Shanghai  :时区   
TZ=Asia/Shanghai   :时区     
N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS=true   :重置文件权限   
WEBHOOK_URL=https://用户名-space名.hf.space     

Space secrets(Private)说明:   
N8N_ENCRYPTION_KEY:  加密密钥【需要保存，如不保存以后重置数据无法恢复】   
RCLONE_CONF:rclone配置内容，可选，用来同步数据  

同步配置目录命令   
```
rclone sync /home/node/.n8n huggingface:/n8n --create-empty-src-dirs
```

### 以下为参考环境变量
EXECUTIONS_DATA_PRUNE=true     //是否开启自动清理运行日志   
EXECUTIONS_DATA_MAX_AGE=168   //几小时后删除运行日志    
EXECUTIONS_DATA_PRUNE_MAX_COUNT=50000   //保留的日志条数   






Check out the configuration reference at https://huggingface.co/docs/hub/spaces-config-reference
