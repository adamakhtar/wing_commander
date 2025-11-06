# frozen_string_literal: true

module MinitestExample
  class Worker
    def initialize
      @helper = Helper.new
    end

    def addition(a, b, raise_error: false)
      @helper.process(a, b, raise_error: raise_error)
    end
  end
end
