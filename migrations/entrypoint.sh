#!/bin/bash

DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"
TESTDBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$TESTDB sslmode=$DBSSL"

goose postgres "$DBSTRING" up
goose postgres "$TESTDBSTRING" up