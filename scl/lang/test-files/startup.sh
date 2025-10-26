#!/bin/bash
echo "Starting application..."
cd /opt/myapp
./myapp --config=/etc/myapp/app.conf
echo "Application started successfully"