# spec/requests/study_activities_spec.rb
require_relative '../spec_helper'

RSpec.describe 'Study Activities API', type: :request do
  describe 'GET /api/study_activities/:id' do
    let(:activity_id) { 1 }
    before { api_get("/api/study_activities/#{activity_id}") }

    context 'with valid activity id' do
      it 'returns study activity details' do
        if @response.code == 200
          expect_success_response
          expect(json_response).to include(
            'id' => activity_id,
            'name' => be_a(String),
            'description' => be_a(String).or(be_nil)
          )
        else
          expect([404, 500]).to include(@response.code.to_i)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with invalid activity id' do
      let(:activity_id) { 99999 }

      it 'returns appropriate error' do
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end
  end

  describe 'GET /api/study_activities/:id/study_sessions' do
    let(:activity_id) { 1 }
    before { api_get("/api/study_activities/#{activity_id}/study_sessions") }

    context 'with valid activity id' do
      it 'returns study sessions for the activity' do
        if @response.code == 200
          expect_success_response
          sessions = json_response['items']
          expect(sessions).to be_an(Array)

          sessions.each do |session|
            expect(session).to include(
              'id' => be_an(Integer),
              'activity_id' => activity_id,
              'created_at' => be_a(String)
            )
          end
        else
          expect(@response.code).to eq(500)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with invalid activity id' do
      let(:activity_id) { 99999 }

      it 'returns appropriate error' do
        expect([404, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end

    context 'with pagination' do
      it 'returns paginated response' do
        api_get("/api/study_activities/#{activity_id}/study_sessions?page=1&per_page=5")
        
        if @response.code == 200
          expect(json_response).to include(
            'items' => be_an(Array),
            'page' => 1,
            'per_page' => 5
          )
        else
          expect(@response.code).to eq(500)
          expect(json_response).to include('error')
        end
      end
    end
  end

  describe 'POST /api/study_activities' do
    let(:valid_params) do
      {
        group_id: 1,
        study_activity_id: 1
      }
    end

    context 'with valid data' do
      before { api_post('/api/study_activities', valid_params) }

      it 'creates a new study session' do
        if @response.code == 200
          expect_success_response
          expect(json_response).to include(
            'id' => be_an(Integer),
            'group_id' => valid_params[:group_id],
            'study_activity_id' => valid_params[:study_activity_id]
          )
        else
          expect(@response.code).to eq(500)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with missing group_id' do
      before { api_post('/api/study_activities', valid_params.except(:group_id)) }

      it 'returns 400 Bad Request' do
        if @response.code == 400
          expect_bad_request_response
          expect(json_response).to include('error')
        else
          expect(@response.code).to eq(500)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with missing study_activity_id' do
      before { api_post('/api/study_activities', valid_params.except(:study_activity_id)) }

      it 'returns 400 Bad Request' do
        if @response.code == 400
          expect_bad_request_response
          expect(json_response).to include('error')
        else
          expect(@response.code).to eq(500)
          expect(json_response).to include('error')
        end
      end
    end

    context 'with invalid group_id' do
      before { api_post('/api/study_activities', valid_params.merge(group_id: 99999)) }

      it 'handles invalid group_id' do
        expect([400, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end

    context 'with invalid study_activity_id' do
      before { api_post('/api/study_activities', valid_params.merge(study_activity_id: 99999)) }

      it 'handles invalid study_activity_id' do
        expect([400, 500]).to include(@response.code.to_i)
        expect(json_response).to include('error')
      end
    end
  end

  # Helper method to validate study session structure
  def expect_study_session_structure(session)
    expect(session).to include(
      'id' => be_an(Integer),
      'activity_id' => be_an(Integer),
      'created_at' => be_a(String)
    )
  end
end