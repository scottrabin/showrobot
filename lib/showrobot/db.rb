module ShowRobot

	# generic datasources (movies and tv shows)
	class Datasource
		# both will need access to a media file
		attr_accessor :mediaFile
	end

	# TV datasources
	class TVDatasource < Datasource
		# returns a list of series matching the given file
		def series_list
			puts "Fetching series data for [ #{@mediaFile.name_guess} ] from #{self.class::DB_NAME} (#{match_query})" if ShowRobot.config[:verbose] and @series_list.nil?

			@series_list ||= yield ShowRobot.fetch(self.class::DATA_TYPE, match_query)
		end

		# returns a list of episodes matching a given series
		def episode_list
			puts "Fetching episode data for [ #{series[:name]} ] from #{self.class::DB_NAME} (#{episode_query})" if ShowRobot.config[:verbose] and @episode_list.nil?

			@episode_list ||= yield ShowRobot.fetch(self.class::DATA_TYPE, episode_query)
		end

		# series accessor methods
		def series
			@series ||= series_list.first
		end
		def series= val
			@series = series_list.find { |obj| obj[:name] == val }
		end

		# Returns the episode data for the specified episode
		def episode(seasonnum = @mediaFile.season, episodenum = @mediaFile.episode)
			episode_list.find { |ep| ep[:season] == seasonnum and ep[:episode] == episodenum }
		end
	end

	# Movie datasources
	class MovieDatasource < Datasource
		def movie_list
			puts "Fetching list of movies matching [ #{@mediaFile.name_guess} (#{@mediaFile.year}) ] from #{self.class::DB_NAME} (#{match_query})" if ShowRobot.config[:verbose] and @movie_list.nil?

			@movie_list ||= yield ShowRobot.fetch(self.class::DATA_TYPE, match_query)
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
