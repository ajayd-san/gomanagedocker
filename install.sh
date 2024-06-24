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

temp_dir=$(mktemp -d)
if [ $? -ne 0 ]; then
    echo -e "${red}❌ Fail install goManageDocker: ${yellow}Unable to create temporary directory${nc}"
    exit 1
fi

package=gomanagedocker
version=1.2
arch=$(uname -m)
os=$(uname -s)

cd "${temp_dir}"

if [[ "$arch" == "x86_64" ]]; then
    arch="amd64"
elif [[ "$arch" == "arm"* ]]; then
    arch="arm64"
else
    echo -e "${red}❌ Fail install goManageDocker: ${yellow}Unsupported architecture${nc}"
    exit 1
fi

if [[ "$os" == "Linux" ]]; then
    os="linux"
elif [[ "$os" == "Darwin" ]]; then
    os="darwin"
else
    echo -e "${red}❌ Fail install goManageDocker: ${yellow}Unsupported operating system${nc}"
    exit 1
fi

file_name=${package}_${os}_${arch}_v${version}

url="https://github.com/ajayd-san/gomanagedocker/releases/download/v${version}/${file_name}.tar.gz"

if command -v curl &> /dev/null; then
    echo -e "${bright_yellow}Downloading ${cyan}${package} v${version} for ${os} (${arch})...${nc}"
    curl -sLO "$url"
else
    echo -e "${bright_yellow}Downloading ${cyan}${package} v${version} for ${os} (${arch})...${nc}"
    wget -q "$url"
fi

echo -e "${bright_yellow}Extracting ${cyan}${package}...${nc}"
tar -xzf "${file_name}.tar.gz"

echo -e "${bright_yellow}Installing ${cyan}${package}...${nc}"

extracted_folder_name=${os}_${arch}_v${version}
cd ${extracted_folder_name}
chmod +x ./gmd

if sudo mv ./gmd /usr/local/bin/; then
  echo -e "🎉 ${bright_green}Installation complete!${nc}"
  echo -e "${bright_cyan}You can type ${white}\"${bright_yellow}gmd${white}\" ${bright_cyan}to start!${nc}"
else
  echo -e "${red}❌ Fail install goManageDocker: ${yellow}Unable to move binary to /usr/local/bin. Do you have sudo permissions?${nc}"
fi

rm -rf "$temp_dir"
