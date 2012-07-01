module ShowRobot

	class MockTV < TVDatasource

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
				obj.map do |item|
					{ :name => item[:series], :source => item }
				end.uniq do |item|
					item[:name]
				end.sort &by_distance(@mediaFile.name_guess, :name)
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
