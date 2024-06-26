#!/bin/bash
cd `dirname $0`
set -e

OUTPUT_CSV=output.csv

appendHeader() {
  local tHeader=$(docker stats --no-stream --format "table {{.Name}},{{.CPUPerc}},{{.MemUsage}},{{.NetIO}},{{.BlockIO}}" | head -n 1 | sed "s# / #,#g")
  sed -i "1i TIME,${tHeader}" ${OUTPUT_CSV}
  exit 0
}

# trap appendHeader 2
Record(){
  docker stats --no-stream --format "table {{.Name}},{{.CPUPerc}},{{.MemPerc}},{{.MemUsage}},{{.NetIO}},{{.BlockIO}}" \
	  | grep -v "NAME" \
	  | xargs -I@ echo $(date "+%T"),@ \
	  | sed "s# / #,#g" >> ${OUTPUT_CSV}
}

if [ -e ${OUTPUT_CSV} ];then
  rm ${OUTPUT_CSV}
  touch ${OUTPUT_CSV}
else
  touch ${OUTPUT_CSV}
fi

while true
do
  Record;
  # sleep 5;
done

