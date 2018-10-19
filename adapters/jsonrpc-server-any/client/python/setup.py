import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="ldbg-jsonrpc",
    version="1.0.1",
    author="Project Riri Staff",
    author_email="lijiahao99131@gmail.com",
    description="Python3 Client SDK for Little Daemon Bot Gateway "
                "jsonrpc-server-any "
                "(https://projectriri.github.io/bot-gateway/docs/Plugins.html#jsonrpc-server-any) "
                "Plugin.",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/projectriri/bot-gateway/tree/master/adapters/jsonrpc-server-any/client/python",
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
)