import { defineConfig, type DefaultTheme } from 'vitepress'

export const zh = defineConfig({
  lang: 'zh-Hans',
  description: '基于 go-kit 封装的微服务框架',

  themeConfig: {
    nav: nav(),

    sidebar: {
      '/guide/': { base: '/guide/', items: sidebarGuide() },
      '/reference/': { base: '/reference/', items: sidebarReference() }
    },

    editLink: {
      pattern: 'https://github.com/fengjx/luchen/tree/master/docs/:path',
      text: '在 GitHub 上编辑此页面'
    },

    footer: {
      message: '基于  Apache-2.0 License 许可发布',
      copyright: `Copyright © ${new Date().getFullYear()}-present 路辰 <a href="http://beian.miit.gov.cn/" target="_blank">粤ICP备15021633号</a>`
    },

    docFooter: {
      prev: '上一页',
      next: '下一页'
    },

    outline: {
      label: '页面导航'
    },

    lastUpdated: {
      text: '最后更新于'
    },

    langMenuLabel: '多语言',
    returnToTopLabel: '回到顶部',
    sidebarMenuLabel: '菜单',
    darkModeSwitchLabel: '主题',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式'
  }
})

function nav(): DefaultTheme.NavItem[] {
  return [
    {
      text: '指南',
      link: '/guide/introduction',
      activeMatch: '/guide/'
    },
    {
      text: '参考',
      link: '/reference/helloworld',
      activeMatch: '/reference/'
    },
    {
      text: '实践案例',
      items: [
        {
          text: "低代码后台系统",
          link: '/glca/introduction',
          activeMatch: '/glca/'
        }
      ],
    },
    {
      text: 'GoDoc',
      link: 'https://pkg.go.dev/github.com/fengjx/luchen',
    }

  ]
}

function sidebarGuide(): DefaultTheme.SidebarItem[] {
  return [
    {
      text: '开始',
      collapsed: false,
      items: [
        { text: '简介', link: 'introduction' },
        { text: '快速开始', link: 'getting-started' },
      ]
    },
    {
      text: '服务端',
      collapsed: false,
      items: [
        { text: 'HTTP server', link: 'http-server' },
        { text: 'gRPC server', link: 'grpc-server' },
        { text: '服务注册', link: 'register' },
        { text: 'gateway（alpha）', link: 'gateway' },
      ]
    },
    {
      text: '客户端',
      collapsed: false,
      items: [
        { text: 'HTTP client', link: 'http-client' },
        { text: 'gRPC client', link: 'grpc-client' },
      ]
    },
    {
      text: '端点',
      collapsed: false,
      items: [
        { text: '端点定义', link: '' },
        { text: '中间件', link: '' },
      ]
    },
    {
      text: '内置功能',
      collapsed: false,
      items: [
        { text: '环境', link: 'env' },
        { text: '配置加载', link: 'config' },
        { text: '日志', link: 'log' },
      ]
    }
  ]
}

function sidebarReference(): DefaultTheme.SidebarItem[] {
  return [
    {
      text: '参考',
      items: [
        { text: 'helloworld', link: 'helloworld' },
        { text: '多协议支持和服务注册', link: 'multi-protocol-server' },
        { text: '网关', link: 'gateway' },
      ]
    }
  ]
}

export const search: DefaultTheme.AlgoliaSearchOptions['locales'] = {
  zh: {
    placeholder: '搜索文档',
    translations: {
      button: {
        buttonText: '搜索文档',
        buttonAriaLabel: '搜索文档'
      },
      modal: {
        searchBox: {
          resetButtonTitle: '清除查询条件',
          resetButtonAriaLabel: '清除查询条件',
          cancelButtonText: '取消',
          cancelButtonAriaLabel: '取消'
        },
        startScreen: {
          recentSearchesTitle: '搜索历史',
          noRecentSearchesText: '没有搜索历史',
          saveRecentSearchButtonTitle: '保存至搜索历史',
          removeRecentSearchButtonTitle: '从搜索历史中移除',
          favoriteSearchesTitle: '收藏',
          removeFavoriteSearchButtonTitle: '从收藏中移除'
        },
        errorScreen: {
          titleText: '无法获取结果',
          helpText: '你可能需要检查你的网络连接'
        },
        footer: {
          selectText: '选择',
          navigateText: '切换',
          closeText: '关闭',
          searchByText: '搜索提供者'
        },
        noResultsScreen: {
          noResultsText: '无法找到相关结果',
          suggestedQueryText: '你可以尝试查询',
          reportMissingResultsText: '你认为该查询应该有结果？',
          reportMissingResultsLinkText: '点击反馈'
        }
      }
    }
  }
}
