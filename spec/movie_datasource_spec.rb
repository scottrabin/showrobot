require 'helper'
require 'mock/moviesource'
require 'text'

describe ShowRobot, 'movie datasource API' do

	before :each do
		@file = ShowRobot::MediaFile.load 'First Movie(2002).avi'
		@db   = ShowRobot.create_datasource :mockmovie
		@db.mediaFile = @file
	end

	describe 'when querying for a list of matches' do
		it 'should return a list of movies' do
			@db.movie_list.should have(8).items
			@db.movie_list.each do |movie|
				movie[:title].should be_an_instance_of(String)
				movie[:year].should be_a_kind_of(Integer)
				movie[:runtime].should be_a_kind_of(Integer)
			end
		end
	end

	describe 'when prioritizing the returned list' do
		describe 'when there is no runtime data' do
			it 'should order them by word distance to title' do
				last_distance = 0
				@db.movie_list.each do |movie|
					distance = Text::Levenshtein.distance(movie[:title], @file.name_guess)
					last_distance.should <= distance
					last_distance = distance
				end
			end
		end
	end
end
