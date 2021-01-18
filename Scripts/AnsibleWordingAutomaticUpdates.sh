#!/bin/bash

/usr/bin/ansible-playbook -i "$1," /opt/3CX-Reporting/Ansible/PlayBooks/WordingAutomaticUpdates.yaml --extra-vars "ansible_user=root ansible_password=$2 ansible_python_interpreter=/usr/bin/python3";
