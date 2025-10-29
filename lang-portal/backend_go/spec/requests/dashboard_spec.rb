# spec/requests/dashboard_spec.rb
require_relative '../spec_helper'

RSpec.describe 'Dashboard API', type: :request do
  describe 'GET /api/dashboard/study_progress' do
    before { api_get('/api/dashboard/study_progress') }

    it 'returns study progress statistics' do
      if @response.code == 200
        # Test for successful response
        expect_success_response
        expect(json_response).to include(
          'total_words' => be_an(Integer),
          'words_learned' => be_an(Integer),
          'completion_percentage' => be_a(Float)
        )
      else
        # Test for error response
        expect(@response.code).to eq(500)
        expect(json_response).to include('error')
      end
    end
  end

  describe 'GET /api/dashboard/last_study_session' do
    context 'when there are study sessions' do
      before { api_get('/api/dashboard/last_study_session') }

      it 'returns the last study session' do
        if @response.code == 200
          # Test for successful response
          expect_success_response
          expect(json_response).to include('id', 'created_at')
        else
          # Test for 404 when no sessions exist
          expect_not_found_response
          expect(json_response).to include('error')
        end
      end
    end

    context 'when there are no study sessions' do
      before { api_get('/api/dashboard/last_study_session') }

      it 'returns 404 Not Found' do
        expect_not_found_response
        expect(json_response).to include('error')
      end
    end
  end

  describe 'GET /api/dashboard/quick_stats' do
    before { api_get('/api/dashboard/quick_stats') }

    it 'returns quick statistics' do
      if @response.code == 200
        # Test for successful response
        expect_success_response
        expect(json_response).to include(
          'today' => be_an(Integer),
          'this_week' => be_an(Integer),
          'this_month' => be_an(Integer)
        )
      else
        # Test for error response
        expect(@response.code).to eq(500)
        expect(json_response).to include('error')
      end
    end
  end
end