module ShowRobot

	EPISODE_PATTERNS = [
		# S01E02
		/s(?<season>\d{1,4}).?(?<episode>\d{1,4})/i,
		# season01episode02
		/season\s*(?<season>\d{1,4})\sepisode\s*(?<episode>\d{1,4})/i,
		# sxe
		/\D(?<season>\d{1,4})x(?<episode>\d{1,4})\D/
	]

	YEAR_PATTERN = /[\[\(]?(\d{4})[\]\)]?/

	# parse the name of a file into best-guess parts
	def parse_filename fileName
		ext  = File.extname(fileName).sub /^\./, ''
		file = File.basename(fileName, '.' + ext)

		# try parsing as a tv show - needs to have season/episode information
		if not EPISODE_PATTERNS.find { |pattern| pattern.match file }.nil?
			{
				:name_guess => file[0, file.index($&)],
				:year		=> ShowRobot.file_get_year(fileName),
				:season		=> $~['season'].to_i,
				:episode	=> $~['episode'].to_i,
				:type		=> :tv,
				:extension  => ext.to_sym
			}
		else
			# probably a movie
			{
				:name_guess => file[0, file.index(YEAR_PATTERN) || file.length],
				:year       => nil_or(:to_i, fileName[YEAR_PATTERN, 1]),
				:type       => :movie,
				:extension  => ext.to_sym
			}
		end

	end

	module_function :parse_filename

end
