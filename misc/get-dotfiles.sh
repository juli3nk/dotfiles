#!/bin/bash

DOTFILES_PATH="${HOME}/bin/dotfiles"
DOTFILES_URL="https://github.com/juliengk/dotfiles"

function get_latest_tag {
	git ls-remote --tags ${DOTFILES_URL}.git | awk -F'/' '{ print $3 }' | sort | tail -n 1
}

KERNEL_RELEASE=$(uname -s)
ARCH=$(uname -m)

if [ "${KERNEL_RELEASE}" == "Linux" ] ; then
	KERNEL_RELEASE="linux"
elif [ "${KERNEL_RELEASE}" == "Darwin" ] ; then
	KERNEL_RELEASE="darwin"
else
	KERNEL_RELEASE="none"
fi

if [[ ! "${KERNEL_RELEASE}" =~ linux|darwin && ! "${ARCH}" =~ i386|x86_64 ]] ; then
	exit 1
fi

LATEST_TAG=`get_latest_tag`

curl -sSL ${DOTFILES_URL}/releases/download/${LATEST_TAG}/dotfiles-${KERNEL_RELEASE}-${ARCH}.tar.bz2 | tar -xjC /tmp/
mv /tmp/dotfiles-${KERNEL_RELEASE}-${ARCH} ${DOTFILES_PATH}
