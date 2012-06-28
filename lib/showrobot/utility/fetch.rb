module ShowRobot

	# for processing xml data
	require 'xml'
	# for fetching data
	require 'open-uri'
	# extra file utilities
	require 'fileutils'

	# Fetches the given site and processes it as the given type
	def fetch type, url
		# determine the location of the cache file
		cache_file = File.join(ShowRobot.config[:cache_dir], url.gsub(/([^a-zA-Z0-9_\.-]+)/) { '%' + $1.unpack('H2' * $1.size).join('%').upcase }.tr(' ', '+') + '.cache')

		# if USE_CACHE is true, attempt to find the file out of the cache
		contents = if ShowRobot.config[:use_cache] && File.exists?(cache_file)
					   puts "Found cache entry for [ #{url} ]" if ShowRobot.config[:verbose]
					   File.read cache_file
				   else # not in cache, fetch from web
					   open(url).read
				   end

		# if USE_CACHE and the cache file doesn't exist, write to it
		if ShowRobot.config[:use_cache] and not File.exists? cache_file
			puts "Creating cache entry for [ #{url} ]" if ShowRobot.config[:verbose]

			# create the cache directory if it doesn't exist
			FileUtils.mkdir_p File.dirname(cache_file) if not File.directory? File.dirname(cache_file)

			# write to the cache file
			File.open(cache_file, 'w') { |f| f.write contents }
		end

		# dispatch on requested type
		case type
		when :xml
			XML::Parser.string(contents).parse
		else
			raise "Invalid datatype to fetch: [ #{type.to_s} ]"
		end
	end

	module_function :fetch

end
