# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.5.5-10.2.14-MariaDB)
# Database: foreign_currency
# Generation Time: 2018-09-12 21:09:07 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table exchange_rates
# ------------------------------------------------------------

DROP TABLE IF EXISTS `exchange_rates`;

CREATE TABLE `exchange_rates` (
  `oid` varchar(50) NOT NULL DEFAULT '',
  `date` date NOT NULL,
  `from` char(3) DEFAULT NULL,
  `to` char(3) DEFAULT NULL,
  `rate` float DEFAULT NULL,
  PRIMARY KEY (`oid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `exchange_rates` WRITE;
/*!40000 ALTER TABLE `exchange_rates` DISABLE KEYS */;

INSERT INTO `exchange_rates` (`oid`, `date`, `from`, `to`, `rate`)
VALUES
	('13213451435','2018-07-06','USD','GBP',0.709),
	('512342345135daf32rfsf432r','2018-07-05','USD','GBP',NULL),
	('5b95740aa059a33537225789','2018-07-08','USD','GBP',0.123),
	('5b95740aa059a33537225790','2018-07-01','IDR','USD',0),
	('5b95740aa059a335372257e1','2018-07-08','IDR','USD',14322),
	('5b95740aa059a335372257e2','2018-07-07','IDR','USD',14347),
	('5b95740aa059a335372257e3','2018-07-04','IDR','USD',14129),
	('5b95740aa059a335372257e4','2018-07-02','IDR','USD',14091),
	('5b95740aa059a335372257e5','2018-07-03','IDR','USD',14222),
	('5b957435a059a335fe0f4912','2018-07-07','USD','GBP',0.7609),
	('5b957435a059a335fe0f4970','2018-07-06','IDR','USD',13233),
	('5b957435a059a335fe0f49800','2018-07-08','JPY','IDR',NULL),
	('5b957435a059a335fe0f49e9','2018-07-05','IDR','USD',14123);

/*!40000 ALTER TABLE `exchange_rates` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tracks
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tracks`;

CREATE TABLE `tracks` (
  `from` char(3) NOT NULL DEFAULT '',
  `to` char(3) NOT NULL DEFAULT '',
  `uid` char(6) NOT NULL DEFAULT 'NUL',
  PRIMARY KEY (`uid`),
  KEY `from` (`from`),
  KEY `to` (`to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `tracks` WRITE;
/*!40000 ALTER TABLE `tracks` DISABLE KEYS */;

INSERT INTO `tracks` (`from`, `to`, `uid`)
VALUES
	('IDR','USD','idrusd'),
	('JPY','IDR','jpyidr'),
	('USD','GBP','usdgbp');

/*!40000 ALTER TABLE `tracks` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
