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
				list = []
				obj.each do |item|
					if list.index { |i| i[:name] == item[:series] }.nil?
						list << Hash[{ :name => item[:series], :source => item }]
					end
				end
				list.sort do |a, b|
					distance = Text::Levenshtein.distance(a[:name], @mediaFile.name_guess) - Text::Levenshtein.distance(b[:name], @mediaFile.name_guess)
					case
					when distance < 0
						-1
					when distance > 0
						1
					else
						0
					end
				end
			end
		end

		def episode_list
			super do |obj|
				# take out the non-matching series
				obj.select do |item|
					item[:series] == series[:name]
				end.collect do |episode|
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
