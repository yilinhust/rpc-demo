CREATE DATABASE TEST;

USE TEST;

DROP TABLE  IF EXISTS `t_user`;

CREATE TABLE `t_user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NULL NOT NULL,
    `password` VARCHAR(200) NULL NOT NULL,
    `nickname` VARCHAR(64) NULL DEFAULT '',
    `profile` VARCHAR(200) NULL DEFAULT '',
    `CreateTime` DATETIME NULL DEFAULT NOW(),
    `LoginTime` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
);

ALTER TABLE t_user ADD INDEX idx_user_username (username);

INSERT t_user SET username='test',password='r7HDvaNzbzG_MpB1R3ogVtCqDj-AAI8vQwkvCVqNb2s',nickname='hust',CreateTime=NOW();
