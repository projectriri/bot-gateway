#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @File  : main.py.py
# @Author: harry
# @Date  : 18-10-7 上午1:39
# @Desc  : An example bot using jsonrpc sdk
import bot_gateway_jsonrpc_client.client as client
import logging.config
from os import path
import config

if __name__ == "__main__":
    log_conf_path = path.join(path.dirname(path.abspath(__file__)), 'logging.conf')
    logging.config.fileConfig(log_conf_path)
    c = client.Client("47.88.223.83", 4700, config.UUID)
    c.connect()
    c.init_channel_default()
    while True:
        c.get_updates()
