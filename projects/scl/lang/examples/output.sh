#!/bin/bash

set -e  # Exit on error

MODE=compile
SETTING=configuration
SUPER_USER=true
ON_ERROR=rollback
TARGET=ec2-user@YOUR_PUBLIC_IP_HERE

main() {
    echo "🔍 EC2 Connection Diagnostics"
    echo "============================="
    echo "Target: "+target"
    echo ""
    echo "⚠️  IMPORTANT: Update the target IP in this script!"
    echo "   1. Go to AWS Console → EC2 → Instances"
    echo "   2. Find your instance"
    echo "   3. Copy the PUBLIC IPv4 address or Public DNS"
    echo "   4. Update the 'target' variable above"
    echo ""
    connection_test
    echo "✅ Connection test completed"
}

connection_test() {
    echo "🧪 Testing basic connection..."
    sysinfo
    echo "✅ If you see system info above, the connection works!"
}

# Execute main function
main
