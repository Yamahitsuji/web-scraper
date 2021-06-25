DROP DATABASE IF EXISTS `wikipedia`;
CREATE DATABASE `wikipedia` DEFAULT COLLATE utf8mb4_general_ci;
USE `wikipedia`;

drop table if exists `article`;
CREATE TABLE IF NOT EXISTS `article` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(190) NOT NULL,
    `url` varchar(190) NOT NULL,
    `latitude` float NOT NULL DEFAULT 0,
    `longitude` float NOT NULL DEFAULT 0,
    `read` text,
    `text` mediumtext NOT NULL,
    `detail_json` mediumtext NOT NULL,
    `created_at` DATETIME,
    `updated_at` DATETIME,
    PRIMARY KEY (`id`),
    KEY `idx_title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='記事一覧';
