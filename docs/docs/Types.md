# 类型定义

本页上定义了一些类型，可能会出现在[包](/docs/Concept.html#包)的体中。

## Common

通用协议类型。

### HTTPRequest

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| method | string | 请求方法："GET"，"POST" 等 |
| url | string | 请求地址 |
| header | { [key: string]: string[] } | HTTP Headers |
| body | byte[] | HTTP Body |

## CMD

解析后的命令。下表中前五个字段都是可选字段，是否存在取决于插件的[配置](/docs/Plugins.html#commander)。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| cmd | [RichTextElement](#richtextelement)[] | 命令：由一组图文消息段组成 |
| cmd_str | string | 命令：纯字符串，去除了其他媒体（如图片、At 等）以后的命令 |
| args | [RichTextElement](#richtextelement)[][] | 命令参数：有多个命令参数，每个命令参数由一组图文消息段组成 |
| args_txt | string[] | 命令参数：多个命令参数，每个命令参数都是去除了其他媒体（如图片、At 等）以后的纯字符串 |
| args_str | string | 命令参数字符串：由 *args_txt* 以空格分割组成个一个字符串 |
| cmd_prefix | string | 命令前缀 |
| message | [Message](#message) | 命令对应的消息 |

## UBM-API

通用机器人消息 API。

### UBM

通用机器人消息。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| type | string | 消息类型 |
| message | [Message](#message) | 聊天消息（消息类型是 message 时需要） |
| notice | [Notice](#notice) | 系统消息：如加群通知等（消息类型是 notice 时需要，仅用于收到的消息） |
| response | [Response](#response) | 响应消息：发出 Message 时的发送结果，或发出 Action 时的响应数据（消息类型是 response 时需要，仅用于收到的消息） |
| action | [Action](#action) | 操作消息：退群等非聊天请求（消息类型是 action 时需要，仅用于发出的消息） |
| self | [User](#user) | 自己的用户（仅用于收到的消息）|

### Message

聊天消息。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| id | string | 消息 ID，只在当前会话中唯一（仅用于收到的消息）|
| from | [User](#user) | 消息来自的用户（仅用于收到的消息）|
| chat | [Chat](#chat) | 会话（仅用于收到的消息）|
| cid | [CID](#cid) | 会话 CID（仅用于发出的消息）|
| is_message_to_me | bool | 是否 at 或回复我（仅用于收到的消息）|
| reply_id | string | 回复的消息 ID |
| edit_id | string | 编辑的消息 ID |
| delete_id | string | 撤回的消息 ID |
| type | string | 消息类型 |
| rich_text | [RichTextElement](#richtextelement)[] | 图文（消息类型是 rich_text 时需要）|
| sticker | [Sticker](#sticker) | 贴纸（消息类型是 sticker 时需要）|
| voice | [Voice](#voice) | 语音（消息类型是 voice 时需要） |
| location | [Location](#location) | 位置（消息类型是 location 时需要）|

### CID

CID 标识一个会话 ID。任何两个不同会话的 CID 都不相同。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| messenger | string | 平台，实际为插件配置中的适配器名 |
| chat_id | string | 平台中的会话 ID |
| chat_type | string | 平台的会话类型 |

### Chat

会话。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| cid | [CID](#cid) | 会话 CID |
| title | string | 会话名（用户昵称或群名）|
| description | string | 会话描述 |

### UID

UID 标识一个用户 ID。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| messenger | string | 平台，实际为插件配置中的适配器名 |
| id | string | 平台中的用户 ID |
| username | string | 用户名（用于在文本中@一个用户的全平台唯一用户名，有些平台可能没有） |

### User

用户。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| display_name | string | 显示名，取决于平台和会话 |
| first_name | string | 名字 |
| last_name | string | 姓（有些平台可能没有）|
| uid | [UID](#uid) | 用户 UID |
| private_chat | [CID](#cid) | 私聊该用户的会话 CID |

### RichTextElement

图文消息段。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| type | string | 类型 |
| text | string | 文本（消息类型是 text 时需要）|
| at | [At](#at) | At（消息类型是 at 时需要）|
| image | [Image](#image) | 图片（消息类型是 image 时需要）|

### At

At.

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| display_name | string | 有些平台允许 at 时自定义显示名 |
| uid | [UID](#uid) | 要 at 的用户 UID |

### Image

图片。

接收时可选参数不一定有，取决于平台。接收时没有 data 字段，只有 url 和
file_id 字段。接收到的 url 不一定可以直接下载，因为可能需要 token，
可以通过发送 [Action](#action) 中的 [GetFile](#getfile) 获取。

发送时只需要提供 data，url，file_id 之一即可：url 为公网可访问地址。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| messenger | string | 平台（仅用于收到的消息） |
| format | string |（可选）图片格式 |
| width | int |（可选）图片宽度 |
| height | int |（可选）图片高度 |
| data | byte[] | 图片 |
| url | string | 图片 URL |
| file_id | string | 文件 ID，取决于平台 |
| file_size | int |（可选）文件大小 |

### Sticker

贴纸。接收到的消息是否拉取图片取决于插件的配置。
发送时发送（id，pack_id）和 image 任一即可。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| messenger | string | 平台（仅用于收到的消息） |
| id | string | 贴纸 ID，有些平台是唯一的，有些平台取决于贴纸包 ID |
| pack_id | string | 贴纸包 ID，有些平台不存在 |
| image | [Image](#image) | 图片 |

### Voice

语音。

接收时可选参数不一定有，取决于平台。接收时没有 data 字段，只有 url 和
file_id 字段。接收到的 url 不一定可以直接下载，因为可能需要 token，
可以通过发送 [Action](#action) 中的 [GetFile](#getfile) 获取。

发送时只需要提供 data，url，file_id 之一即可：url 为公网可访问地址。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| messenger | string | 平台（仅用于收到的消息） |
| format | string |（可选）音频格式 |
| duration | int |（可选）音频时长 |
| data | byte[] | 音频 |
| url | string | 音频 URL |
| file_id | string | 文件 ID，取决于平台 |
| file_size | int |（可选）文件大小 |

### Location

位置。

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| content | string | 详细地址 |
| latitude | float64 | 纬度 |
| longitude | float64 | 经度 |
| title | string | 标题 |

