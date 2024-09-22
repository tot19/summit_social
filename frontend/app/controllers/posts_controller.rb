# app/controllers/posts_controller.rb
class PostsController < ApplicationController
  def index
    service = GoApiService.new
    begin
      @posts = service.get_posts
      Rails.logger.debug "Received posts: #{@posts.inspect}"
    rescue StandardError => e
      @error = "Failed to fetch posts from the backend: #{e.message}"
      @posts = []
    end
  end
end
