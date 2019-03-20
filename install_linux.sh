#!/bin/bash

os=$(uname -s | awk '{print tolower($0)}')
processor=$(uname -p)

if [ "$processor" == "x86_64" ]; then
  github_release_arch="amd64"
fi
echo "os=$os"
echo "processor=$processor"
echo "github_release_arch=$github_release_arch"

echo -e "\n\n===================================================="
echo "Downloading plugin ..."
mkdir -p ~/.packer.d/plugins

latest_version=$(curl -s https://api.github.com/repos/wata727/packer-post-processor-amazon-ami-management/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")' )
echo "latest_version=$latest_version"

#wget https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/download/$latest_version/packer-post-processor-amazon-ami-management_${os}_${github_release_arch}.zip -P /tmp/
curl -s -L -o packer-post-processor-amazon-ami-management.zip https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/download/$latest_version/packer-post-processor-amazon-ami-management_${os}_${github_release_arch}.zip
retVal=$?
if [ $retVal -ne 0 ]; then
  echo "Failed to download the plugin"
  exit $retVal
fi


echo -e "\n\n===================================================="
echo "Unpacking and Installing plugin ..."
cd ~/.packer.d/plugins
unzip -f -j /tmp/packer-post-processor-amazon-ami-management_${os}_${github_release_arch}.zip -d ~/.packer.d/plugins

echo -e "\n\n===================================================="
echo "Current list of Packer plugins"
ls -tlr ~/.packer.d/plugins
