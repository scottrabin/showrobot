module ShowRobot

	EPISODE_PATTERNS = [
		# S01E02
		/s(?<season>\d{1,4}).?(?<episode>\d{1,4})/i,
		# season01episode02
		/season\s*(?<season>\d{1,4})\sepisode\s*(?<episode>\d{1,4})/i
	]

	# parse the name of a file into best-guess parts
	def parse_filename fileName

		# try parsing as a tv show - needs to have season/episode information
		if not EPISODE_PATTERNS.find { |pattern| pattern.match fileName }.nil?
			{
				:name_guess => fileName[0, fileName.index($&)],
				:year		=> ShowRobot.file_get_year(fileName),
				:season		=> $~['season'],
				:episode	=> $~['episode']
			}
		else
			# probably a movie
		end

	end

	module_function :parse_filename

end
