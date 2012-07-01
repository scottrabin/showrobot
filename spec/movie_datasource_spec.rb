require 'helper'
require 'mock/moviesource'

describe ShowRobot, 'movie datasource API' do

	before :each do
		@file = ShowRobot::MediaFile.load 'Movie(2002).avi'
		@db   = ShowRobot.create_datasource :mockmovie
		@db.mediaFile = @file
	end

	describe 'when querying for a list of matches' do
		it 'should return a list of movies' do
			@db.movie_list.should have(8).items
			@db.movie_list.each do |movie|
				movie[:name].should be_an_instance_of(String)
				movie[:year].should be_a_kind_of(Integer)
				movie[:runtime].should be_a_kind_of(Integer)
			end
		end
	end
end
