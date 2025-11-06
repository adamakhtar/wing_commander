# frozen_string_literal: true

module MinitestExample
  class Helper
    def process(a, b, raise_error:)
      if raise_error
        raise "Error in Helper#process"
      else
        a + b
      end
    end
  end
end
