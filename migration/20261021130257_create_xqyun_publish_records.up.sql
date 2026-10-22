-- 小鹊云：页面 / 内容发布记录（站点编排 → 发布与下架）
-- 对应 docs/MySQL数据库表设计.md §3.16、docs/初版上线方案.md §5.2

CREATE TABLE `xqyun_publish_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `target_type` VARCHAR(20) NOT NULL COMMENT '目标类型：page / content',
  `target_id` BIGINT UNSIGNED NOT NULL COMMENT '对象 ID',
  `locale` VARCHAR(10) NOT NULL COMMENT '发布语言',
  `status` VARCHAR(20) NOT NULL DEFAULT 'success' COMMENT 'success / failed',
  `message` VARCHAR(500) NULL DEFAULT NULL COMMENT '失败原因或备注',
  `operator_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '操作人',
  `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发布时间（Unix 秒）',
  PRIMARY KEY (`id`),
  KEY `idx_target_created` (`target_type`, `target_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发布记录';
