#!/bin/bash

# testing user must be a superuser role

DB_NAME=${1:-coupons_test}
DB_PASSWORD=${2:-testingtesting123}
DB_USER=${3:-testing}

ROOT_DIR_PATH=$(cd $(dirname $0)/.. && pwd)

for FILE in `ls -1 ${ROOT_DIR_PATH}/db/migrations/*up.sql`
do
  echo "Executing: `basename ${FILE}`"
  psql "dbname='${DB_NAME}' password='${DB_PASSWORD}' user='${DB_USER}'" -f ${FILE} 2>allfiles.err.log
done