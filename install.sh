#!/bin/bash

## Modified script from github.com/yoruokot/superfile. Thanks!!

green='\033[0;32m'
red='\033[0;31m'
yellow='\033[0;33m'
blue='\033[0;34m'
purple='\033[0;35m'
cyan='\033[0;36m'
white='\033[0;37m'
bright_red='\033[1;31m'
bright_green='\033[1;32m'
bright_yellow='\033[1;33m'
bright_blue='\033[1;34m'
bright_purple='\033[1;35m'
bright_cyan='\033[1;36m'
bright_white='\033[1;37m'
nc='\033[0m' # No Color

echo -e '
\033[1;36m
░██████╗░░█████╗░███╗░░░███╗░█████╗░███╗░░██╗░█████╗░░██████╗░███████╗██████╗░░█████╗░░█████╗░██╗░░██╗███████╗██████╗░
 ██╔════╝░██╔══██╗████╗░████║██╔══██╗████╗░██║██╔══██╗██╔════╝░██╔════╝██╔══██╗██╔══██╗██╔══██╗██║░██╔╝██╔════╝██╔══██╗
 ██║░░██╗░██║░░██║██╔████╔██║███████║██╔██╗██║███████║██║░░██╗░█████╗░░██║░░██║██║░░██║██║░░╚═╝█████═╝░█████╗░░██████╔╝
 ██║░░╚██╗██║░░██║██║╚██╔╝██║██╔══██║██║╚████║██╔══██║██║░░╚██╗██╔══╝░░██║░░██║██║░░██║██║░░██╗██╔═██╗░██╔══╝░░██╔══██╗
 ╚██████╔╝╚█████╔╝██║░╚═╝░██║██║░░██║██║░╚███║██║░░██║╚██████╔╝███████╗██████╔╝╚█████╔╝╚█████╔╝██║░╚██╗███████╗██║░░██║
 ░╚═════╝░░╚════╝░╚═╝░░░░░╚═╝╚═╝░░╚═╝╚═╝░░╚══╝╚═╝░░╚═╝░╚═════╝░╚══════╝╚═════╝░░╚════╝░░╚════╝░╚═╝░░╚═╝╚══════╝╚═╝░░╚═╝
'

package=gomanagedocker
arch=$(uname -m)
os=$(uname -s)

if [[ "$arch" == "x86_64" ]]; then
    arch="amd64"
elif [[ "$arch" == "arm"* || "$arch" == "aarch64" ]]; then
    arch="arm64"
else
    echo -e "${red}❌ Fail install goManageDocker: ${yellow}Unsupported architecture: ${nc}${arch}"
    exit 1
fi

if [[ "$os" == "Linux" ]]; then
    os="linux"
elif [[ "$os" == "Darwin" ]]; then
    os="darwin"
else
    echo -e "${red}❌ Fail install goManageDocker: ${yellow}Unsupported operating system${nc}${arch}"
    exit 1
fi

# allow specifying different destination directory
DIR="${DIR:-"$HOME/.local/bin"}"

GITHUB_LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/ajayd-san/gomanagedocker/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
FILE_NAME=${package}_${os}_${arch}_${GITHUB_LATEST_VERSION}
GITHUB_URL="https://github.com/ajayd-san/gomanagedocker/releases/download/${GITHUB_LATEST_VERSION}/${FILE_NAME}.tar.gz"

PARENT_FOLDER="${os}_${arch}_${GITHUB_LATEST_VERSION}"
# install/update the local binary
echo -e "${bright_yellow}Downloading ${cyan}${package} v${version} for ${os} (${arch})...${nc}"
curl -L -o gomanagedocker.tar.gz $GITHUB_URL

echo -e "${bright_yellow}Extracting ${cyan}${package}...${nc}"
tar xzvf gomanagedocker.tar.gz "${PARENT_FOLDER}/gmd" -O >gmd

if install -Dm 755 gmd -t "$DIR" && rm gmd gomanagedocker.tar.gz; then
    echo -e "🎉 ${bright_green}Installation complete!${nc}"
    echo -e "${bright_cyan}You can type ${white}\"${bright_yellow}gmd${white}\" ${bright_cyan}to start!${nc}"
else
    echo -e "${red}❌ Fail install goManageDocker to ${DIR} ${nc}"
fi
