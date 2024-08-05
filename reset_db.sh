set -e

export PGPASSWORD=password
if [[ -z $PGPASSWORD ]]; then
    echo "Please set the PGPASSWORD environment variable"
    exit 1
fi


dropdb -U gooo -h 127.0.0.1 -f gooo_development
psql -U gooo -h 127.0.0.1 -d postgres -c "CREATE DATABASE gooo_development;"
psql -U gooo -h 127.0.0.1 -d gooo_development -f fixtures/sql/schema.sql
