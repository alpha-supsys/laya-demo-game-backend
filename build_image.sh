#!/bin/sh

version=$1

docker build -t ${hub}/alpha_supsys/laya-demo-game-backend:${version} .