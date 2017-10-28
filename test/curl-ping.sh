#!/bin/bash

curl -XPOST http://localhost:9000/ping?appid=1 -d $'s3cr3t\ngoogle.com.'
dig @127.0.0.1 -p 5454 google.com.
