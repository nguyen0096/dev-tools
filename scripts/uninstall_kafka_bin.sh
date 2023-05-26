#!/bin/bash
set -o xtrace

kafkaVersion=kafka_2.13-3.4.0

# uninstall function to remove kafka
# delete all symbolic links into /usr/local/kafka/bin
# delete /usr/local/kafka
uninstall() {
    rm -rf "/usr/local/$kafkaVersion"
    perl -0777 -i -pe "s/# $sectionName.begin\\n.*\\n# $sectionName.end//s" ~/.zprofile
}

uninstall