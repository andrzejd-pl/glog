# Glog

Mini blog written in golang.

## How to run

```shell script
docker run --name some-api andrzejd/glog:latest
```

## Example configuration

```json
{
  "database_dsn": "user:passwd@tcp(ip-database)/database-name"
}
```

See: [config.example.json](config.example.json)

## Database

Type: MySQL (>= 8.0.18) / MariaDB (>= 10.3)

Install SQL:
```mysql
create database glog;
use glog;
CREATE TABLE `Posts`
(
    `PostId`              int(11)      NOT NULL AUTO_INCREMENT,
    `PostUUID`            bigint(20)   NOT NULL DEFAULT (uuid_short()),
    `PostTitle`           varchar(255) NOT NULL,
    `PostContent`         longtext,
    `PostInsertTimestamp` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`PostId`)
);
```

See: [install.sql](install/install.sql)

## Endpoints

- `GET /api/posts` - return all posts from database

## TODO

- [x] ~~add custom configuration~~
- [ ] add categories to posts
- [ ] add admin panel
- [ ] add frontend
- [ ] add pagination