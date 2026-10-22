-- 小鹊云业务数据库（与 service/config.yaml 中的 database 保持一致）
-- 执行方式：mysql -u root -p < deploy/init.sql

CREATE DATABASE IF NOT EXISTS `xqyun`
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE `xqyun`;

-- 业务表由 GORM 迁移或后续 schema 脚本创建。

-- 可选：创建独立应用账号，避免应用直接使用 root。
-- CREATE USER IF NOT EXISTS 'xqyun'@'%' IDENTIFIED BY 'change-me';
-- GRANT ALL PRIVILEGES ON `xqyun`.* TO 'xqyun'@'%';
-- FLUSH PRIVILEGES;
