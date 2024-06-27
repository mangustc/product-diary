#!/bin/sh

mkdir ./libs
gcc -fPIC -shared icu.c `pkg-config --libs --cflags icu-uc icu-io` -o libSqliteIcu.so

mv ./libSqliteIcu.so ./libs
