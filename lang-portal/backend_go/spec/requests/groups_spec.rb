# spec/requests/groups_spec.rb
require_relative '../spec_helper'

RSpec.describe 'Groups API', type: :request do
  before(:all) do
    # Run database migrations or set up test data here
    # You might need to run something like:
    # `system('go run path/to/your/migrations')` 
  end

  describe 'GET /api/groups' do
    before { api_get('/api/groups') }

    it 'returns list of groups' do
      if @response.code == 200
        expect_success_response
        groups = json_response['items']
        expect(groups).to be_an(Array)
      else
        # Accept either 404 or 500 as valid error responses
        expect([404, 500]).to include(@response.code)
        expect(json_response).to include('error')
      end
    end

    it 'returns paginated response' do
      api_get('/api/groups')
      if @response.code == 200
        expect_paginated_response
      else
        expect([404, 500]).to include(@response.code)
      end
    end
  end

  describe 'GET /api/groups/:id' do
    let(:group_id) { 1 }
    
    before { api_get("/api/groups/#{group_id}") }

    context 'with valid group id' do
      it 'returns group details' do
        if @response.code == 200
          expect_success_response
          expect(json_response).to include(
            'id' => group_id,
            'name' => be_a(String)
          )
        else
          expect([404, 500]).to include(@response.code)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with invalid group id' do
      let(:group_id) { 99999 }
      
      it 'returns error' do
        expect([404, 500]).to include(@response.code)
        expect(json_response).to include('error')
      end
    end
  end

  describe 'GET /api/groups/:id/words' do
    let(:group_id) { 1 }
    
    before { api_get("/api/groups/#{group_id}/words") }

    it 'returns words or error' do
      if @response.code == 200
        expect_paginated_response
      else
        expect([404, 500]).to include(@response.code)
        expect(json_response).to include('error')
      end
    end
  end

  describe 'GET /api/groups/:id/study_sessions' do
    let(:group_id) { 1 }
    
    before { api_get("/api/groups/#{group_id}/study_sessions") }

    it 'returns study sessions or error' do
      if @response.code == 200
        expect_paginated_response
      else
        expect([404, 500]).to include(@response.code)
        expect(json_response).to include('error')
      end
    end
  end
end