#!/bin/bash

echo "🔍 EC2 Connection Checker"
echo "========================"
echo ""

# Check if Swan.pem exists and has correct permissions
if [ ! -f "Swan.pem" ]; then
    echo "❌ Swan.pem not found"
    echo "💡 Copy your key file: cp ~/.ssh/Swan.pem ."
    exit 1
fi

# Set correct permissions
chmod 400 Swan.pem
echo "✅ Swan.pem permissions set correctly"

echo ""
echo "🧪 Testing different connection methods..."
echo ""

# Test 1: Private IP (if you're inside AWS VPC)
echo "Test 1: Private IP (172.31.7.120)"
echo "-----------------------------------"
timeout 10 ssh -i "Swan.pem" -o ConnectTimeout=5 -o StrictHostKeyChecking=no ec2-user@172.31.7.120 'echo "Private IP works!"' 2>/dev/null && {
    echo "✅ Private IP connection successful!"
    echo "🎯 Use this in your SCL script: target := \"ec2-user@172.31.7.120\";"
    exit 0
} || {
    echo "❌ Private IP connection failed"
}

echo ""

# Instructions for finding public IP
echo "🔍 To find your EC2 public IP:"
echo "------------------------------"
echo "1. Go to AWS Console: https://console.aws.amazon.com/ec2/"
echo "2. Click 'Instances' in the left sidebar"
echo "3. Find your instance in the list"
echo "4. Look for these values:"
echo "   • Public IPv4 address (e.g., 3.15.123.456)"
echo "   • Public IPv4 DNS (e.g., ec2-3-15-123-456.us-east-2.compute.amazonaws.com)"
echo ""

echo "🧪 Manual test commands:"
echo "------------------------"
echo "# Test with public IP (replace with your actual IP):"
echo "ssh -i \"Swan.pem\" ec2-user@YOUR_PUBLIC_IP"
echo ""
echo "# Test with public DNS (replace with your actual DNS):"
echo "ssh -i \"Swan.pem\" ec2-user@ec2-xx-xx-xx-xx.compute-1.amazonaws.com"
echo ""

echo "🔧 Common issues and solutions:"
echo "------------------------------"
echo "1. ❌ Connection timeout:"
echo "   → Check EC2 instance is running (not stopped)"
echo "   → Check Security Group allows SSH (port 22) from your IP"
echo ""
echo "2. ❌ Permission denied:"
echo "   → Make sure you're using the correct key pair"
echo "   → Check key permissions: chmod 400 Swan.pem"
echo ""
echo "3. ❌ Host key verification failed:"
echo "   → Add -o StrictHostKeyChecking=no to ssh command"
echo ""

echo "📋 Security Group Requirements:"
echo "------------------------------"
echo "Your EC2 Security Group must allow:"
echo "• Type: SSH"
echo "• Protocol: TCP"
echo "• Port: 22"
echo "• Source: Your IP address or 0.0.0.0/0 (less secure)"
echo ""

echo "🎯 Once you find the working connection, update your SCL script:"
echo "target := \"ec2-user@YOUR_WORKING_IP_OR_DNS\";"