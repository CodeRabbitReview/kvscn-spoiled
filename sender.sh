#!/bin/bash
source_file=$1
host=$2

if [ "$host" == "m" ]
then
  host="host.docker.internal"
else
  host="172.17.0.1"
fi

file_data="["

while read line; do
  curl -d "$line" -H "Content-Type: application/json" -k -X POST https://"$host":8080/api/
  file_data+="${line},"
done < "$source_file"

file_data=${file_data::-1}
file_data+="]"
file_data=${file_data//[[:blank:]]/}

server_data=$(curl -k https://"$host"/api/)

if [ "$file_data" != "$server_data" ]
then
  echo "Failed to store records"
else
  echo "All records stored"
fi
