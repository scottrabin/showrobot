module ShowRobot

	class TheTVDB < Datasource
		DB_NAME = "The TVDB"

		def match_query
			"http://www.thetvdb.com/api/GetSeries.php?seriesname=#{ShowRobot.url_encode @mediaFile.name_guess}&language=en"
		end

		# Returns a list of series related to the media file
		def series
			super()

			Hash[ShowRobot.fetch(:xml, match_query).find('//Series').collect { |series| [series.find('SeriesName').first.content, series] }]
		end

		def episode
			puts "  Fetching #{@mediaFile.name_guess} from #{DB_NAME} (#{match_query})" if ShowRobot::VERBOSE

		end

	end

	add_datasource :tvdb, TheTVDB
end
