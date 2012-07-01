require 'helper'
require 'mock/moviesource'

describe ShowRobot, 'movie datasource API' do

	before :each do
		@file = ShowRobot::MediaFile.load 'First Movie.avi'
		@db   = ShowRobot.create_datasource :mockmovie
		@db.mediaFile = @file
	end

	describe 'when querying for a list of matches' do
		it 'should return a list of movies' do
			@db.movie_list.should have(12).items
			@db.movie_list.each do |movie|
				movie[:title].should be_an_instance_of(String)
				movie[:year].should be_a_kind_of(Integer)
				movie[:runtime].should be_a_kind_of(Integer)
			end
		end
	end

	describe 'when prioritizing the returned list' do
		describe 'when there is no runtime or year data' do
			it 'should order them by word distance to title' do
				last_distance = 0
				@db.movie_list.each do |movie|
					distance = word_distance movie[:title], @file.name_guess
					last_distance.should <= distance
					last_distance = distance
				end
			end
		end

		describe 'when there is year data but no runtime data' do
			it 'should filter the movie list by matching year' do
				@file.year = 2002
				@db.movie_list.should have(2).items
				@db.movie_list.each do |movie|
					movie[:year].should eq(2002)
				end
			end
		end

		describe 'when there is runtime data but no year data' do
			it 'should order the movie list by word distance to title them by closest runtime' do
				@file.runtime = 65
				@db.movie_list.should have(12).items
				last_distance, last_runtime_diff = 0, 0
				@db.movie_list.each do |movie|
					distance = word_distance movie[:title], @file.name_guess
					last_distance.should <= distance
					if last_distance == distance
						last_runtime_diff.should <= (movie[:runtime] - @file.runtime)
						last_runtime_diff = movie[:runtime] - @file.runtime
					else
						last_runtime_diff = 0
					end
					last_distance = distance
				end
			end
		end

		describe 'when there is runtime and year data' do
			it 'should filter the list by year and sort by word distance then by closest runtime' do
				@file.year = 2003
				@file.runtime = 87
				# verify filtering
				@db.movie_list.should have(2).items
				last_distance, last_runtime_diff = 0, 0
				@db.movie_list.each do |movie|
					# verify filtering, part 2
					movie[:year].should eq(2003)
					# verify distance
					distance = word_distance movie[:title], @file.name_guess
					last_distance.should <= distance
					if last_distance == distance
						last_runtime_diff.should <= (movie[:runtime] - @file.runtime)
						last_runtime_diff = movie[:runtime] - @file.runtime
					else
						last_runtime_diff = 0
					end
					last_distance = distance
				end
			end
		end
	end
end
