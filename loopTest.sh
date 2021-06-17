#!/bin/bash

#!/bin/bash
e=0
c=0
for i in {1..20}
do
  echo "Run #$i"

  (battlesnake play --width 7 --height 7 --name Control --url "http://0.0.0.0:8080" --name Experiment --url "http://0.0.0.0:8080/experiment" &> sneks.log)

  if [[ "$(tail -1 sneks.log)" == *"Experiment is the winner"* ]]; then
   ((e=e+1))
  fi

  if [[ "$(tail -1 sneks.log)" == *"Control is the winner"* ]]; then
   ((c=c+1))
  fi

  echo "$(tail -1 sneks.log)"
done

echo "Score: Experiment $e - Control $c"
