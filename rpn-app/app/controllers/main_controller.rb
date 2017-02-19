require 'HTTParty'

class MainController < ApplicationController
    def index; end

    def parse
        @apiUrl = 'http://127.0.0.1:5000/parse'
        @data = {}

        @data[:count] = params[:input].lines.first.gsub("\r\n", '')
        @data[:expressions] = []

        lineno = 1

        params[:input].each_line do |line|
            @data[:expressions].push(value: line.gsub("\r\n", '')) if lineno > 1
            lineno += 1
        end

        request = HTTParty.post(@apiUrl, headers: { 'Content-Type' => 'application/json' }, body: @data.to_json)

        @response = request.parsed_response

        @responseString = ''

        @response['results'].each do |result|
            @responseString << result['value'].to_s << ', ' << result['time'].to_s << "\r\n"
        end
    end
end
