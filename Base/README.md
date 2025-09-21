---
title: Base
emoji: 🏃
colorFrom: yellow
colorTo: purple
sdk: docker
pinned: false
app_port: 5700
---


## 以Code Serer为核心的用nginx做为反代来部署多个服务
目前已经添加服务   
n8n地址:/   
Code Server:/coder/     
code server可通过n8n的shell节点执行：pm2 start/stop code-server来启动或停止


## 如使用Rclone方法需要设置:

Space variables (Public)说明:   
GENERIC_TIMEZONE=Asia/Shanghai  :时区   
TZ=Asia/Shanghai   :时区     
N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS=true   :重置文件权限   
WEBHOOK_URL=https://用户名-space名.hf.space/    !!!!!不要漏了最后的反斜杠     

Space secrets(Private)说明:   
N8N_ENCRYPTION_KEY:  加密密钥【需要保存，如不保存以后重置数据无法恢复】   
RCLONE_CONF:rclone配置内容，可选，用来同步数据  
ADMIN_PASSWORD:Code Server登陆密码    

同步配置目录命令   
```
rclone sync /home/coder/.n8n huggingface:/n8n --create-empty-src-dirs
```

### 以下为参考环境变量
EXECUTIONS_DATA_PRUNE=true     //是否开启自动清理运行日志   
EXECUTIONS_DATA_MAX_AGE=168   //几小时后删除运行日志    
EXECUTIONS_DATA_PRUNE_MAX_COUNT=50000   //保留的日志条数   
