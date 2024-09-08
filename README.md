# spaces
Huggingface的一些space代码   
由于Huggingface的权限问题 Dockerfile均以非root用户执行，因此项目的Dockerfile并不能完美适应Dockerhub.    
以青龙面板为例：非root用户在后台添加依赖时会失败，root用户则不影响.    
因此当你用于Dockerhub的时候最好对Dockerfile进行修改以确保足够的权限来满足需求.
