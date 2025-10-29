require 'httparty'
require 'json'

module ApiHelpers
  BASE_URL = 'http://localhost:8080'  

  def api_get(path, headers = {})
    @response = HTTParty.get("#{BASE_URL}#{path}", headers: headers)
  end

  def api_post(path, body = {}, headers = {})
    headers['Content-Type'] ||= 'application/json'
    @response = HTTParty.post("#{BASE_URL}#{path}", body: body.to_json, headers: headers)
  end

  def api_put(path, body = {}, headers = {})
    headers['Content-Type'] ||= 'application/json'
    @response = HTTParty.put("#{BASE_URL}#{path}", body: body.to_json, headers: headers)
  end

  def api_delete(path, headers = {})
    @response = HTTParty.delete("#{BASE_URL}#{path}", headers: headers)
  end

  def expect_paginated_response
    expect(json_response).to have_key('items')
    expect(json_response).to have_key('pagination')
    expect(json_response['pagination']).to have_key('current_page')
    expect(json_response['pagination']).to have_key('items_per_page')
    expect(json_response['pagination']).to have_key('total_pages')
    expect(json_response['pagination']).to have_key('total_items')
  end

  def expect_word_structure(word)
    expect(word).to have_key('id')
    expect(word).to have_key('japanese')
    expect(word).to have_key('romaji')
    expect(word).to have_key('english')
    expect(word).to have_key('correct_count')
    expect(word).to have_key('wrong_count')
  end

  def expect_group_structure(group)
    expect(group).to have_key('id')
    expect(group).to have_key('name')
    expect(group).to have_key('word_count')
  end

  def expect_study_session_structure(session)
    expect(session).to have_key('id')
    expect(session).to have_key('activity_name')
    expect(session).to have_key('group_name')
    expect(session).to have_key('start_time')
    expect(session).to have_key('review_items_count')
  end

  def expect_dashboard_progress_structure(progress)
    expect(progress).to have_key('total_words_studied')
    expect(progress).to have_key('total_words_available')
  end

  def expect_dashboard_stats_structure(stats)
    expect(stats).to have_key('success_rate')
    expect(stats).to have_key('total_study_sessions')
    expect(stats).to have_key('total_active_groups')
    expect(stats).to have_key('study_streak_days')
  end
end

RSpec.shared_examples 'a paginated endpoint' do |path, method = :api_get, params = {}|
  it 'returns paginated response' do
    send(method, "#{path}?#{params.to_query}")
    expect_success_response
    expect_paginated_response
  end
end

def expect_paginated_response
  expect(json_response).to have_key('items')
  expect(json_response).to have_key('pagination')
  expect(json_response['pagination']).to include(
    'current_page' => be_a(Integer),
    'items_per_page' => be_a(Integer),
    'total_pages' => be_a(Integer),
    'total_items' => be_a(Integer)
  )
end

RSpec.shared_examples "returns proper error for invalid id" do |endpoint_type|
  context "with invalid id" do
    let(:invalid_id) { 99999 }

    it "returns 404 Not Found" do
      subject
      expect_not_found_response
      expect(json_response).to have_key('error')
    end
  end
end

RSpec.shared_examples "returns proper error for bad request" do
  it "returns 400 Bad Request" do
    subject
    expect_bad_request_response
    expect(json_response).to have_key('error')
  end
end

RSpec.configure do |config|
  config.include ApiHelpers
end
