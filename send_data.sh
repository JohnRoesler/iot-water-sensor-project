#!/usr/bin/env sh
while IFS=, read -r timestamp symbol volume temperature
do
  data=$(printf '{"timeStamp": %s, "symbol": "%s", "volume": %s, "temperature": %s}' "$timestamp" "$symbol" "$volume" "$temperature")
  curl localhost:8888 -v -d "$data"
done < input.csv
