#!/bin/bash


source /opt/3CX-Reporting/Scripts/DB_CONFIG_FILE

daTeY=$(date +%Y)
daTem=$(date +%m)
daTed=$(date +%d)
daTe=$(date +%Y_%m_%d_%H_%M)

mysqldump -P $DBport -h $DBhost -u $DBuser -p$DBpass --extended-insert=FALSE 3cxReporting > /tmp/$(echo $daTe)_3cxReporting.sql

mkdir -p $BaseDIR/$daTeY/$daTem/$daTed
tar cf - /tmp/$(echo $daTe)_3cxReporting.sql | 7za a -si -t7z -m0=lzma -mx=9 -mfb=64 -md=32m -ms=on $BaseDIR/$daTeY/$daTem/$daTed/$(echo $daTe)_3cxReporting.tar.7z

rm /tmp/$(echo $daTe)_3cxReporting.sql

