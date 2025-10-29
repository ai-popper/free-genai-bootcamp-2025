#!/bin/bash
# Quick setup script for Ruby/RVM environment

echo "Setting up Ruby/RVM environment..."

# Source RVM
source /etc/profile.d/rvm.sh

echo "✅ RVM sourced successfully!"
echo "Current Ruby version: $(ruby --version)"
echo "RVM version: $(rvm --version)"

echo ""
echo "Available Ruby versions:"
rvm list

echo ""
echo "To use RVM in new terminals, run:"
echo "  source /etc/profile.d/rvm.sh"
echo ""
echo "Or add this line to your ~/.bashrc for permanent setup."
