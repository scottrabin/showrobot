module ShowRobot

	class TheTVDB < Datasource
		DB_NAME		= "The TVDB"
		DATA_TYPE	= :xml

		def match_query
			"http://www.thetvdb.com/api/GetSeries.php?seriesname=#{ShowRobot.url_encode @mediaFile.name_guess}&language=en"
		end

		# Returns a list of series related to the media file
		def series
			super do |xml|
				xml.find('//Series').collect { |series| {:name => series.find('SeriesName').first.content, :source => series} }
			end
		end

		def episode

		end

	end

	add_datasource :tvdb, TheTVDB
end
