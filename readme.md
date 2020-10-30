# 服务端框架

> 使用 golang 高性能语言实现

---

### 目录结构

```
├─api                      # 接口逻辑
|   ├─example.go           # 示例
|   ├─...                  # 更多
├─app                      # 框架核心
|   ├─app.go               # 框架初始化
|   ├─cron.go              # 自动任务管理
|   ├─db.go                # 数据库加载
|   ├─helper.go            # 全局函数
|   ├─logger.go            # 日志处理
|   ├─middleware.go        # 中间件管理
|   ├─redis.go             # redis初始化
├─config                   # 配置
|   ├─develop              # 开发环境配置
|   |   ├─db.json          # 数据库配置
|   ├─prod                 # 线上配置
├─console                  # 自动任务管理
├─dto                      # dto管理
├─exception                # 异常管理
|   ├─LogicException.go    # 逻辑异常
|   ├─...                  # 更多异常
├─logs                     # 日志管理
├─model                    # 模型层
|   ├─mysql                # 数据库模型目录
├─probotuf                 # proto文件生成管理
├─router                   # 路由管理
|   ├─config.go            # 路由初始化配置
|   ├─User.go              # 用户模块路由初始化
|   ├─...                  # 更多模块
├─util                     # 第三方包管理
├─websocket                # 海龟汤游戏服务
├─main.go                  # 入口文件
├─.gitignore               # git提交忽略配置
├─.env                     # 环境文件
```

----

### .env 环境文件配置

```
#develop 测试开发环境，对应的配置为config/develop
#prod 生产环境，对应配置为config/prod
env=develop
```

### 配置目录下的文件示例

> 数据库配置 db.json
```json
{
  "database": "database_name",
  "port": 3306,
  "charset": "utf8",
  "protocol": "tcp",
  "master": {
    "ip": "127.0.0.1",
    "username": "root",
    "password": "root"
  },
  "slaves": []
}
```

> 小程序配置 applet.json

```json
[
  {
    "appId": "your appId",
    "secret": "your secret"
  },
  {
    "appId": "your appId",
    "secret": "your secret"
  }
]
```

> redis配置 redis.json

```json
{
  "ip": "127.0.0.1",
  "port": "6379",
  "password": "",
  "database": 0
}
```

> websocket配置 websocket.json

```json
{
  "addr": ":9504",
  "port": "9504",
  "ip": "127.0.0.1",
  "path": "/wss"
}
```

> mongo配置 mongo.json

```json
{
  "ip": "127.0.0.1",
  "port": "27017",
  "username": "root",
  "password": "root"
}
```