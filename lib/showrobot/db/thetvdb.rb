module ShowRobot

	class TheTVDB < Datasource
		DB_NAME		= "The TVDB"
		DATA_TYPE	= :xml

		def match_query
			"http://www.thetvdb.com/api/GetSeries.php?seriesname=#{ShowRobot.url_encode @mediaFile.name_guess}&language=en"
		end

		def episode_query
			lang = 'en' # TODO
			"http://www.thetvdb.com/api/#{ShowRobot.config[:tvdb_api_key]}/series/#{series[:source].find('seriesid').first.content}/all/#{lang}.xml"
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
				xml.find('//Episode').collect do |episode|
					{
						:series		=> series[:name],
						:title		=> episode.find('EpisodeName').first.content,
						:season		=> episode.find('SeasonNumber').first.content.to_i,
						:episode	=> episode.find('EpisodeNumber').first.content.to_i,
						:episode_ct	=> episode.find('Combined_episodenumber').first.content.to_i
					}
				end
			end
		end

	end

	add_datasource :tvdb, TheTVDB
end
