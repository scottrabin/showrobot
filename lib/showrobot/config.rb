module ShowRobot

	require 'yaml'

	# configuration defaults
	@config = {
		# base directory against which other files are loaded (i.e. cache/config/logs)
		:basepath  => ENV['HOME'] + '/.showrobot',
		# verbosity of output
		:verbose   => false,
		# whether or not to use the cache store for data fetching
		:use_cache => true,
		# where the cache directory is located
		:cache_dir => ENV['HOME'] + '/.showrobot/cache',
		# default configuration file
		:config_file => ENV['HOME'] + '/.showrobot/config.yml'
	}
	def self.config
		@config
	end

	# Configure via hash
	def self.configure(opts = {})
		opts.each { |k, v| @config[k.to_sym] = v } if not opts.nil?
	end

	def self.load_config(file = @config[:config_file])
		begin
			config = YAML::load(IO.read(file))
		rescue Errno::ENOENT
			puts :warning, "YAML configuration file could not be found. Using defaults."
		rescue Psych::SyntaxError
			puts :warning, "YAML configuration file contains invalid syntax. Using defaults"
		end

		configure config
	end

	# load the default configuration file
	load_config(File.dirname(__FILE__) + '/../../config.yml')
	# and the environment configuration file
	load_config

end
