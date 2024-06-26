#!/bin/bash
echo "Create zip file from module"
cd /tmp/mmpf
DATE=`date "+%Y%m%d"`
zip -r mmpf_modules_${DATE}.zip *

echo "Delivery"
mkdir -p /app/delivery
cd /tmp/mmpf
cp mmpf_modules_${DATE}.zip /app/install/unzip_module.sh /app/delivery
