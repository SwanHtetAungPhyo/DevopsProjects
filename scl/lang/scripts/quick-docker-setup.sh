#!/bin/bash

echo "🐳 Quick Docker Setup for SCL SSH Testing"
echo "=========================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Build and start container
echo "🔨 Building and starting SSH test container..."
docker-compose -f docker-test/docker-compose.yml up -d --build

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 15

# Check if container is running
if docker ps | grep -q scl-ssh-test; then
    echo "✅ Container is running!"
    
    # Test SSH connectivity
    echo "🔑 Testing SSH connectivity..."
    
    # Remove old host key
    ssh-keygen -f ~/.ssh/known_hosts -R "[localhost]:2222" 2>/dev/null || true
    
    # Add new host key
    ssh-keyscan -p 2222 localhost >> ~/.ssh/known_hosts 2>/dev/null
    
    # Test connection with password
    echo "🧪 Testing SSH connection (you may need to enter 'testpass')..."
    if command -v sshpass > /dev/null; then
        sshpass -p 'testpass' ssh -p 2222 -o StrictHostKeyChecking=no testuser@localhost 'echo "✅ SSH connection successful!"'
    else
        echo "💡 Install sshpass for automatic testing: brew install sshpass (macOS) or apt-get install sshpass (Linux)"
        echo "🔐 Manual test: ssh -p 2222 testuser@localhost (password: testpass)"
    fi
    
    echo ""
    echo "🚀 Ready to test! Run:"
    echo "   ./scl examples/ssh_demo.scl --verbose"
    echo ""
    echo "🛑 To stop: docker-compose -f docker-test/docker-compose.yml down"
    
else
    echo "❌ Container failed to start. Check Docker logs:"
    echo "   docker-compose -f docker-test/docker-compose.yml logs"
fi