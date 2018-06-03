#!/usr/bin/env bash
count=0
for dp in "$@"
do
dir_path=$dp
count=$(($count+1))
cd $dir_path
for i in `ls -1v ${pwd}`; do
    # 再遍历一次防止变成 1、10、2、3 这样的顺序
    echo "file '${i}'" >> "ff.txt"
done

filename=`basename $dp`
ffmpeg -f concat -i ff.txt -c copy ${filename}.mp4

rm ff.txt
rm *.flv
rm *.ts
rm *.f4v
cd -
done
