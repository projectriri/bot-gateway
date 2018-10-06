#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @File  : client.py
# @Author: harry
# @Date  : 18-10-7 上午1:27
# @Desc  : jsonrpc client wrapper

import socket
import logging
import json


class Client:
    def __init__(self, host: str, port: int, uuid: str, buffer_size: int = 1024):
        self.host = host
        self.port = port
        self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.buffer_size = buffer_size
        self.uuid = uuid
        self.id = 0

    def connect(self):
        logging.info("connecting to server {}".format(self.host + ":" + str(self.port)))
        if self._socket.connect_ex((self.host, self.port)) != 0:
            logging.critical("failed to connect, exiting")
            exit(-1)
        logging.info("connect success")

    def _invoke(self, method: str, params: object) -> (str, bool):
        request = {
            "method": method,
            "params": [params],
            "jsonrpc": "2.0",
            "id": self.id
        }
        self.id += 1
        # send request
        try:
            req = json.dumps(request)
            logging.info("invoke: {}".format(req))
            self._socket.sendall(bytes(req, encoding="utf-8"))
        except socket.error as e:
            logging.error("failed to invoke method {} with params {}, error: {}".format(method, params, e))
            return None, False

        # get response
        response_list = []
        while True:
            try:
                data = self._socket.recv(self.buffer_size)
            except socket.timeout:
                logging.error("failed to invoke method {} with params {}, error: timeout".format(method, params))
                return None, False
            if not data:
                break
            response_list.append(str(data, encoding="utf-8"))
            if len(data) < self.buffer_size:
                break

        rst = ''.join(response_list)
        logging.info("invoke response: {}".format(rst))
        return rst, True

    def init_channel(self, producer: bool, consumer: bool, accept: [dict]) -> bool:
        res, ok = self._invoke("Broker.InitChannel", {
            "uuid": self.uuid,
            "producer": producer,
            "consumer": consumer,
            "accept": accept
        })
        if ok:
            rst = json.loads(res)
            if rst["result"]["code"] == 10001:
                logging.info("channel init success")
                return True
        return False

    def init_channel_default(self) -> bool:
        return self.init_channel(True, True, [
            {
                "from": ".*",
                "to": ".*",
                "formats": [
                    {
                        "api": "ubm-api",
                        "version": "1.0",
                        "method": "receive"
                    }
                ]
            }
        ])

    def get_updates(self) -> (list, bool):
        res, ok = self._invoke("Broker.GetUpdates", {
            "uuid": self.uuid,
            "timeout": "10s",
            "limit": 100
        })
        if ok:
            rst = json.loads(res)
            if rst["result"]["code"] == 10000 or rst["result"]["code"] == 10004:
                if rst["result"]["packets"] is None:
                    rst["result"]["packets"] = []
                logging.info("receiving updates: {}".format(rst["result"]["packets"]))
                return rst["result"]["packets"], True
        return None, False
