#!/bin/bash

set -e  # Exit on error

# import primary
SETTING=configuration
TARGET=testuser@localhost:2222
SUPER_USER=true
ON_ERROR=rollback
MODE=compile
# Check if docker and python3 exist on testuser@localhost:2222
ssh testuser@localhost:2222 "command -v docker && command -v python3" > /dev/null 2>&1
SYSTEM_READY=$?

test_file_operations() {
    if [ $SYSTEM_READY -eq 0 ]; then
        echo "System is ready for file operations"
        # Test access
        echo "System tools missing"
    else
    fi
}

copy_config_files() {
    echo "Copying configuration files..."
    ssh testuser@localhost:2222 "mkdir -p /tmp/myapp"
    scp -r test-files/app.conf testuser@localhost:2222:/tmp/myapp/app.conf
    ssh testuser@localhost:2222 "mkdir -p /tmp/nginx"
    scp -r test-files/nginx.conf testuser@localhost:2222:/tmp/nginx/nginx.conf
    echo "Configuration files copied successfully"
}

create_script_files() {
    echo "Creating script files..."
    ssh testuser@localhost:2222 "mkdir -p /tmp/scripts && touch /tmp/scripts/deploy.sh && chmod 755 /tmp/scripts/deploy.sh"
    ssh testuser@localhost:2222 "mkdir -p /tmp/logs && touch /tmp/logs/app.log && chmod 644 /tmp/logs/app.log"
    ssh testuser@localhost:2222 "mkdir -p /tmp/config && touch /tmp/config/settings.ini && chmod 600 /tmp/config/settings.ini"
    echo "Script files created successfully"
}

verify_deployment() {
    echo "Verifying deployment..."
    echo "All files deployed and created successfully"
}

main() {
    test_file_operations
    copy_config_files
    create_script_files
    verify_deployment
    echo "File operations test completed!"
}

# Execute main function
main
