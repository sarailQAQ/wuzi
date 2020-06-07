# 服务器接口简述

## 端口

使用8080端口

## API

### 1. /login

登录功能，使用json格式传输数据，如

```json
{
	username: "customer",
	password: "123456"
}
```



以返回json格式数据，如

```json
{
    "code": 0,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0aGhibXoiLCJleHAiOiIxNTkxNTI0NjgxIiwidXNlciI6ImNncCIsImlkIjoxLCJpYXQiOiIxNTkxNTEzODgxIn0=.v6BD4U1ID+y9pHu0GjCUYc4S17mYq8ZuljzIeiU6hgE=",
        "username": "customer"
    }
}
```

后续请求中需在header设置以token为关键字的值，以便服务器确认登录状态

### 2. /register

注册功能，请求格式与登录相同，返回json格式，如

```json
{
    "code": 0,
    "data": {}
}
```

### 3. /ws/: room

开启五子棋房间，其中room为你的房间名。服务端与客户端建立WebSocket，之后服务端与客户端使用json格式交互数据。

服务端不会保存棋盘，只会保存下棋的步骤（每个人每一步的落点）。胜负由客户端判断，客户端判断后服务器会对判断结果进行检查，如果结果存在异常则以无胜者的方式退出。

游戏正常结束后，所有步骤会存到数据库中。

## 代码运行环境

服务器需要安装mysql，

客户端无要求