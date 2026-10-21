-- 小鹊云：全站设置（站点编排 → 站点配置 / 头部 / 页脚）
-- 按 key + locale 存储，值统一用 JSON；locale 用空字符串表示全局配置。
-- 对应 docs/MySQL数据库表设计.md §3.12、docs/初版上线方案.md §5.2

CREATE TABLE `xqyun_settings` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `key` VARCHAR(100) NOT NULL COMMENT '设置项，如 brand、header、footer',
  `locale` VARCHAR(10) NOT NULL DEFAULT '' COMMENT '语言；空字符串表示全局配置',
  `value` JSON NOT NULL COMMENT '设置值',
  `updated_by` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '最后修改人',
  `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（Unix 秒）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key_locale` (`key`, `locale`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='全站设置';
