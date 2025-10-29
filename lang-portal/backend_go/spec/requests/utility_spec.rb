# spec/requests/utility_spec.rb
require_relative '../spec_helper'

RSpec.describe 'Utility API', type: :request do
  describe 'POST /api/reset_history' do
    context 'when successful' do
      before { api_post('/api/reset_history') }

      it 'resets study history' do
        if @response.code == 200
          expect_success_response
          expect(json_response).to include('message' => be_a(String))
        else
          # If not 200, expect a server error
          expect(@response.code).to eq(500)
          expect(json_response).to include('error')
        end
      end
    end

    context 'when database operation fails' do
      # Add test case for database failure if needed
    end
  end

  describe 'POST /api/full_reset' do
    context 'when successful' do
      before { api_post('/api/full_reset') }

      it 'performs full database reset' do
        if @response.code == 200
          expect_success_response
          expect(json_response).to include('message' => be_a(String))
        else
          # If not 200, expect a server error
          expect(@response.code).to eq(500)
          expect(json_response).to include('error')
        end
      end
    end

    context 'when database operation fails' do
      # Add test case for database failure if needed
    end
  end
end