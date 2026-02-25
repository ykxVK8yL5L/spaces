---
title: baihu
emoji: 🐯
colorFrom: yellow
colorTo: purple
sdk: docker
pinned: false
app_port: 8052
---
 
Space secrets(Private)说明:    
RCLONE_CONF:rclone配置内容


## 注意 
同步命令：[$REMOTE_FOLDER为自己设置的rclone远程目录]
```
rclone sync /app $REMOTE_FOLDER --exclude="/baihu" --exclude "/docker-entrypoint.sh"
```

