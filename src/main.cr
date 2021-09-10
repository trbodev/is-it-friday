require "kemal"
require "./pages"

ENV["HOST"] ||= "0.0.0.0"
ENV["PORT"] ||= "3000"

error 404 do |ctx|
  ctx.response.content_type = "text/plain"
  "Not Found"
end

get "/" do |ctx|
  Pages::Home.new.to_s
end

get "/plain" do |ctx|
  ctx.response.content_type = "text/plain"
  ctx.response.headers["Content-Disposition"] = "filename=\"friday.txt\""
  Time.utc.friday? ? "Yes" : "No"
end

get "/boolean" do |ctx|
  ctx.response.content_type = "text/plain"
  Time.utc.friday? ? "true" : "false"
end

get "/json" do |ctx|
  ctx.response.content_type = "application/json"
  ctx.response.headers["Content-Disposition"] = "filename=\"friday.json\""
  Time.utc.friday? ? "true" : "false"
end

get "/yaml" do |ctx|
  ctx.response.content_type = "text/yaml"
  ctx.response.headers["Content-Disposition"] = "attachment; filename=\"friday.yaml\""
  "friday: #{Time.utc.friday? ? "true" : "false"}"
end

get "/xml" do |ctx|
  ctx.response.content_type = "text/xml"
  ctx.response.headers["Content-Disposition"] = "filename=\"friday.xml\""
  "<?xml version=\"1.0\" encoding=\"UTF-8\"?><friday>#{Time.utc.friday? ? "true" : "false"}</friday>"
end

get "/binary" do |ctx|
  ctx.response.content_type = "text/plain"
  Time.utc.friday? ? "1" : "0"
end

Kemal.config.powered_by_header = false

logging Kemal.config.env === "development"

Kemal.run do |config|
  config.powered_by_header = false
  config.server.not_nil!.bind_tcp(ENV["HOST"], ENV["PORT"].to_i)
end
