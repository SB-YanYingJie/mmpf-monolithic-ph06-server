#!/bin/bash

if [ $# != 1 ]; then
  echo "Specify zip file name as argument"
  exit 1
fi

ZIP_FILE=$1
if [ ! -e $ZIP_FILE ]; then
  echo "File:${ZIP_FILE} not exists"
  exit 1
fi

echo "Create User"
mkdir -p /home/mmpf

echo "Unzip:$ZIP_FILE..."
unzip -u $ZIP_FILE -d /home/mmpf/`basename $ZIP_FILE .zip`
wait

LINK_TO=/home/mmpf/mmpf_modules
if [ -e $LINK_TO ]; then
  echo "Unlink symbolic link:$LINK_TO..."
  unlink $LINK_TO
fi

echo "Create symbolic link..."
ln -s /home/mmpf/`basename $ZIP_FILE .zip` $LINK_TO
