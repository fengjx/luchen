import { defineConfig } from 'vitepress'
import { zh, search as zhSearch } from './zh'

export default defineConfig({
  title: 'luchen',

  lastUpdated: true,
  cleanUrls: true,
  ignoreDeadLinks: true,
  markdown: {
    math: true,
    codeTransformers: [
      // We use `[!!code` in demo to prevent transformation, here we revert it back.
      {
        postprocess(code) {
          return code.replace(/\[\!\!code/g, '[!code')
        }
      }
    ]
  },

  sitemap: {
    hostname: 'http://luchen.fun'
  },

  /* prettier-ignore */
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/luchen-logo.svg' }],
    ['link', { rel: 'icon', type: 'image/png', href: '/luchen-logo.png' }],
    ['meta', { name: 'theme-color', content: '#5f67ee' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:locale', content: 'en' }],
    ['meta', { name: 'og:site_name', content: 'VitePress' }],
    ['meta', { name: 'og:image', content: 'https://vitepress.dev/vitepress-og.jpg' }],
    ['meta', { name: 'keywords', content: 'luchen go-kit golang go 微服务 网关' }],
    ['script', { src: 'https://cdn.usefathom.com/script.js', 'data-site': 'AZBRSFGG', 'data-spa': 'auto', defer: '' }]
  ],

  themeConfig: {
    logo: { src: '/luchen-logo.svg', width: 24, height: 24 },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/fengjx/luchen' }
    ],

    search: {
      provider: 'local',
      options: {
        locales: { ...zhSearch }
      }
    },
    // carbonAds: { code: 'CEBDT27Y', placement: 'vuejsorg' }
  },

  locales: {
    root: { label: '中文', ...zh },
  }
})
