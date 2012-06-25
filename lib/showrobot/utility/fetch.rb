module ShowRobot

	# for processing xml data
	require 'xml'
	# for fetching data
	require 'open-uri'

	# Fetches the given site and processes it as the given type
	def fetch type, url
		case type
		when :xml
			XML::Parser.string(fetch_cache(url)).parse
		else
			raise "Invalid datatype to fetch: [ #{type.to_s} ]"
		end
	end

	def fetch_cache url
		cache_file = File.join(ShowRobot::CACHE_DIRECTORY, url.gsub(/([^a-zA-Z0-9_\.-]+)/) { '%' + $1.unpack('H2' * $1.size).join('%').upcase }.tr(' ', '+'))
		in_cache = ShowRobot::USE_CACHE && File.exists?(cache_file)

		if in_cache
			puts "Found cache entry for [ #{url} ]" if ShowRobot.verbose
			File.read cache_file
		else
			contents = open(url).read

			if ShowRobot::USE_CACHE
				puts "Created cache entry for [ #{url} ] at [ #{cache_file} ]" if ShowRobot.verbose
				File.open(cache_file, 'w') { |f| f.write(contents) }
			end
			contents
		end
	end
	module_function :fetch, :fetch_cache

end
