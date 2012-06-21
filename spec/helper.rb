require 'rubygems'
require 'bundler'
begin
  Bundler.setup(:default, :development)
rescue Bundler::BundlerError => e
  $stderr.puts e.message
  $stderr.puts "Run `bundle install` to install missing gems"
  exit e.status_code
end

require 'showrobot'

module ShowRobotHelper
	def with_file fileName, &block
		describe "with file #{fileName}" do
			subject do
				ShowRobot::MediaFile.load fileName
			end

			instance_eval &block
		end
	end
end

RSpec.configure do |c|
	c.extend ShowRobotHelper
end
