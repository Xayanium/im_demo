### 用户表

```text
{
    "username": "xxx",
    "password": "xxx",
    "nickname": "xxx",
    "sex": 1,  // 1-男，2-女
    "email": "xxx",
    "avatar": "xxx",
    "created_at": 1, // 创建时间
    "updated_at": 1, // 更新时间
}
```

### 消息表
```text
{
    "message_id": "xxx", // 消息唯一标识
    "send_user_id": "xxx", // 发送者的唯一标识
    "recv_user_ids": "[xxx,xxx,xxx]", // 接收者的唯一标识
    "room_id": "xxx", // 房间唯一标识
    "data": "xxx", // 发送的数据
    "create_at": 1, // 创建时间
    "updated_at": 1, // 更新时间
}
```


### 房间表
```text
{
    "room_id": "xxx", // 房间唯一标识
    "room_name": "xxx", // 房间名称
    "info": "xxx", // 房间简介
    "user_id": "xxx", // 房间创建者的唯一标识
    "create_at": 1,
    "update_at": 1,
}
```


### 用户与房间关联表
```text
{
    "user_id": "xxx", // 用户唯一标识
    "room_id": "xxx", // 房间唯一标识
    "create_at": 1, // 创建时间
    "updated_at": 1, // 更新时间
}
```