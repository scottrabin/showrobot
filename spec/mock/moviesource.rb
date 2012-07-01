module ShowRobot
	require 'text'

	class MockMovie < Datasource

		DB_NAME     = 'Mock Movie Database'
		DATA_TYPE   = :yml

		def match_query
			File.expand_path './spec/mock/moviedatabase.yml'
		end

		def match_list

		end
	end

	add_datasource :mockmovie, MockMovie
end
