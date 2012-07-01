module ShowRobot

	class MockMovie < MovieDatasource

		DB_NAME     = 'Mock Movie Database'
		DATA_TYPE   = :yml

		def match_query
			File.expand_path './spec/mock/moviedatabase.yml'
		end

		def movie_list
			super do |obj|
				obj
			end
		end
	end

	add_datasource :mockmovie, MockMovie
end
