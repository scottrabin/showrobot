module ShowRobot

	require 'pathname'
	require 'fileutils'

	# ShowRobot constants

	# Determines if ShowRobot will attempt to use the query cache
	USE_CACHE = true unless defined?(::ShowRobot::USE_CACHE)
	# Determines the directory ShowRobot will use for the query cache
	if ::ShowRobot::USE_CACHE
		if not defined?(::ShowRobot::CACHE_DIRECTORY)
			CACHE_DIRECTORY = [File.dirname(__FILE__) + '/..', ENV['HOME'], '.'].find { |path| File.writable? path } + '/.showrobot_cache' rescue nil
		end

		if CACHE_DIRECTORY.nil?
			$stderr.puts "Cannot write to any default directory for ShowRobot cache"
		else
			FileUtils.mkdir_p CACHE_DIRECTORY
		end
	end

	class << self
		# verbosity of output
		attr_accessor :verbose

	end

end

require 'showrobot/media_file'
require 'showrobot/db'
(
	# load all utility files
	Dir[File.dirname(__FILE__) + '/showrobot/utility/*.rb'] +
	# load all media file parsers
	Dir[File.dirname(__FILE__) + '/showrobot/video/*.rb'] +
	# load all database shims
	Dir[File.dirname(__FILE__) + '/showrobot/db/*.rb']
).each { |file| require file }
