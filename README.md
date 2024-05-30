# PgCat-Experiments

This is mostly me just tinkering around.  It's also an overengineered pattern for optional cli flags i've screwed around on other projects that i don't wanna forget about. 

I pointed this at pgcat k8s deployment in google cloud on cloudsql.

## Take Aways
- Query splitting works as expected
- Transactions as stated do not support stored procedures.
- Binary binary_parameters are very much required even in session.  Seemingly this is a hold out from pg-bouncer.  
- Without the binary_parameters even in session pooling we seem to see about 10% error ratio and it will attempt to periodically insert on the replicas.
- No support for env variables at this time.  So the user/pw ends up in plaintext in a tomlfile somewhere.  This probably needs to be put as a k8s secret and not a config map.  Not really ideal.
- During a writer failover pgcat continues to try and connect.
- probably use auth_query to do authentication to hide the users, but you need access to pg_shadow

## Cloudsql - Take Aways
- Postgres cloudsql does seem generally better than mysql.
- The primary HA replica is not addressable in postgres.  This seems true with alloydb as well
- Promoting a replica is not infact a takeover but a 'fork this replica and make it writable'
- Replicas even seemingly small ones (>10GB) take 15 minutes to spin up.
- Manually triggered failover takes about 30s on a n2-highmem-2.  GUI indicated this took 30s, practice looks more like
2s.  I've seen mysqls take 15m in cloudsql when sufficiently large. 
- Primary is not addressable during a failover.
- Seems like async replication the cloudsql isnt always connected?  Or at least it doesnt show up in pg_stat_replication.  While i expected it to async commit i didnt expect it to be intermittently connected.

### Useful commands
`SELECT pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) AS replication_lag_bytes FROM pg_stat_replication`
`SELECT * FROM pg_stat_wal_receiver;`
10:10:51
