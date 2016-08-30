#!/bin/bash
set -e

ARCHS[0]=386
ARCHS[1]=amd64
#ARCHS[2]=arm

OSes[0]=linux
OSes[1]=windows

#VERSION=$(git describe --tags)
VERSION=$TRAVIS_TAG
DIR=$(dirname $(cd $(dirname "${BASH_SOURCE[0]}") && pwd))
NAME=$(basename $DIR)
TMP_DIR="$DIR/tmp"

mkdir -p $DIR/release

for ARCH in "${ARCHS[@]}"; do
    for OS in "${OSes[@]}"; do
        echo "Building $OS $ARCH release..."
        cd $DIR
        env GOOS=$OS GOARCH=$ARCH go build

        if [[ $OS == "windows" ]]; then
            case "$ARCH" in
            "amd64") ARCH="64bit";;
            "386") ARCH="32bit";;
            esac

            FILE=$NAME.exe
        else
            FILE=$NAME
        fi

        CURRENT=$NAME-$OS-$ARCH-$VERSION
        CURRENT_DIR=$TMP_DIR/$CURRENT
        mkdir -p $CURRENT_DIR
        cd $CURRENT_DIR
        mv $DIR/$FILE ./
        cp $DIR/config.example.yml ./config.yml
        cp $DIR/README.md ./

        if [[ $OS == "windows" ]]; then
            zip -r "$DIR/release/$CURRENT.zip" *
        else
            tar -zcvf "$DIR/release/$CURRENT.tar.gz" *
        fi

    done
done

rm -rf $DIR/tmp