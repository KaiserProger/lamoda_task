#!/bin/bash

echo "SELECT 'CREATE DATABASE $TESTDB' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$TESTDB')\gexec" | psql -U $DBUSER $DBNAME