#!/usr/bin/env python

import os
import sys
import json
import subprocess
import signal
import pycurl
import re
import csv
from subprocess import Popen, PIPE

from io import BytesIO
import StringIO

headers = {}
myFile = open('/tmp/gitlab_users.csv', 'w')
record = [['WWID', 'Name', 'Email']]

def header_function(header_line):
    # HTTP standard specifies that headers are encoded in iso-8859-1.
    # On Python 2, decoding step can be skipped.
    # On Python 3, decoding step is required.
    header_line = header_line.decode('iso-8859-1')

    # Header lines include the first status line (HTTP/1.x ...).
    # We are going to ignore all lines that don't have a colon in them.
    # This will botch headers that are split on multiple lines...
    if ':' not in header_line:
        return

    # Break the header line into header name and value.
    name, value = header_line.split(':', 1)

    # Remove whitespace that may be present.
    # Header lines include the trailing newline, and there may be whitespace
    # around the colon.
    name = name.strip()
    value = value.strip()

    # Header names are case insensitive.
    # Lowercase name here.
    name = name.lower()

    headers[name] = value

URL = 'https://<my gitlab server>/api/v4/users'
Headers = 'Private-Token: <my token>'

# Setup curl values for GitLab API
buffer = BytesIO()
ch = pycurl.Curl()
ch.setopt(ch.HEADER, True)
ch.setopt(ch.URL, URL)
ch.setopt(ch.HTTPHEADER, [Headers])
ch.setopt(ch.NOBODY, True)
ch.setopt(ch.HEADERFUNCTION, header_function)
ch.setopt(ch.WRITEFUNCTION, lambda x: None)

ch.perform()
ch.close()

Total_Users = headers['x-total']
Num_pages = int(Total_Users)//100
remainder = int(Total_Users)%100

print ('Total users = ' + str(Total_Users) )
if (remainder > 0 ):
    Num_pages += 2

print ('Total pages needed = ' + str(Num_pages) )

writer = csv.writer(myFile)
writer.writerows(record)
for x in range(1, Num_pages):
    URL = 'https://<my gitlab server>/api/v4/users?per_page=100&page=' + str(x)
    Headers = 'Private-Token: <my token>'

    pages = pycurl.Curl()
    data = BytesIO()
    pages.setopt(pages.URL, URL)
    pages.setopt(pages.HTTPHEADER, [Headers])
    pages.setopt(pages.WRITEFUNCTION, data.write)

    pages.perform()
    pages.close()

    Results = data.getvalue()
    dict = json.loads(Results)

    for id in dict:
        getName = str(id["name"])
        getUserid  = str(id["username"])
        getEmail = str(id["email"])
        record[0][0] = getUserid
        record[0][1] = getName
        record[0][2] = getEmail
        writer.writerows(record)

myFile.close()

Command = 'echo "GitLab production user report in CSV format" | mutt -s "Production Gitlab monthly user reporting" -a "/tmp/gitlab_users.csv" -- myemail@domain.com 2> /tmp/email.log'
outPut = Popen( Command , shell=True, stderr=PIPE, stdout=PIPE)

print ( outPut.communicate()[0] )
