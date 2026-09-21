#!/bin/sh
set -eu

sudo pacman --needed --noconfirm -S \
  fontconfig \
  freetype2 \
  libx11 \
  libxcb \
  libxft \
  libxinerama \
  xclip \
  xdotool
