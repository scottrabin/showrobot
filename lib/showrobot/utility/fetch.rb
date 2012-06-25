module ShowRobot

	# for processing xml data
	require 'xml'
	# for fetching data
	require 'open-uri'

	# Fetches the given site and processes it as the given type
	def fetch type, url
		case type
		when :xml
			XML::Parser.string(open(url).read).parse
		else
			raise "Invalid datatype to fetch: [ #{type.to_s} ]"
		end
	end
	module_function :fetch

end
