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
    def __init__(self, host: str, port: int):
        self.host = host
        self.port = port
        self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.id = 0

    @staticmethod
    def _to_rpc_method(method: str):
        return "Broker." + method.title().replace("_", "")

    def __getattr__(self, item):
        return lambda obj: self._invoke(self._to_rpc_method(item), obj)

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
            "id": self.id
        }
        self.id = (self.id + 1) % 65536
        try:
            req = json.dumps(request)
            logging.debug("invoke: {}".format(req))
            self._socket.sendall(bytes(req, encoding="utf-8"))
        except socket.error as e:
            logging.error("failed to invoke method {} with params {}, error: {}".format(method, params, e))
            return None, False

        response = bytes()
        while True:
            try:
                data = self._socket.recv(1)
            except socket.timeout:
                logging.error("failed to invoke method {} with params {}, error: timeout".format(method, params))
                return None, False
            if not data:
                break
            response += data
            if data == b'\n':
                break

        rst = str(response, encoding="utf-8")
        logging.debug("invoke response: {}".format(rst))
        rst = json.loads(rst)
        return rst["result"], True
