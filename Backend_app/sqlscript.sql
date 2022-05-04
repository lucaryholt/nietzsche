SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema distributedsystems
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS systemsintegration;

USE `systemsintegration`;

DROP TABLE IF EXISTS `verifications`;
CREATE TABLE IF NOT EXISTS `verifications` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `cpr` VARCHAR(100) NOT NULL,
  `iat` VARCHAR(100) NULL DEFAULT NULL,
  `exp` VARCHAR(100) NULL DEFAULT NULL,
  `code` int(5) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `games_id_UNIQUE` (`id` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;