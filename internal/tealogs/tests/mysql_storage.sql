CREATE TABLE `accessLogs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `log` text COLLATE utf8mb4_bin,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;