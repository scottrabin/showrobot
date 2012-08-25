require 'helper'
require 'mock/moviesource'

describe ShowRobot, 'movie datasource API' do

	MOVIE = {
		:filename   => 'Fourth Movie: The Subtitle.(2003).avi',
		:title      => 'Fourth Movie: The Subtitle',
		:runtime    => 92,
		:year       => '2003',
		:extension  => 'avi'
	}

	before :each do
		@file = ShowRobot::MediaFile.load MOVIE[:filename]
		@db   = ShowRobot.create_datasource :mockmovie
		@db.mediaFile = @file
	end

	describe 'when querying for a list of matches' do
		it 'should return a list of movies' do
			# verify the movie list is an Array
			@db.movie_list.should be_a_kind_of(Array)
			# verify that all items have the right properties
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
			it 'should order the movie list by matching year' do
				@file.year = 2002
				@db.movie_list.should have(12).items
				cutoff = false
				@db.movie_list.each do |movie|
					# only cut off after the year stops being equal
					cutoff = true if not cutoff and movie[:year] != 2002
					if cutoff
						movie[:year].should_not == 2002
					end
				end
			end
		end

		describe 'when there is runtime data but no year data' do
			it 'should order the movie list by word distance to title them by closest runtime' do
				# year gets automatically read from the parsed value, so it needs to be force-overridden
				class << @file
					def year; nil; end
					def runtime; 65; end
				end
				@db.movie_list.should have(12).items
				last_distance, last_runtime_diff = 0, 0
				@db.movie_list.each do |movie|
					distance = word_distance movie[:title], @file.name_guess
					last_distance.should <= distance
					if last_distance == distance
						runtime_diff = (movie[:runtime] - @file.runtime).abs
						last_runtime_diff.should <= runtime_diff
						last_runtime_diff = runtime_diff
					else
						last_runtime_diff = 0
					end
					last_distance = distance
				end
			end
		end

		describe 'when there is runtime and year data' do
			it 'should sort the list by matching year, then by word distance, then by closest runtime' do
				@file.year = 2003
				@file.runtime = 87
				# verify filtering
				@db.movie_list.should have(12).items
				last_distance, last_runtime_diff = 0, 0
				cutoff = false
				@db.movie_list.each do |movie|
					if not cutoff and movie[:year] != 2003
						cutoff = true
						# reset distance
						last_distance = 0
					end
					# verify filtering, part 2
					movie[:year].should_not == 2003 if cutoff
					# verify distance
					distance = word_distance movie[:title], @file.name_guess
					last_distance.should <= distance
					if last_distance == distance
						last_runtime_diff.should <= (movie[:runtime] - @file.runtime).abs
						last_runtime_diff = (movie[:runtime] - @file.runtime).abs
					else
						last_runtime_diff = 0
					end
					last_distance = distance
				end
			end
		end
	end
end
