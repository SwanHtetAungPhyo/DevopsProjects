#!/bin/bash

echo "🧪 SCL Language Test"
echo "==================="

# Build SCL
echo "Building SCL..."
go build -o scl
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"

# Test compile mode
echo ""
echo "Testing compile mode..."
./scl examples/demo.scl --verbose
if [ $? -eq 0 ]; then
    echo "✅ Compile mode works"
    echo "📄 Generated bash script: output.sh"
else
    echo "❌ Compile mode failed"
    exit 1
fi

# Check if Docker is available
if command -v docker &> /dev/null; then
    echo ""
    echo "Setting up Docker SSH environment..."
    ./setup-docker.sh
    
    # Wait for SSH to be ready
    sleep 3
    
    # Test SSH connection
    if ssh -p 2222 -o ConnectTimeout=5 testuser@localhost "echo 'SSH works'" &> /dev/null; then
        echo "✅ SSH connection working"
        
        # Test interpret mode
        echo ""
        echo "Testing interpret mode..."
        sed 's/mode := interpret/mode := compile/g' examples/demo.scl | sed 's/mode := compile/mode := interpret/g' > /tmp/demo-interpret.scl
        ./scl /tmp/demo-interpret.scl --verbose
        
        if [ $? -eq 0 ]; then
            echo "✅ Interpret mode works"
            
            # Verify files
            echo ""
            echo "Verifying files on remote server..."
            ssh -p 2222 testuser@localhost "ls -la /tmp/myapp/ /tmp/scripts/ 2>/dev/null | head -5"
        else
            echo "❌ Interpret mode failed"
        fi
        
        # Cleanup
        docker-compose -f docker-test/docker-compose.yml down &> /dev/null
        rm -f /tmp/demo-interpret.scl
    else
        echo "❌ SSH connection failed"
    fi
else
    echo "⚠️  Docker not available - skipping SSH tests"
fi

echo ""
echo "🎉 Test completed!"
echo ""
echo "Next steps:"
echo "  ./scl examples/demo.scl           # Test compile mode"
echo "  ./setup-docker.sh                # Setup SSH environment"
echo "  ./scl examples/demo.scl --verbose # Test interpret mode"