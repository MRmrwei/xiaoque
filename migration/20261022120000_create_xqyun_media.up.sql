-- 媒体素材表：当前用于后台视频素材库，后续可扩展到图片 / 文件。
CREATE TABLE `xqyun_media` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `kind` VARCHAR(20) NOT NULL DEFAULT 'video' COMMENT '素材类型：video / image / file',
  `original_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
  `filename` VARCHAR(255) NOT NULL COMMENT '存储文件名',
  `url` VARCHAR(500) NOT NULL COMMENT '访问地址',
  `mime` VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'MIME 类型',
  `size` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
  `created_by` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '上传人',
  `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '上传时间（Unix 秒）',
  `deleted_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '软删除时间（Unix 秒），0 表示未删除',
  PRIMARY KEY (`id`),
  KEY `idx_kind_created` (`kind`, `created_at`),
  KEY `idx_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='媒体素材';
