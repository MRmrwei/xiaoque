-- 小鹊云：表单留言（后台「线索 → 表单留言」）
-- 对应 docs/MySQL数据库表设计.md §3.13、docs/初版上线方案.md §5.2

CREATE TABLE `xqyun_form_submissions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `form_key` VARCHAR(50) NOT NULL COMMENT '表单标识，如 contact',
  `locale` VARCHAR(10) NULL DEFAULT NULL COMMENT '提交时的语言：zh-CN / en',
  `name` VARCHAR(100) NULL DEFAULT NULL COMMENT '姓名',
  `email` VARCHAR(255) NULL DEFAULT NULL COMMENT '邮箱',
  `phone` VARCHAR(50) NULL DEFAULT NULL COMMENT '电话',
  `company` VARCHAR(200) NULL DEFAULT NULL COMMENT '公司',
  `message` TEXT NULL COMMENT '留言正文',
  `data` JSON NOT NULL COMMENT '完整表单数据，兼容后续新增字段',
  `status` VARCHAR(20) NOT NULL DEFAULT 'new' COMMENT 'new / processing / done / spam',
  `ip` VARCHAR(45) NULL DEFAULT NULL COMMENT '提交 IP',
  `user_agent` VARCHAR(500) NULL DEFAULT NULL COMMENT '浏览器信息',
  `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '提交时间（Unix 秒）',
  `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '处理时间（Unix 秒）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='表单留言';
