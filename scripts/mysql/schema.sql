SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL';
SHOW WARNINGS;

DROP SCHEMA IF EXISTS `hasd_covid` ;
CREATE SCHEMA IF NOT EXISTS `hasd_covid` ;
USE `hasd_covid` ;

DROP TABLE IF EXISTS `hasd_covid`.`category` ;

CREATE  TABLE IF NOT EXISTS `hasd_covid`.`category` (
  `category_skey` INT NOT NULL AUTO_INCREMENT ,
  `description` VARCHAR(100) NOT NULL ,
  PRIMARY KEY (`category_skey`) ,
  UNIQUE INDEX `category_skey_UNIQUE` (`category_skey` ASC) )
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;
SHOW WARNINGS;

DROP TABLE IF EXISTS `hasd_covid`.`school` ;
CREATE  TABLE IF NOT EXISTS `hasd_covid`.`school` (
  `school_skey` INT NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(100) NOT NULL ,
  PRIMARY KEY (`school_skey`),
  UNIQUE INDEX `school_skey_UNIQUE` (`school_skey` ASC) )
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;
SHOW WARNINGS;

DROP TABLE IF EXISTS `hasd_covid`.`metric` ;

CREATE  TABLE IF NOT EXISTS `hasd_covid`.`metric` (
    `metric_skey` INT NOT NULL AUTO_INCREMENT,
    `category_skey` INT NOT NULL ,
    `school_skey` INT NOT NULL ,
    `active_cases` int NOT NULL default 0,
    `total_positive_cases` int NOT NULL default 0,
    `total_probable_cases` int NOT NULL default 0,
    `resolved` int NOT NULL default 0,
  PRIMARY KEY (`metric_skey`) ,
  UNIQUE INDEX `metric_metric_skey_UNIQUE` (`metric_skey` ASC) ,
  KEY `metric_category_skey` (`category_skey`),
  KEY `metric_school_skey` (`school_skey`),
  CONSTRAINT `metric_category_skey` FOREIGN KEY (`category_skey`) REFERENCES `category` (`category_skey`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `metric_school_skey` FOREIGN KEY (`school_skey`) REFERENCES `school` (`school_skey`) ON DELETE NO ACTION ON UPDATE NO ACTION
)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;

SHOW WARNINGS;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
