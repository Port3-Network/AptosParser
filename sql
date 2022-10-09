DROP DATABASE IF EXISTS aptos_sync;
CREATE DATABASE IF NOT EXISTS aptos_sync;
USE aptos_sync;
SET sql_mode="NO_ENGINE_SUBSTITUTION";

-- NOTE: aptos address length -> 66
-- ----------------------------
-- Table sysconfig -> system config
-- ----------------------------
DROP TABLE IF EXISTS `sysconfig`;

CREATE TABLE `sysconfig` (
    `id` int NOT NULL AUTO_INCREMENT,
    `cfg_name` varchar(64) NOT NULL COMMENT 'config name',
    `cfg_val` varchar(64) NOT NULL COMMENT 'config value',
    `cfg_type` varchar(16) NOT NULL COMMENT 'config type',
    `cfg_comment` varchar(128) NOT NULL DEFAULT '' COMMENT 'config comment',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `cfg_name` (`cfg_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table block -> block detail
-- ----------------------------
DROP TABLE IF EXISTS `block`;

CREATE TABLE `block` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `height` bigint NOT NULL COMMENT 'block height',
    `hash` char(66) NOT NULL COMMENT 'block hash',
    `block_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `first_version` char(64) NOT NULL DEFAULT '0' COMMENT 'first version of current block',
    `last_version` char(64) NOT NULL DEFAULT '0' COMMENT 'last version of current block',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `height` (`height`),
    KEY `hash` (`hash`),
    KEY `block_time` (`block_time`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table transaction -> all tx record
-- ----------------------------
DROP TABLE IF EXISTS `transaction`;

CREATE TABLE `transaction` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `success` tinyint NOT NULL DEFAULT 0 COMMENT 'vm state',
    `sequence_number` bigint NOT NULL DEFAULT 0 COMMENT 'sequence number',
    `gas_used` varchar(32) NOT NULL COMMENT 'gas used',
    `gas_price` varchar(32) NOT NULL COMMENT 'gas price',
    `gas` varchar(32) NOT NULL COMMENT 'gas',
    `type` char(24) NOT NULL DEFAULT '' COMMENT 'tx type',
    `sender` char(66) NOT NULL COMMENT 'from',
    `receiver` char(66) NOT NULL COMMENT 'to',
    `tx_value` varchar(42) NOT NULL DEFAULT '0' COMMENT 'value',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `hash` (`hash`),
    KEY `version` (`version`),
    KEY `sender` (`sender`),
    KEY `receiver` (`receiver`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table payload -> tx log, which function to call
-- ----------------------------
DROP TABLE IF EXISTS `payload`;

CREATE TABLE `payload` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `sequence_number` int NOT NULL COMMENT 'sequence_number',
    `sender` char(66) NOT NULL COMMENT 'tx sender',
    `payload_func` char(128) NOT NULL COMMENT 'call function, payload function',
    `payload_type` char(128) NOT NULL COMMENT 'call type, payload type',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table record_coin -> publish pkg record
-- ----------------------------
DROP TABLE IF EXISTS `record_coin`;

-- resource = sender::module_name::contract_name
CREATE TABLE `record_coin` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `sender` char(66) NOT NULL DEFAULT '' COMMENT 'tx sender',
    `module_name` char(128) NOT NULL DEFAULT '' COMMENT '',
    `contract_name` char(128) NOT NULL DEFAULT '' COMMENT '',
    `resource` char(128) NOT NULL DEFAULT '' COMMENT 'resource name',
    `name` text NOT NULL DEFAULT '' COMMENT 'contract name',
    `symbol` text NOT NULL DEFAULT '' COMMENT 'contract symbol',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `resource` (`resource`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table history_coin -> coin transfer histories
-- ----------------------------
DROP TABLE IF EXISTS `history_coin`;

CREATE TABLE `history_coin` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `sender` char(66) NOT NULL COMMENT 'tx sender',
    `receiver` char(66) NOT NULL COMMENT 'tx receiver',
    `resource` char(128) NOT NULL COMMENT 'coin resource',
    `amount` varchar(128) NOT NULL COMMENT 'tx amount',
    `action` tinyint NOT NULL COMMENT '0: unknow, 1: mint, 2: transfer, 3:burn',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `hash` (`hash`),
    KEY `version` (`version`),
    KEY `index_sender_receiver` (`sender`, `receiver`),
    KEY `index_time` (`tx_time`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table collection -> collection
-- ----------------------------
DROP TABLE IF EXISTS `collection`;

-- resource = sender::module_name::contract_name
CREATE TABLE `collection` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `sender` char(66) NOT NULL DEFAULT '' COMMENT 'tx sender',
    `creator` char(66) NOT NULL DEFAULT '' COMMENT 'collection owner',
    `name` char(66) NOT NULL DEFAULT '' COMMENT 'collection name',
    `description` text NOT NULL DEFAULT '' COMMENT 'collection description',
    `uri` text NOT NULL DEFAULT '' COMMENT 'collection uri',
    `maximum` char(128) NOT NULL DEFAULT '' COMMENT 'collection maximum',
    `type` char(128) NOT NULL DEFAULT '' COMMENT 'collection type',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table record_token -> publish pkg record
-- ----------------------------
DROP TABLE IF EXISTS `record_token`;

-- resource = sender::module_name::contract_name
CREATE TABLE `record_token` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `sender` char(66) NOT NULL DEFAULT '' COMMENT 'tx sender',
    `creator` char(66) NOT NULL DEFAULT '' COMMENT 'token creator',
    `collection` char(255) NOT NULL DEFAULT '' COMMENT 'collection',
    `name` char(255) NOT NULL DEFAULT '' COMMENT 'token name',
    `description` text NOT NULL DEFAULT '' COMMENT 'token description',
    `uri` text NOT NULL DEFAULT '' COMMENT 'token uri',
    `maximum` char(128) NOT NULL DEFAULT '' COMMENT 'collection maximum',
    `type` char(128) NOT NULL DEFAULT '' COMMENT 'collection type',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `token_data` (`creator`, `collection`, `name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table asset_token -> owner of each nft
-- ----------------------------
DROP TABLE IF EXISTS `asset_token`;

-- resource = sender::module_name::contract_name
CREATE TABLE `asset_token` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `owner` char(66) NOT NULL DEFAULT '' COMMENT 'owner',
    `creator` char(66) NOT NULL DEFAULT '' COMMENT 'token creator',
    `collection` char(255) NOT NULL DEFAULT '' COMMENT 'collection',
    `name` char(255) NOT NULL DEFAULT '' COMMENT 'token name',
    `amount` char(66) NOT NULL DEFAULT '' COMMENT 'nft amount',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `owner_token_data` (`owner`, `creator`, `collection`, `name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;

-- ----------------------------
-- Table history_token -> token transfer histories
-- ----------------------------
DROP TABLE IF EXISTS `history_token`;

CREATE TABLE `history_token` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'auto increment',
    `version` char(32) NOT NULL COMMENT 'tx version',
    `hash` char(66) NOT NULL COMMENT 'tx hash',
    `tx_time` bigint NOT NULL DEFAULT 0 COMMENT 'block timestamp',
    `sender` char(66) NOT NULL COMMENT 'tx sender',
    `receiver` char(66) NOT NULL COMMENT 'tx receiver',
    `creator` char(66) NOT NULL DEFAULT '' COMMENT 'token creator',
    `collection` char(255) NOT NULL DEFAULT '' COMMENT 'collection',
    `name` char(255) NOT NULL DEFAULT '' COMMENT 'token name',
    `amount` varchar(66) NOT NULL COMMENT 'nft amount',
    `action` tinyint NOT NULL COMMENT '0: unknow, 1: mint, 2: transfer, 3:burn',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `hash` (`hash`),
    KEY `index_sender_receiver` (`sender`, `receiver`),
    KEY `index_time` (`tx_time`),
    KEY `token_data` (`creator`, `collection`, `name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 AUTO_INCREMENT = 1;
