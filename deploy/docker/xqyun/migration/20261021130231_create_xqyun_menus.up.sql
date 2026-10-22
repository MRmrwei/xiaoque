-- 小鹊云：菜单容器（站点编排 → 顶部 / 底部导航）
-- 对应 docs/MySQL数据库表设计.md §3.10、docs/初版上线方案.md §5.2

CREATE TABLE `xqyun_menus` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` VARCHAR(50) NOT NULL COMMENT '菜单编码：header / footer / footer_bottom',
  `name` VARCHAR(100) NOT NULL COMMENT '后台显示名称',
  `description` VARCHAR(255) NULL DEFAULT NULL COMMENT '说明',
  `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（Unix 秒）',
  `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（Unix 秒）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单容器';
