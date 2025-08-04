# 在huggingface上部署Alist会秒封号 请勿使用！！！


# spaces
Huggingface的一些space代码   
由于Huggingface的权限问题 Dockerfile均以非root用户执行，因此项目的Dockerfile并不能完美适应Dockerhub.    
以青龙面板为例：非root用户在后台添加依赖时会失败，root用户则不影响.    
因此当你用于Dockerhub的时候最好对Dockerfile进行修改以确保足够的权限来满足需求.


# 推荐适用于Rclone的Webdav：   
https://infini-cloud.net/en/    
推荐码:QU2MJ

#### 如需在脚本中使用rclone配置文件 需要在github项目的settings/secrets/actions处添加变量名为RCLONE_CONF的secrets  


python文件为使用抱脸sdk创建space代码：   
青龙面板使用方法[token和userid为必填参数]:   
```
python qinglong.py --token="" --userid="" --admin="" --password="" --image="" --rclone_conf_path=""
```
