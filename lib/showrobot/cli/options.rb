# Define various command line option sets here

BASIC_USAGE = <<EOS
ShowRobot is a utility for matching media files against
TV and movie databases to categorize your media collection.
EOS

OptionSet = {}
OptionSet[:config] = Proc.new do
	# Use the specified configuration file
	opt :config_file, "Use the specified configuration file for ShowRobot", :type => String, :short => :f

	# Causes ShowRobot to display additional debugging information
	opt :verbose, "Show detailed information messages", :default => false
end

OptionSet[:cache] = Proc.new do
	# Cache control
	opt :cache_dir, "Use the specified directory for caching server queries", :type => String, :short => :c
	opt :no_cache, "Do not use caching for any queries", :default => false, :short => nil
	opt :clear_cache, "Clear the cache and start fresh with queries", :default => false, :short => :C
end

OptionSet[:database] = Proc.new do
	# Media database identification options
	opt :movie_database, "Use the specified database for movie matching", :type => String, :short => nil
	opt :tv_database, "Use the specified database for TV show matching", :type => String, :short => nil

	# prompts user to choose amongst various options instead of deferring to the first
	opt :prompt, "Prompt for media file matches", :default => false, :short => :p
end

OptionSet[:file] = Proc.new do
	# File type control
	opt :force_movie, "Force the specified file(s) to be matched against the movie database", :default => false, :short => :M
	opt :force_tv, "Force the specified file(s) to be matched against the tv database", :default => false, :short => :T
end

def apply_options options={}
	require 'showrobot'

	options.each do |key, value|
		# if the value is nil, skip
		next if not options.include?((key.to_s + '_given').to_sym) or value.nil?

		case key
		when :config_file
			ShowRobot.configure :config_file => options[:config_file]
			ShowRobot.load_config
		when :cache_dir, :tv_database, :movie_database, :verbose
			ShowRobot.configure key, value
		when :no_cache
			ShowRobot.configure :use_cache, false
		when :clear_cache
			puts "Clearing cache in [ #{ShowRobot.config[:cache_dir]} ]" if ShowRobot.config[:verbose]
			File.delete(*Dir[ShowRobot.config[:cache_dir] + '/*.cache'])
		end
	end
end
