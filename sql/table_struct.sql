CREATE TABLE IF NOT EXISTS `donations` (
  `donate_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL,
  `ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `current_donations` int(10) unsigned NOT NULL DEFAULT '0',
  `prev_donations` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`donate_id`,`current_donations`,`ts`,`member_id`),
  KEY `idx_ts` (`ts`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=6445 ;


CREATE TABLE IF NOT EXISTS `members` (
  `member_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `tag` varchar(45) NOT NULL,
  `name` varchar(60) NOT NULL,
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `active` int(11) NOT NULL,
  `exited` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `alert_sent_donations` int(11) DEFAULT '0',
  `current_donations` int(10) unsigned DEFAULT '0',
  `last_donation_time` timestamp NULL DEFAULT NULL,
  `prev_donations` int(10) unsigned DEFAULT NULL,
  `alerted_discord` int(10) unsigned DEFAULT '0',
  `current_rec` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`member_id`),
  UNIQUE KEY `idx_tag_name` (`tag`,`name`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=451 ;


CREATE TABLE IF NOT EXISTS `receive` (
  `receive_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` int(10) unsigned NOT NULL,
  `ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `current` int(11) NOT NULL,
  `prev` int(11) NOT NULL,
  PRIMARY KEY (`receive_id`),
  KEY `idx_ts` (`ts`),
  KEY `idx_member` (`member_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=6510 ;

