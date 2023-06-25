# YIM 聊天推送系统

该项目可以业务于 推送、聊天等系统，鉴于goim 不怎么维护，所以复刻了此项目。欢迎感兴趣的伙伴一起维护(v: -Xy961025 ＋注明来意)。

效果:
http://120.27.141.27/h123/
http://120.27.141.27/h124/

Comet:

有状态服务，目前只支持ws通信(支持多节点部署)

Job:

消息分发层

Logic:

logic路由层，目前实现单聊



流程:

客户端连接ws -> auth- > 建立连接

logic发送消息 -> job分发 -> comet发送消息到客户端



需要环境:

zookeeper kafka redis golang环境（建议>=1.17） 





启动:

**comet:**

**docker build -t comet:v1 . -f Dockerfile_comet**

**docker run -d  -p 8081:8081 -p 8083:8083 -it comet:v1  -config ./cmd/comet/dev/comet.yaml** 



**job:**

**docker build -t job:v1 . -f Dockerfile_job**

**docker run -d   job:v1  -config ./cmd/job/dev/job.yaml** 



**logic:**

**docker build -t logic:v1 . -f Dockerfile_logic**

**docker run -d  -p 8080:8080 -p 9091:9091 -it logic:v1  -config ./cmd/logic/dev/logic.yaml** 



消息协议：

```
| parameter             | is required  | type     | comment|
| :-----                | :---         | :---     | :---       |
| package length        | true  | int32 bigendian | package length |
| header Length         | true  | int16 bigendian | header length |
| ver                   | true  | int16 bigendian | Protocol version |
| operation             | true |  int32 bigendian | Operation |
| seq                   | true |  int32 bigendian | jsonp callback |
| body                  | false | binary          | $(package lenth) - $(header length) |


| operation     | comment | 
| :-----     | :---  |
| 2 | Client send heartbeat|
| 3 | Server reply heartbeat|
| 7 | authentication request |
| 8 | authentication response |
```


目前完成:
1. 单聊、群聊

后续规划:
1. 离线消息存储漫游
2. ack机制
3. 等等..


