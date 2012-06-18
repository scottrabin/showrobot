module ShowRobot

	# determines the season and episode from the file name
	def get_season_episode fileName
		case fileName
		when /s(\d{1,4}).?e(\d{1,4})/i
			[$1, $2, fileName.index($&), $&.length]
		when /season\s*(\d{1,4})\s*episode\s*(\d{1,4})/i
			[$1, $2, fileName.index($&), $&.length]
		end
	end
	module_function :get_season_episode

end
