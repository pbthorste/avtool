# avtool
Decrypt ansible vault files in golang.

The tool uses the same command line parameters as the `ansible-vault` command.

# Why?
This is useful if you are writing code in golang and need to manipulate Ansible vault files.

Also useful if you want to decrypt an Ansible vault file, and are not keen on installing python.

# Downloading
See the [Releases](https://github.com/pbthorste/avtool/releases)

# Code
The code is derived from code in the ansible project:

https://github.com/ansible/ansible

# Building
Install golang 1.7. To build, use the Makefile, run:
```bash
make ci
```
The tool will be built and placed in the `build/` folder.
