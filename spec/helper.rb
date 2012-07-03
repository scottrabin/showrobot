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
	CMD_RENAME_TO = /--> \[ (.*?) \]/

	module Extend
		def with_file fileName, &block
			describe "with file #{fileName}" do
				subject do
					ShowRobot::MediaFile.load fileName
				end

				instance_eval &block
			end
		end
	end

	module Include
		def cli cmd
			r = `#{File.expand_path(File.dirname(__FILE__) + '/../bin/showrobot')} #{cmd}`
			$?.should == 0
			r
		end
	end
end

RSpec.configure do |c|
	c.extend ShowRobotHelper::Extend
	c.include ShowRobotHelper::Include
end
