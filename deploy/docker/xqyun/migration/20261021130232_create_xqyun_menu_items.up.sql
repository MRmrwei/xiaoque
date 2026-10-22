-- 小鹊云：菜单项（站点编排 → 顶部 / 底部导航）
-- 同一张表承载顶部菜单项和页脚栏目 / 链接。
-- 对应 docs/MySQL数据库表设计.md §3.11、docs/初版上线方案.md §5.2

CREATE TABLE `xqyun_menu_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `menu_id` BIGINT UNSIGNED NOT NULL COMMENT '所属菜单',
  `parent_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '父项；页脚中顶级项是栏目，子项是链接',
  `label` JSON NOT NULL COMMENT '多语言文案：{"zh-CN":"...","en":"..."}',
  `link_type` VARCHAR(20) NOT NULL DEFAULT 'page' COMMENT 'page / content / external / none',
  `link_value` VARCHAR(500) NULL DEFAULT NULL COMMENT '页面路径、内容 slug 或外链',
  `target` VARCHAR(20) NOT NULL DEFAULT '_self' COMMENT '打开方式：_self / _blank',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '同级排序，越小越靠前',
  `visible` TINYINT NOT NULL DEFAULT 1 COMMENT '1 显示，0 隐藏',
  `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（Unix 秒）',
  `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（Unix 秒）',
  PRIMARY KEY (`id`),
  KEY `idx_menu_parent_sort` (`menu_id`, `parent_id`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单项与页脚链接';
