#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @File  : main.py.py
# @Author: harry
# @Date  : 18-10-7 上午1:39
# @Desc  : An example bot using jsonrpc sdk

# import from relative path
import ldbg_jsonrpc.client as client
# import from library
# from ldbg_jsonrpc import client
import logging.config
from os import path
import config

if __name__ == "__main__":
    log_conf_path = path.join(path.dirname(path.abspath(__file__)), 'logging.conf')
    logging.config.fileConfig(log_conf_path)
    c = client.Client(config.Host, config.Port)
    c.connect()
    c.init_channel({
        "uuid": config.UUID,
        "producer": True,
        "consumer": True,
        "accept": [{
            "from": ".*",
            "to": ".*",
            "formats": [{
                "api": "ubm-api",
                "version": "1.0",
                "method": "receive"
            }]
        }]
    })
    while True:
        res, ok = c.get_updates({
            "uuid": config.UUID,
            "timeout": "10s",
            "limit": 100,
        })
        if ok:
            print(res)
