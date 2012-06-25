module ShowRobot

	class Datasource
		def initialize fileName
			@mediaFile = fileName
		end

		# returns a list of series matching the given file
		def series
			puts "Fetching series data for [ #{@mediaFile.name_guess} ] from #{self.class::DB_NAME} (#{match_query})" if ShowRobot::VERBOSE and @series.nil?

			@series ||= Hash[(yield ShowRobot.fetch(self.class::DATA_TYPE, match_query)).collect { |hash| [hash[:name], hash[:source]] }]
		end

		def episodes
			puts "Fetching #{@mediaFile.name_guess} from #{self.class::DB_NAME} (#{match_query})" if ShowRobot::VERBOSE
		end
	end

	class << self
		DATASOURCES = {}
		def add_datasource sym, klass
			DATASOURCES[sym] = klass
		end

		def datasource_for sym
			DATASOURCES[sym.to_sym]
		end

		def url_encode(s)
			s.to_s.gsub(/[^a-zA-Z0-9_\-.]/n){ sprintf("%%%02X", $&.unpack("C")[0]) }
		end
	end

end
