#!/bin/bash -e

os=$(uname -s | awk '{print tolower($0)}')
processor=$(uname -m)

if [ "$processor" == "x86_64" ]; then
  github_release_arch="amd64"
fi
echo "os=$os"
echo "processor=$processor"
echo "github_release_arch=$github_release_arch"

echo -e "\n\n===================================================="
echo "Downloading plugin ..."
mkdir -p ~/.packer.d/plugins

get_latest_release() {
  curl --silent "https://api.github.com/repos/wata727/packer-post-processor-amazon-ami-management/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

echo "Looking up the latest version ..."
latest_version=$(get_latest_release)
echo "Downloading latest version of $latest_version"
curl -L -o /tmp/packer-post-processor-amazon-ami-management.zip "https://github.com/wata727/packer-post-processor-amazon-ami-management/releases/download/${latest_version}/packer-post-processor-amazon-ami-management_${os}_${github_release_arch}.zip"
retVal=$?
if [ $retVal -ne 0 ]; then
  echo "Failed to download the plugin"
  exit $retVal
else
  echo "Download was successfully"
fi


echo -e "\n\n===================================================="
echo "Unpacking and Installing plugin ..."
cd ~/.packer.d/plugins
unzip -u -j /tmp/packer-post-processor-amazon-ami-management.zip -d ~/.packer.d/plugins

echo -e "\n\n===================================================="
echo "Current list of Packer plugins"
ls -tlr ~/.packer.d/plugins
