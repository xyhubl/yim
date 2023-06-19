# YIM 聊天推送系统

该项目可以业务于 推送、聊天等系统，鉴于goim 不怎么维护，所以复刻了此项目。欢迎感兴趣的伙伴一起维护(v: -Xy961025 ＋注明来意)。

Comet:

有状态服务，目前只支持ws通信(支持多节点部署)

Job:

消息分发层

Logic:

logic路由层，目前实现单聊



需要环境:

zookeeper kafka redis golang环境（建议>=1.17） 





启动:

**comet:**

**docker build -t comet:v1 . -f Dockerfile_comet**

**docker run -d  -p 8081:8081 -p 8083:8083 -it comet:v1  -config ./cmd/comet/dev/comet.yaml** 



**job:**

**docker build -t job:v1 . -f Dockerfile_job**

**docker run -d   job:v1  -config ./cmd/job/dev/job.yaml** 



**comet:**

**docker build -t logic:v1 . -f Dockerfile_logic**

**docker run -d  -p 8080:8080 -p 9091:9091 -it logic:v1  -config ./cmd/logic/dev/logic.yaml** 



后续规划:

1. 群聊、广播、事件订阅
2. 离线消息存储漫游
3. ack机制
4. 服务发现注册
5. 等等..



