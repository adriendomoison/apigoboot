CREATE USER apigoboot;
ALTER ROLE apigoboot WITH CREATEDB;
CREATE DATABASE apigoboot OWNER apigoboot;
ALTER USER apigoboot WITH ENCRYPTED PASSWORD 'apigoboot';

CREATE USER apigoboot_test;
ALTER ROLE apigoboot_test WITH CREATEDB;
CREATE DATABASE apigoboot_test OWNER apigoboot_test;
ALTER USER apigoboot_test WITH ENCRYPTED PASSWORD 'apigoboot_test';
