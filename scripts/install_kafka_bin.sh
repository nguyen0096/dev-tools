#!/bin/bash
set -o xtrace

kafkaVersion=kafka_2.13-3.4.0
zipFilename=$kafkaVersion.tgz

# check if file exists
if [ -f "$zipFilename" ]; then
    printf "Removing file with same name\n"
    rm "$zipFilename"
fi
printf "Download binary file\n"
# download binary file, exit with code 1 if failed
curl -O "https://downloads.apache.org/kafka/3.4.0/$zipFilename"

# check if folder with same name exists
if [ -d "$kafkaVersion" ]; then
    printf "Removing folder with same name\n"
    rm -rf "$kafkaVersion"
fi
printf "Extracting\n"
tar -xzf "$zipFilename"
printf "Removing zip file\n"
rm "$zipFilename"

printf "Move Kafka to /usr/local/kafka\n"
# check if folder exists
if [ -d "/usr/local/$kafkaVersion" ]; then
    printf "Removing folder with same name\n"
    rm -rf "/usr/local/$kafkaVersion"
fi
sudo mv "$kafkaVersion" "/usr/local/$kafkaVersion"

# add kafka bin to PATH
printf "Add kafka bin to PATH\n"
# check if sudo is required, if not, clean up and exit
if [ "$EUID" -ne 0 ]; then
    printf "Please run as root\n"
    rm -rf "/usr/local/$kafkaVersion"
    exit 1
fi

# detect devtool section syntax to check if it is already added
sectionName=devtool.kafka_bin
if grep -q "# $sectionName.begin" ~/.zprofile; then
    printf "Section already exists, replace it\n"
    perl -0777 -i -pe "s/# $sectionName.begin\\n.*\\n# $sectionName.end/# $sectionName.begin\\nexport PATH=\\\$PATH:\\/usr\\/local\\/$kafkaVersion\\/bin\\n# $sectionName.end/s" ~/.zprofile
else
    printf "Adding kafka bin path\n"
    echo "\n# $sectionName.begin" >> ~/.zprofile
    echo "export PATH=\$PATH:/usr/local/$kafkaVersion/bin" >> ~/.zprofile
    echo "# $sectionName.end" >> ~/.zprofile
fi

# Below aproach is not working, since bash script need jar files in lib folder added into CLASSPATH
# printf "Create alias for all bin files in /usr/local/$kafkaVersion/bin\n"
# # Loop though all files in /usr/local/kafka/bin
# for file in /usr/local/$kafkaVersion/bin/*; do
#     # if file is not a regular file, skip it
#     [[ -f $file ]] || continue
    
#     filename=$(basename -- "$file")

#     # check if sudo is required, if not, clean up and exit
#     if [ "$EUID" -ne 0 ]; then
#         printf "Please run as root\n"
#         rm -rf /usr/local/kafka
#         exit 1
#     fi

#     # delete existing file
#     rm "/usr/local/bin/$filename"

#     # Create symbol link for each file in /usr/local/bin
#     ln -s "$file" "/usr/local/bin/$filename"
# done