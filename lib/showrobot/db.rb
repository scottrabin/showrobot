module ShowRobot

	class Datasource
		def initialize fileName
			@mediaFile = fileName
		end

		# returns a list of series matching the given file
		def series_list
			puts "Fetching series data for [ #{@mediaFile.name_guess} ] from #{self.class::DB_NAME} (#{match_query})" if ShowRobot::VERBOSE and @series_list.nil?

			@series_list ||= yield ShowRobot.fetch(self.class::DATA_TYPE, match_query)
		end

		def episode_list
			series = series_list.first if series.nil?

			puts "Fetching episode data for [ #{series[:name]} ] from #{self.class::DB_NAME} (#{episode_query})" if ShowRobot::VERBOSE and @episode_list.nil?

			@episode_list ||= yield ShowRobot.fetch(self.class::DATA_TYPE, episode_query)
		end

		attr_writer :series
		def series
			@series ||= series_list.first
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
