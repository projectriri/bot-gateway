# 使用教程

## 功能简介

小恶魔网关能够将不同即时通讯平台的 API 转换为统一的 API，
从而使得机器人开发者无须关心多个平台的 API 适配，
用统一的 API 来调用不同平台，从而更加专注于业务代码。

小恶魔网关同时拥有可靠的消息路由功能，
你的机器人可以分为多个程序，使用不同的语言编写，
而无须定义不同程序之间的调用协议：它们可以统一接入网关来实现通信。

## 下载安装

### 下载二进制程序

1. 从 [GitHub Release](https://github.com/projectriri/bot-gateway/releases)
中下载已经编译好的程序压缩包，解压到任意目录
2. 复制 `config.toml.example` 到 `config.toml`
3. 删除 `lib` 目录中不需要的插件
4. 复制 `conf.d` 目录下所有 `*.example` 文件到 `*`

::: tip
如果你只需要很少的官方插件，你也可以选择分别下载它们，然后将插件放置在
`lib` 目录下，将插件配置文件放置在 `conf.d` 目录下
:::

### 从源码编译安装

::: warning
请使用 go >= 1.10
:::

1. 克隆[源码仓库](https://github.com/projectriri/bot-gateway)
2. 编辑 `plugins.txt`，加入所有想要的插件，删除不想要的插件
3. `make`

## 配置文件

### 根配置文件

根配置文件应命名为 `config.toml` 放置在程序运行目录下。

| 配置项名称 | 默认配置文件中的值 | 说明 |
| --- | --- | --- |
| buffer_size | 65536 | 缓冲区大小：单个[频道](/docs/Concept.html#频道)消息最大缓存条数，可以根据内存大小修改 |
| log_level | info | 日志等级，可选 debug、info、warn、error、fatal |
| plugin_dir | lib | 插件加载的目录 |
| plugin_conf_dir | conf.d | 插件配置文件加载的目录 |

### 插件配置文件

插件配置文件应命名为“插件文件名.配置扩展名”，放置在 *plugin_conf_dir* 目录下：
例如某插件 `chajian.so` 接受格式为 toml 的配置，那么配置将从 *plugin_conf_dir*
目录下的 `chajian.toml` 加载。

插件配置文件的定义见[插件文档](/docs/Plugins.html)。

## 插件列表

**客户型**插件表示插件是对外表现为一个客户端，这类插件一启动就向相应的服务器发起连接，
从服务器接收消息，向服务器发送消息。

**服务型**插件表示插件对外表现为一个服务端，这类插件启动时将启动一个服务程序，
外界程序可以通过连接到这个插件来向网关发送和接收消息。

**自耦型**插件表示插件不对外通信，它从网关接收消息，处理后直接送回网关。

| 插件名 | 类型 | 功能描述 |
| --- | --- | --- |
| [longpolling-client-tgbot](/docs/Plugins.html#longpolling-client-tgbot) | 客户型 | 从 Telegram Bot API 以长轮询形式拉取更新 |
| [http-client-tgbot](/docs/Plugins.html#http-client-tgbot) | 客户型 | 调用 Telegram Bot API 发送消息 |
| [jsonrpc-server-any](/docs/Plugins.html#jsonrpc-server-any) | 服务型 | 提供 JSON RPC 服务 |
| [tgbot-ubm-conv](/docs/Plugins.html#tgbot-ubm-conv) | 自耦型 | 将消息在 [Telegram-Bot-API](/docs/Formats.html) 和 [UBM-API](/docs/Formats.html) 两种格式之间转换 |
