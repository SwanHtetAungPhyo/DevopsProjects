#!/bin/bash

echo "🐳 Setting up Docker SSH test environment for SCL..."
echo ""

# Build and start Docker container
cd docker-test
echo "Building Docker container..."
docker-compose build

echo "Starting SSH test server..."
docker-compose up -d

echo "Waiting for SSH server to start..."
sleep 5

# Setup SSH keys
echo "Setting up SSH key authentication..."

# Generate SSH key if it doesn't exist
if [ ! -f ~/.ssh/id_rsa ]; then
    echo "Generating SSH key..."
    ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""
fi

# Copy SSH key to container
docker exec scl-ssh-test mkdir -p /home/testuser/.ssh
docker cp ~/.ssh/id_rsa.pub scl-ssh-test:/home/testuser/.ssh/authorized_keys
docker exec scl-ssh-test chown testuser:testuser /home/testuser/.ssh/authorized_keys
docker exec scl-ssh-test chmod 600 /home/testuser/.ssh/authorized_keys

# Add to known_hosts
ssh-keyscan -p 2222 localhost >> ~/.ssh/known_hosts 2>/dev/null

cd ..

echo ""
echo "✅ Docker SSH test environment is ready!"
echo ""
echo "🧪 Test SSH connection:"
echo "  ssh -p 2222 testuser@localhost"
echo ""
echo "🚀 Test SCL with Docker:"
echo "  # Update target in examples/docker-ssh-test.scl to: [\"testuser@localhost:2222\"]"
echo "  ./scl examples/docker-ssh-test.scl --verbose"
echo ""
echo "🛑 Stop test environment:"
echo "  docker-compose -f docker-test/docker-compose.yml down"