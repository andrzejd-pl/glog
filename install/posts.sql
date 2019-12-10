CREATE TABLE `Posts`
(
    `PostId`              int(11)      NOT NULL AUTO_INCREMENT,
    `PostUUID`            bigint(20)   NOT NULL DEFAULT (uuid_short()),
    `PostTitle`           varchar(255) NOT NULL,
    `PostContent`         longtext,
    `PostInsertTimestamp` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`PostId`)
)
