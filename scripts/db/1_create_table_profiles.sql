CREATE TABLE `profiles` (
  `id`     varchar(64) NOT NULL,
  `handle` varchar(32) NOT NULL,
  `avatar` varchar(256),
  UNIQUE(`handle`),
  PRIMARY KEY (`id`)
);