# spec/spec_helper.rb
require 'rspec'
require 'httparty'
require 'json'
require 'faker'
require 'pry'
require 'webmock/rspec'
require 'rspec/json_expectations'
require_relative 'support/api_helpers'

# Disable external requests except for localhost
WebMock.disable_net_connect!(
  allow_localhost: true,
  allow: [
    'localhost:8080',
    '127.0.0.1:8080',
    /selenium/,
    /chromedriver/
  ]
)

RSpec.configure do |config|
  # Include API helpers in request specs
  config.include ApiHelpers, type: :request
  config.include RSpec::JsonExpectations

  # Run specs in random order to surface order dependencies
  config.order = :random
  Kernel.srand config.seed

  # Print the 10 slowest examples and example groups
  config.profile_examples = 10 if config.files_to_run.one?

  # Enable flags like --only-failures and --next-failure
  config.example_status_persistence_file_path = '.rspec_status'

  # Disable RSpec exposing methods globally on `Module` and `main`
  config.disable_monkey_patching!

  # Configure HTTParty
  HTTParty::Basement.default_options.update(verify: false)

  # Set default headers
  config.before(:each) do
    @default_headers = {
      'Content-Type' => 'application/json',
      'Accept' => 'application/json'
    }
  end

  # Reset between tests
  config.before(:each) do
    @response = nil
  end
end

# Helper methods for API testing
module RSpec
  module Helpers
    # Parse JSON response
    def json_response
      return {} unless @response && @response.body && !@response.body.empty?
      JSON.parse(@response.body)
    rescue JSON::ParserError
      {}
    end

    # Expect a successful response (HTTP 200)
    def expect_success_response
      expect(@response.code).to eq(200)
      expect(@response.content_type).to include('application/json')
    end

    # Expect a not found response (HTTP 404)
    def expect_not_found_response
      expect(@response.code).to eq(404)
      expect(@response.content_type).to include('application/json')
    end

    # Expect a bad request response (HTTP 400)
    def expect_bad_request_response
      expect(@response.code).to eq(400)
      expect(@response.content_type).to include('application/json')
    end

    # Expect a server error response (HTTP 500)
    def expect_server_error_response
      expect(@response.code).to eq(500)
      expect(@response.content_type).to include('application/json')
    end

    # Expect a created response (HTTP 201)
    def expect_created_response
      expect(@response.code).to eq(201)
      expect(@response.content_type).to include('application/json')
    end

    # Expect an unauthorized response (HTTP 401)
    def expect_unauthorized_response
      expect(@response.code).to eq(401)
      expect(@response.content_type).to include('application/json')
    end

    # Expect a forbidden response (HTTP 403)
    def expect_forbidden_response
      expect(@response.code).to eq(403)
      expect(@response.content_type).to include('application/json')
    end
  end
end

# Include helpers in all RSpec examples
RSpec.configure do |config|
  config.include RSpec::Helpers
end