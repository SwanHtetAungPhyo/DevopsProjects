#!/bin/bash

echo "🐳 Ubuntu Docker Setup for SCL Testing"
echo "======================================="
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "🔨 Building Ubuntu test container..."
docker-compose -f docker-ubuntu-test/docker-compose.yml up -d --build

echo "⏳ Waiting for services to start..."
sleep 15

# Check if container is running
if docker ps | grep -q scl-ubuntu-test; then
    echo "✅ Container is running!"
    
    # Test SSH connectivity
    echo "🔑 Testing SSH connectivity..."
    
    # Remove old host key
    ssh-keygen -f ~/.ssh/known_hosts -R "[localhost]:2222" 2>/dev/null || true
    
    # Add new host key
    ssh-keyscan -p 2222 localhost >> ~/.ssh/known_hosts 2>/dev/null
    
    # Test connection with password
    echo "🧪 Testing SSH connection..."
    if command -v sshpass > /dev/null; then
        sshpass -p 'ubuntu123' ssh -p 2222 -o StrictHostKeyChecking=no ubuntu@localhost 'echo "✅ SSH connection successful!"'
    else
        echo "💡 Testing SSH manually (password: ubuntu123)..."
        timeout 10 ssh -p 2222 -o ConnectTimeout=5 ubuntu@localhost 'echo "✅ SSH connection successful!"' || {
            echo "⚠️  SSH test failed, but container should be running"
            echo "💡 Manual test: ssh -p 2222 ubuntu@localhost (password: ubuntu123)"
        }
    fi
    
    echo ""
    echo "🎯 Container Information:"
    echo "========================"
    docker ps | grep scl-ubuntu-test
    echo ""
    
    echo "👥 Available Users:"
    echo "   • ubuntu:ubuntu123 (sudo access)"
    echo "   • devops:devops123 (sudo access)"
    echo "   • admin:admin123 (sudo access)"
    echo "   • root:rootpass"
    echo ""
    
    echo "🌐 Port Mappings:"
    echo "   • SSH: localhost:2222 → container:22"
    echo "   • HTTP: localhost:8080 → container:80"
    echo "   • HTTPS: localhost:8443 → container:443"
    echo ""
    
    echo "🧪 Manual Testing:"
    echo "   • SSH: ssh -p 2222 ubuntu@localhost"
    echo "   • Web: curl http://localhost:8080"
    echo "   • Shell: docker exec -it scl-ubuntu-test bash"
    echo ""
    
    echo "🚀 Ready to run SCL demo!"
    echo "   ./scl examples/ubuntu_docker_demo.scl --verbose"
    echo ""
    
    echo "🛑 To stop:"
    echo "   docker-compose -f docker-ubuntu-test/docker-compose.yml down"
    
else
    echo "❌ Container failed to start. Check Docker logs:"
    echo "   docker-compose -f docker-ubuntu-test/docker-compose.yml logs"
fi