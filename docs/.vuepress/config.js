module.exports = {
  title: 'Vgo文档',
  description: 'Vgo文档',
  base: '/vgo/',
  locales: {
    '/': {
      lang: 'zh-CN',
      title: 'Vgo文档',
      description: 'Vgo文档',
    },
  },
  themeConfig: {
    logo: '/logo-admin-new.png',
    repo: 'https://github.com/vera-byte/vgo',
    docsDir: 'docs',
    editLinks: true,
    editLinkText: '在 GitHub 上编辑此页',
    nav: [
      { text: '首页', link: '/' },
      { text: 'vAdmin官网', link: 'https://v-js.com' },
      { text: 'GoFrame官网', link: 'https://goframe.org' },
    ],
    lastUpdated: '上次更新',
    sidebar: [
      "/",
      "/introduction",
      "/feedback",
      "/development",
      "/cli",
      "/quick_start",
      "/config",
      "/changelog",
      "/known_issues",
    ],
  },
}
