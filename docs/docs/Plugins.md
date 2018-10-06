# 插件文档

::: warning
到目前为止，本页上所有插件都只支持 toml 格式的配置文件。
:::

## longpolling-client-tgbot

这个插件通过[长轮询](https://core.telegram.org/bots/api#getupdates)从 [Telegram Bot API](https://core.telegram.org/bots/api)
拉取[更新](https://core.telegram.org/bots/api#update)，将更新送入网关中。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| token |（必填）| [Telegram Bot Token](https://core.telegram.org/bots/api#authorizing-your-bot) |
| adaptor_name | Telegram | 适配器名称：如果你开了两个这个插件（如你用两个功能完全相同机器人账号进行负载均衡），可以通过适配器名称来区分 |
| timeout | 60 | [长轮询超时时间](https://core.telegram.org/bots/api#getupdates) |
| limit | 100 | [长轮询单次拉取的消息上限](https://core.telegram.org/bots/api#getupdates) |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID，可为空

这是一个[生产者](/docs/Concept.html#生产者)，它生产的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| *adaptor_name* | | telegram-bot-api | latest | update | http |

它生产的[包](/docs/Concept.html#包)的体的类型为 [APIResponse](/docs/Other.html#apiresponse)，其中 Result 为 [Update](https://core.telegram.org/bots/api#update)[]

## http-client-tgbot

这个插件用于请求 [Telegram Bot API](https://core.telegram.org/bots/api)。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| token |（必填）| [Telegram Bot Token](https://core.telegram.org/bots/api#authorizing-your-bot) |
| adaptor_name | Telegram | 适配器名称：如果你开了两个这个插件（如你用两个功能完全相同机器人账号进行负载均衡），可以通过适配器名称来区分 |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID

这是一个[消费者](/docs/Concept.html#消费者)，它接受的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| .* | *adaptor_name* | telegram-bot-api | latest | apirequest | http |

它接受的[包](/docs/Concept.html#包)的体的类型为 [HTTPRequest](/docs/Types.html#httprequest)。

这是一个[生产者](/docs/Concept.html#生产者)，它生产的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| *adaptor_name* | *回复的包的 from* | telegram-bot-api | latest | apiresponse | http |

它生产的[包](/docs/Concept.html#包)的体的类型为 [APIResponse](/docs/Other.html#apiresponse)。

## websocket-client-cqhttp

这个插件与 [CoolQ HTTP API](https://cqhttp.cc) 建立正向 WebSocket 连接，
接收事件上报和进行 API 请求。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| adaptor_name | QQ | 适配器名称：如果你开了两个这个插件（如你用两个功能完全相同机器人账号进行负载均衡），可以通过适配器名称来区分 |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID |
| cqhttp_access_token | | CoolQ HTTP API 的 [access_token](https://cqhttp.cc/docs/4.4/#/Configuration) |
| cqhttp_websocket_addr | ws://localhost:6700 | CoolQ HTTP API 正向 WebSocket 地址 |
| cqhttp_version | latest | CoolQ HTTP API 版本 |

这是一个[消费者](/docs/Concept.html#消费者)，它接受的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| .* | *adaptor_name* | coolq-http-api | latest | apirequest | websocket |

它接受的[包](/docs/Concept.html#包)的体的类型为 [APIRequest](https://cqhttp.cc/docs/4.4/#/WebSocketAPI?id=api-%E6%8E%A5%E5%8F%A3)。

这是一个[生产者](/docs/Concept.html#生产者)，它生产的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| *adaptor_name* | | coolq-http-api | latest | event | websocket |

它生产的[包](/docs/Concept.html#包)的体的类型为 [Update](https://cqhttp.cc/docs/4.4/#/Post?id=%E4%B8%8A%E6%8A%A5%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F)。

这是一个[生产者](/docs/Concept.html#生产者)，它生产的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| *adaptor_name* | *回复的包的 from* | coolq-http-api | latest | apiresponse | websocket |

它生产的[包](/docs/Concept.html#包)的体的类型为 [APIResponse](https://cqhttp.cc/docs/4.4/#/API?id=%E5%93%8D%E5%BA%94%E8%AF%B4%E6%98%8E)。

## tgbot-ubm-conv

这个插件能够将包在 *telegram-bot-api* 格式和 *ubm-api* 格式之间转换。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| telegram_adaptors | Telegram | 转换的包可能来自的适配器名（正则表达式） |
| adaptor_name | TGBot-UBM-Converter | 适配器名称 |
| fetch_file | true | 是否拉取 [file_path](https://core.telegram.org/bots/api#file) 作为 URL |
| api_response_timeout | 5s | 拉取 [file_path](https://core.telegram.org/bots/api#file) 和 [self](https://core.telegram.org/bots/api#getme) 的最长等待时间 |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID，可为空 |

**支持转换的格式**

| from.api | from.version | from.method | from.protocol | to.api | to.version | to.method | to.protocol |
| --- | --- | --- | --- | --- | --- | --- | --- |
| telegram-bot-api | latest | update | http | ubm-api | 1.0 | receive |（不限）|
| ubm-api | 1.0 | send |（不限）| telegram-bot-api | latest | apirequest | http |

## cqhttp-ubm-conv

这个插件能够将包在 *coolq-http-api* 格式和 *ubm-api* 格式之间转换。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| qq_adaptors | QQ | 转换的包可能来自的适配器名（正则表达式） |
| adaptor_name | CQHTTP-UBM-Converter | 适配器名称 |
| api_response_timeout | 5s | 拉取 [self](https://cqhttp.cc/docs/4.4/#/API?id=get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF) 的最长等待时间 |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID，可为空 |

**支持转换的格式**

| from.api | from.version | from.method | from.protocol | to.api | to.version | to.method | to.protocol |
| --- | --- | --- | --- | --- | --- | --- | --- |
| coolq-http-api | >=3 | event | websocket | ubm-api | 1.0 | receive |（不限）|
| ubm-api | 1.0 | send |（不限）| coolq-http-api | >=3 | apirequest | websocket |

::: tip
头中定义的格式与体中数据结构的关系可以通过查询[格式对照表](/docs/Formats.html)找到。
:::

## jsonrpc-server-any

这个插件提供一个基于 TCP 的 JSON RPC 服务。用户可以通过这个插件建立频道，自定义路由规则，发送和接收任意格式的包。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| port | 4700 | 服务监听端口 |
| channel_lifetime | 2h | 用户断线后频道最长保存时间 |
| garbage_collection_interval | 5m | 垃圾回收检查周期 |

### Broker.InitChannel

创建频道。

**请求参数**

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| uuid | string | 频道 UUID |
| producer | bool | 是否建立生产者频道 |
| consumer | bool | 是否建立消费者频道 |
| accept | [RoutingRule](/docs/Concept.html#路由规则)[] | 消费者频道接受的路由规则（只在 consumer 为 true 时有效）|

**响应内容**

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| uuid | string | 建立的频道的 UUID（如果请求的 UUID 不合法，将会生成一个） |
| code | int | 响应码 |

### Broker.Send

发送包。

**请求参数**

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| uuid | string | 频道 UUID |
| packet | [Packet](/docs/Concept.html#包) | 要发送的包 |

**响应内容**

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| code | int | 响应码 |

### Broker.GetUpdates

通过长轮询接收包。

**请求参数**

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| uuid | string | 频道 UUID |
| timeout | string | 没有收到消息时的最大超时时长，例如："5s" |
| limit | int | 单次请求最多返回的包数量（1-100）|

**响应内容**

| 字段名 | 数据类型 | 说明 |
| --- | --- | --- |
| code | int | 响应码 |
| packets | [Packet](/docs/Concept.html#包)[] | 包 |

### 响应码

| 响应吗 | 说明 |
| --- | --- |
| 10000 | 收到长轮询更新包 |
| 10001 | 成功建立了频道 |
| 10002 | 包发送成功 |
| 10004 | 长轮询没有收到消息超时返回 |
| 10042 | 建立频道时 producer 和 consumer 都是 false |
| 10044 | 要发往/接收的频道不存在 |
| 10048 | 要发往的频道不是一个生产者 |

## commander

这个插件接收 *ubm-api* 格式的包，将符合命令格式要求的消息解析为命令，并送回到网关中。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| command_prefix | ["/"] | 命令前缀，是一个字符串数组。 |
| response_mode | 31 | 用于配置解析后的命令包含哪些内容。 |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID，可为空 |

命令解析器工作时，对于一条图文消息，会首先在 *command_prefix* 中从前向后依次匹配。
如果一条消息的前缀匹配 *command_prefix* 中的元素，它会被当作一条命令，同时命令前缀会被删除。

:::tip
空的 *command_prefix* 数组或者数组中包含空字符串将导致所有图文消息被当作命令。
:::

**响应模式**

| 位 | 5 | 4 | 3 | 2 | 1 |
| --- | --- | --- | --- | --- | --- |
| 名称 | args_str | args_txt | args | cmd_str | cmd |

*response_mode* 是一个范围在 0~31 的整数，它的各位代表启用 [CMD](/docs/Types.html#cmd) 中的各项。
如果某位设置为 1 那么插件生产的命令报文中则会有此项。
例如设置 *response_mode* 为 26，即 11010，则表示启用 args_str、args_txt、cmd_str 三项。

**命令和命令参数**

命令和命令参数之间以 [Unicode 空白字符](https://golang.org/pkg/unicode/#IsSpace) 分隔。
使用引号（`'` 或 `"` 或 `` ` ``）引起来的内容不会被断开。
转义字符 `\` 后面所接的字符会直接作为命令或命令参数的内容，失去其特殊含义（如空格、引号、`\` 本身等）。

这是一个[消费者](/docs/Concept.html#消费者)，它接受的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| .* | .* | ubm-api | 1.0 | receive | |

它接受的[包](/docs/Concept.html#包)的体的类型为 [UBM](/docs/Types.html#ubm)。

这是一个[生产者](/docs/Concept.html#生产者)，它生产的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| *接受的包的 from* | *接受的包的 to* | cmd | 1.0 | cmd |  |

它生产的[包](/docs/Concept.html#包)的体的类型为 [CMD](/docs/Types.html#cmd)。
