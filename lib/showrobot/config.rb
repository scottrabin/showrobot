module ShowRobot
	require 'yaml'
	require 'pathname'

	class << self
		attr_accessor :configuration

		def config
			self.configuration ||= Configuration.new
		end

		def configure
			yield config
		end
	end

	class Configuration

		# base directory against which other files are loaded (i.e. cache/config/logs)
		attr_accessor :basepath

		# verbosity of output
		attr_accessor :verbose

		# whether or not to use the cache store for data fetching
		attr_accessor :use_cache

		# where the cache directory is located
		attr_accessor :cache_dir

		# API key for The TVDB
		attr_accessor :tvdb_api_key

		# API key for The Movie DB
		attr_accessor :tmdb_api_key

		# Default database for TV shows
		attr_accessor :tv_database

		# Default database for movies
		attr_accessor :movie_database

		def initialize
			# provide a default basepath
			@basepath = Pathname.new(ENV['HOME'] + '/.showrobot')

			# load the default configuration file to provide default values for the above
			load_from(File.dirname(__FILE__) + '/../../config.yml')
		end

		# load a set of configuration options from a file
		def load_from file
			YAML::load(IO.read(file)).each do |key, value|
				send(key + '=', value)
			end
		rescue Errno::ENOENT
			puts :warning, "No configuration file found at [ #{file} ]."
		rescue Psych::SyntaxError
			puts :warning, "YAML configuration file at [ #{file} ] contains invalid syntax."
		end
	end

end
