module.exports = {
    title: 'Bot-Gateway',
    description: 'The one and only one Bot API you need!',
    base: '/bot-gateway/',
    locales: {
        // The key is the path for the locale to be nested under.
        // As a special case, the default locale can use '/' as its path.
        '/': {
            lang: 'zh-CN',
            title: '小恶魔机器人网关',
            description: '用同一套代码服务多个平台的机器人',
        },
        '/en/': {
            lang: 'en',
            title: 'Bot-Gateway',
            description: 'The one and only one Bot API you need!',
        }
    },
    themeConfig: {
        sidebar: 'auto',
        sidebarDepth: 3,
        locales: {
            '/': {
                // 多语言下拉菜单的标题
                selectText: '选择语言',
                // 该语言在下拉菜单中的标签
                label: '简体中文',
                // 编辑链接文字
                editLinkText: '在 GitHub 上编辑此页',
                sidebar: [
                    '/ChangeLog',
                    '/docs/',
                    '/docs/Concept',
                    '/docs/Types',
                    '/docs/Plugins',
                    '/docs/Formats',
                    '/docs/Other',
                ]
            },
            '/en/': {
                selectText: 'Languages',
                label: 'English',
                editLinkText: 'Edit this page on GitHub',
                sidebar: [
                    '/en/ChangeLog',
                    '/en/docs/',
                    '/en/docs/Concept',
                    '/en/docs/Types',
                    '/en/docs/Plugins',
                    '/en/docs/Formats',
                    '/en/docs/Other',
                ]
            }
        },
        repo: 'projectriri/bot-gateway',
        editLinks: true,
    },
    configureWebpack: (config, isServer) => {
        config.module.rules.push({
            test: /\.(webp)$/,
            loader: 'url-loader',
            options: {
                limit: 8192
            }
        })
    }
}