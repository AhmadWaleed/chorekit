/* *****************************************************************************
 // Setup the preferences
 // ****************************************************************************/
SET
  NAMES utf8 COLLATE 'utf8_unicode_ci';

SET
  foreign_key_checks = 1;

SET
  time_zone = '+00:00';

SET
  sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET
  default_storage_engine = InnoDB;

SET
  CHARACTER SET utf8;

/* *****************************************************************************
 // Remove old database
 // ****************************************************************************/
DROP DATABASE IF EXISTS chorekit;

/* *****************************************************************************
 // Create new database
 // ****************************************************************************/
CREATE DATABASE chorekit DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

USE chorekit;

/* *****************************************************************************
 // Create the tables
 // ****************************************************************************/
CREATE TABLE `users` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(191),
  `email` VARCHAR(191) DEFAULT NULL,
  `password` VARCHAR(191),
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE = InnoDB AUTO_INCREMENT = 2, CHARSET = utf8 COLLATE = utf8_unicode_ci;

CREATE TABLE `servers` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(191),
  `ip` VARCHAR(191),
  `user` VARCHAR(191),
  `port` bigint DEFAULT NULL,
  `ssh_public_key` longtext,
  `ssh_private_key` longtext,
  `status` VARCHAR(50),
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 3;

CREATE TABLE `tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(191),
  `env` VARCHAR(191),
  `script` longtext,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 2;

CREATE TABLE `task_servers` (
  `task_id` bigint unsigned NOT NULL,
  `server_id` bigint unsigned NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`task_id`, `server_id`),
  KEY `fk_task_servers_server` (`server_id`),
  CONSTRAINT `fk_task_servers_server` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`),
  CONSTRAINT `fk_task_servers_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE = InnoDB;

CREATE TABLE `buckets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `parallel` tinyint(1) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB;

CREATE TABLE `bucket_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_id` bigint unsigned DEFAULT NULL,
  `bucket_id` bigint unsigned DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_bucket_tasks_task` (`task_id`),
  CONSTRAINT `fk_bucket_tasks_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE = InnoDB;

CREATE TABLE `task_runs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_id` bigint unsigned DEFAULT NULL,
  `output` longtext,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_tasks_runs_run` (`task_id`),
  CONSTRAINT `fk_tasks_runs_run` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 9;

CREATE TABLE `bucket_runs` (
  `bucket_id` bigint unsigned NOT NULL,
  `task_run_id` bigint unsigned NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`bucket_id`, `task_run_id`),
  KEY `fk_bucket_runs_run` (`task_run_id`),
  CONSTRAINT `fk_bucket_runs_bucket` FOREIGN KEY (`bucket_id`) REFERENCES `buckets` (`id`),
  CONSTRAINT `fk_bucket_runs_task_run` FOREIGN KEY (`task_run_id`) REFERENCES `task_runs` (`id`)
) ENGINE = InnoDB;
