#!/bin/bash

path="."
date=$(date +%Y-%m-%d)      
filename=$date.zip
log=$path$filename
days=15 # Set the age limit of files to delete (in days)

log_message() {
    echo "$(date +%Y-%m-%d)  - $1" >> "$log"
}

find "$path" -type f -mtime +$days -exec rm -rf {} \; -print | while read line; do
    log_message "[/] $line deleted."
done

exit 0