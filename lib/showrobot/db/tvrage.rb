module ShowRobot

	class TVRage < Datasource
		DB_NAME = 'TVRage'

		def match_query
			"http://services.tvrage.com/feeds/full_search.php?show=#{ShowRobot.url_encode @mediaFile.name_guess}"
		end

		def episode_query
			"http://services.tvrage.com/feeds/full_show_info.php?sid=#{show_id}"
		end

		def fetch
			puts "  Fetching #{@mediaFile.name_guess} from #{DB_NAME} (#{match_query})" if ShowRobot::VERBOSE

			doc = XML::Parser.string(open(match_query).read).parse

			# TODO - make sure this matches
			bestMatch = doc.find('//show').first
			@showName = bestMatch.find('name').first.content
			@showId = bestMatch.find('showid').first.content

			puts "    --> Query: [ #{@mediaFile.name_guess} ] Match: [ #{@showName} ]@#{@showId}" if ShowRobot::VERBOSE
		end

		def episode
			puts "  Fetching episode list for #{show_name}@#{show_id} (#{episode_query})" if ShowRobot::VERBOSE

			puts fetch(:xml, episode_query).find("//episode[seasonnum/.=#{@mediaFile.season} and epnum/.=#{@mediaFile.episode}]")
		end

		def show_name
			if @showName.nil?
				fetch
			end
			@showName
		end

		def show_id
			if @showId.nil?
				fetch
			end
			@showId
		end
	end

	add_datasource :tvrage, TVRage
end
