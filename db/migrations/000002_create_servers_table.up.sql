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
) ENGINE = InnoDB CHARSET = utf8 COLLATE = utf8_unicode_ci;