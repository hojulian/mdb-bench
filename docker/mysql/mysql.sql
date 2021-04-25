CREATE DATABASE IF NOT EXISTS test;

USE test;

CREATE TABLE IF NOT EXISTS `voyages` (
    `number` varchar(191) NOT NULL,
    `schedule_id` bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `carrier_movements` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `departure_location` longtext,
    `arrival_location` longtext,
    `departure_time` datetime(3) DEFAULT NULL,
    `arrival_time` datetime(3) DEFAULT NULL,
    `schedule_refer` bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_carrier_movements_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `schedules` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_schedules_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `locations` (
    `un_locode` varchar(191) NOT NULL,
    `name` longtext,
    PRIMARY KEY (`un_locode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cargos` (
    `tracking_id` varchar(191) NOT NULL,
    `origin` longtext,
    `route_specification_id` bigint(20) unsigned DEFAULT NULL,
    `itinerary_id` bigint(20) unsigned DEFAULT NULL,
    `delivery_id` bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (`tracking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `handling_histories` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_handling_histories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `handling_events` (
    `tracking_id` varchar(191) NOT NULL,
    `activity_id` bigint(20) unsigned DEFAULT NULL,
    `handling_history_refer` bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (`tracking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `handling_activities` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `type` bigint(20) DEFAULT NULL,
    `location` longtext,
    `voyage_number` longtext,
    PRIMARY KEY (`id`),
    KEY `idx_handling_activities_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `route_specifications` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `origin` longtext,
    `destination` longtext,
    `arrival_deadline` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_route_specifications_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `itineraries` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_itineraries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `deliveries` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `itinerary_id` bigint(20) unsigned DEFAULT NULL,
    `route_specification_id` bigint(20) unsigned DEFAULT NULL,
    `routing_status` bigint(20) DEFAULT NULL,
    `transport_status` bigint(20) DEFAULT NULL,
    `next_expected_activity_id` bigint(20) unsigned DEFAULT NULL,
    `last_event_id` varchar(64) DEFAULT NULL,
    `last_known_location` longtext,
    `current_voyage` longtext,
    `eta` datetime(3) DEFAULT NULL,
    `is_misdirected` tinyint(1) DEFAULT NULL,
    `is_unloaded_at_destination` tinyint(1) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_deliveries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `legs` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `voyage_number` longtext,
    `load_location` longtext,
    `unload_location` longtext,
    `load_time` datetime(3) DEFAULT NULL,
    `unload_time` datetime(3) DEFAULT NULL,
    `itinerary_refer` bigint(20) unsigned DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_legs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
