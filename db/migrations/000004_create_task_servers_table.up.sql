CREATE TABLE `task_servers` (
  `task_id` bigint unsigned NOT NULL,
  `server_id` bigint unsigned NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`task_id`, `server_id`),
  KEY `fk_task_servers_server` (`server_id`),
  CONSTRAINT `fk_task_servers_server` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`),
  CONSTRAINT `fk_task_servers_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE = InnoDB CHARSET = utf8 COLLATE = utf8_unicode_ci;