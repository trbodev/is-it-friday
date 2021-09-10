require "kemal"

ENV["PORT"] ||= "3000"

error 404 do |ctx|
  ctx.response.content_type = "text/plain"
  "Not Found"
end

get "/" do |ctx|
  ctx.response.content_type = "text/plain"
  Time.utc.friday? ? "Yes" : "No"
end

Kemal.config.powered_by_header = false

logging Kemal.config.env === "development"

Kemal.run do |config|
  config.powered_by_header = false
  config.server.not_nil!.bind_tcp(ENV["HOST"], ENV["PORT"].to_i)
end
