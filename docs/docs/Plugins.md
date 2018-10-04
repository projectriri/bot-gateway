# 插件文档

::: warning
到目前为止，本页上所有插件都只支持 toml 格式的配置文件。
:::

## longpolling-client-tgbot

这个插件通过[长轮询](https://core.telegram.org/bots/api#getupdates)从 [Telegram Bot API](https://core.telegram.org/bots/api)
拉取 [更新](https://core.telegram.org/bots/api#update)，将更新送入网关中。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| token |（必填）| [Telegram Bot Token](https://core.telegram.org/bots/api#authorizing-your-bot) |
| adaptor_name | Telegram | 适配器名称：如果你开了两个这个插件（如你用两个功能完全相同机器人账号进行负载均衡），可以通过适配器名称来区分 |
| timeout | 60 | [长轮询超时时间](https://core.telegram.org/bots/api#getupdates) |
| limit | 100 | [长轮询单次拉取的消息上限](https://core.telegram.org/bots/api#getupdates) |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID

这是一个[生产者](/docs/Concept.html#生产者)，它生产的[包](/docs/Concept.html#包)的头为

| from | to | format.api | format.version | format.method | format.protocol |
| --- | --- | --- | --- | --- | --- |
| *adaptor_name* | | telegram-bot-api | latest | update | http |

它生产的[包](/docs/Concept.html#包)的体的类型为 [APIResponse](/docs/Other.html#apiresponse)，其中 Result 为 [Update](https://core.telegram.org/bots/api#update)[]

## http-client-tgbot

这个插件用于请求[Telegram Bot API](https://core.telegram.org/bots/api)。

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

## tgbot-ubm-conv

这个插件能够将包在 *telegram-bot-api* 格式和 *ubm-api* 格式之间转换。

**配置文件格式**

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| telegram_adaptors | Telegram | 转换的包可能来自的适配器名（正则表达式） |
| adaptor_name | TGBot-UBM-Converter | 适配器名称 |
| fetch_file | true | 是否拉取 [file_path](https://core.telegram.org/bots/api#file) 作为 URL |
| fetch_file_timeout | 5s | 拉取 [file_path](https://core.telegram.org/bots/api#file) 的最长等待时间 |
| channel_uuid | | 插件用于注册[频道](/docs/Concept.html#频道)的 UUID |

**支持转换的格式**

| from.api | from.version | from.method | from.protocol | to.api | to.version | to.method | to.protocol |
| --- | --- | --- | --- | --- | --- | --- | --- |
| telegram-bot-api | latest | update | http | ubm-api | 1.0 | receive |（不限）|
| ubm-api | 1.0 | send |（不限）| telegram-bot-api | latest | apirequest | http |

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