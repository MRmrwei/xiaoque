-- 小鹊云：页面主表（站点编排 → 页面内容）
-- 对应 docs/MySQL数据库表设计.md §3.6、docs/初版上线方案.md §5.2

CREATE TABLE `xqyun_pages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `translation_group_id` CHAR(26) NOT NULL COMMENT '翻译组 ID；中英文记录共用',
  `locale` VARCHAR(10) NOT NULL COMMENT '语言：zh-CN / en',
  `path` VARCHAR(255) NOT NULL COMMENT '页面路径，如 /about、/en/about-us',
  `name` VARCHAR(200) NOT NULL COMMENT '后台显示名称',
  `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT 'draft 未发布，published 已发布',
  `draft_content` JSON NOT NULL COMMENT '当前草稿：SEO + 区块列表',
  `published_content` JSON NULL COMMENT '最近一次发布的线上内容；未发布为空',
  `version` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发布序号，每次发布 +1',
  `created_by` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '最后修改人',
  `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（Unix 秒）',
  `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（Unix 秒）',
  `deleted_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '软删除时间（Unix 秒），0 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_locale_path` (`locale`, `path`),
  KEY `idx_group` (`translation_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='页面主表';
