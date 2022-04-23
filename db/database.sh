#!/bin/bash

echo "Starting db init"

sqlite3 main.db <<EOF
CREATE TABLE user (id INTEGER PRIMARY KEY, name TEXT, email TEXT, avatarUrl TEXT, timeCreated INTEGER);
EOF