# 服务端框架

> 使用 golang 高性能语言实现

#### 技术相关
- go >= 1.14
- protobuf
- websocket
- mongodb
- nginx

#### 特点
- @todo

---

### 目录结构

<details>
  <summary>展开</summary>
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
</details>

----

### .env 环境文件配置

```
#develop 测试开发环境，对应的配置为config/develop
#prod 生产环境，对应配置为config/prod
env=develop
```

----

### 配置目录下的文件示例

<details>
 <summary>数据库配置 db.json</summary>
 
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
</details>


<details>
  <summary>小程序配置 applet.json</summary>

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
</details>

<details>
  <summary>redis配置 redis.json</summary>

  ```json
  {
    "ip": "127.0.0.1",
    "port": "6379",
    "password": "",
    "database": 0
  }
  ```
</details>

<details>
  <summary>websocket配置 websocket.json</summary>

  ```json
  {
    "addr": ":9504",
    "port": "9504",
    "ip": "127.0.0.1",
    "path": "/wss"
  }
  ```
</details>

<details>
  <summary>mongo配置 mongo.json</summary>

  ```json
  {
    "ip": "127.0.0.1",
    "port": "27017",
    "username": "root",
    "password": "root"
  }
  ```
</details>

----

### 关于 http 调用

> 根域名 https://api.sunanzhi.com

**请求返回格式 `json`**

```json
{
    "_code": 0,
    "_data": null,
    "_message": "success"
}
```

**_code说明**

value | desc
----- | ----
0 | 正常
100 | 逻辑异常

----

**接口列表**

<details>
  <summary> 小程序登录 </summary>
  
  #### url:/auth/appletLogin

  > 请求参数

  key | type | desc
  --- | ---- | ----
  code | string | wx.Login 获取的code
  encryptedData | string | 小程序获取的加密数据
  iv | string | 小程序获取的iv
  appId | string | 小程序的appId

  > 返回数据说明

  key | type | desc
  --- | ---- | ----
  accessToken | string | token
  expire | int | 过期时间
  userInfo | array | 用户基本信息
  userInfo.role | int | 角色 0：普通用户 1：管理员 2：超级管理员
  userInfo.userId | int | 用户id
  userInfo.username | string | 用户名

  > 返回数据示例

  ```json

  {
      "accessToken": "db029874-c31a-4e3c-924f-b1918099dd73",
      "expire": 604800,
      "userInfo": {
          "role": 0,
          "userId": 1,
          "username": "sunanzhi"
      }
  }
  ```
</details>

----

### @todo 关于 wss 链接

> wss://api.sunanzhi.com/wss?Authorization=yourtokenvalue

