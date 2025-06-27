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
hexo地址:/   
Code Server:/coder/   


## 如使用Rclone方法需要设置:

Space variables (Public)说明:   

Space secrets(Private)说明:   
ADMIN_PASSWORD: coder server管理密码   
RCLONE_CONF: rclone配置内容，可选，用来同步数据  
 

同步配置目录命令   
```
rclone sync /home/coder/blog huggingface:/hexo --create-empty-src-dirs
```

### 以下为参考环境变量
  
