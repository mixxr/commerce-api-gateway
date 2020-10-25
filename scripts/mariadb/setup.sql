--CREATE DATABASE dcgw;

USE dcgw;

CREATE TABLE IF NOT EXISTS table_ (
   id BIGINT NOT NULL AUTO_INCREMENT,
   created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deflang CHAR(2),
   ncols SMALLINT NOT NULL DEFAULT 0,
   nrows SMALLINT NOT NULL DEFAULT 0,
   owner VARCHAR(32) NOT NULL,
   name VARCHAR(32) NOT NULL,
   descr VARCHAR(256) NOT NULL DEFAULT '',
   tags VARCHAR(256) NOT NULL DEFAULT '',
   PRIMARY KEY ( id ),
   UNIQUE KEY ( owner, name )
);

CREATE INDEX table_index_desc ON table_( descr );
CREATE INDEX table_index_tags ON table_( tags );
