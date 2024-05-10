-- I suspect this is bad practice in postgres.... because the creation of the user ends up in WAL logs
-- Its a bad idea in mysql anyhoo, but i'm lazy

CREATE USER replicator WITH REPLICATION ENCRYPTED PASSWORD 'replicator';
SELECT pg_create_physical_replication_slot('replication_slot');

CREATE TABLE example_table (id SERIAL PRIMARY KEY, name VARCHAR(100), age INT);
