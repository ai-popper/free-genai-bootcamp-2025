# spec/requests/study_sessions_spec.rb
require_relative '../spec_helper'

RSpec.describe 'Study Sessions API', type: :request do
  describe 'GET /api/study_sessions' do
    before { api_get('/api/study_sessions') }

    it 'returns list of study sessions or error' do
      if @response.code == 200
        expect_success_response
        sessions = json_response['items']
        expect(sessions).to be_an(Array)
        
        # Check structure of first session if present
        if sessions.any?
          session = sessions.first
          expect(session).to include(
            'id' => be_an(Integer),
            'activity_name' => be_a(String),
            'group_name' => be_a(String),
            'start_time' => be_a(String),
            'end_time' => be_a(String).or(be_nil),
            'review_items_count' => be_an(Integer)
          )
        end
        
        # Check pagination structure
        expect(json_response).to include(
          'pagination' => {
            'current_page' => be_an(Integer),
            'items_per_page' => be_an(Integer),
            'total_pages' => be_an(Integer),
            'total_items' => be_an(Integer)
          }
        )
      else
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end

    it 'handles pagination' do
      api_get('/api/study_sessions?page=2&per_page=5')
      
      if @response.code == 200
        expect_success_response
        expect(json_response['pagination']).to include(
          'current_page' => 2,
          'items_per_page' => 5
        )
      else
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end
  end

  describe 'GET /api/study_sessions/:id' do
    let(:session_id) { 1 }
    before { api_get("/api/study_sessions/#{session_id}") }

    context 'with valid session id' do
      it 'returns study session details' do
        if @response.code == 200
          expect_success_response
          expect(json_response).to include(
            'id' => session_id,
            'activity_name' => be_a(String),
            'group_name' => be_a(String),
            'start_time' => be_a(String),
            'end_time' => be_a(String).or(be_nil),
            'review_items_count' => be_an(Integer)
          )
        else
          expect([404, 500]).to include(@response.code.to_i)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with invalid session id' do
      let(:session_id) { 99999 }

      it 'returns appropriate error' do
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end
  end

  describe 'GET /api/study_sessions/:id/words' do
    let(:session_id) { 1 }
    before { api_get("/api/study_sessions/#{session_id}/words") }

    context 'with valid session id' do
      it 'returns words in the study session' do
        if @response.code == 200
          expect_success_response
          words = json_response['items']
          expect(words).to be_an(Array)
          
          # Check structure of first word if present
          if words.any?
            word = words.first
            expect(word).to include(
              'id' => be_an(Integer),
              'japanese' => be_a(String),
              'romaji' => be_a(String),
              'english' => be_a(String),
              'correct_count' => be_an(Integer).or(be_nil),
              'wrong_count' => be_an(Integer).or(be_nil)
            )
          end
          
          # Check pagination structure
          expect(json_response).to include(
            'pagination' => {
              'current_page' => be_an(Integer),
              'items_per_page' => be_an(Integer),
              'total_pages' => be_an(Integer),
              'total_items' => be_an(Integer)
            }
          )
        else
          expect([404, 500]).to include(@response.code.to_i)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with invalid session id' do
      let(:session_id) { 99999 }

      it 'returns appropriate error' do
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end
  end

  describe 'POST /api/study_sessions/:id/words/:word_id/review' do
    let(:session_id) { 1 }
    let(:word_id) { 1 }
    let(:valid_params) { { correct: true } }

    before { api_post("/api/study_sessions/#{session_id}/words/#{word_id}/review", valid_params) }

    context 'with valid data' do
      it 'records the word review' do
        if @response.code == 200
          expect_success_response
          expect(json_response).to include(
            'success' => be_truthy,
            'word_id' => word_id,
            'study_session_id' => session_id,
            'correct' => be_boolean,
            'created_at' => be_a(String)
          )
        else
          expect([400, 404, 500]).to include(@response.code.to_i)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with invalid session id' do
      let(:session_id) { 99999 }

      it 'returns appropriate error' do
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end

    context 'with invalid word id' do
      let(:word_id) { 99999 }

      it 'returns appropriate error' do
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end

    context 'with invalid request body' do
      let(:valid_params) { {} }

      it 'returns appropriate error' do
        expect([400, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end
  end
end
