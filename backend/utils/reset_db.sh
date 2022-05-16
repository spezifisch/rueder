#!/bin/bash -ex

# disconnect users sitting on the db
# https://stackoverflow.com/questions/5408156/how-to-drop-a-postgresql-database-if-there-are-active-connections-to-it
# https://stackoverflow.com/questions/43099116/error-the-input-device-is-not-a-tty
cat <<-EOF | docker-compose exec -T db psql -U postgres -w
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE datname = 'rueder_development'
  AND pid <> pg_backend_pid();

DROP DATABASE rueder_development WITH (FORCE);
EOF

docker-compose exec dev soda create
docker-compose exec dev soda migrate up
