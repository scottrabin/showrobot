module ShowRobot

	class TheMovieDB < MovieDatasource
		DB_NAME     = "The Movie Database"
		DATA_TYPE   = :xml
		API_VERSION = 2.1

		def match_query opts = {}
			# TODO
			opts[:lang] = 'en'
			"http://api.themoviedb.org/#{API_VERSION}/Movie.search/#{opts[:lang] || ShowRobot.config[:lang]}/xml/#{ShowRobot.config[:tmdb_api_key]}/#{ShowRobot.url_encode @mediaFile.name_guess}"
		end

		def movie_list
			super do |xml|
				xml.find('//movie').collect do |movie|
					{
						:title   => movie.find('name').first.content,
						:runtime => nil, # TheMovieDB does not return this data
						:year    => nil_or(:to_i, movie.find('released').first.content[/^\d{4}/])
					}
				end
			end
		end
	end

	add_datasource :tmdb, TheMovieDB
end
