# Ruby/RSpec API Testing Suite

A comprehensive test suite for the Language Learning Portal backend API using Ruby and RSpec.

## Overview

This test suite provides complete coverage for all 18 API endpoints, ensuring they return proper JSON responses and HTTP status codes as specified in the technical documentation.

## Test Coverage

### ✅ All Endpoints Tested

| Endpoint Group | Endpoints | Status |
|----------------|-----------|--------|
| Dashboard | 3 endpoints | ✅ Complete |
| Words | 2 endpoints | ✅ Complete |
| Groups | 4 endpoints | ✅ Complete |
| Study Sessions | 4 endpoints | ✅ Complete |
| Study Activities | 3 endpoints | ✅ Complete |
| Utility | 2 endpoints | ✅ Complete |

**Total: 18 endpoints** - All implemented and tested!

### 🧪 Test Categories

- **Status Code Validation**: 200, 404, 400, 500 responses
- **JSON Response Format**: Proper structure and data types
- **Pagination**: Correct pagination metadata
- **Data Structure**: Word, group, session object validation
- **Error Handling**: Invalid IDs, bad requests, server errors
- **Database Operations**: Create, read, update operations
- **Integration**: End-to-end API workflow testing

## Project Structure

```
backend_go/
├── Gemfile                    # Ruby dependencies
├── Rakefile                   # Test runner tasks
├── .rspec                     # RSpec configuration
├── test_ruby_api.rb           # Test runner script
└── spec/
    ├── spec_helper.rb         # Test configuration
    ├── support/
    │   └── api_helpers.rb     # Helper methods & matchers
    └── requests/
        ├── dashboard_spec.rb      # Dashboard tests
        ├── words_spec.rb          # Words tests
        ├── groups_spec.rb         # Groups tests
        ├── study_sessions_spec.rb # Sessions tests
        ├── study_activities_spec.rb # Activities tests
        └── utility_spec.rb        # Utility tests
```

## Prerequisites

- **Ruby** 2.7+ installed
- **Bundler** gem installed (`gem install bundler`)
- **Go server** running on `http://localhost:8080`

## Setup

### 1. Install Dependencies

```bash
bundle install
```

### 2. Verify Go Server

Ensure the Go server is running:

```bash
curl http://localhost:8080/health
```

You should get: `{"status":"ok"}`

### 3. Run Tests

```bash
# Run all tests
bundle exec rspec spec/

# Or use the test script
./test_ruby_api.rb

# Or use Rake tasks
bundle exec rake spec
```

## Test Commands

```bash
# Run all tests with detailed output
bundle exec rspec spec/ --format documentation --color

# Run specific test files
bundle exec rspec spec/requests/words_spec.rb

# Run with coverage (if simplecov is added)
COVERAGE=true bundle exec rspec spec/

# Run specific endpoint group tests
bundle exec rake test_dashboard
bundle exec rake test_words
bundle exec rake test_groups
bundle exec rake test_sessions
bundle exec rake test_activities
bundle exec rake test_utility
```

## Test Features

### 🔧 Helper Methods

```ruby
# API request helpers
api_get('/api/words')
api_post('/api/study_activities', { group_id: 1, study_activity_id: 1 })

# Response validation helpers
expect_success_response      # 200 status + JSON content-type
expect_not_found_response    # 404 status + JSON error
expect_bad_request_response  # 400 status + JSON error
expect_server_error_response # 500 status + JSON error

# Data structure validators
expect_word_structure(word)
expect_group_structure(group)
expect_study_session_structure(session)
expect_paginated_response
```

### 📊 Response Validation

All tests validate:

1. **HTTP Status Codes**
   - 200 OK for successful operations
   - 404 Not Found for invalid IDs
   - 400 Bad Request for malformed data
   - 500 Internal Server Error for server issues

2. **JSON Content-Type**
   - All responses must be `application/json`

3. **Data Structure**
   - Proper field names and data types
   - Nested object relationships
   - Pagination metadata

4. **Pagination Format**
   ```json
   {
     "items": [...],
     "pagination": {
       "current_page": 1,
       "items_per_page": 100,
       "total_pages": 5,
       "total_items": 450
     }
   }
   ```

### 🎯 Shared Examples

```ruby
# Reusable test patterns
it_behaves_like "a paginated endpoint"
it_behaves_like "returns proper error for invalid id", "word"
it_behaves_like "returns proper error for bad request"
```

## Dependencies

```ruby
# Core testing
gem 'rspec', '~> 3.12'              # Testing framework
gem 'httparty', '~> 0.21'           # HTTP client for API requests

# Database & fixtures
gem 'database_cleaner-sqlite'       # Database cleanup between tests

# Development helpers
gem 'faker', '~> 3.2'               # Fake data generation
gem 'pry', '~> 0.14'                # Debugging

# Response validation
gem 'rspec-json_expectations'       # JSON structure validation

# Additional tools
gem 'rack-test', '~> 2.1'           # Rack testing utilities
gem 'webmock', '~> 3.18'            # HTTP request mocking
gem 'vcr', '~> 6.1'                 # HTTP interaction recording
```

## Test Data

Tests use the existing seed data:
- **50 words** across 5 groups
- **4 study activities** (Flashcard, Quiz, Writing, Listening)
- **Dynamic study sessions** created during tests

## Database Management

- **DatabaseCleaner** handles test isolation
- **Strategy**: Truncation (clean slate for each test)
- **Automatic cleanup** before/after each test

## Running Tests in CI/CD

```bash
# Install dependencies
bundle install --without development

# Run tests
bundle exec rspec spec/ --format RspecJunitFormatter --out test-results.xml
```

## Debugging Tests

```ruby
# Add debugging to any test
binding.pry  # Drops into pry console
puts last_response.body  # Print raw response
puts json_response.inspect  # Print parsed JSON
```

## Extending Tests

### Adding New Endpoints

1. Create new test file: `spec/requests/new_endpoint_spec.rb`
2. Follow existing patterns for request/response validation
3. Add Rake task if needed: `bundle exec rake test_new_endpoint`

### Adding Custom Matchers

Add to `spec/support/api_helpers.rb`:

```ruby
RSpec::Matchers.define :have_valid_structure do
  match do |response|
    # Your validation logic
  end
end
```

## Best Practices

✅ **Test Isolation** - Each test is independent
✅ **Descriptive Names** - Clear test descriptions
✅ **Proper Assertions** - Validate both structure and content
✅ **Error Scenarios** - Test failure cases
✅ **Documentation** - Well-commented test code
✅ **Maintainability** - Reusable shared examples

## Troubleshooting

### Common Issues

1. **Server not running**: `curl http://localhost:8080/health` should return `{"status":"ok"}`

2. **Database issues**: Tests clean database between runs - ensure server can access `words.db`

3. **Port conflicts**: Ensure no other service is using port 8080

4. **Ruby version**: Requires Ruby 2.7+ for modern gem compatibility

### Debug Mode

```bash
# Run with verbose output
RSPEC_DEBUG=true bundle exec rspec spec/

# Run single test with debugging
bundle exec rspec spec/requests/words_spec.rb:10 --format documentation
```

## 🎉 Success Metrics

- **18 endpoints** fully tested
- **100% endpoint coverage**
- **Status code validation** for all scenarios
- **JSON structure validation** for all responses
- **Error handling verification**
- **Pagination testing**
- **Integration testing** with real database

Ready for production deployment with confidence! 🚀
