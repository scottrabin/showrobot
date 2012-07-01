module ShowRobot

	class Datasource
		def initialize
		end

		# returns a list of series matching the given file
		def series_list
			puts "Fetching series data for [ #{@mediaFile.name_guess} ] from #{self.class::DB_NAME} (#{match_query})" if ShowRobot.config[:verbose] and @series_list.nil?

			@series_list ||= yield ShowRobot.fetch(self.class::DATA_TYPE, match_query)
		end

		def episode_list
			puts "Fetching episode data for [ #{series[:name]} ] from #{self.class::DB_NAME} (#{episode_query})" if ShowRobot.config[:verbose] and @episode_list.nil?

			@episode_list ||= yield ShowRobot.fetch(self.class::DATA_TYPE, episode_query)
		end

		attr_accessor :mediaFile

		attr_writer :series
		def series
			@series ||= series_list.first
		end

		# Returns the episode data for the specified episode
		def episode(seasonnum = @mediaFile.season, episodenum = @mediaFile.episode)
			episode_list.find { |ep| ep[:season] == seasonnum and ep[:episode] == episodenum }
		end
	end

	class << self
		DATASOURCES = {}
		def add_datasource sym, klass
			DATASOURCES[sym] = klass
		end

		def create_datasource sym
			datasource_for(sym).new
		end

		def datasource_for sym
			DATASOURCES[sym.to_sym]
		end

		def url_encode(s)
			s.to_s.gsub(/[^a-zA-Z0-9_\-.]/n){ sprintf("%%%02X", $&.unpack("C")[0]) }
		end
	end

end
