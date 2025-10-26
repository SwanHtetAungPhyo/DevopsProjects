#!/bin/bash

echo "🚀 SCL EC2 Setup Script"
echo "======================="
echo ""

# Check if Swan.pem exists
if [ ! -f "Swan.pem" ]; then
    echo "❌ Swan.pem not found in current directory"
    echo "💡 Please copy your Swan.pem file to this directory first"
    echo ""
    echo "Commands to run:"
    echo "  cp ~/.ssh/Swan.pem ."
    echo "  chmod 400 Swan.pem"
    exit 1
fi

# Set proper permissions on the key
echo "🔐 Setting proper permissions on Swan.pem..."
chmod 400 Swan.pem

# Test SSH connection
echo "🧪 Testing SSH connection to EC2..."
echo "💡 This will test the connection (you may need to accept the host key)"

ssh -i "Swan.pem" -o ConnectTimeout=10 -o StrictHostKeyChecking=no ec2-user@172.31.7.120 'echo "✅ SSH connection successful!"' || {
    echo "❌ SSH connection failed"
    echo ""
    echo "🔧 Troubleshooting:"
    echo "  1. Make sure your EC2 instance is running"
    echo "  2. Check the security group allows SSH (port 22)"
    echo "  3. Verify the private IP address: 172.31.7.120"
    echo "  4. Make sure you're using the correct key pair"
    echo ""
    echo "🧪 Manual test command:"
    echo "  ssh -i \"Swan.pem\" ec2-user@172.31.7.120"
    echo ""
    exit 1
}

echo ""
echo "✅ EC2 SSH connection verified!"
echo ""
echo "🚀 Ready to deploy with SCL!"
echo ""
echo "📋 Next steps:"
echo "  1. Review the deployment script: examples/ec2_deployment.scl"
echo "  2. Customize the configuration variables if needed"
echo "  3. Run the deployment:"
echo "     ./scl examples/ec2_deployment.scl --verbose"
echo ""
echo "🔧 Configuration in ec2_deployment.scl:"
echo "  • Target: ec2-user@172.31.7.120"
echo "  • Environment: production"
echo "  • SSL: enabled"
echo "  • Monitoring: enabled"
echo "  • Backup: enabled"
echo ""
echo "💡 The script will:"
echo "  ✓ Analyze your EC2 system"
echo "  ✓ Install and configure Nginx"
echo "  ✓ Set up security (firewall, fail2ban)"
echo "  ✓ Deploy a web application"
echo "  ✓ Configure SSL certificates"
echo "  ✓ Set up monitoring and backups"
echo ""