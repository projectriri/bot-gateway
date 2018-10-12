# Getting Started

## Introduction

A bot developer, with the use of The Little Daemon Gateway,
will get emancipated from supporting various Instant Messaging platforms
and be able to concentrate on the application,
thanks to the ability of The Little Daemon Gateway
to convert different API of different IM platforms into an unified API.

Meanwhile, The Little Daemon Gateway is blessed with a feature that
it can route messages reliably, which means that you may divide your
bot into different parts, written in different kind of languages.
No need to worry about how the programs communicate with each other,
all they need is just connecting to the gateway.

## Comparision with Other Solutions

### [Bot-SDK](https://github.com/projectriri/bot-sdk)

Bot-SDK is an alternative solution to use one API across many platforms.
Compared with Bot-SDK:


+ Bot-SDK 的接口在 SDK 上层封装，其本身的开发成本相比小恶魔网关更低。
+ Bot-SDK 使用时可以直接调用 SDK 的相关实例和方法，使用起来更方便。
+ 要使用 Bot-SDK，用户必须使用 Go 语言开发应用。而小恶魔网关是基于协议的，
用户可以使用任何语言开发。得益于可扩展的插件系统，小恶魔网关能够支持 RPC、HTTP、WebSocket 等多种通信协议。
+ 小恶魔网关在允许用户用一套 API 服务多个聊天平台的同时，还有路由功能。
即当你有多个 Bot 应用程序时，他们可以接入网关，使用共同的聊天平台账号。
而使用 Bot-SDK 你要为每个应用分配一组聊天平台账号。同时借助小恶魔网关你也可以实现 Bot 应用程序之间的通信。

### [Telegram Bot Gateway](https://gitlab.com/FiveYellowMice/telegram-bot-gateway)

Telegram Bot Gateway 是专门针对 Telegram Bot API 的多个 Bot 使用同一账号的解决方案。
小恶魔网关与 Telegram Bot Gateway 相比：

+ Telegram Bot Gateway 是专门针对 Telegram 多个 Bot 使用同一账号的应用情景设计的：
例如当用户使用的命令多个 Bot 都能够响应时，可以弹出键盘询问用户想要发送给哪个 Bot；
而小恶魔网关虽然具备消息路由功能，却并不是为了让多个互不相关的 Bot 使用同一账号设计的。
小恶魔网关设计上是为了使 Bot 应用程序之间能够紧密联系和配合，共同提供作为 1 个 Bot 的功能。
+ 小恶魔网关在消息路由之外还具备同时服务多个聊天平台，以及进行 API 转换的功能。
+ 由于 Telegram Bot Gateway 设计上可以作为一个透明代理，小恶魔网关可以在 Telegram Bot Gateway 下运行。

## Install

### Download Binary File

1. 从 [GitHub Release](https://github.com/projectriri/bot-gateway/releases)
中下载已经编译好的程序压缩包，解压到任意目录
2. 复制 `config.toml.example` 到 `config.toml`
3. 删除 `lib` 目录中不需要的插件
4. 复制 `conf.d` 目录下所有 `*.example` 文件到 `*`

::: tip
如果你只需要很少的官方插件，你也可以选择分别下载它们，然后将插件放置在
`lib` 目录下，将插件配置文件放置在 `conf.d` 目录下
:::

### Build from Source

::: warning
请使用 go >= 1.10
:::

1. 克隆[源码仓库](https://github.com/projectriri/bot-gateway)
2. 编辑 `plugins.txt`，加入所有想要的插件，删除不想要的插件
3. `make`

## Configuration

### Root Configuration

根配置文件应命名为 `config.toml` 放置在程序运行目录下。

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| buffer_size | 65536 | 缓冲区大小：单个[频道](/docs/Concept.html#频道)消息最大缓存条数，可以根据内存大小修改 |
| log_level | info | 日志等级，可选 debug、info、warn、error、fatal |
| plugin_dir | lib | 插件加载的目录 |
| plugin_conf_dir | conf.d | 插件配置文件加载的目录 |
| enable_little_daemon | true | 是否启用小恶魔 |
| little_daemon_name | LittleDaemon | 小恶魔在网关中的适配器名称 |
| little_daemon_channel_uuid | | 小恶魔用于注册[频道](/docs/Concept.html#频道)的 UUID，可为空 |

“小恶魔”是用来报告网关状态的内置机器人。
如果小恶魔被启用，当网关收到命令 `ping` 或 `status` 时，会报告网关当前状态。
启用小恶魔你需要在程序运行目录放置 `locale.yml`，程序源码中包含一个该文件的[示例](https://github.com/projectriri/bot-gateway/blob/master/locale.yml)。

:::tip
要让小恶魔正常工作你必须启用插件 [commander](/docs/Plugins.html#commander)，并根据需要正确配置命令前缀。
:::

### Plugin Configurations

插件配置文件应命名为“插件文件名.配置扩展名”，放置在 *plugin_conf_dir* 目录下：
例如某插件 `chajian.so` 接受格式为 toml 的配置，那么配置将从 *plugin_conf_dir*
目录下的 `chajian.toml` 加载。

插件配置文件的定义见[插件文档](/docs/Plugins.html)。

## List of Plugins

**客户型**插件表示插件是对外表现为一个客户端，这类插件一启动就向相应的服务器发起连接，
从服务器接收消息，向服务器发送消息。

**服务型**插件表示插件对外表现为一个服务端，这类插件启动时将启动一个服务程序，
外界程序可以通过连接到这个插件来向网关发送和接收消息。

**自耦型**插件表示插件不对外通信，它从网关接收消息，处理后直接送回网关。

| 插件名 | 类型 | 功能描述 |
| --- | --- | --- |
| [longpolling-client-tgbot](/docs/Plugins.html#longpolling-client-tgbot) | 客户型 | 从 Telegram Bot API 以长轮询形式拉取更新 |
| [http-client-tgbot](/docs/Plugins.html#http-client-tgbot) | 客户型 | 调用 Telegram Bot API 发送消息 |
| [websocket-client-cqhttp](/docs/Plugins.html#websocket-client-cqhttp) | 客户型 | 从 CoolQ HTTP API 正向 WebSocket 连接上接收事件上报和 API 调用 |
| [jsonrpc-server-any](/docs/Plugins.html#jsonrpc-server-any) | 服务型 | 提供 JSON RPC 服务 |
| [tgbot-ubm-conv](/docs/Plugins.html#tgbot-ubm-conv) | 自耦型 | 将消息在 [Telegram-Bot-API](/docs/Formats.html) 和 [UBM-API](/docs/Formats.html) 两种格式之间转换 |
| [cqhttp-ubm-conv](/docs/Plugins.html#cqhttp-ubm-conv) | 自耦型 | 将消息在 [CoolQ-HTTP-API](/docs/Formats.html) 和 [UBM-API](/docs/Formats.html) 两种格式之间转换 |
| [commander](/docs/Plugins.html#commander) | 自耦型 | 命令解析器 |
