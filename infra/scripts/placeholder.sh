#!/bin/bash 



echo "Hello from the pulumi " >> /var/www/html/index.html


sudo yum update -y 


sudo install -y httpd


systemctl start httpd

systemctl enable httpd
