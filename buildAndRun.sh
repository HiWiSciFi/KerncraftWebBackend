#!/bin/bash
rm -r ~/.local
pip3 install Kerncraft
go build -o ./out/___go_build_github_com_HiWiSciFi_KerncraftWebBackend github.com/HiWiSciFi/KerncraftWebBackend
./out/___go_build_github_com_HiWiSciFi_KerncraftWebBackend
