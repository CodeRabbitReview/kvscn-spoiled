#!/bin/bash
source_file=$1
stopServer=$2
host=$3
statusCode=0

if [ "hostType" == "d" ]
then
  host="host.docker.internal"
else
  host="172.17.0.1"
fi
sleep 10

file_data="["

while read line; do
  curl -d "$line" -H "Content-Type: application/json" -k -X POST https://"$host":8080/api/
  file_data+="${line},"
done < "$source_file"

file_data=${file_data::-1}
file_data+="]"
file_data=${file_data//[[:blank:]]/}

server_data=$(curl -k https://"$host":8080/api/)

if [ "$file_data" != "$server_data" ]
then
  echo "Failed to store records"
  statusCode=1
else
  echo "All records stored"
fi

if test "$stopServer" == "ok"
then
        fuser -k 8080/tcp
        docker stop storage_server
        exit "$statusCode"
fi