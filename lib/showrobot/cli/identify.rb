# helper function for identifying a file
def identify media, database, prompt
	if media.is_movie?

	else
		db = ShowRobot.create_datasource database
		db.mediaFile = media
		
		db.series = if prompt
						select db.series_list, "Select a series for [ #{media.fileName} ]" do |i, item|
							sprintf " %3d) %s", i, item[:name]
						end
					else
						db.series_list.first
					end

		if prompt
			select db.episode_list, "Select an episode for [ #{media.fileName} ]" do |i, item|
				sprintf " %3d) [%02d.%02d] %s", i, item[:season], item[:episode], item[:title]
			end
		else
			db.episode media.season, media.episode
		end
	end
end
