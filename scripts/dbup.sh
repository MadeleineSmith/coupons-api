#!/bin/bash

ROOT_DIR_PATH=$(cd $(dirname $0)/.. && pwd)

for FILE in `ls -1 ${ROOT_DIR_PATH}/db/migrations/*up.sql`
do
  cat ${FILE} >> /tmp/allfiles.sql
done

DB_NAME=$1
DB_USER=$2

psql -d ${DB_NAME} -U ${DB_USER} -f /tmp/allfiles.sql 2>allfiles.err.log

rm /tmp/allfiles.sql