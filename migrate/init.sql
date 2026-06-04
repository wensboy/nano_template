-- @author: wendisx
-- @date: 2026-06-03 pm
-- @description: stable migrate definations

SET NAMES utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS `test` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `test`;

-- user
DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
    `id`         INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `created_at` DATETIME(3)  NULL                       COMMENT '创建时间',
    `updated_at` DATETIME(3)  NULL                       COMMENT '更新时间',
    `deleted_at` DATETIME(3)  NULL                       COMMENT '逻辑删除时间',
    `username`   VARCHAR(255) NOT NULL                   COMMENT '用户名, 逻辑唯一',
    `password`   VARCHAR(255) NOT NULL                   COMMENT '用户密码'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- user profile
DROP TABLE IF EXISTS `user_profiles`;
CREATE TABLE IF NOT EXISTS `user_profiles` (
    `id`         INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `created_at` DATETIME(3)  NULL                       COMMENT '创建时间',
    `updated_at` DATETIME(3)  NULL                       COMMENT '更新时间',
    `deleted_at` DATETIME(3)  NULL                       COMMENT '逻辑删除时间',
    `user_id`    INT UNSIGNED NOT NULL                   COMMENT '关联: 用户ID',
    `avatar`     VARCHAR(255)                            COMMENT '头像地址',
    `nickname`   VARCHAR(255)                            COMMENT '昵称',
    `email`      VARCHAR(255)                            COMMENT '邮箱',
    `phone`      VARCHAR(255)                            COMMENT '手机号',
    `signature`  TEXT                                    COMMENT '个性签名'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- role
DROP TABLE IF EXISTS `roles`;
CREATE TABLE IF NOT EXISTS `roles` (
    `id`         INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `created_at` DATETIME(3)  NULL                       COMMENT '创建时间',
    `updated_at` DATETIME(3)  NULL                       COMMENT '更新时间',
    `deleted_at` DATETIME(3)  NULL                       COMMENT '逻辑删除时间',
    `name`       VARCHAR(191) NOT NULL                   COMMENT '角色名称',
    `level`      INT UNSIGNED NOT NULL                   COMMENT '角色级别',
    `description` VARCHAR(191) NULL                      COMMENT '角色描述',
    `state`      INT UNSIGNED NOT NULL                   COMMENT '状态 - 0代表未启用, >=1代表其他启用状态'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- user role
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE IF NOT EXISTS `user_roles` (
    `id`         INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `created_at` DATETIME(3)  NULL                       COMMENT '创建时间',
    `updated_at` DATETIME(3)  NULL                       COMMENT '更新时间',
    `deleted_at` DATETIME(3)  NULL                       COMMENT '逻辑删除时间',
    `user_id`    INT UNSIGNED NOT NULL                   COMMENT '关联用户ID',
    `role_id`    INT UNSIGNED NOT NULL                   COMMENT '关联角色ID',
    `description` VARCHAR(191) NULL                      COMMENT '关联描述 -> 备注',
    `state`      INT UNSIGNED NOT NULL                   COMMENT '状态 - 0代表未启用, >=1代表其他启用状态'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;