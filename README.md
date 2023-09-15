# GRTCMPRO
v1.0版本地址 https://github.com/aoweichen/GoRealTimeChatWithVUE.git
v1.1版本更新了整体服务端的构架的构架。  
1. 添加代理服务端，实现了登录、注册、发送验证码以及鉴权功能，其他的核心服务将会先鉴权后转发到相关服务集群。
2. 添加GRPC微服务集群系统，使用consul作为GRPC服务注册和服务发现，实现登录、注册、发送验证码以及鉴权的功能。
3. 添加HTTP服务集群，包含主要的IM功能。
（注释：由于没有合适的前端实现，所以只用Apifox进行了单测。）

![image](https://github.com/aoweichen/GRTCMPRO/assets/73885370/12a35004-305a-4c35-8198-a9fd4a19f74b)

![image](https://github.com/aoweichen/GRTCMPRO/assets/73885370/232e97a5-13c1-4e2b-a7c7-bc27a1590f8a)


