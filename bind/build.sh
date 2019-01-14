#!/bin/sh

export ANDROID_HOME=/Users/mannix/Android/sdk

gomobile bind -target=ios github.com/GaoMjun/ladder/bind/ladderclient
gomobile bind -target=android github.com/GaoMjun/ladder/bind/ladderclient