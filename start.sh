#!/bin/bash
folder="/root/logs/"
file="/root/logs/clickhouse-import.log"
if [ ! -d "$folder" ]; then
  mkdir "$folder"
fi
if [ ! -f "$file" ]; then
  touch "$file"
fi
nohup /root/linux_clickhouse_sinker -conf /root/conf >> $file 2>&1 & echo $! > /root/pidfile.txt && sleep 10
PID=$(cat /root/pidfile.txt)
if [ -n "$PID" -a -e /proc/$PID ]; then
    echo "process exists"
    tail -f -n 100 $file
fi
tail -n 100 $file
echo "start faild"
# tail -f -n 100 $file
