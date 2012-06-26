module ShowRobot

	class TheTVDB < Datasource
		DB_NAME		= "The TVDB"
		DATA_TYPE	= :xml
		API_KEY		= 'BA864DEE427E384A' # TODO

		def match_query
			"http://www.thetvdb.com/api/GetSeries.php?seriesname=#{ShowRobot.url_encode @mediaFile.name_guess}&language=en"
		end

		def episode_query
			lang = 'en' # TODO
			"http://www.thetvdb.com/api/#{API_KEY}/series/#{series[:source].find('seriesid').first.content}/all/#{lang}.xml"
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
						:name		=> episode.find('EpisodeName').first.content,
						:season		=> episode.find('SeasonNumber').first.content.to_i,
						:episode	=> episode.find('EpisodeNumber').first.content.to_i,
						:episode_ct	=> episode.find('Combined_episodenumber').first.content.to_i
					}
				end
			end
		end

		# Returns the episode data for the specified episode
		def episode seasonnum, episodenum
			episode_list.find { |ep| ep[:season] == seasonnum and ep[:episode] == episodenum }
		end

	end

	add_datasource :tvdb, TheTVDB
end
