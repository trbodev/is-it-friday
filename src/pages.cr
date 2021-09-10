require "ecr"

module Pages
  class Home
    def initialize()
    end
    ECR.def_to_s("views/index.html")
  end
end
