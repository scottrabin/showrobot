module ShowRobot

	class TheTVDB < Datasource
		DB_NAME		= "The TVDB"
		DATA_TYPE	= :xml

		def match_query
			"http://www.thetvdb.com/api/GetSeries.php?seriesname=#{ShowRobot.url_encode @mediaFile.name_guess}&language=en"
		end

		def episode_query
			"http://www.thetvdb.com/?tab=seasonall&id=#{series[:source].find('seriesid').first.content}"
		end

		# Returns a list of series related to the media file
		def series_list
			super do |xml|
				xml.find('//Series').collect { |series| {:name => series.find('SeriesName').first.content, :source => series} }
			end
		end

		# Returns a list of episodes related to the media file from a given series
		def episode_list
			super do |xml|

			end
		end

		# Returns the episode data for the specified episode
		def episode seasonnum, episodenum

		end

	end

	add_datasource :tvdb, TheTVDB
end
