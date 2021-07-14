-- MySQL dump 10.19  Distrib 10.3.29-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: 3cxReporting
-- ------------------------------------------------------
-- Server version       10.3.29-MariaDB-0+deb10u1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `config`
--

DROP TABLE IF EXISTS `config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `PasswordStateURL` varchar(255) NOT NULL,
  `PasswordStateAPIKey` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `hosts`
--

DROP TABLE IF EXISTS `hosts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hosts` (
  `name` varchar(255) DEFAULT NULL,
  `fqdn` varchar(255) DEFAULT NULL,
  `main_ip` varchar(15) DEFAULT NULL,
  `os_name` varchar(80) DEFAULT NULL,
  `os_version` varchar(40) DEFAULT NULL,
  `system` varchar(40) DEFAULT NULL,
  `kernel` varchar(40) DEFAULT NULL,
  `arch_hardware` varchar(12) DEFAULT NULL,
  `arch_userspace` varchar(12) DEFAULT NULL,
  `virt_type` varchar(20) DEFAULT NULL,
  `virt_role` varchar(20) DEFAULT NULL,
  `cpu_type` varchar(60) DEFAULT NULL,
  `vcpus` int(11) DEFAULT NULL,
  `ram` float DEFAULT NULL,
  `disk_total` float DEFAULT NULL,
  `disk_free` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `servers`
--

DROP TABLE IF EXISTS `servers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `servers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) NOT NULL,
  `Location` varchar(255) NOT NULL,
  `URL` varchar(255) NOT NULL,
  `UserName` varchar(255) NOT NULL,
  `PwStateID` int(11) NOT NULL,
  `CustomerID` int(11) DEFAULT NULL,
  `Version` varchar(255) DEFAULT NULL,
  `FQDN` varchar(255) DEFAULT NULL,
  `CallRecordingUsage` varchar(255) DEFAULT NULL,
  `RecordingUsedSpace` bigint(20) DEFAULT NULL,
  `RecordingQuota` bigint(20) DEFAULT NULL,
  `RecordingQuotaSold` bigint(20) DEFAULT NULL,
  `MaxSimCalls` int(11) DEFAULT NULL,
  `ExtTotal` int(11) DEFAULT NULL,
  `vcpus` varchar(255) DEFAULT NULL,
  `OSram` int(11) DEFAULT NULL,
  `OSswap` int(11) DEFAULT NULL,
  `OSDiskSpace` int(11) DEFAULT NULL,
  `AutoUpdate` tinyint(1) DEFAULT NULL,
  `License` varchar(255) DEFAULT NULL,
  `LicenseKey` varchar(255) DEFAULT NULL,
  `LicenseExpiration` varchar(255) DEFAULT NULL,
  `ResellerName` varchar(255) DEFAULT NULL,
  `bespoke` varchar(255) DEFAULT NULL,
  `AnsibleUpdates` tinyint(1) DEFAULT 1,
  `lastUpdated` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-07-14 10:02:25
