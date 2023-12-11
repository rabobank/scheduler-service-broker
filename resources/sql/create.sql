CREATE TABLE `schedulables`
(
    `guid`       varchar(36) not null,
    `created_at` timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`guid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;


CREATE TABLE `jobs`
(
    `guid`       varchar(36)   not null,
    `appguid`    varchar(36)   NOT NULL,
    `spaceguid`  varchar(36)   NOT NULL,
    `state`      varchar(32)   NOT NULL DEFAULT 'CREATED',
    `name`       varchar(255)  NOT NULL,
    `command`    varchar(4096) NOT NULL,
    `memoryinmb` int(1)                 DEFAULT '0',
    `diskinmb`   int(1)                 DEFAULT '0',
    PRIMARY KEY (`guid`),
    UNIQUE KEY `name_in_space` (`name`, `spaceguid`),
    CONSTRAINT `fk_jobs_are_schedulables` FOREIGN KEY (`guid`) REFERENCES `schedulables` (`guid`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;


CREATE TABLE `calls`
(
    `guid`       varchar(36)  not null,
    `appguid`    varchar(36)  NOT NULL,
    `spaceguid`  varchar(36)  NOT NULL,
    `state`      varchar(255) NOT NULL DEFAULT 'CREATED',
    `name`       varchar(255) NOT NULL,
    `url`        text         NOT NULL,
    `authheader` text,
    PRIMARY KEY (`guid`),
    UNIQUE KEY `name_in_space` (`name`, `spaceguid`),
    CONSTRAINT `fk_calls_are_schedulables` FOREIGN KEY (`guid`) REFERENCES `schedulables` (`guid`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;


CREATE TABLE `schedules`
(
    `guid`             varchar(36)  not null,
    `schedulable_guid` varchar(36)           DEFAULT NULL,
    `expression_type`  varchar(255) NOT NULL,
    `expression`       varchar(255)          DEFAULT NULL,
    `enabled`          tinyint(1)   NOT NULL,
    `created_at`       timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`guid`),
    KEY `schedule_fk_schedulable` (`schedulable_guid`),
    CONSTRAINT `schedule_fk_schedulable` FOREIGN KEY (`schedulable_guid`) REFERENCES `schedulables` (`guid`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `histories`
(
    `guid`                 varchar(36) NOT NULL,
    `scheduled_time`       timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `execution_start_time` timestamp   NULL     DEFAULT NULL,
    `execution_end_time`   timestamp   NULL     DEFAULT NULL,
    `message`              varchar(255)         DEFAULT NULL,
    `state`                varchar(64) NOT NULL DEFAULT 'PENDING',
    `schedule_guid`        varchar(36)          DEFAULT NULL,
    `task_guid`            varchar(36)          DEFAULT NULL,
    `created_at`           timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `counted_job`          tinyint(1)           DEFAULT '0',
    PRIMARY KEY (`guid`),
    KEY `history_fk_schedule_guid` (`schedule_guid`),
    CONSTRAINT `history_fk_schedule_guid` FOREIGN KEY (`schedule_guid`) REFERENCES `schedules` (`guid`) ON DELETE SET NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8