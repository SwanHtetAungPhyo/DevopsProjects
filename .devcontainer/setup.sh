#!/bin/bash
set -e

echo "=== Installing Pulumi ==="
curl -fsSL https://get.pulumi.com | sh
export PATH=$PATH:$HOME/.pulumi/bin

echo "=== Installing Go dependencies ==="
cd /workspaces/$(basename "$PWD")
go mod download || echo "No go.mod found, skipping..."

echo "=== Installing additional Go tools ==="
go install github.com/rs/zerolog/log@latest || true

echo "=== Setting up permissions ==="
sudo chown -R ubuntu:ubuntu /home/ubuntu/.pulumi || true
sudo chown -R ubuntu:ubuntu /home/ubuntu/go || true

echo "=== Verifying installations ==="
echo "Go version:"
go version
echo ""
echo "Pulumi version:"
pulumi version
echo ""
echo "AWS CLI version:"
aws --version
echo ""

echo "=== Setup Complete ==="