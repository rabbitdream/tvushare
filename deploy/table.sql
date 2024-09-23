CREATE DATABASE tvushare;

drop table tvu_social_media_account;

drop table oauth_token;

CREATE TABLE `tvu_social_media_account` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `tvu_account` VARCHAR(1024) DEFAULT '' COMMENT 'tvu account',
    `social_media_user_id` VARCHAR(128) DEFAULT '' COMMENT 'social media user id',
    `social_media_user_name` VARCHAR(128) DEFAULT '' COMMENT 'social media user name',
    `social_media_avatar_url` VARCHAR(128) DEFAULT '' COMMENT 'social media avatar url',
    `social_media_platform` VARCHAR(32) DEFAULT '' COMMENT 'social media platform',
    `social_media_token_id` INT(10) DEFAULT 0 COMMENT 'social media oauth2.0 token id',
    `created_time` INT(11) NOT NULL COMMENT 'created_time',
    `updated_time` INT(11) NOT NULL COMMENT 'updated_time',
    PRIMARY KEY (`id`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = 'tvu account to social media platform token and account';

CREATE TABLE `oauth_token` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `access_token` VARCHAR(256) NOT NULL COMMENT 'the token that authorizes and authenticates the requests.',
    `token_secret` VARCHAR(256) DEFAULT '' COMMENT 'the token secret of oauth1',
    `token_type` VARCHAR(64) DEFAULT '' COMMENT 'the type of token.',
    `refresh_token` VARCHAR(1024) DEFAULT '' COMMENT 'a token that is used by the application  (as opposed to the user) to refresh the access token if it expires.',
    `expires_in` INT(10)  DEFAULT 0 COMMENT 'how many seconds later the token expires',
    `scope` VARCHAR(1024) DEFAULT '' COMMENT 'the token scope',
    `created_time` INT(11) NOT NULL COMMENT 'created_time',
    `updated_time` INT(11) NOT NULL COMMENT 'updated_time',
    PRIMARY KEY (`id`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = 'social medias platform oauth2.0 token';


SELECT tsma.* FROM tvu_social_media_account AS tsma INNER JOIN oauth_token AS ot ON tsma.social_media_token_id = ot.id WHERE tsma.tvu_account = 'calebpan@tvunetworks.com' AND (ot.updated_time + ot.expires_in > 1726816171 OR ot.expires_in = 0)
