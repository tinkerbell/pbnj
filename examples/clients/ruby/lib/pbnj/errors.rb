module PBnJ
  class Error < RuntimeError; end

  class ResponseError < Error
    attr_reader :response

    def initialize(response)
      @response = response
      super(message)
    end

    def message
      "PBnJ Error { status: #{response&.try(:status)}, body: #{response&.try(:body)} }"
    end
  end
end
