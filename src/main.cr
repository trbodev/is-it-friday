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

Kemal.run ENV["PORT"].to_i
