module ShowRobot
	require 'text'

	class MockTV < Datasource

		DB_NAME     = 'Mock Database'
		DATA_TYPE   = :yml

		def match_query
			File.expand_path './spec/mock/tvdatabase.yml'
		end

		def episode_query
			File.expand_path './spec/mock/tvdatabase.yml'
		end

		def series_list
			super do |obj|
				obj.map { |item| { :name => item[:series], :source => item }}.uniq { |item| item[:name] }.sort do |a, b|
					Text::Levenshtein.distance(a[:name], @mediaFile.name_guess) - Text::Levenshtein.distance(b[:name], @mediaFile.name_guess)
				end
			end
		end

		def episode_list
			super do |obj|
				# take out the non-matching series
				obj.select { |item| item[:series] == series[:name] }.collect do |episode|
					{
						:series => episode[:series],
						:title => episode[:title],
						:season => episode[:season],
						:episode => episode[:episode],
						:episode_ct => episode[:combined_ep]
					}
				end
			end
		end

	end

	add_datasource :mocktv, MockTV
end
