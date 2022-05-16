#!/bin/bash -ex

cat <<-EOF | docker-compose exec -T db psql -U postgres -d rueder_development -w
UPDATE feeds
SET fetched_at = '2020-01-01 00:00:00'
EOF

docker-compose restart worker
