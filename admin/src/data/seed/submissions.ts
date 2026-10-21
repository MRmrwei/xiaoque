import type { ContactSubmissionRecord } from '@/types/admin'

/**
 * 演示用联系表单留言。
 * 数据只存在浏览器 localStorage，不会提交到任何后端。
 */
export const submissions: ContactSubmissionRecord[] = [
  {
    id: 'sub_1001',
    formKey: 'contact',
    pageName: '联系我们',
    pagePath: '/contact',
    locale: 'zh',
    name: '赵一鸣',
    email: 'zhaoyiming@example.com',
    phone: '138-0013-8000',
    company: '远山制造',
    message:
      '我们正在升级集团官网，想了解区块化编辑、中英双语内容管理和媒体库的实施周期。方便安排一次 30 分钟的线上演示吗？',
    data: {
      name: '赵一鸣',
      email: 'zhaoyiming@example.com',
      phone: '138-0013-8000',
      company: '远山制造',
      message:
        '我们正在升级集团官网，想了解区块化编辑、中英双语内容管理和媒体库的实施周期。方便安排一次 30 分钟的线上演示吗？'
    },
    status: 'new',
    ip: '203.0.113.24',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/141.0',
    createdAt: '2026-10-19T09:42:00+08:00',
    updatedAt: '2026-10-19T09:42:00+08:00'
  },
  {
    id: 'sub_1002',
    formKey: 'contact',
    pageName: '联系我们',
    pagePath: '/contact',
    locale: 'zh',
    name: '李然',
    email: 'liran@example.com',
    phone: '139-1188-2600',
    company: '星河数据',
    message:
      '英文站目前用独立域名，想确认迁移到 /en 前缀后 hreflang、canonical 和旧链接 301 应该怎么配置。',
    data: {
      name: '李然',
      email: 'liran@example.com',
      phone: '139-1188-2600',
      company: '星河数据',
      message:
        '英文站目前用独立域名，想确认迁移到 /en 前缀后 hreflang、canonical 和旧链接 301 应该怎么配置。'
    },
    status: 'processing',
    ip: '198.51.100.18',
    userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15',
    createdAt: '2026-10-18T16:08:00+08:00',
    updatedAt: '2026-10-19T10:12:00+08:00'
  },
  {
    id: 'sub_1003',
    formKey: 'contact',
    pageName: 'Contact',
    pagePath: '/en/contact',
    locale: 'en',
    name: 'Emily Chen',
    email: 'emily.chen@example.com',
    phone: '+1 415 555 0132',
    company: 'Northwind Studio',
    message:
      'We are a design studio looking for a bilingual site builder for our clients. Could you share the onboarding process and pricing for three sites?',
    data: {
      name: 'Emily Chen',
      email: 'emily.chen@example.com',
      phone: '+1 415 555 0132',
      company: 'Northwind Studio',
      message:
        'We are a design studio looking for a bilingual site builder for our clients. Could you share the onboarding process and pricing for three sites?'
    },
    status: 'done',
    ip: '192.0.2.77',
    userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/140.0',
    createdAt: '2026-10-17T11:26:00+08:00',
    updatedAt: '2026-10-18T09:05:00+08:00'
  },
  {
    id: 'sub_1004',
    formKey: 'contact',
    pageName: '联系我们',
    pagePath: '/contact',
    locale: 'zh',
    name: '王珂',
    email: 'wangke@example.com',
    phone: '137-2200-9104',
    company: '木野品牌',
    message:
      '官网内容更新频率比较高，主要是新闻和客户案例。想了解后台发布后多久生效，以及能不能限制编辑人员直接发布。',
    data: {
      name: '王珂',
      email: 'wangke@example.com',
      phone: '137-2200-9104',
      company: '木野品牌',
      message:
        '官网内容更新频率比较高，主要是新闻和客户案例。想了解后台发布后多久生效，以及能不能限制编辑人员直接发布。'
    },
    status: 'new',
    ip: '203.0.113.91',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edge/140.0',
    createdAt: '2026-10-19T14:31:00+08:00',
    updatedAt: '2026-10-19T14:31:00+08:00'
  },
  {
    id: 'sub_1005',
    formKey: 'contact',
    pageName: '联系我们',
    pagePath: '/contact',
    locale: 'zh',
    name: '推广专员',
    email: 'seo-boost@example.net',
    phone: '400-000-0000',
    company: '全球推广',
    message: '专业 SEO 优化，快速提升关键词排名，首月五折，欢迎联系。',
    data: {
      name: '推广专员',
      email: 'seo-boost@example.net',
      phone: '400-000-0000',
      company: '全球推广',
      message: '专业 SEO 优化，快速提升关键词排名，首月五折，欢迎联系。'
    },
    status: 'spam',
    ip: '198.51.100.203',
    userAgent: 'Mozilla/5.0 (compatible; ExampleBot/2.1; +http://example.net/bot)',
    createdAt: '2026-10-18T22:17:00+08:00',
    updatedAt: '2026-10-19T08:20:00+08:00'
  },
  {
    id: 'sub_1006',
    formKey: 'contact',
    pageName: 'Contact',
    pagePath: '/en/contact',
    locale: 'en',
    name: 'Daniel Wu',
    email: 'daniel.wu@example.com',
    phone: '+44 20 7946 0188',
    company: 'Harbour & Co.',
    message:
      'Could you share a rough timeline for migrating 40 existing pages? We need both Chinese and English versions and a staging site for review.',
    data: {
      name: 'Daniel Wu',
      email: 'daniel.wu@example.com',
      phone: '+44 20 7946 0188',
      company: 'Harbour & Co.',
      message:
        'Could you share a rough timeline for migrating 40 existing pages? We need both Chinese and English versions and a staging site for review.'
    },
    status: 'done',
    ip: '192.0.2.140',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Firefox/142.0',
    createdAt: '2026-10-16T18:44:00+08:00',
    updatedAt: '2026-10-17T10:30:00+08:00'
  },
  {
    id: 'sub_1007',
    formKey: 'contact',
    pageName: '联系我们',
    pagePath: '/contact',
    locale: 'zh',
    name: '陈默',
    email: 'chenmo@example.com',
    phone: '136-7788-4120',
    company: '青屿教育',
    message:
      '想了解媒体库的图片压缩规则和存储方式。我们每篇新闻大概 8 到 10 张图，担心后台上传太慢。',
    data: {
      name: '陈默',
      email: 'chenmo@example.com',
      phone: '136-7788-4120',
      company: '青屿教育',
      message:
        '想了解媒体库的图片压缩规则和存储方式。我们每篇新闻大概 8 到 10 张图，担心后台上传太慢。'
    },
    status: 'processing',
    ip: '203.0.113.156',
    userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/141.0',
    createdAt: '2026-10-19T08:05:00+08:00',
    updatedAt: '2026-10-19T09:18:00+08:00'
  },
  {
    id: 'sub_1008',
    formKey: 'contact',
    pageName: 'Contact',
    pagePath: '/en/contact',
    locale: 'en',
    name: 'Olivia Park',
    email: 'olivia.park@example.com',
    phone: '+82 2 555 0147',
    company: 'Lumen Labs',
    message:
      'Do you support custom domains and SSO for the admin console? We are evaluating vendors for a launch in January.',
    data: {
      name: 'Olivia Park',
      email: 'olivia.park@example.com',
      phone: '+82 2 555 0147',
      company: 'Lumen Labs',
      message:
        'Do you support custom domains and SSO for the admin console? We are evaluating vendors for a launch in January.'
    },
    status: 'new',
    ip: '192.0.2.201',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/140.0',
    createdAt: '2026-10-19T11:53:00+08:00',
    updatedAt: '2026-10-19T11:53:00+08:00'
  }
]
