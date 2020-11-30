#! /bin/bash

for i in $(find . -iname "*_test.go" -exec dirname {} \; | uniq | grep -v example)
do
    go test -race -cover $i;
    if [ $? != 0 ]
    then
      return 1
    fi
done