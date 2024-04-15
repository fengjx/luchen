---
layout: home

title: luchen
titleTemplate: 基于go-kit封装的微服务框架

hero:
  name: luchen
  text: 基于 go-kit 封装的微服务框架
  tagline: 开始这个项目的时候，我女儿刚出生，取名【路辰】，所以将项目名命名为 luchen。
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: GitHub
      link: https://github.com/fengjx/luchen
  image:
      src: /luchen-logo.svg
      alt: luchen

features:
  - icon: 🖥
    title: 单体服务
    details: 秉承go-kit的简单，自己选择使用微服务还是单体服务。
  - icon: 🚀
    title: 微服务
    details: 在单体服务的基础上，只需要增加一个 Register 即可完成服务注册。
  - icon: ⚙
    title: 网关支持
    details: 实现了支持静态路由和动态服务发现网关服务，通过插件化很容易对功能进行扩展。
---
<style>
:root {
  --vp-home-hero-name-color: transparent;
  --vp-home-hero-name-background: -webkit-linear-gradient(120deg, #bd34fe 30%, #41d1ff);

  --vp-home-hero-image-background-image: linear-gradient(-45deg, #bd34fe 50%, #47caff 50%);
  --vp-home-hero-image-filter: blur(44px);
}

@media (min-width: 640px) {
  :root {
    --vp-home-hero-image-filter: blur(56px);
  }
}

@media (min-width: 960px) {
  :root {
    --vp-home-hero-image-filter: blur(68px);
  }
}
</style>
