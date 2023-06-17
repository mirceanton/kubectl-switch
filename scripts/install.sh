#!/bin/bash

cd /tmp
wget https://raw.githubusercontent.com/mirceanton/kswitcher/v1.0.0/kswitcher.py
wget https://raw.githubusercontent.com/mirceanton/kswitcher/v1.0.0/requirements.txt
pip install -r requirements.txt
chmod +x kswitcher.py
sudo mv kswitcher.py /usr/local/bin/kswitcher
