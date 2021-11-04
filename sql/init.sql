CREATE DATABASE IF NOT EXISTS imloader;

-- create the users for each database
CREATE USER IF NOT EXISTS 'imloader'@'%' IDENTIFIED BY 'password';
GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON `imloader`.* TO 'imloader'@'%';

FLUSH PRIVILEGES;

USE imloader;


CREATE TABLE IF NOT EXISTS `images` (
    `token` VARCHAR(255) NOT NULL,
    `data` LONGBLOB NOT NULL,
    PRIMARY KEY (`token`)
) ENGINE = InnoDB CHARSET=utf8 COLLATE utf8_general_ci;
