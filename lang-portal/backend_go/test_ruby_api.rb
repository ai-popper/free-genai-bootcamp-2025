#!/bin/bash

# Test script for Language Learning Portal API
# This script sets up Ruby environment and runs RSpec tests

set -e

echo "=== Language Learning Portal API Tests ==="
echo ""

# Check if Ruby is installed
if ! command -v ruby &> /dev/null; then
    echo "❌ Ruby is not installed. Please install Ruby first."
    exit 1
fi

# Check if bundler is installed
if ! command -v bundle &> /dev/null; then
    echo "Installing bundler..."
    gem install bundler
fi

# Install dependencies
echo "📦 Installing Ruby dependencies..."
bundle install

# Check if Go server is running
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Go server is running on port 8080"
else
    echo "❌ Go server is not running. Please start it first:"
    echo "   cd backend_go"
    echo "   go run cmd/server/main.go"
    echo ""
    echo "Then run this test script again."
    exit 1
fi

echo ""
echo "🧪 Running RSpec tests..."
echo ""

# Run all tests
bundle exec rspec spec/ --format documentation --color

echo ""
echo "=== Tests Complete ==="

# Check if tests passed
if [ $? -eq 0 ]; then
    echo "✅ All tests passed!"
    echo ""
    echo "🎉 API testing setup is complete and working!"
    echo ""
    echo "Available commands:"
    echo "  bundle exec rspec spec/                    # Run all tests"
    echo "  bundle exec rspec spec/requests/words_spec.rb    # Run specific test file"
    echo "  bundle exec rake test_dashboard            # Run dashboard tests"
    echo "  bundle exec rake test_words                # Run words tests"
    echo "  bundle exec rake test_groups               # Run groups tests"
    echo "  bundle exec rake test_sessions             # Run sessions tests"
    echo "  bundle exec rake test_activities           # Run activities tests"
    echo "  bundle exec rake test_utility              # Run utility tests"
else
    echo "❌ Some tests failed. Please check the output above."
    exit 1
fi
