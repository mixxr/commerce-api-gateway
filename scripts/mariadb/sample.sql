INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows) VALUES
('it','mike','ssn_ca','Security Social Number, State of California','ssn,ca,california,wellfare',3,3),
('it','mike','ssn_ny','Security Social Number, State of New York','ssn,ny,new york,wellfare',3,0),
('it','anthony','ssn_ny','Security Social Number, State of New York','ssn,ny,new york,wellfare',4,3),
('it','mike','cars_diesel_vw','VW Diesel Cars','car,diesel,volkswagen',4,0);

-- mike ssn_ca
CREATE TABLE mike_ssn_ca_colnames (
   id BIGINT NOT NULL AUTO_INCREMENT,
   created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   lang CHAR(2),
   colname1 VARCHAR(32) NOT NULL,
   colname2 VARCHAR(32) NOT NULL,
   colname0 VARCHAR(32) NOT NULL,
   PRIMARY KEY ( id )
);
INSERT INTO mike_ssn_ca_colnames (lang,colname1,colname2,colname0) VALUES
('it','nome','cognome','ssn'),
('en','name','surname','ssn');

CREATE TABLE mike_ssn_ca_values (
   id BIGINT NOT NULL AUTO_INCREMENT,
   created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   colvalue1 VARCHAR(256) NOT NULL,
   colvalue2 VARCHAR(256) NOT NULL,
   colvalue0 VARCHAR(256) NOT NULL,
   PRIMARY KEY ( id )
);
INSERT INTO mike_ssn_ca_values (colvalue1,colvalue2,colvalue0) VALUES
('mike','douglàs','3897428934EWREW'),
('äbel','òmar ópël','3897428934EWREW'),
('anthony','di martino','234234FSAFSADF');

-- anthony ssn_ca
CREATE TABLE anthony_ssn_ca_colnames (
   id BIGINT NOT NULL AUTO_INCREMENT,
   created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   lang CHAR(2),
   colname1 VARCHAR(32) NOT NULL,
   colname2 VARCHAR(32) NOT NULL,
   colname3 VARCHAR(32) NOT NULL,
   colname0 VARCHAR(32) NOT NULL,
   PRIMARY KEY ( id )
);
INSERT INTO anthony_ssn_ca_colnames (lang,colname1,colname2,colname3,colname0) VALUES
('it','nome','cognome','sesso','ssn'),
('en','name','surname','gender','ssn');

CREATE TABLE anthony_ssn_ca_values (
   id BIGINT NOT NULL AUTO_INCREMENT,
   created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   colvalue1 VARCHAR(256) NOT NULL,
   colvalue2 VARCHAR(256) NOT NULL,
   colvalue3 VARCHAR(256) NOT NULL,
   colvalue0 VARCHAR(256) NOT NULL,
   PRIMARY KEY ( id )
);
INSERT INTO anthony_ssn_ca_values (colvalue1,colvalue2,colvalue3,colvalue0) VALUES
('mike','douglàs','male','3897428934EWREW'),
('äbel','òmar ópël','male','3897428934EWREW'),
('zoe','di martino','female','93749823ASFSAFD');

-- end