import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Kite CLI',
  description: 'A blazing fast command-line interface for Zerodha Kite',
  base: '/kite-cli/',
  
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/kite-cli/logo.svg' }],
    ['meta', { name: 'theme-color', content: '#387ED1' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'Kite CLI' }],
    ['meta', { property: 'og:description', content: 'A blazing fast command-line interface for Zerodha Kite' }],
  ],

  themeConfig: {
    logo: '/logo.svg',
    
    nav: [
      { text: 'Guide', link: '/' },
      { text: 'GitHub', link: 'https://github.com/kavinkumar999/kite-cli' }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Introduction', link: '/' },
          { text: 'Installation', link: '#installation' },
          { text: 'Setup', link: '#setup' },
        ]
      },
      {
        text: 'Usage',
        items: [
          { text: 'Multi-Account Support', link: '#multi-account-support' },
          { text: 'Trading', link: '#trading' },
          { text: 'Portfolio & Holdings', link: '#portfolio-holdings' },
          { text: 'Market Data', link: '#market-data' },
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'Order Flags', link: '#order-flags' },
          { text: 'Configuration', link: '#configuration' },
          { text: 'Troubleshooting', link: '#troubleshooting' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/kavinkumar999/kite-cli' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024-present Kavin Kumar'
    },

    search: {
      provider: 'local'
    },

    outline: {
      level: [2, 3]
    }
  }
})
