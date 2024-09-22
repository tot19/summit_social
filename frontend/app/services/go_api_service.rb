# app/services/go_api_service.rb
require "httparty"

class GoApiService
  include HTTParty
  base_uri ENV["GO_BACKEND_URL"]

  def get_posts
    response = self.class.get("/posts")
    if response.success?
      JSON.parse(response.body)
    else
      raise "API request failed with status #{response.code}"
    end
  end

  def get_post(id)
    response = self.class.get("/api/posts/#{id}")
    if response.success?
      JSON.parse(response.body)
    else
      raise "API request failed with status #{response.code}"
    end
  end
end
