-- 小鹊云：后台账号
-- 对应 docs/MySQL数据库表设计.md §3.1

CREATE TABLE `xqyun_users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` VARCHAR(64) NOT NULL COMMENT '登录名',
  `password_hash` VARCHAR(255) NOT NULL COMMENT 'bcrypt 哈希',
  `name` VARCHAR(100) NOT NULL COMMENT '显示名称',
  `email` VARCHAR(255) NULL DEFAULT NULL COMMENT '邮箱',
  `mobile` VARCHAR(32) NULL DEFAULT NULL COMMENT '手机号',
  `avatar` MEDIUMTEXT NULL COMMENT '头像：data URL 或外部地址',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1 启用，0 禁用',
  `last_login_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '最近登录时间（Unix 秒）',
  `last_login_ip` VARCHAR(45) NULL DEFAULT NULL COMMENT '最近登录 IP，兼容 IPv6',
  `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（Unix 秒）',
  `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（Unix 秒）',
  `deleted_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '软删除时间（Unix 秒），0 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台账号';

-- 初始超级管理员：账号 admin，显示名称「小阙」，初始密码 123456。
-- 密码通过 service/app/services/auth/password.go 的 HashPassword 生成（bcrypt）。
-- 首次登录后请立即修改密码。
INSERT INTO `xqyun_users` (`username`, `password_hash`, `name`, `status`, `created_at`, `updated_at`)
VALUES (
  'admin',
  '$2a$10$E2BArJj/J3BSGLY/iW82eenYIZxHlIYfqmK9UcjjS5F1ABRRGF/aS',
  '小阙',
  1,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP()
);
